package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	jsoniter "github.com/json-iterator/go"

	"github.com/Masterminds/sprig"
	"github.com/joho/godotenv"
)

func main() {
	db, err := LoadDatabase("../database/database.json")
	if err != nil {
		panic(err)
	}
	files, err := ioutil.ReadDir("../src")
	if err != nil {
		panic(err)
	}
	fmt.Println("------")
	if BuildingForProduction() {
		fmt.Println("Hydrating for production")
	} else {
		fmt.Println("Hydrating for developement")
	}
	fmt.Println("------")
	os.Mkdir("../artifacts", 0777) // TODO: 0777 is evil
	os.Mkdir("../artifacts/phase_1", 0777)
	for _, lang := range []string{"fr", "en"} {
		println("[hydrator]     language: " + lang)
		os.Mkdir("../artifacts/phase_1/"+lang, 0777)
		for _, file := range files {
			if file.IsDir() && file.Name() == "using" {
				files, err := ioutil.ReadDir("../src/using")
				if err != nil {
					panic(err)
				}
				os.Mkdir("../artifacts/phase_1/"+lang+"/using", 0777)
				for _, file := range files {
					if file.Name() == "_technology.pug" {
						templateString, err := ReadFile("../src/using/_technology.pug")
						if err != nil {
							panic(err)
						}
						for _, tech := range AvailableTechnologies {
							println("[hydrator]     hydrating: 'using/_technology.pug' @ " + tech.URLName)
							HydrateDynamicFileWithLang(HydrationArgs{
								templateString: templateString,
								db:             db,
								lang:           lang,
								tech:           tech,
								templateName:   "using/_technology.pug",
							})
						}
					}
				}
			} else if strings.HasSuffix(file.Name(), ".pug") {
				templateString, err := ReadFile("../src/" + file.Name())
				if file.Name() == "_work.pug" {
					if err != nil {
						panic(err)
					}
					for _, work := range db.Works {
						println("[hydrator]     hydrating: '" + file.Name() + "' @ " + work.ID)
						HydrateDynamicFileWithLang(HydrationArgs{
							templateString: templateString,
							db:             db,
							lang:           lang,
							work:           work,
							templateName:   file.Name(),
						})
					}
				} else if file.Name() == "_tag.pug" {
					templateString, err := ReadFile("../src/" + file.Name())
					if err != nil {
						panic(err)
					}
					for _, tag := range tags {
						println("[hydrator]     hydrating: '" + file.Name() + "' @ " + tag.URLName)
						HydrateDynamicFileWithLang(HydrationArgs{
							templateString: templateString,
							db:             db,
							lang:           lang,
							tag:            tag,
							templateName:   file.Name(),
						})
					}
				} else {
					println("[hydrator]     hydrating: '" + file.Name() + "'")
					currentPath := strings.TrimSuffix(file.Name(), ".pug")
					if currentPath == "index" {
						currentPath = ""
					}
					// replaced, err := ExecuteTemplate(string(templateString), "src/"+file.Name(), db, WorkOneLang{}, "/"+currentPath, Tag{}, lang)
					replaced, err := ExecuteTemplate(ExecuteTemplateArgs{
						templateString: string(templateString),
						currentPath:    "/" + currentPath,
						db:             db,
						templateName:   file.Name(),
						lang:           lang,
					})
					if err != nil {
						panic(err)
					}
					file, err := os.Create("../artifacts/phase_1/" + lang + "/" + file.Name())
					if err != nil {
						panic(err)
					}
					_, err = file.WriteString(replaced)
					if err != nil {
						panic(err)
					}
				}

			}
		}
	}
}

type HydrationArgs struct {
	templateString []byte
	templateName   string
	db             Database
	work           Work
	tag            Tag
	tech           Technology
	lang           string
}

func HydrateDynamicFileWithLang(args HydrationArgs) {
	// replaced, err := ExecuteTemplate(string(args.templateString), args.templateName, args.db, GetOneLang(args.lang, args.work)[0], "/"+args.work.ID, args.tag, lang)
	replaced, err := ExecuteTemplate(ExecuteTemplateArgs{
		templateString: string(args.templateString),
		templateName:   args.templateName,
		db:             args.db,
		currentWork:    GetOneLang(args.lang, args.work)[0],
		currentPath:    "/" + args.work.ID,
		currentTag:     args.tag,
		currentTech:    args.tech,
	})
	if err != nil {
		panic(err)
	}
	var pathIdentifier string
	if args.work.ID != "" {
		pathIdentifier = args.work.ID
	} else if args.tag.URLName != "" {
		pathIdentifier = args.tag.URLName
	} else {
		pathIdentifier = "using/" + args.tech.URLName
	}
	file, err := os.Create("../artifacts/phase_1/" + args.lang + "/" + pathIdentifier + ".pug")
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(replaced)
	if err != nil {
		panic(err)
	}
}


