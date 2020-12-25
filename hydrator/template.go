package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/davecgh/go-spew/spew"
	"github.com/jaytaylor/html2text"
)

// GetTemplateFuncMap returns the funcmap used to hydrate files
func GetTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		// "into" transforms to HTML structures
		"intoGallery": intoGallery,
		"intoLayout":  intoLayout,
		// "get" gets a Go value (string, map[string]string, etc.)
		"getColorsMap":       getColorsMap,
		"getSummary":         getSummary,
		"getThumbnailSource": getThumbnailSource,
		"getYears":           getYears,
		// "with" filters a []WorkOneLang
		"withTag":       withTag,
		"withTech":      withTech,
		"withWIPStatus": withWIPStatus,
		// Nice, cosy aliases for filters
		"withCreatedYear": withCreatedYear,
		"tagged":          withTag,
		"madeWith":        withTech,
		"createdIn":       withCreatedYear,
		"finished": func(ws []WorkOneLang) []WorkOneLang {
			return withWIPStatus(false, ws)
		},
		"unfinished": func(ws []WorkOneLang) []WorkOneLang {
			return withWIPStatus(true, ws)
		},
		// reduces a []WorkOneLang down to a single WorkOneLang
		"latest": latest,
		// functions acting on paths
		"asset": asset,
		"media": media,
		// lookups for tags & technologies
		"lookupTag":  lookupTag,
		"lookupTech": lookupTech,
		// debugging
		"log": log,
	}
}

type TemplateData struct {
	Age               uint8
	KnownTags         [len(KnownTags)]Tag
	KnownTechnologies [len(KnownTechnologies)]Technology
	Works             []WorkOneLang
	MusicTag          Tag
	// Template data for _-prefixed .pug files: relevant struct instance of what's being hydrated
	CurrentTag  Tag
	CurrentTech Technology
	CurrentWork WorkOneLang
}

func GetAge() uint8 {
	// TODO Do it dynamically
	return 17
}

func intoGallery(ws []WorkOneLang) string {
	templateContent, err := ReadFile("../src/.gallery.pug")

	// Hydrate .gallery.pug with ws
	tmpl := template.Must(
		template.New(".gallery.pug").Funcs(GetTemplateFuncMap()).Funcs(sprig.TxtFuncMap()).Funcs(template.FuncMap{
			"tindent":  IndentWithTabs,
			"tnindent": IndentWithTabsNewline,
		}).Parse(string(templateContent)))
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, struct {
		GivenWorks []WorkOneLang
		KnownTags  [len(KnownTags)]Tag
	}{
		GivenWorks: ws,
		KnownTags:  KnownTags,
	})
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func intoLayout(w WorkOneLang) string {
	return w.BuildLayout()
}

// getColorsMap returns a mapping of "primary", "secondary", etc to the color values,
// with an added "#" prefix if needed
func getColorsMap(w WorkOneLang) map[string]string {
	colorsMap := make(map[string]string, 3)
	if w.Metadata.Colors.Primary != "" {
		colorsMap["primary"] = AddOctothorpeIfNeeded(w.Metadata.Colors.Primary)
	}
	if w.Metadata.Colors.Secondary != "" {
		colorsMap["secondary"] = AddOctothorpeIfNeeded(w.Metadata.Colors.Secondary)
	}
	if w.Metadata.Colors.Tertiary != "" {
		colorsMap["tertiary"] = AddOctothorpeIfNeeded(w.Metadata.Colors.Tertiary)
	}
	return colorsMap
}

// getSummary summarizes the given work's first description paragraph
func getSummary(w WorkOneLang) string {
	if len(w.Paragraphs) == 0 {
		return ""
	}
	plainText, err := html2text.FromString(w.Paragraphs[0].Content, html2text.Options{OmitLinks: true})
	if err != nil {
		panic(err)
	}
	return SummarizeString(plainText, 150)
}

