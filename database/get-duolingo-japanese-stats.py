#!/usr/bin/env python
import requests
import re
import json
from pathlib import Path

EXTRACT_LEAGUE = re.compile(r"(\w+) League · \d+ weeks?")
EXTRACT_STREAK = re.compile(r"(\d+) days streak")
EXTRACT_CROWNS = re.compile(r"<b>\s*[cC]rowns\s*:\s*</b>\s*(\d+)/(\d+)")

content = requests.get("https://duome.eu/ewen_lbh/en/ja").text

crowns_done, crowns_total = (int(d) for d in (EXTRACT_CROWNS.search(content).groups()))

data = {
    "league": EXTRACT_LEAGUE.search(content).group(1),
    "streak": int(EXTRACT_STREAK.search(content).group(1)),
    "crowns": {
        "done": crowns_done,
        "total": crowns_total,
        "progress": crowns_done/crowns_total,
    }
}

print(f"Writing {data}")

(Path(__file__).parent / "duolingo-japanese.json").write_text(json.dumps(data))
