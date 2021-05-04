package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/Joker/jade"
	"github.com/Masterminds/sprig"
	"github.com/joho/godotenv"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

// getAbsPath returns the absolute path of basename,
// joining the absolute path of src/ and the given basename
func getAbsPath(basename string) string {
	if strings.HasPrefix(basename, "/") {
		return basename
	}
	absdir, err := filepath.Abs("src")
	if err != nil {
		panic(err)
	}
	return path.Join(absdir, basename)
}

func main() {
	//
	// Preparing dist directory
	//
	err := os.MkdirAll("dist/fr/using", 0777)
	if err != nil {
		printerr("Couldn't create directories for writing", err)
		return
	}
	os.MkdirAll("dist/en/using", 0777)
	//
	// Loading files
	//
	db, err := LoadDatabase("database/database.json")
	if err != nil {
		printerr("Could not load the database", err)
		return
	}
	translations, err := LoadTranslations()
	if err != nil {
		printerr("Couldn't load the translation files", err)
	}
	data := GlobalData{translations, db}
	//
	// Watch mode
	//
	if len(os.Args) >= 2 && os.Args[1] == "watch" {
		StartWatcher(db)
	} else {
		data.BuildAll()

		// Save the updated .po file
		translations.SavePO("i18n/fr.po")

		// Final newline
		println("")
	}
}

// GlobalData holds data that is used throughout the whole build process
type GlobalData struct {
	Translations
	Database
}

// BuildAll builds every page
func (data *GlobalData) BuildAll() (built []string) {
	files, err := ioutil.ReadDir("src")
	if err != nil {
		printerr("Could not read src/", err)
		return
	}

	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), "_") || !strings.HasSuffix(file.Name(), ".pug") || file.Name() == "gallery.pug" {
			continue
		}
		built = append(built, data.BuildRegularPage(file.Name())...)
	}

	// Process the technologies index file
	// FIXME: I have to dot it separately since it's in src/using/
	// and not just src/
	built = append(built, data.BuildRegularPage("using/index.pug")...)
	built = append(built, data.BuildWorkPages()...)
	built = append(built, data.BuildTechPages()...)
	return
}

// BuildTechPages builds all technology pages using using/_technology.pug
func (data *GlobalData) BuildTechPages() (built []string) {
	techTemplate := BuildTemplate(getAbsPath("using/_technology.pug"))
	if techTemplate == "" {
		return
	}
	templ, err := data.ParseTemplate(
		"using/_technology.pug",
		techTemplate,
	)
	if err != nil {
		PrintTemplateErrorMessage("parsing template", "using/_technology.pug", techTemplate, err, 2)
		return
	}
	for _, tech := range KnownTechnologies {
		for _, language := range []string{"fr", "en"} {
			content, err := data.ExecuteTemplate(
				templ,
				language,
				CurrentlyHydrated{tech: tech},
			)
			if err != nil {
				PrintTemplateErrorMessage("executing template", "using/_technology<"+tech.URLName+">", techTemplate, err, 2)
				continue
			}
			content = data.TranslateHydrated(content, language)
			fmt.Printf("\r\033[KTranslated using/%s into %s", tech.URLName, language)
			built = append(built, WriteDistFile("using/"+tech.URLName, content, language))
		}
	}
	return
}

// BuildTagPages builds all tag pages using _tag.py
func (data *GlobalData) BuildTagPages() (built []string) {
	tagTemplate := BuildTemplate(getAbsPath("_tag.pug"))
	if tagTemplate == "" {
		return
	}
	templ, err := data.ParseTemplate(
		"_tag.pug",
		tagTemplate,
	)
	if err != nil {
		PrintTemplateErrorMessage("parsing template", "_tag.pug", tagTemplate, err, 2)
		return
	}

	for _, tag := range KnownTags {
		for _, language := range []string{"fr", "en"} {
			content, err := data.ExecuteTemplate(
				templ,
				language,
				CurrentlyHydrated{tag: tag},
			)
			if err != nil {
				PrintTemplateErrorMessage("executing template", NameOfDynamicTemplate(templ, CurrentlyHydrated{tag: tag}), tagTemplate, err, 2)
				continue
			}
			content = data.TranslateHydrated(content, language)
			fmt.Printf("\r\033[KTranslated %s into %s", tag.URLName(), language)
			built = append(built, WriteDistFile(tag.URLName(), content, language))
		}
	}
	return
}

