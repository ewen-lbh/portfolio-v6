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

func main() {
	db, err := LoadDatabase("../database/database.json")
	if err != nil {
		panic(err)
	}
	files, err := ioutil.ReadDir("../src")
	if err != nil {
		panic(err)
	}
	for _, lang := range []string{"fr", "en"} {
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".pug") {
				templateString, err := ReadFile("../src/" + file.Name())
				if file.Name() == "_work.pug" {
					if err != nil {
						panic(err)
					}
					for _, work := range db.Works {
						println("[hydrator]     hydrating: '" + file.Name() + "' @ " + work.ID)
						HydrateWorkFileWithLang(templateString, "src/"+file.Name(), db, work, lang)
					}
				} else {
					println("[hydrator]     hydrating: '" + file.Name() + "'")
					currentPath := strings.TrimSuffix(file.Name(), ".pug")
					if currentPath == "index" {
						currentPath = ""
					}
					replaced, err := ExecuteTemplate(string(templateString), "src/"+file.Name(), db, WorkOneLang{}, "/"+currentPath, lang)
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

func HydrateWorkFileWithLang(templateString []byte, templateName string, db Database, work Work, lang string) {
	replaced, err := ExecuteTemplate(string(templateString), templateName, db, GetOneLang(lang, work)[0], "/"+work.ID, lang)
	if err != nil {
		panic(err)
	}
	os.MkdirAll("../artifacts/phase_1/"+lang, 0o0666)
	file, err := os.Create("../artifacts/phase_1/" + lang + "/" + work.ID + ".pug")
	if err != nil {
		panic(err)
	}
	_, err = file.WriteString(replaced)
	if err != nil {
		panic(err)
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
	parsedDate, err := iso8601.ParseString(
		strings.ReplaceAll(
			strings.Replace(createdDate, "????", "9999", 1), "?", "1",
		),
	)
	if err != nil {
		fmt.Printf("Error while parsing creation date of %v:\n", work.ID)
		panic(err)
	}
	return parsedDate
}

func (db *Database) LatestWork() Work {
	latest := db.Works[0]
	for _, work := range db.Works {
		if work.Created().Year() == 9999 {
			continue
		}
		if work.Created().After(latest.Created()) {
			latest = work
		}
	}
	return latest
}

func GetOneLang(lang string, works ...Work) []WorkOneLang {
	result := make([]WorkOneLang, 0, len(works))
	for _, work := range works {
		var title string
		paragraphs := make([]Paragraph, 0)
		media := make([]Media, 0)
		links := make([]Link, 0)
		footnotes := make([]Footnote, 0)
		if len(work.Title[lang]) > 0 {
			title = work.Title[lang]
		} else {
			title = work.Title["default"]
		}
		if len(work.Paragraphs[lang]) > 0 {
			paragraphs = work.Paragraphs[lang]
		} else {
			paragraphs = work.Paragraphs["default"]
		}
		if len(work.Media[lang]) > 0 {
			media = work.Media[lang]
		} else {
			media = work.Media["default"]
		}
		if len(work.Links[lang]) > 0 {
			links = work.Links[lang]
		} else {
			links = work.Links["default"]
		}
		if len(work.Footnotes[lang]) > 0 {
			footnotes = work.Footnotes[lang]
		} else {
			footnotes = work.Footnotes["default"]
		}
		result = append(result, WorkOneLang{
			ID:         work.ID,
			Metadata:   work.Metadata,
			Title:      title,
			Paragraphs: paragraphs,
			Media:      media,
			Links:      links,
			Footnotes:  footnotes,
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

func ExecuteTemplate(templateString string, templateName string, db Database, currentWork WorkOneLang, currentPath string, lang string) (string, error) {
	tmpl := template.Must(template.New(templateName).Funcs(template.FuncMap{
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
			panic("cannot find tag with name " + tagName + ", look at /mnt/d/projects/portfolio-next/hydrator/tags.go")
		},
	}).Funcs(sprig.TxtFuncMap()).Funcs(template.FuncMap{
		"tindent":  IndentWithTabs,
		"tnindent": IndentWithTabsNewline,
	}).Parse(templateString))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, TemplateData{
		Tags:                   tags,
		LatestWork:             GetOneLang(lang, db.LatestWork())[0],
		SortingOptions:         []string{"date"}, //TODO: more sorting options
		WorksByYear:            db.WorksByYearOneLang(lang),
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
