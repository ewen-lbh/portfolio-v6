from gettext import GNUTranslations
import os, re
from pathlib import Path
import bs4
from polib import mofile
from typing import Any
from termcolor import colored, cprint
import sys

def verbose() -> bool:
	return False # TODO

def p(indent_level: int, *args: Any, **kwargs: Any):
	# TODO: indent level proportional to verbosity level
	if not verbose() and indent_level >= 1:
		return
	args = list(args)
	# if indent_level == 0:
	# 	args = ["\n"] + args
	print(indent_level * "\t", *args, **kwargs)


os.chdir("..")
to_translate_dir = Path("artifacts/phase_2/")
translated_dir = Path("dist/")
translate_strings_pattern = re.compile(r"\[# '([^']+)' \| translate #\]")

catalog = {e.msgid: e.msgstr for e in mofile("messages/fr.mo")}

longest_msgid_length = max(len(k) for k in catalog.keys())


def translate(lang: str):
	p(0, f"Translating to {lang}")
	translate_files(lang, (to_translate_dir / lang).iterdir())


def translate_files(lang: str, files):
	if lang != "en":
		other_lang = "en"
		translation = lambda msgid: catalog.get(msgid) or msgid
	else:
		other_lang = "fr"
		translation = lambda msgid: msgid

	for filepath in files:
		Path(
			str(filepath).replace(str(to_translate_dir), str(translated_dir))
		).parent.mkdir(exist_ok=True)
		if filepath.is_dir():
			translate_files(lang, filepath.iterdir())
		else:
			p(1, f"Translating {filepath}")
			with open(filepath, "r") as file:
				raw = file.read()
				raw = raw.replace("[# LANGUAGE CODE #]", lang)
				raw = raw.replace("[# OTHER LANGUAGE CODE #]", other_lang)
				if (matches := translate_strings_pattern.finditer(raw)) :
					for match in matches:
						p(2, f"{match.group(0)} ~> {translation(match.group(1))}")
						raw = raw.replace(match.group(0), translation(match.group(1)))
				html = bs4.BeautifulSoup(raw, features="html.parser")
				for tag in html("i18n") + html(translate="translate"):
					del tag["translate"]
					del tag["translate-context"]
					if tag.name == "i18n":
						# TODO: Find a way to _remove_ the tag and insert the content
						# into the parent tag as a text node directly, this is ugly af.
						tag.name = "span"
					contents_before = tag.decode_contents()
					tag.clear()
					tag.append(bs4.BeautifulSoup(translation(contents_before), features="html.parser"))
					contents_after = tag.decode_contents()
					p(2, f"{contents_before:<{min(longest_msgid_length, 20)}} ~> {contents_after}")
					if not contents_after:
						cprint(f"WARN: Translation of {contents_before!r} is empty")
			with open(
				str(filepath).replace(str(to_translate_dir), str(translated_dir)), "w"
			) as file:
				file.write(str(html))


translate("en")
translate("fr")
