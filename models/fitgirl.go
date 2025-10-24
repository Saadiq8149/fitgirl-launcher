package models

type FitgirlRepack struct {
	Name         string
	OriginalSize string
	RepackSize   string
	Screenshots  []string
	Sources      []string
	CoverImage   string
}

type PopularRepack struct {
	Link       string
	CoverImage string
}

type PopularRepacks struct {
	Repacks []PopularRepack
}

type FitgirlPage struct {
	Results []string
	Page    int
	Total   int
}
