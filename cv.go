package main

type Header struct {
	Header_format                                  []string
	Name, Email, Location, Phone, Linkedin, Github string
	Name_size                                      int
}

type SectionEntry struct {
	Section_type                                     string
	Dates, Description, Location, Title, Institution string
	Bulletpoints                                     []string
	Points                                           map[string]string
}

type Config struct {
	Font_size               int
	Font_scale              float64
	Page_margin             float64
	Font                    string
	Section_order           []string
	Project_header_order    []string
	Experience_header_order []string
	Education_header_order  []string
}

type CV struct {
	Config  Config
	Header  Header
	Skills  map[string]string
	Section map[string][]SectionEntry // all sections can live here
}
