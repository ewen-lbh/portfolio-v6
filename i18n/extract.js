#!/usr/bin/env node

const { GettextExtractor, HtmlExtractors } = require("gettext-extractor")

let extractor = new GettextExtractor()

extractor
	.createHtmlParser([
		HtmlExtractors.elementContent('i18n', {
			attributes: {
				textPlural: "plural",
				context: "context",
				comment: "comment"
			}
		}),
		HtmlExtractors.elementContent('[i18n]', {
			attributes: {
				textPlural: "i18n-plural",
				context: "i18n-context",
				comment: "i18n-comment"
			}
		})
	])
	.parseFilesGlob("dist/{fr,en}/**.html")

extractor.savePotFile("i18n/fr.new.po")
extractor.printStats()