func GetAge() uint8 {
	// TODO Do it dynamically
	return 17
}

func BuildingForProduction() bool {
	err := godotenv.Load("./.env")
	if err != nil {
		panic("Could not load the .env file")
	}
	return os.Getenv("ENVIRONMENT") != "dev"
}

// AssetURL returns the full URL for a given asset.
func AssetURL(assetPath string) string {
	var urlScheme string
	if !BuildingForProduction() {
		urlScheme = "file://" + os.Getenv("LOCAL_PROJECTS_DIR") + "portfolio-next/assets/%s"
	} else {
		urlScheme = "https://assets.ewen.works/%s"
	}
	return fmt.Sprintf(urlScheme, assetPath)
}

// MediaURL returns the full URL for a given media.
func MediaURL(mediaPath string) string {
	var urlScheme string
	if !BuildingForProduction() {
		urlScheme = "file://" + os.Getenv("LOCAL_PROJECTS_DIR") + "/%s"
	} else {
		urlScheme = "https://media.ewen.works/%s"
	}
	return fmt.Sprintf(urlScheme, mediaPath)
}

type ExecuteTemplateArgs struct {
	templateString string
	templateName   string
	db             Database
	currentWork    WorkOneLang
	currentPath    string
	currentTag     Tag
	currentTech    Technology
	lang           string
}

func ExecuteTemplate(args ExecuteTemplateArgs) (string, error) {
	tmpl := template.Must(template.New(args.templateName).Funcs(template.FuncMap{
		"summarize": func(s string) string {
			var runesCount = 0
			for index := range s {
				runesCount++
				if runesCount > 150 {
					return s[:index] + "…"
				}
			}
			return s
		},
		"lookupTag": func(tagName string) Tag {
			for _, tag := range tags {
				if tag.DisplayName == tagName {
					return tag
				}
			}
			panic("cannot find tag with name " + tagName + ", look at /home/ewen/projects/portfolio-next/hydrator/tags.go")
		},
		"lookupTech": func(techName string) Technology {
			for _, tech := range AvailableTechnologies {
				if tech.URLName == strings.TrimSpace(strings.ToLower(techName)) {
					return tech
				}
			}
			panic("cannot find technology with name " + techName + ", look at /home/ewen/projects/portfolio-next/hydrator/technologies.go")
		},
		"asset": AssetURL,
		"media": func(mediaPath string) string {
			if strings.HasPrefix(mediaPath, "/") {
				return MediaURL(strings.TrimPrefix(mediaPath, "/"))
			}
			return MediaURL(args.currentWork.ID + "/" + mediaPath)
		},
	}).Funcs(sprig.TxtFuncMap()).Funcs(template.FuncMap{
		"tindent":  IndentWithTabs,
		"tnindent": IndentWithTabsNewline,
	}).Parse(args.templateString))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, TemplateData{
		Tags:                   tags,
		LatestWork:             GetOneLang(args.lang, args.db.LatestWork())[0],
		SortingOptions:         []string{"date"}, //TODO: more sorting options
		WorksByYear:            args.db.WorksByYearOneLang(args.lang),
		CurrentWork:            args.currentWork,
		CurrentWorkBuiltLayout: args.currentWork.BuildLayout(),
		CurrentTag:             args.currentTag,
		CurrentTagWorks:        GetOneLang(args.lang, args.db.WorksOfTag(args.currentTag)...),
		CurrentTech:            args.currentTech,
		CurrentTechWorks:       GetOneLang(args.lang, args.db.WorksMadeWith(args.currentTech)...),
		CurrentPath:            args.currentPath,
		Age:                    GetAge(),
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func LoadDatabase(filename string) (Database, error) {
	var db []Work
	json := jsoniter.ConfigFastest
	SetJSONNamingStrategy(LowerCaseWithUnderscores)
	content, err := ReadFile(filename)
	if err != nil {
		return Database{}, err
	}
	err = json.Unmarshal(content, &db)
	return Database{Works: db}, nil
}

func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return []byte{}, err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}
	return content, nil
}

func WriteToFile(filename string, content []byte) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	return err
}
