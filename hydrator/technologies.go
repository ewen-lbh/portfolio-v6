package main

// Represents something that a work was made with
// used for the /using/_technology path
type Technology struct {
	URLName      string   // (unique) identifier used in the URL
	DisplayName  string   // name displayed to the user
	Aliases      []string // aliases pointing to the canonical URL (built from URLName)
	Author       string   // What company is behind the tech? (to display i.e. 'Adobe Photoshop' instead of 'Photoshop')
	LearnMoreURL string   // The technology's website
}

// ReferredToBy returns whether the given name refers to the tech
func (t *Technology) ReferredToBy(name string) bool {
	return StringsLooselyMatch(name, t.URLName, t.DisplayName) || StringsLooselyMatch(name, t.Aliases...)
}

var KnownTechnologies = [...]Technology{
	{
		URLName:      "aftereffects",
		DisplayName:  "After Effects",
		Author:       "Adobe",
		LearnMoreURL: "https://www.adobe.com/fr/products/aftereffects.html",
	},
	{
		URLName:      "asciidoc",
		DisplayName:  "asciidoc",
		LearnMoreURL: "",
	},
	{
		URLName:      "assembly",
		DisplayName:  "Assembly",
		LearnMoreURL: "https://apps.apple.com/us/app/assembly-graphic-design-art/id1024210402",
	},
	{
		URLName:      "coffeescript",
		DisplayName:  "CoffeeScript",
		LearnMoreURL: "https://coffeescript.org/",
	},
	{
		URLName:      "css",
		DisplayName:  "CSS",
		LearnMoreURL: "https://developer.mozilla.org/fr/docs/Web/CSS",
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
		URLName:      "djangorestframework",
		DisplayName:  "Django REST Framework",
		Author:       "encode",
		LearnMoreURL: "https://www.django-rest-framework.org",
	},
	{
		URLName:      "django",
		DisplayName:  "Django",
		LearnMoreURL: "https://www.djangoproject.com",
	},
	{
		URLName:      "docopt",
		DisplayName:  "Docopt",
		LearnMoreURL: "https://github.com/docopt",
	},
	{
		URLName:      "figma",
		DisplayName:  "Figma",
		Author:       "Google",
		LearnMoreURL: "https://figma.com",
	},
	{
		URLName:      "fishshell",
		DisplayName:  "Fish Shell",
		Aliases:      []string{"fish"},
		LearnMoreURL: "https://fishshell.com",
	},
	{
		URLName:      "flstudio",
		DisplayName:  "FL Studio",
		Aliases:      []string{"fruityloops"},
		Author:       "Image-Line",
		LearnMoreURL: "https://www.image-line.com/flstudio/",
	},
	{
		URLName:      "gimp",
		DisplayName:  "GIMP",
		Author:       "GNU",
		LearnMoreURL: "https://www.gimp.org/",
	},
	{
		URLName:      "go",
		DisplayName:  "Go",
		Author:       "Google",
		LearnMoreURL: "https://go.dev",
	},
	{
		URLName:      "html",
		DisplayName:  "HTML",
		LearnMoreURL: "https://developer.mozilla.org/fr/docs/Web/HTML",
	},
	{
		URLName:      "illustrator",
		DisplayName:  "Illustrator",
		Author:       "Adobe",
		LearnMoreURL: "https://www.adobe.com/fr/products/illustrator.html",
	},
	{
		URLName:      "indesign",
		DisplayName:  "InDesign",
		Author:       "Adobe",
		LearnMoreURL: "https://www.adobe.com/fr/products/indesign.html",
	},
	{
		URLName:      "javascript",
		DisplayName:  "JavaScript",
		Aliases:      []string{"js"},
		Author:       "Mozilla",
		LearnMoreURL: "https://developer.mozilla.org/fr/docs/Web/JavaScript",
	},
	{
		URLName:      "json",
		DisplayName:  "JSON",
		LearnMoreURL: "https://www.json.org/",
	},
	{
		URLName:      "latex",
		DisplayName:  "LaTeX",
		LearnMoreURL: "https://www.latex-project.org/",
	},
	{
		URLName:      "livescript",
		DisplayName:  "LiveScript",
		LearnMoreURL: "https://livescript.net",
	},
	{
		URLName:      "lua",
		DisplayName:  "Lua",
		LearnMoreURL: "https://lua.org",
	},
	{
		URLName:      "markdown",
		DisplayName:  "Markdown",
		LearnMoreURL: "https://daringfireball.net/projects/markdown/",
	},
	{
		URLName:     "nestjs",
		DisplayName: "NestJS",
	},
	{
		URLName:      "nim",
		DisplayName:  "Nim",
		LearnMoreURL: "https://nim-lang.org/",
	},
	{
		URLName:      "nuxt",
		DisplayName:  "Nuxt",
		Aliases:      []string{"nuxtjs"},
		LearnMoreURL: "https://nuxtjs.org",
	},
	{
		URLName:      "oclif",
		DisplayName:  "Oclif",
		Author:       "Heroku",
		LearnMoreURL: "https://oclif.io/",
	},
	{
		URLName:      "photoshop",
		DisplayName:  "Photoshop",
		Author:       "Adobe",
		LearnMoreURL: "https://www.adobe.com/fr/products/photoshop.html",
	},
	{
		URLName:      "php",
		DisplayName:  "PHP",
		LearnMoreURL: "https://www.php.net/",
	},
	{
		URLName:      "plantuml",
		DisplayName:  "PlantUML",
		LearnMoreURL: "https://plantuml.com/",
	},
	{
		URLName:      "postgresql",
		DisplayName:  "PostGreSQL",
		LearnMoreURL: "https://www.postgresql.org/",
	},
	{
		URLName:      "premierepro",
		DisplayName:  "Premiere Pro",
		Author:       "Adobe",
		LearnMoreURL: "https://www.adobe.com/fr/products/premiere.html",
	},
	{
		URLName:      "pug",
		DisplayName:  "Pug",
		LearnMoreURL: "https://pugjs.org",
	},
	{
		URLName:      "pychemin",
		DisplayName:  "PyChemin",
		Author:       "ewen-lbh",
		LearnMoreURL: "https://ewen.works/pychemin",
	},
	{
		URLName:      "python",
		DisplayName:  "Python",
		Author:       "PSF",
		LearnMoreURL: "https://python.org",
	},
	{
		URLName:      "rubyonrails",
		DisplayName:  "Ruby On Rails",
		LearnMoreURL: "https://rubyonrails.org",
	},
	{
		URLName:      "ruby",
		DisplayName:  "Ruby",
		LearnMoreURL: "https://www.ruby-lang.org",
	},
	{
		URLName:      "rust",
		DisplayName:  "Rust",
		Author:       "Mozilla",
		LearnMoreURL: "https://www.rust-lang.org",
	},
	{
		URLName:      "sapper",
		DisplayName:  "Sapper",
		Author:       "Svelte",
		LearnMoreURL: "https://sapper.svelte.dev",
	},
	{
		URLName:      "sass",
		DisplayName:  "SASS",
		LearnMoreURL: "https://sass-lang.com/",
	},
	{
		URLName:     "shell",
		DisplayName: "Shell",
	},
	{
		URLName:      "stylus",
		DisplayName:  "Stylus",
		LearnMoreURL: "https://stylus-lang.com",
	},
	{
		URLName:      "svelte",
		DisplayName:  "Svelte",
		LearnMoreURL: "https://svelte.dev/",
	},
	{
		URLName:      "toml",
		DisplayName:  "TOML",
		LearnMoreURL: "https://toml.io",
	},
	{
		URLName:      "typescript",
		DisplayName:  "TypeScript",
		Author:       "Microsoft",
		LearnMoreURL: "https://www.typescriptlang.org/",
	},
	{
		URLName:      "vue",
		DisplayName:  "Vue",
		Aliases:      []string{"vuejs"},
		LearnMoreURL: "https://vuejs.org",
	},
	{
		URLName:      "webpack",
		DisplayName:  "Webpack",
		LearnMoreURL: "https://webpack.js.org/",
	},
	{
		URLName:      "yaml",
		DisplayName:  "YAML",
		LearnMoreURL: "https://yaml.org/",
	},
	{
		URLName:      "manim",
		DisplayName:  "Manim",
		Author:       "3Blue1Brown",
		LearnMoreURL: "https://www.manim.community/",
	},
	{
		URLName:      "lark",
		DisplayName:  "Lark",
		LearnMoreURL: "https://lark-parser.readthedocs.io/",
	},
}
