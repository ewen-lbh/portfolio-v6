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

func (work *WorkOneLang) BuildLayout() string {
	layout := loadLayout(work.Metadata.Layout)
	usedCounts := usedCounts{}
	var built string
	for _, layoutRow := range layout {
		var row string
		for _, layoutElement := range layoutRow {
			var element string
			if layoutElement.IsSpacer {
				element = `div.spacer`
			} else if layoutElement.IsLink {
				data := work.Links[usedCounts.l]
				usedCounts.l++
				element = fmt.Sprintf(`a(href="%v" id="%v" title="%v") %v`, data.URL, data.ID, data.Title, data.Name)
			} else if layoutElement.IsMedia {
				data := work.Media[usedCounts.m]
				usedCounts.m++
				mediaGeneralContentType := strings.Split(data.ContentType, "/")[0]
				switch mediaGeneralContentType {
				case "video":
					element = fmt.Sprintf(`video(src="%v" id="%v" title="%v") %v`, data.Source, data.ID, data.Title, data.Alt)
				case "audio":
					element = fmt.Sprintf(`audio(src="%v" id="%v" title="%v") %v`, data.Source, data.ID, data.Title, data.Alt)
				case "image":
					element = fmt.Sprintf(`img(src="%v" id="%v" title="%v" alt="%v")`, data.Source, data.ID, data.Title, data.Alt)
				default:
					element = fmt.Sprintf(`a(href="%v" id="%v" title="%v") %v`, data.Source, data.ID, data.Title, data.Alt)
				}
			} else if layoutElement.IsParagraph {
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

// TODO: handle no layout defined explicitly

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
