package main

import (
	"io/ioutil"

	"github.com/metal3d/go-slugify"
	"gopkg.in/yaml.v2"
)

// Tag represents a tag
type Tag struct {
	Singular string   `yaml:"singular"` // Plural form display name
	Plural   string   `yaml:"plural"`   // Singular form display name
	Aliases  []string `yaml:"name"`     // Works with a tag name in this array will be considered as tagged by the Tag
}

// URLName computes the identifier to use in the tag's page's URL
func (t Tag) URLName() string {
	return slugify.Marshal(t.Plural)
}

// ReferredToBy returns whether the given name refers to the tag
func (t *Tag) ReferredToBy(name string) bool {
	return StringsLooselyMatch(name, t.Plural, t.Singular, t.URLName()) || StringsLooselyMatch(name, t.Aliases...)
}

// LoadTags loads the tags from the given yaml file into a []Tag
func LoadTags(filename string) (tags []Tag, err error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(raw, &tags)
	return
}
