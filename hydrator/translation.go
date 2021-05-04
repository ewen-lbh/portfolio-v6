package main

import (
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	po "github.com/chai2010/gettext-go/po"
	"github.com/snapcore/go-gettext"
	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

// Translations holds both the gettext catalog from the .mo file
// and a po file object used to update the .po file (e.g. when discovering new translatable strings)
type Translations struct {
	poFile po.File
	moFile gettext.Catalog
}

func getLanguageCode(french bool) string {
	if french {
		return "fr"
	}
	return "en"
}

// TranslateToLanguage translates the given html node to french or english, removing translation-related attributes
func (t *Translations) TranslateToLanguage(french bool, root *html.Node) string {
	// Open files
	doc := goquery.NewDocumentFromNode(root)
	doc.Find("i18n, [i18n], [i18n-context]").Each(func(_ int, element *goquery.Selection) {
		element.RemoveAttr("i18n")
		msgContext, _ := element.Attr("i18n-context")
		element.RemoveAttr("i18n-context")
		if french {
			innerHTML, _ := element.Html()
			innerHTML = strings.TrimSpace(innerHTML)
			element.SetHtml(t.GetTranslation(innerHTML))
			if !t.IsInPOFile(innerHTML) {
				t.AppendTranslation(po.Message{
					MsgId:      innerHTML,
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

// IsInPOFile checks whether the given msgid is in the PO file
func (t *Translations) IsInPOFile(msgid string) bool {
	for _, message := range t.poFile.Messages {
		if message.MsgId == msgid {
			return true
		}
	}
	return false
}

// LoadTranslations reads from i18n/fr.{m,p}o to load both translation files
func LoadTranslations() (Translations, error) {
	messagesFile, err := os.Open("i18n/fr.mo")
	if err != nil {
		return Translations{}, err
	}
	moFile, err := gettext.ParseMO(messagesFile)
	if err != nil {
		return Translations{}, err
	}
	poFile, err := po.LoadFile("i18n/fr.po")
	if err != nil {
		return Translations{}, err
	}
	return Translations{
		moFile: moFile,
		poFile: *poFile,
	}, nil
}

// SavePO writes the .po file to the disk, with its potential modifications
func (t *Translations) SavePO(path string) {
	t.poFile.Save(path)
}

// GetTranslation returns the msgstr corresponding to msgid from the .mo file
// If not found, it returns the given msgid
func (t *Translations) GetTranslation(msgid string) string {
	return t.moFile.Gettext(msgid)
}

// AppendTranslation adds a new message to the .po file
// use SavePO to write the changes
func (t *Translations) AppendTranslation(msg po.Message) {
	t.poFile.Messages = append(t.poFile.Messages, msg)
}
