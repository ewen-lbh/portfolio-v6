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
	"github.com/snapcore/go-gettext"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

func IsVerbose() bool {
	// return len(os.Args) > 1 && os.Args[1] == "--verbose"
	return false // TODO
}

func VerboseLog(s string, fmtArgs ...interface{}) {
	if IsVerbose() {
		fmt.Printf(s+"\n", fmtArgs...)
	}
}

func main() {
	//
	// Loading files
	//
	db, err := LoadDatabase("database/database.json")
	if err != nil {
		printerr("Could not load the database", err)
		return
	}
	messagesFile, err := os.Open("messages/fr.mo")
	if err != nil {
		printerr("Could not open the .mo file", err)
		return
	}
	messages, err := gettext.ParseMO(messagesFile)
	if err != nil {
		printerr("Could not parse the .mo file", err)
		return
	}
	files, err := ioutil.ReadDir("src")
	getAbsPath := func(basename string) string {
		absdir, err := filepath.Abs("src")
		if err != nil {
			panic(err)
		}
		return path.Join(absdir, basename)
	}
	if err != nil {
		printerr("Could not read src/", err)
		return
	}
	//
	// Preparing dist directory
	//
	err = os.MkdirAll("dist/fr/using", 0777)
	if err != nil {
		printerr("Couldn't create directories for writing", err)
		return
	}
	os.MkdirAll("dist/en/using", 0777)
	//
	// Processing regular pages
	//
	for _, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), "_") || !strings.HasSuffix(file.Name(), ".pug") || file.Name() == "gallery.pug" {
			continue
		}
		absFilepath := getAbsPath(file.Name())
		//
		// Build the template
		//
		templateContent := BuildTemplate(absFilepath)
		for _, language := range []string{"fr", "en"} {
			//
			// Execute the template
			//
			content, err := ExecuteTemplate(db, &messages, language, absFilepath, templateContent, CurrentlyHydrated{})
			if err != nil {
				continue
			}
			content = TranslateHydrated(content, language, &messages)
			fmt.Printf("\r\033[KTranslated %s into %s", file.Name(), language)
			WriteDistFile(file.Name(), content, language, &messages)
		}
	}
	//
	// Processing works
	//
	workTemplate := BuildTemplate(getAbsPath("_work.pug"))
	if workTemplate != "" {

		for _, work := range db.Works {
			for _, language := range []string{"fr", "en"} {
				content, err := ExecuteTemplate(
					db, &messages, language,
					"_work<"+work.ID+">",
					workTemplate,
					CurrentlyHydrated{work: work.InLanguage(language)},
				)
				if err != nil {
					continue
				}
				content = TranslateHydrated(content, language, &messages)
				fmt.Printf("\r\033[KTranslated %s into %s", work.ID, language)
				WriteDistFile(work.ID, content, language, &messages)
			}
		}
	}
	//
	// Processing tags
	//
	tagTemplate := BuildTemplate(getAbsPath("_tag.pug"))
	if tagTemplate != "" {
		for _, tag := range KnownTags {
			for _, language := range []string{"fr", "en"} {
				content, err := ExecuteTemplate(
					db, &messages, language,
					"_tag<"+tag.URLName()+">",
					tagTemplate,
					CurrentlyHydrated{tag: tag},
				)
				if err != nil {
					continue
				}
				content = TranslateHydrated(content, language, &messages)
				fmt.Printf("\r\033[KTranslated %s into %s", tag.URLName(), language)
				WriteDistFile(tag.URLName(), content, language, &messages)
			}
		}
	}
	//
	// Processing technologies
	//
	// Process the index file
	absFilepath := getAbsPath("using/index.pug")
	templateContent := BuildTemplate(absFilepath)
	for _, language := range []string{"fr", "en"} {
		// Execute the template
		content, err := ExecuteTemplate(db, &messages, language, absFilepath, templateContent, CurrentlyHydrated{})
		if err != nil {
			continue
		}
		content = TranslateHydrated(content, language, &messages)
		fmt.Printf("\r\033[KTranslated using/index.pug into %s", language)
		WriteDistFile("using/index.pug", content, language, &messages)
	}
	// Process all the technologies
	techTemplate := BuildTemplate(getAbsPath("using/_technology.pug"))
	if techTemplate != "" {
		for _, tech := range KnownTechnologies {
			for _, language := range []string{"fr", "en"} {
				content, err := ExecuteTemplate(
					db, &messages, language,
					"using/_technology<"+tech.URLName+">",
					techTemplate,
					CurrentlyHydrated{tech: tech},
				)
				if err != nil {
					continue
				}
				content = TranslateHydrated(content, language, &messages)
				fmt.Printf("\r\033[KTranslated using/%s into %s", tech.URLName, language)
				WriteDistFile("using/"+tech.URLName, content, language, &messages)
			}
		}
	}
	// Final newline
	println("")
}

