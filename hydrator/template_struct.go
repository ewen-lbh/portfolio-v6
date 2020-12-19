package main

type TemplateData struct {
	Tags                   [len(tags)]Tag
	SortingOptions         []string
	WorksByYear            map[int][]WorkOneLang
	LatestWork             WorkOneLang
	CurrentWork            WorkOneLang
	CurrentWorkBuiltLayout string
	CurrentPath            string
	CurrentTag             Tag
	CurrentTagWorks        []WorkOneLang
	CurrentTech            Technology
	CurrentTechWorks       []WorkOneLang
	Age                    uint8
}
