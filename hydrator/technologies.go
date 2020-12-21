package main

// Represents something that a work was made with
// used for the /using/_technology path
type Technology struct {
	URLName     string
	DisplayName string
	Aliases     []string
	Author      string
}

var KnownTechnologies = [...]Technology{
	{
		URLName:     "aftereffects",
		DisplayName: "After Effects",
		Author:      "Adobe",
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
		Author:      "Microsoft",
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
		Author:      "encode",
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
		Author:      "Google",
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
		Author:      "Image-Line",
	},
	{
		URLName:     "gimp",
		DisplayName: "GIMP",
		Author:      "GNU",
	},
	{
		URLName:     "go",
		DisplayName: "Go",
		Author:      "Google",
	},
	{
		URLName:     "html",
		DisplayName: "HTML",
		Author:      "W3C & WHATWG",
	},
	{
		URLName:     "illustrator",
		DisplayName: "Illustrator",
		Author:      "Adobe",
	},
	{
		URLName:     "indesign",
		DisplayName: "InDesign",
		Author:      "Adobe",
	},
	{
		URLName:     "javascript",
		DisplayName: "JavaScript",
		Aliases:     []string{"js"},
		Author:      "Mozilla",
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
		Author:      "Heroku",
	},
	{
		URLName:     "photoshop",
		DisplayName: "Photoshop",
		Author:      "Adobe",
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
		Author:      "Adobe",
	},
	{
		URLName:     "pug",
		DisplayName: "Pug",
	},
	{
		URLName:     "pychemin",
		DisplayName: "PyChemin",
		Author:      "ewen-lbh",
	},
	{
		URLName:     "python",
		DisplayName: "Python",
		Author:      "PSF",
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
		Author:      "Mozilla",
	},
	{
		URLName:     "sapper",
		DisplayName: "Sapper",
		Author:      "Svelte",
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
		Author:      "Microsoft",
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
		URLName:     "manim",
		DisplayName: "Manim",
		Author:      "3Blue1Brown",
	},
	{
		URLName:     "lark",
		DisplayName: "Lark",
	},
}
