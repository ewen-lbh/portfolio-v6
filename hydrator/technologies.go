package main

// Represents something that a work was made with
// used for the /using/_technology path
type Technology struct {
	URLName     string
	DisplayName string
	Aliases     []string
}

var KnownTechnologies = [...]Technology{
	{
		URLName:     "aftereffects",
		DisplayName: "After Effects",
	},
	{
		URLName:     "asciidoc",
		DisplayName: "asciidoc",
	},
	{
		URLName:     "assembly",
		DisplayName: "Assembly",
	},
	{
		URLName:     "coffeescript",
		DisplayName: "CoffeeScript",
	},
	{
		URLName:     "css",
		DisplayName: "CSS",
	},
	{
		URLName:     "c#",
		DisplayName: "C#",
		Aliases:     []string{"cs", "csharp"},
	},
	{
		URLName:     "c++",
		DisplayName: "C++",
		Aliases:     []string{"cpp", "cplusplus"},
	},
	{
		URLName:     "c",
		DisplayName: "C",
	},
	{
		URLName:     "djangorestframework",
		DisplayName: "Django REST Framework",
	},
	{
		URLName:     "django",
		DisplayName: "Django",
	},
	{
		URLName:     "docopt",
		DisplayName: "Docopt",
	},
	{
		URLName:     "figma",
		DisplayName: "Figma",
	},
	{
		URLName:     "fishshell",
		DisplayName: "Fish Shell",
		Aliases:     []string{"fish"},
	},
	{
		URLName:     "flstudio",
		DisplayName: "FL Studio",
		Aliases:     []string{"fruityloops"},
	},
	{
		URLName:     "gimp",
		DisplayName: "GIMP",
	},
	{
		URLName:     "go",
		DisplayName: "Go",
	},
	{
		URLName:     "html",
		DisplayName: "HTML",
	},
	{
		URLName:     "illustrator",
		DisplayName: "Illustrator",
	},
	{
		URLName:     "indesign",
		DisplayName: "InDesign",
	},
	{
		URLName:     "javascript",
		DisplayName: "JavaScript",
		Aliases:     []string{"js"},
	},
	{
		URLName:     "json",
		DisplayName: "JSON",
	},
	{
		URLName:     "latex",
		DisplayName: "LaTeX",
	},
	{
		URLName:     "livescript",
		DisplayName: "LiveScript",
	},
	{
		URLName:     "lua",
		DisplayName: "Lua",
	},
	{
		URLName:     "markdown",
		DisplayName: "Markdown",
	},
	{
		URLName:     "nestjs",
		DisplayName: "NestJS",
	},
	{
		URLName:     "nim",
		DisplayName: "Nim",
	},
	{
		URLName:     "nuxt",
		DisplayName: "Nuxt",
		Aliases:     []string{"nuxtjs"},
	},
	{
		URLName:     "oclif",
		DisplayName: "Oclif",
	},
	{
		URLName:     "photoshop",
		DisplayName: "Photoshop",
	},
	{
		URLName:     "php",
		DisplayName: "PHP",
	},
	{
		URLName:     "plantuml",
		DisplayName: "PlantUML",
	},
	{
		URLName:     "postgresql",
		DisplayName: "PostGreSQL",
	},
	{
		URLName:     "premierepro",
		DisplayName: "Premiere Pro",
	},
	{
		URLName:     "pug",
		DisplayName: "Pug",
	},
	{
		URLName:     "pychemin",
		DisplayName: "PyChemin",
	},
	{
		URLName:     "python",
		DisplayName: "Python",
	},
	{
		URLName:     "rubyonrails",
		DisplayName: "Ruby On Rails",
	},
	{
		URLName:     "ruby",
		DisplayName: "Ruby",
	},
	{
		URLName:     "rust",
		DisplayName: "Rust",
	},
	{
		URLName:     "sapper",
		DisplayName: "Sapper",
	},
	{
		URLName:     "sass",
		DisplayName: "SASS",
	},
	{
		URLName:     "shell",
		DisplayName: "Shell",
	},
	{
		URLName:     "stylus",
		DisplayName: "Stylus",
	},
	{
		URLName:     "svelte",
		DisplayName: "Svelte",
	},
	{
		URLName:     "toml",
		DisplayName: "TOML",
	},
	{
		URLName:     "typescript",
		DisplayName: "TypeScript",
	},
	{
		URLName:     "vue",
		DisplayName: "Vue",
		Aliases:     []string{"vuejs"},
	},
	{
		URLName:     "webpack",
		DisplayName: "Webpack",
	},
	{
		URLName:     "yaml",
		DisplayName: "YAML",
	},
	{
		URLName: "manim",
		DisplayName: "ManimCE",
	},
	{
		URLName: "lark",
		DisplayName: "Lark",
	},
}
