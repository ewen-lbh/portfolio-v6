package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/Masterminds/sprig"
	"github.com/relvacode/iso8601"
)

var tags = [...]Tag{
	{
		DisplayName: "card",
		URLName:     "cards",
	},
	{
		DisplayName: "cover art",
		URLName:     "cover-arts",
	},
	{
		DisplayName: "game",
		URLName:     "games",
	},
	{
		DisplayName: "graphism",
		URLName:     "graphism",
	},
	{
		DisplayName: "poster",
		URLName:     "posters",
	},
	{
		DisplayName: "automation",
		URLName:     "automation",
	},
	{
		DisplayName: "web",
		URLName:     "web",
	},
	{
		DisplayName: "intro",
		URLName:     "intros",
	},
	{
		DisplayName: "music",
		URLName:     "music",
	},
	{
		DisplayName: "app",
		URLName:     "apps",
	},
	{
		DisplayName: "book",
		URLName:     "books",
	},
	{
		DisplayName: "api",
		URLName:     "APIs",
	},
	{
		DisplayName: "program",
		URLName:     "programs",
	},
	{
		DisplayName: "cli",
		URLName:     "CLIs",
	},
	{
		DisplayName: "motion design",
		URLName:     "motion-design",
	},
	{
		DisplayName: "compositing",
		URLName:     "compositing",
	},
	{
		DisplayName: "visual identity",
		URLName:     "visual-identities",
	},
	{
		DisplayName: "banner",
		URLName:     "banners",
	},
	{
		DisplayName: "illustration",
		URLName:     "illustrations",
	},
	{
		DisplayName: "logo",
		URLName:     "logos",
	},
	{
		DisplayName: "typography",
		URLName:     "typography",
	},
	{
		DisplayName: "drawing",
		URLName:     "drawings",
	},
	{
		DisplayName: "iconography",
		URLName:     "iconography",
	},
	{
		DisplayName: "site",
		URLName:     "sites",
	},
	{
		DisplayName: "language",
		URLName:     "languages",
	},
}

func main() {
	db, err := LoadDatabase("../database/works.json")
	if err != nil {
		panic(err)
	}
	files, err := ioutil.ReadDir("../src")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".pug") {
			templateString, err := ReadFile("../src/" + file.Name())
			if err != nil {
				panic(err)
			}
			if file.Name() == "_work.pug" {
				for _, work := range db.Works {
					println("[hydrator]     hydrating: '" + file.Name() + "' @ " + work.ID)
					replaced, err := ExecuteTemplate(string(templateString), db, GetOneLang("en-US", work)[0], "/"+work.ID)
					if err != nil {
						panic(err)
					}
					file, err := os.Create("../trans/src/" + work.ID + ".pug")
					if err != nil {
						panic(err)
					}
					_, err = file.WriteString(replaced)
					if err != nil {
						panic(err)
					}
				}
			} else {
				println("[hydrator]     hydrating: '" + file.Name() + "'")
				currentPath := strings.TrimSuffix(file.Name(), ".pug")
				if currentPath == "index" {
					currentPath = ""
				}
				replaced, err := ExecuteTemplate(string(templateString), db, WorkOneLang{}, "/"+currentPath)
				if err != nil {
					panic(err)
				}
				file, err := os.Create("../trans/src/" + file.Name())
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

type TemplateData struct {
	Tags                   [len(tags)]Tag
	SortingOptions         []string
	WorksByYear            map[int][]WorkOneLang
	LatestWork             WorkOneLang
	CurrentWork            WorkOneLang
	CurrentWorkBuiltLayout string
	CurrentPath            string
}

type Tag struct {
	URLName     string
	DisplayName string
}

type Database struct {
	Works []Work
}

func (work *Work) Created() time.Time {
	createdDate := work.Metadata.Created
	parsedDate, err := iso8601.ParseString(strings.ReplaceAll(createdDate, "?", "1"))
	if err != nil {
		fmt.Printf("Error while parsing creation date of %v:\n", work.ID)
		panic(err)
	}
	return parsedDate
}

func (db *Database) LatestWork() Work {
	latest := db.Works[0]
	for _, work := range db.Works {
		if work.Created().After(latest.Created()) {
			latest = work
		}
	}
	return latest
}

func GetOneLang(lang string, works ...Work) []WorkOneLang {
	result := make([]WorkOneLang, 0, len(works))
	for _, work := range works {
		result = append(result, WorkOneLang{
			ID:         work.ID,
			Metadata:   work.Metadata,
			Title:      work.Title[lang],
			Paragraphs: work.Paragraphs[lang],
			Media:      work.Media[lang],
			Links:      work.Links[lang],
			Footnotes:  work.Footnotes[lang],
		})
	}
	return result
}

func (db *Database) WorksByYearOneLang(lang string) map[int][]WorkOneLang {
	worksByYear := make(map[int][]WorkOneLang)
	for _, work := range db.Works {
		year := work.Created().Year()
		worksByYear[year] = append(worksByYear[year], GetOneLang(lang, work)[0])
	}
	return worksByYear
}

func ExecuteTemplate(templateString string, db Database, currentWork WorkOneLang, currentPath string) (string, error) {
	tmpl := template.Must(template.New("whatever").Funcs(template.FuncMap{
		"summarize": func(s string) string {
			var runesCount = 0
			for index := range s {
				runesCount++
				if runesCount > 25 {
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
			panic("Cannot find tag with name " + tagName)
		},
	}).Funcs(sprig.TxtFuncMap()).Funcs(template.FuncMap{
		"tindent":  IndentWithTabs,
		"tnindent": IndentWithTabsNewline,
	}).Parse(templateString))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, TemplateData{
		Tags:                   tags,
		LatestWork:             GetOneLang("en-US", db.LatestWork())[0],
		SortingOptions:         []string{"date"}, //TODO: more sorting options
		WorksByYear:            db.WorksByYearOneLang("en-US"),
		CurrentWork:            currentWork,
		CurrentWorkBuiltLayout: currentWork.BuildLayout(),
		CurrentPath:            currentPath,
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
