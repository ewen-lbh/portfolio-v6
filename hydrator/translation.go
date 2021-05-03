package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/snapcore/go-gettext"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
	po "github.com/chai2010/gettext-go/po"
)

func getLanguageCode(french bool) string {
	if french {
		return "fr"
	}
	return "en"
}

// TranslateToLanguage translates the given html node to french or english, removing translation-related attributes
func TranslateToLanguage(french bool, root *html.Node, catalog *gettext.Catalog, poFile *po.File) string {
	// Open files
	doc := goquery.NewDocumentFromNode(root)
	doc.Find("i18n, [i18n], [i18n-context]").Each(func(_ int, element *goquery.Selection) {
		element.RemoveAttr("i18n")
		msgContext, _ := element.Attr("i18n-context")
		element.RemoveAttr("i18n-context")
		if french {
			innerHTML, _ := element.Html()
			innerHTML = strings.TrimSpace(innerHTML)
			element.SetHtml(catalog.Gettext(innerHTML))
			if !IsMsgidInPoFile(poFile, innerHTML) {
				poFile.Messages = append(poFile.Messages, po.Message{
					MsgId: innerHTML,
					MsgContext: msgContext,
				})
			}
		}
	})
	htmlString, _ := doc.Html()
	htmlString = strings.ReplaceAll(htmlString, "<i18n>", "")
	htmlString = strings.ReplaceAll(htmlString, "</i18n>", "")
	htmlString = strings.ReplaceAll(htmlString, "[# LANGUAGE CODE #]", getLanguageCode(french))
	htmlString = strings.ReplaceAll(htmlString, "[# OTHER LANGUAGE CODE #]", getLanguageCode(!french))
	return gohtml.Format(htmlString)
}

func IsMsgidInPoFile(poFile *po.File, msgid string) bool {
	for _, message := range poFile.Messages {
		if message.MsgId == msgid {
			return true
		}
	}
	return false
}

