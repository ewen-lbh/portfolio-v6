package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Technology represents something that a work was made with
// used for the /using/_technology path
type Technology struct {
	URLName      string   `yaml:"slug"`          // (unique) identifier used in the URL
	DisplayName  string   `yaml:"name"`          // name displayed to the user
	Aliases      []string `yaml:"aliases"`       // aliases pointing to the canonical URL (built from URLName)
	Author       string   `yaml:"by"`            // What company is behind the tech? (to display i.e. 'Adobe Photoshop' instead of 'Photoshop')
	LearnMoreURL string   `yaml:"learn more at"` // The technology's website
	Description  string   `yaml:"description"`   // A short description of the technology
}

// ReferredToBy returns whether the given name refers to the tech
func (t *Technology) ReferredToBy(name string) bool {
	return StringsLooselyMatch(name, t.URLName, t.DisplayName) || StringsLooselyMatch(name, t.Aliases...)
}

// LoadTechnologies loads the technologies from the given yaml file into a []Technology
func LoadTechnologies(filename string) (technologies []Technology, err error) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(raw, &technologies)
	return
}
