from gettext import GNUTranslations
import os, re
from pathlib import Path
import bs4
from polib import mofile

os.chdir("..")
to_translate_dir = Path("artifacts/phase_2/")
translated_dir = Path("dist/")
translate_strings_pattern = re.compile(r"\[# '([^']+)' \| translate #\]")

catalog = { e.msgid:e.msgstr for e in mofile("messages/fr.mo") }

def translate(lang: str):
    if lang != 'en':
        translation = lambda msgid: catalog.get(msgid) or msgid
    else:
        translation = lambda msgid: msgid
    for filepath in [ Path(f) for f in os.listdir(to_translate_dir / lang)]:
        with open(to_translate_dir / lang / filepath, 'r') as file:
            raw = file.read()
            if (matches := translate_strings_pattern.finditer(raw)):
                for match in matches:
                    raw = raw.replace(match.group(0), translation(match.group(1)))
            html = bs4.BeautifulSoup(raw, features="html.parser")
            for tag in html("i18n") + html(translate="translate"):
                del tag["translate"]
                if tag.name == "i18n":
                    # TODO: Find a way to _remove_ the tag and insert the content 
                    # into the parent tag as a text node directly, this is ugly af.
                    tag.name = "span" 
        with open(translated_dir / lang / filepath, 'w') as file:
            file.write(str(html))

translate("en")
translate("fr")
