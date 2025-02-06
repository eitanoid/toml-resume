package main

import "strings"

type Header struct {
	Details   [][]string
	Name      string
	Name_size int
}

type SectionEntry struct {
	Section_type                                     string
	Dates, Description, Location, Title, Institution string
	Bulletpoints                                     []string
	Points                                           [][2]string
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
	CV_Builder strings.Builder
	Config     Config
	Header     Header
	Skills     map[string]string
	Section    map[string][]SectionEntry // all sections can live here
}

func (cv *CV) WriteString(str string) {
	cv.CV_Builder.WriteString(str)
}

func (cv *CV) String() string {
	return cv.CV_Builder.String()
}
