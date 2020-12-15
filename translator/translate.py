from gettext import GNUTranslations
import os, re
from pathlib import Path
import bs4
from polib import mofile
from typing import Any

def p(indent_level: int, *args: Any, **kwargs: Any):
    args = list(args)
    if indent_level == 0:
        args = ['\n'] + args
    print(indent_level * "  ", *args, **kwargs)

os.chdir("..")
to_translate_dir = Path("artifacts/phase_2/")
translated_dir = Path("dist/")
translate_strings_pattern = re.compile(r"\[# '([^']+)' \| translate #\]")

catalog = { e.msgid:e.msgstr for e in mofile("messages/fr.mo") }

def translate(lang: str):
    p(0, f"Translating in {lang}")
    if lang != 'en':
        translation = lambda msgid: catalog.get(msgid) or msgid
    else:
        translation = lambda msgid: msgid
    
    for filepath in [ Path(f) for f in os.listdir(to_translate_dir / lang)]:
        p(1, f"Translating {filepath}")
        with open(to_translate_dir / lang / filepath, 'r') as file:
            raw = file.read()
            raw = raw.replace("[# LANGUAGE CODE #]", lang)
            if (matches := translate_strings_pattern.finditer(raw)):
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
                    p(2, f"{tag.string} ~> {translation(tag.string)}")
                tag.string = translation(tag.string)
        with open(translated_dir / lang / filepath, 'w') as file:
            file.write(str(html))

translate("en")
translate("fr")