// BuildWorkPages builds all work pages using _work.pug
func (data *GlobalData) BuildWorkPages() (built []string) {
	workTemplate := BuildTemplate(getAbsPath("_work.pug"))
	if workTemplate == "" {
		return
	}
	templ, err := data.ParseTemplate("_work.pug", workTemplate)
	if err != nil {
		PrintTemplateErrorMessage("parsing template", "_work.pug", workTemplate, err, 2)
		return
	}
	for _, work := range data.Works {
		for _, language := range []string{"fr", "en"} {
			content, err := data.ExecuteTemplate(
				templ,
				language,
				CurrentlyHydrated{work: work.InLanguage(language)},
			)
			if err != nil {
				PrintTemplateErrorMessage("executing template", NameOfDynamicTemplate(templ, CurrentlyHydrated{work: work.InLanguage(language)}), workTemplate, err, 2)
				continue
			}
			content = data.TranslateHydrated(content, language)
			fmt.Printf("\r\033[KTranslated %s into %s", work.ID, language)
			built = append(built, WriteDistFile(work.ID, content, language))
		}
	}
	return
}

// BuildRegularPage builds a given page that isn't dynamic (i.e. does not require object data,
// as opposed to work, tag and tech pages)
func (data *GlobalData) BuildRegularPage(filepath string) (built []string) {
	absFilepath := getAbsPath(filepath)
	templateContent := BuildTemplate(absFilepath)
	if templateContent == "" {
		return
	}
	templ, err := data.ParseTemplate(absFilepath, templateContent)
	if err != nil {
		PrintTemplateErrorMessage("parsing template", absFilepath, templateContent, err, 2)
		return
	}
	for _, language := range []string{"fr", "en"} {
		//
		// Execute the template
		//
		content, err := data.ExecuteTemplate(templ, language, CurrentlyHydrated{})
		fmt.Printf("\r\033[KHydrated %s", absFilepath)
		if err != nil {
			PrintTemplateErrorMessage("executing template", absFilepath, templateContent, err, 2)
			continue
		}
		content = data.TranslateHydrated(content, language)
		fmt.Printf("\r\033[KTranslated %s into %s", GetPathRelativeToSrcDir(absFilepath), language)
		built = append(built, WriteDistFile(GetPathRelativeToSrcDir(absFilepath), content, language))
	}
	return
}

// CurrentlyHydrated represents a Tag, Technology or WorkOneLang
type CurrentlyHydrated struct {
	tag  Tag
	tech Technology
	work WorkOneLang
}

// BuildingForProduction returns true if the environment file declares ENVIRONMENT to not "dev"
func BuildingForProduction() bool {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Could not load the .env file")
	}
	return os.Getenv("ENVIRONMENT") != "dev"
}

// BuildTemplate converts a .pug template to an HTML one
func BuildTemplate(absFilepath string) string {
	raw, err := ioutil.ReadFile(absFilepath)
	if err != nil {
		printerr("Could not read file "+absFilepath, err)
	}

	// Fix `extends` statement
	// From joker/jade's point of view, the current work dir is just the project's root,
	// thus (project root)/layout.pug does not exist.
	// Fix that by adding src/ in front
	// Joker/jade also requires the .pug extension
	extendsPattern := regexp.MustCompile(`(?m)^extends (.+)$`)
	raw = extendsPattern.ReplaceAllFunc(raw, func(line []byte) []byte {
		// printfln("transforming %s", line)
		extendsArgument := strings.TrimPrefix(string(line), "extends ")
		if strings.HasPrefix(extendsArgument, "src/") {
			return line
		}
		return []byte(fmt.Sprintf("extends src/%s.pug", extendsArgument))
	})

	template, err := jade.Parse(absFilepath, raw)
	if err != nil {
		PrintTemplateErrorMessage("converting template to HTML", absFilepath, string(raw), err, 1)
		return ""
	}

	return template
}

