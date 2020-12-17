package main

// Abbreviation represents an abbreviation declaration in a description.md file
type Abbreviation struct {
	Name       string
	Definition string
}

// Footnote represents a footnote declaration in a description.md file
type Footnote struct {
	Name    string
	Content string
}

// Paragraph represents a paragraph declaration in a description.md file
type Paragraph struct {
	ID      string
	Content string
}

// Link represents an (isolated) link declaration in a description.md file
type Link struct {
	ID    string
	Name  string
	Title string
	URL   string
}

// Work represents a complete work, with analyzed mediae
type Work struct {
	ID       string
	Metadata struct {
		Created  string
		Tags     []string
		Layout   []interface{}
		MadeWith []string `json:"made with"`
	}
	Title      map[string]string
	Paragraphs map[string][]Paragraph
	Media      map[string][]Media
	Links      map[string][]Link
	Footnotes  map[string][]Footnote
}

type WorkOneLang struct {
	ID       string
	Metadata struct {
		Created  string
		Tags     []string
		Layout   []interface{}
		MadeWith []string `json:"made with"`
	}
	Title      string
	Paragraphs []Paragraph
	Media      []Media
	Links      []Link
	Footnotes  []Footnote
}

// MediaEmbedDeclaration represents media embeds. (abusing the ![]() syntax to extend it to any file)
// Only stores the info extracted from the syntax, no filesystem interactions.
type MediaEmbedDeclaration struct {
	Alt    string
	Title  string
	Source string
}

// ParsedDescription represents a work, but without analyzed media. All it contains is information from the description.md file
type ParsedDescription struct {
	Metadata               map[string]interface{}
	Title                  map[string]string
	Paragraphs             map[string][]Paragraph
	MediaEmbedDeclarations map[string][]MediaEmbedDeclaration
	Links                  map[string][]Link
	Footnotes              map[string][]Footnote
}

// ImageDimensions represents metadata about a media as it's extracted from its file
type ImageDimensions struct {
	Width       int
	Height      int
	AspectRatio float32
}

// Media represents a media object inserted in the work object's ``media`` array.
type Media struct {
	ID          string
	Alt         string
	Title       string
	Source      string
	ContentType string
	Size        uint64 // In bytes
	Dimensions  ImageDimensions
	Duration    uint // In seconds
	Online      bool // Whether the media is hosted online (referred to by an URL)
}
