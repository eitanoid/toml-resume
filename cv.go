package main

type Header struct {
	Header_format                                  []string
	Name, Email, Location, Phone, Linkedin, Github string
	Name_size                                      int
}

type Experience struct {
	Dates, Location, Title, Institution string
	Bulletpoints                        []string
}

type Education struct {
	Dates, Location, Title, Institution string
	Bulletpoints                        []string
}

type Project struct {
	Title, Dates, Description string
	Bulletpoints              []string
}

type Config struct {
	Font_size               int
	Font_name               string
	Margin_size             int
	Cv_order                []string
	Project_header_order    []string
	Experience_header_order []string
	Education_header_order  []string
	Section_titles          []string
}

type CV struct {
	Config     Config
	Header     Header
	Experience []Experience
	Education  []Education
	Project    []Project
	Skills     map[string]string
}