// getThumbnailSource gets the source URL of the work's first media
func getThumbnailSource(w WorkOneLang) string {
	if len(w.Media) == 0 {
		return ""
	}
	return media(w.Media[0].Source)
}

func getYears(ws []WorkOneLang) []int {
	years := make([]int, 0)
	for _, work := range ws {
		years = append(years, work.Created().Year())
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years)))
	return years
}

// withTag returns an array of works that have tag in their tags
func withTag(tag Tag, ws []WorkOneLang) []WorkOneLang {
	filtered := make([]WorkOneLang, 0)
	for _, work := range ws {
		_, findError := FindInArrayLax(work.Metadata.Tags, tag.URLName)
		if findError == nil {
			filtered = append(filtered, work)
		}
	}
	if len(filtered) == 0 && IsVerbose() {
		fmt.Printf("WARNING: No works from %v have the %s tag", ws, tag.URLName)
	}
	return filtered
}

// withTag returns an array of works that have tech in their "made with" technologies list
func withTech(tech Technology, ws []WorkOneLang) []WorkOneLang {
	filtered := make([]WorkOneLang, 0)
	for _, work := range ws {
		_, findError := FindInArrayLax(work.Metadata.MadeWith, tech.URLName)
		if findError == nil {
			filtered = append(filtered, work)
		}
	}
	return filtered
}

func withWIPStatus(wipStatus bool, ws []WorkOneLang) []WorkOneLang {
	filtered := make([]WorkOneLang, 0)
	for _, work := range ws {
		if work.IsWIP() == wipStatus {
			filtered = append(filtered, work)
		}
	}
	return filtered
}

func withCreatedYear(createdYear int, ws []WorkOneLang) []WorkOneLang {
	filtered := make([]WorkOneLang, 0)
	for _, work := range ws {
		if work.Created().Year() == createdYear {
			filtered = append(filtered, work)
		}
	}
	return filtered
}

func latest(ws []WorkOneLang) WorkOneLang {
	if len(ws) == 0 {
		panic("cannot get the latest element of an empty array")
	}
	latest := ws[0]
	for _, work := range ws {
		if work.Created().Year() == 65535 {
			continue
		}
		if work.Created().After(latest.Created()) {
			latest = work
		}
	}
	return latest
}

func log(o interface{}) string {
	spew.Dump(o)
	if !BuildingForProduction() {
		return fmt.Sprintf("logged %v to stdout", o)
	}
	return ""
}

// asset returns the full URL for a given asset (ie a website's static asset like an icon)
func asset(assetPath string) string {
	var urlScheme string
	if !BuildingForProduction() {
		urlScheme = "file://" + os.Getenv("LOCAL_PROJECTS_DIR") + "portfolio/assets/%s"
	} else {
		urlScheme = "https://assets.ewen.works/%s"
	}
	return fmt.Sprintf(urlScheme, assetPath)
}

// media returns the full URL for a given media (ie a work's media URL)
func media(mediaPath string) string {
	var urlScheme string
	if !BuildingForProduction() {
		urlScheme = "file://" + os.Getenv("LOCAL_PROJECTS_DIR") + "/%s"
	} else {
		urlScheme = "https://media.ewen.works/%s"
	}
	return fmt.Sprintf(urlScheme, mediaPath)
}

// lookupTag returns the tag with DisplayName name
func lookupTag(name string) Tag {
	for _, tag := range KnownTags {
		if StringsLooselyMatch(tag.DisplayName, name) {
			return tag
		}
	}
	panic("cannot find tag with display name " + name + ", look at /home/ewen/projects/portfolio/hydrator/tags.go")
}

// lookupTech returns the tech with DisplayName name
func lookupTech(name string) Technology {
	for _, tech := range KnownTechnologies {
		if StringsLooselyMatch(tech.DisplayName, name) {
			return tech
		}
	}
	panic("cannot find tech with display name " + name + ", look at /home/ewen/projects/portfolio/hydrator/technologies.go")
}