type CurrentlyHydrated struct {
	tag  Tag
	tech Technology
	work WorkOneLang
}

func BuildingForProduction() bool {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Could not load the .env file")
	}
	return os.Getenv("ENVIRONMENT") != "dev"
}

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

func PrintTemplateErrorMessage(whileDoing string, templateName string, templateContent string, err error, lineNumberSliceIndex int) {
	_lineIndex, intParseError := strconv.ParseInt(strings.Split(err.Error(), ":")[lineNumberSliceIndex], 10, 64)
	if intParseError != nil {
		fmt.Printf("While %s %s: %s", whileDoing, templateName, err.Error())
		return
	}
	lineIndex := int(_lineIndex)
	fmt.Printf("While %s %s:%d: %s\n", whileDoing, templateName, lineIndex, strings.SplitN(err.Error(), ":", lineNumberSliceIndex+1+1)[lineNumberSliceIndex+1])
	lineIndex -= 1 // Lines start at 1, arrays of line are indexed from 0
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

func ExecuteTemplate(db Database, catalog *gettext.Catalog, language string, templateName string, templateContent string, currentlyHydrated CurrentlyHydrated) (string, error) {
	tmpl := template.New(templateName)
	tmpl = tmpl.Funcs(GetTemplateFuncMap(language, catalog))
	tmpl = tmpl.Funcs(sprig.TxtFuncMap())
	tmpl, err := tmpl.Parse(gohtml.Format(templateContent))

	if err != nil {
		PrintTemplateErrorMessage("parsing template", templateName, templateContent, err, 2)
		return "", err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, TemplateData{
		KnownTags:         KnownTags,
		KnownTechnologies: KnownTechnologies,
		Works:             GetOneLang(language, db.Works...),
		Age:               GetAge(),
		CurrentTag:        currentlyHydrated.tag,
		CurrentTech:       currentlyHydrated.tech,
		CurrentWork:       currentlyHydrated.work,
	})

	if err != nil {
		PrintTemplateErrorMessage("executing template", templateName, templateContent, err, 2)
		return "", err
	} else {
		fmt.Printf("\r\033[KHydrated %s", templateName)
	}
	return buf.String(), nil
}

func TranslateHydrated(content string, language string, messages *gettext.Catalog) string {
	parsedContent, err := html.Parse(strings.NewReader(content))
	if err != nil {
		printerr("An error occured while parsing the hydrated HTML for translation", err)
		return ""
	}
	return TranslateToLanguage(language == "fr", parsedContent, messages)
}

func WriteDistFile(fileName string, content string, language string, messages *gettext.Catalog) {
	distFilePath := fmt.Sprintf("dist/%s/%s", language, strings.TrimSuffix(fileName, ".pug")+".html")
	distFile, err := os.Create(distFilePath)
	if err != nil {
		printerr("Could not create the destination file "+distFilePath, err)
		return
	}
	_, err = distFile.WriteString(content)
	if err != nil {
		printerr("Could not write to the destination file "+distFilePath, err)
		return
	}
	fmt.Printf("\r\033[KWritten %s", distFilePath)
}
