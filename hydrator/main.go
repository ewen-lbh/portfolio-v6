package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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
	db, err := LoadDatabase("database")
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
		_, err := data.BuildAll("/home/ewen/projects/portfolio/src")

		if err != nil {
			printerr("While building", err)
		}

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

// BuildAll builds all page in the given directory, recursively.
func (data *GlobalData) BuildAll(in string) (built []string, err error) {
	err = filepath.WalkDir(in, func(path string, entry fs.DirEntry, err error) error {
		if strings.Contains(path, "/mixins/") {
			return nil
		}
		if !(strings.HasSuffix(path, ".pug") || strings.HasSuffix(path, ".html")) {
			return nil
		}
		if err != nil {
			return err
		}
		if strings.HasPrefix(entry.Name(), ":work") || strings.HasPrefix(entry.Name(), ":if(work") {
			built = append(built, data.BuildWorkPages(path)...)
		} else if strings.HasPrefix(entry.Name(), ":tag") {
			built = append(built, data.BuildTagPages(path)...)
		} else if strings.HasPrefix(entry.Name(), ":technology") {
			built = append(built, data.BuildTechPages(path)...)
		} else if strings.HasPrefix(entry.Name(), ":") {
			return fmt.Errorf("Dynamic path %s uses unknown variable %s", path, strings.TrimPrefix(entry.Name(), ":"))
		} else {
			built = append(built, data.BuildRegularPage(path)...)
		}
		return err
	})
	return
}

// BuildTechPages builds all technology pages using `using`
func (data *GlobalData) BuildTechPages(using string) (built []string) {
	templateHTML, err := ConvertTemplateIfNeeded(using)
	if err != nil {
		printerr("couldn't convert technology template", err)
		return
	}
	for _, tech := range data.Technologies {
		built = append(built, data.BuildPage(using, templateHTML, &Hydration{tech: tech})...)
	}
	return
}

// BuildTagPages builds all tag pages using the given filename
func (data *GlobalData) BuildTagPages(using string) (built []string) {
	templateHTML, err := ConvertTemplateIfNeeded(using)
	if err != nil {
		printerr("couldn't convert template", err)
		return
	}
	for _, tag := range data.Tags {
		built = append(built, data.BuildPage(using, templateHTML, &Hydration{tag: tag})...)
	}
	return
}

// BuildWorkPages builds all work pages using the given filepath
func (data *GlobalData) BuildWorkPages(using string) (built []string) {
	templateHTML, err := ConvertTemplateIfNeeded(using)
	if err != nil {
		printerr("couldn't build work pages", err)
		return
	}
	for _, work := range data.Works {
		built = append(built, data.BuildPage(using, templateHTML, &Hydration{work: work})...)
	}
	return
}

// BuildRegularPage builds a given page that isn't dynamic (i.e. does not require object data,
// as opposed to work, tag and tech pages)
func (data *GlobalData) BuildRegularPage(filepath string) (built []string) {
	templateContent, err := ConvertTemplateIfNeeded(filepath)
	if err != nil {
		printerr("could not convert the template", err)
		return
	}
	return data.BuildPage(filepath, templateContent, &Hydration{})
}

// BuildPage builds a single page
func (data *GlobalData) BuildPage(using string, templateHTML string, hydration *Hydration) (built []string) {
	for _, language := range []string{"fr", "en"} {
		hydration.language = language
		templ, err := data.ParseTemplate(
			language,
			using,
			templateHTML,
		)
		if err != nil {
			PrintTemplateErrorMessage("parsing template", using, templateHTML, err, 2)
			return
		}
		content, err := data.ExecuteTemplate(
			templ,
			language,
			*hydration,
		)
		if err != nil {
			PrintTemplateErrorMessage("executing template", NameOfTemplate(templ, *hydration), templateHTML, err, 2)
			continue
		}
		content = data.TranslateHydrated(content, language)
		fmt.Printf("\r\033[KTranslated %s into %s", hydration.Name(), language)
		outPath := hydration.GetDistFilepath(using)
		os.MkdirAll(filepath.Dir(outPath), 0777)
		ioutil.WriteFile(outPath, []byte(content), 0777)
		fmt.Printf("\r\033[KRendered %s as %s", using, outPath)
		built = append(built, outPath)
	}
	return
}
