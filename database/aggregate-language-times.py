#!/usr/bin/env python
from pathlib import Path
import json
import yaml

here = Path(__file__).parent
wakatime = json.loads((here / 'wakatime.json').read_text())
technologies = yaml.load((here / 'technologies.yaml').read_text(), Loader=yaml.SafeLoader)
per_language_totals = {}
per_project_totals = {}

# Maps Wakatime language names to the portfolio's slug names
LANGUAGE_NAMES_MAP = { t['name'] : t['slug'] for t in technologies } | { "fish": "fishshell", "TeX": "latex", "JSON with Comments": "json", "Gettext Catalog": "gettext", "ca65 assembler": "assembler", "Vue.js": "vue", "SCSS": "sass" }

print(LANGUAGE_NAMES_MAP)

def slug_of_language(wakatime_name: str) -> str:
	for name, slug in LANGUAGE_NAMES_MAP.items():
		if name.lower().strip() == wakatime_name.lower().strip():
			return slug
	return wakatime_name

def sort_dict_per_value(o: dict) -> dict:
	return {k: v for k, v in sorted(o.items(), key=lambda item: item[1], reverse=True)}

for day in wakatime["days"]:
	for language in day["languages"]:
		name = slug_of_language(language["name"])
		if name not in per_language_totals:
			per_language_totals[name] = 0

		per_language_totals[name] += language["seconds"] + language["minutes"] * 60 + language["hours"] * 3600

	for project in day["projects"]:
		if project["name"] not in per_project_totals:
			per_project_totals[project["name"]] = 0

		per_project_totals[project["name"]] += project["grand_total"]["total_seconds"]

Path(here / "wakatime-aggregated.json").write_text(
	json.dumps({
		"languages": sort_dict_per_value(per_language_totals),
		"projects": sort_dict_per_value(per_project_totals),
	}, indent=2)
)