// PrintTemplateErrorMessage prints a nice error message with a preview of the code where the error occured
func PrintTemplateErrorMessage(whileDoing string, templateName string, templateContent string, err error, lineNumberSliceIndex int) {
	_lineIndex, intParseError := strconv.ParseInt(strings.Split(err.Error(), ":")[lineNumberSliceIndex], 10, 64)
	if intParseError != nil {
		fmt.Printf("While %s %s: %s", whileDoing, templateName, err.Error())
		return
	}
	lineIndex := int(_lineIndex)
	printfln("While %s %s:%d: %s", whileDoing, templateName, lineIndex, strings.SplitN(err.Error(), ":", lineNumberSliceIndex+1+1)[lineNumberSliceIndex+1])
	lineIndex-- // Lines start at 1, arrays of line are indexed from 0
	lines := strings.Split(gohtml.FormatWithLineNo(templateContent), "\n")
	var lineIndexOffset int
	if len(lines) >= lineIndex+5+1 {
		if lineIndex >= 5 {
			lines = lines[lineIndex-5 : lineIndex+5]
			lineIndexOffset = lineIndex - 5
		} else {
			lines = lines[0 : lineIndex+5]
		}
	}
	for i, line := range lines {
		if i+lineIndexOffset == lineIndex {
			fmt.Print("→ ")
		} else {
			fmt.Print("  ")
		}
		fmt.Println(line)
	}
}

// ParseTemplate parses a given (HTML) template.
func (data *GlobalData) ParseTemplate(templateName string, templateContent string) (*template.Template, error) {
	tmpl := template.New(templateName)
	tmpl = tmpl.Funcs(sprig.TxtFuncMap())
	tmpl, err := tmpl.Parse(gohtml.Format(templateContent))

	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

// NameOfDynamicTemplate returns the name given to a template that is applied to multiple objects, e.g. _work.pug<portfolio>
func NameOfDynamicTemplate(tmpl *template.Template, currentlyHydrated CurrentlyHydrated) string {
	currentlyHydratedName := firstNonEmpty(currentlyHydrated.work.ID, currentlyHydrated.tag.URLName(), currentlyHydrated.tech.URLName)
	if currentlyHydratedName != "" {
		return fmt.Sprintf("%s<%s>", tmpl.Name(), currentlyHydratedName)
	}
	return tmpl.Name()
}

// ExecuteTemplate executes a parsed HTML template to hydrate it with data, potentially with a tag, tech or work.
func (data *GlobalData) ExecuteTemplate(tmpl *template.Template, language string, currentlyHydrated CurrentlyHydrated) (string, error) {
	// Inject Funcs now, since they depend on language
	tmpl = tmpl.Funcs(GetTemplateFuncMap(language, &data.moFile))

	var buf bytes.Buffer

	err := tmpl.Execute(&buf, TemplateData{
		KnownTags:         KnownTags,
		KnownTechnologies: KnownTechnologies,
		Works:             GetOneLang(language, data.Works...),
		Age:               GetAge(),
		CurrentTag:        currentlyHydrated.tag,
		CurrentTech:       currentlyHydrated.tech,
		CurrentWork:       currentlyHydrated.work,
	})

	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// TranslateHydrated translates an hydrated HTML page, removing i18n tags and attributes
// and replacing translatable content with their translations
func (data *Translations) TranslateHydrated(content string, language string) string {
	parsedContent, err := html.Parse(strings.NewReader(content))
	if err != nil {
		printerr("An error occured while parsing the hydrated HTML for translation", err)
		return ""
	}
	return data.TranslateToLanguage(language == "fr", parsedContent)
}

// WriteDistFile writes the given content to the dist/ equivalent of the given fileName and returns that equivalent's path
func WriteDistFile(fileName string, content string, language string) string {
	distFilePath := fmt.Sprintf("dist/%s/%s", language, strings.TrimSuffix(fileName, ".pug")+".html")
	distFile, err := os.Create(distFilePath)
	if err != nil {
		printerr("Could not create the destination file "+distFilePath, err)
		return ""
	}
	_, err = distFile.WriteString(content)
	if err != nil {
		printerr("Could not write to the destination file "+distFilePath, err)
		return ""
	}
	fmt.Printf("\r\033[KWritten %s", distFilePath)
	return distFilePath
}
