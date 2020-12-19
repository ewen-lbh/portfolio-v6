package main

import (
	"github.com/relvacode/iso8601"
	"strings"
	"fmt"
	"time"
)

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

func (db *Database) WorksOfTag(tag Tag) []Work {
	worksOfTag := make([]Work, 0)
	for _, work := range db.Works {
		_, findError := FindInArrayLax(work.Metadata.Tags, tag.URLName)
		if findError == nil {
			worksOfTag = append(worksOfTag, work)
		}
	}
	return worksOfTag
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

func (db *Database) WorksMadeWith(tech Technology) []Work {
	worksMadeWith := make([]Work, 0)
	for _, work := range db.Works {
		_, findError := FindInArrayLax(work.Metadata.MadeWith, tech.URLName)
		if findError == nil {
			worksMadeWith = append(worksMadeWith, work)
		}
	}
	return worksMadeWith
}
