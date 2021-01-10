package main

import "strings"

// Tag represents a tag
type Tag struct {
	URLName     string   // Which name to use as an identifier in the URL
	DisplayName string   // Which name to display to the user
	Aliases     []string // Works with a tag name in this array will be considered as tagged by the Tag
}

// TODO: Define "Plural" and "Singular" in KnownTags and figure out URLName from "Plural" by slugifying it.

// PluralDisplayName computes the plural display name. This is just the URLName without dashes.
func (t *Tag) PluralDisplayName() string {
	return strings.ReplaceAll(t.URLName, "-", " ")
}

// ReferredToBy returns whether the given name refers to the tag
func (t *Tag) ReferredToBy(name string) bool {
	return StringsLooselyMatch(name, t.URLName, t.DisplayName) || StringsLooselyMatch(name, t.Aliases...)
}

// KnownTags defines which tags are valid. Each Tag will get its correspoding page generated from _tag.pug.
var KnownTags = [...]Tag{
	{
		// TODO: deprecate
		DisplayName: "school",
		URLName:     "school",
	},
	{
		DisplayName: "science",
		URLName:     "science",
	},
	{
		DisplayName: "card",
		URLName:     "cards",
	},
	{
		// TODO: deprecate
		DisplayName: "cover art",
		URLName:     "cover-arts",
	},
	{
		DisplayName: "game",
		URLName:     "games",
	},
	{
		DisplayName: "graphism",
		URLName:     "graphism",
	},
	{
		DisplayName: "poster",
		URLName:     "posters",
	},
	{
		DisplayName: "automation",
		URLName:     "automation",
	},
	{
		DisplayName: "web",
		URLName:     "web",
	},
	{
		DisplayName: "intro",
		URLName:     "intros",
	},
	{
		DisplayName: "music",
		URLName:     "music",
	},
	{
		DisplayName: "app",
		URLName:     "apps",
	},
	{
		DisplayName: "book",
		URLName:     "books",
	},
	{
		DisplayName: "api",
		URLName:     "APIs",
	},
	{
		DisplayName: "program",
		URLName:     "programs",
	},
	{
		DisplayName: "cli",
		URLName:     "CLIs",
	},
	{
		DisplayName: "motion design",
		URLName:     "motion-design",
	},
	{
		DisplayName: "compositing",
		URLName:     "compositing",
	},
	{
		DisplayName: "visual identity",
		URLName:     "visual-identities",
		Aliases:     []string{"logo", "logos", "banner", "banners"},
	},
	{
		DisplayName: "illustration",
		URLName:     "illustrations",
	},
	{
		DisplayName: "typography",
		URLName:     "typography",
	},
	{
		DisplayName: "drawing",
		URLName:     "drawings",
	},
	{
		DisplayName: "icons",
		URLName:     "icons",
	},
	{
		DisplayName: "site",
		URLName:     "sites",
	},
	{
		DisplayName: "language",
		URLName:     "languages",
	},
	{
		DisplayName: "math",
		URLName:     "math",
		Aliases:     []string{"maths", "mathematics"},
	},
}
