package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snapcore/go-gettext"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

func getLanguageCode(french bool) string {
	if french {
		return "fr"
	}
	return "en"
}

// TranslateToLanguage translates the given html node to french or english, removing translation-related attributes
func TranslateToLanguage(french bool, root *html.Node, catalog *gettext.Catalog) string {
	// Open files
	doc := goquery.NewDocumentFromNode(root)
	doc.Find("i18n, [i18n], [i18n-context]").Each(func(_ int, element *goquery.Selection) {
		element.RemoveAttr("translate")
		element.RemoveAttr("translate-context")
		if french {
			innerHTML, _ := element.Html()
			element.SetHtml(catalog.Gettext(innerHTML))
		}
	})
	htmlString, _ := doc.Html()
	htmlString = strings.ReplaceAll(htmlString, "<i18n>", "")
	htmlString = strings.ReplaceAll(htmlString, "</i18n>", "")
	htmlString = strings.ReplaceAll(htmlString, "[# LANGUAGE CODE #]", getLanguageCode(french))
	htmlString = strings.ReplaceAll(htmlString, "[# OTHER LANGUAGE CODE #]", getLanguageCode(!french))
	return gohtml.Format(htmlString)
}
