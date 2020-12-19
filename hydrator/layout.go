package main

import (
	"fmt"
	"reflect"
	"strings"
)

type LayoutElement struct {
	IsParagraph bool
	IsMedia     bool
	IsLink      bool
	IsSpacer    bool
}

type usedCounts struct {
	p int
	m int
	l int
}

type Layout = [][]LayoutElement

// LayoutElementRepr returns the strings representation of a Layout
func layoutElementRepr(layoutElement LayoutElement) string {
	if layoutElement.IsLink {
		return "l"
	}
	if layoutElement.IsMedia {
		return "m"
	}
	if layoutElement.IsParagraph {
		return "p"
	}
	if layoutElement.IsSpacer {
		return " "
	}
	panic("unexpected layoutElement: is neither link nor media nor parargraph nor spacer")
}

func layoutRepr(layout Layout) string {
	repr := ""
	for _, row := range layout {
		for _, element := range row {
			repr += layoutElementRepr(element)
		}
		repr += "\n"
	}
	return repr
}

func (work *WorkOneLang) BuildLayout() string {
	var layout Layout
	if len(work.Metadata.Layout) >= 1 {
		layout = loadLayout(work.Metadata.Layout)
	} else {
		layout = autoLayout(work)
	}
	usedCounts := usedCounts{}
	var built string
	for _, layoutRow := range layout {
		var row string
		for _, layoutElement := range layoutRow {
			var element string
			if layoutElement.IsSpacer {
				element = `div.spacer`
			} else if layoutElement.IsLink {
				if len(work.Links) <= usedCounts.l {
					panic(fmt.Sprintf(`Not enough Links to satisfy the given layout:

· Layout is:
%v

· work has only %d links
`, layoutRepr(layout), usedCounts.l))
				}
				data := work.Links[usedCounts.l]
				usedCounts.l++
				element = fmt.Sprintf(`a(href="%v" id="%v" title="%v") %v`, data.URL, data.ID, data.Title, data.Name)
			} else if layoutElement.IsMedia {
				if len(work.Media) <= usedCounts.m {
					panic(fmt.Sprintf(`Not enough Media to satisfy the given layout:

· Layout is:
%v

· work has only %d media
`, layoutRepr(layout), usedCounts.m))
				}
				data := work.Media[usedCounts.m]
				usedCounts.m++
				mediaGeneralContentType := strings.Split(data.ContentType, "/")[0]
				if data.ContentType == "application/pdf" {
					mediaGeneralContentType = "pdf"
				}
				switch mediaGeneralContentType {
				case "video":
					element = fmt.Sprintf(`video(src="%v" id="%v" title="%v") %v`, data.Source, data.ID, data.Title, data.Alt)
				case "audio":
					element = fmt.Sprintf(`audio(src="%v" id="%v" title="%v") %v`, data.Source, data.ID, data.Title, data.Alt)
				case "image":
					element = fmt.Sprintf(`img(src="%v" id="%v" title="%v" alt="%v")`, data.Source, data.ID, data.Title, data.Alt)
				case "pdf":
					element = fmt.Sprintf(`.pdf-frame-container: iframe.pdf-frame(src="%v" id="%v" title="%v" width="100%%" height="100%%") %v`, data.Source, data.ID, data.Title, data.Alt)
				default:
					element = fmt.Sprintf(`a(href="%v" id="%v" title="%v") %v`, data.Source, data.ID, data.Title, data.Alt)
				}
			} else if layoutElement.IsParagraph {
				if len(work.Paragraphs) <= usedCounts.p {
					panic(fmt.Sprintf(`Not enough Paragraphs to satisfy the given layout:

· Layout is:
%v

· work has only %d paragraphs
`, layoutRepr(layout), usedCounts.p))
				}
				data := work.Paragraphs[usedCounts.p]
				usedCounts.p++
				// element = fmt.Sprintf(`<p id="%v">%v</p>`, data.ID, data.Content)
				element = "p(id=\"" + data.ID + "\")." + IndentWithTabsNewline(2, data.Content)
			}
			row += "\t" + element + "\n"
		}
		built += "div.row\n" + row + "\n"
	}
	return built
}

func autoLayout(work *WorkOneLang) Layout {
	layout := make(Layout, 0)
	for range work.Paragraphs {
		layout = append(layout, []LayoutElement{{IsParagraph: true}})
	}
	for range work.Media {
		layout = append(layout, []LayoutElement{{IsMedia: true}})
	}
	for range work.Links {
		layout = append(layout, []LayoutElement{{IsLink: true}})
	}
	return layout
}

func loadLayout(layout []interface{}) Layout {
	loaded := make([][]LayoutElement, 0)
	for _, layoutRowMaybeSlice := range layout {
		loadedRow := make([]LayoutElement, 0)
		var layoutRow []interface{}
		if reflect.TypeOf(layoutRowMaybeSlice).Kind() != reflect.Slice {
			layoutRow = []interface{}{layoutRowMaybeSlice}
		} else {
			layoutRow = layoutRowMaybeSlice.([]interface{})
		}
		for _, layoutElement := range layoutRow {
			loadedRow = append(loadedRow, loadLayoutElement(layoutElement))
		}
		loaded = append(loaded, loadedRow)
	}
	return loaded
}

func loadLayoutElement(layoutElement interface{}) LayoutElement {
	return LayoutElement{
		IsParagraph: layoutElement == "p",
		IsMedia:     layoutElement == "m",
		IsLink:      layoutElement == "l",
		IsSpacer:    layoutElement == nil,
	}
}
