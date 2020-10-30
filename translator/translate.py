import os, re
from pathlib import Path
import bs4
import gettext


translations_dir = Path("../trans/fr/")
translate_strings_pattern = re.compile(r"\[# '([^']+)' | translate #\]")
for filepath in [ Path(f) for f in os.listdir(translations_dir)]:
    print(f"for file {filepath}:")
    with open(translations_dir / filepath, 'r') as file:
        raw = file.read()
        html = bs4.BeautifulSoup(raw, features="html.parser")
        for tag in html("i18n") + html(translate="translate"):
            del tag["translate"]
            if tag.name == "i18n":
                for tag in tag.parent.children:
                    print(tag)
            print(tag)