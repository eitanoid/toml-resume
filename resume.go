package main

import "strings"

type Resume struct {
	Builder strings.Builder
	Data    *RawData
}

func NewResume() *Resume {
	return &Resume{
		Builder: strings.Builder{},
		Data:    &RawData{},
	}
}

type Config struct {
	FontSize               int      `toml:"font_size"`
	FontScale              float64  `toml:"font_scale"`
	PageMargin             float64  `toml:"page_margin"`
	FontName               string   `toml:"font"`
	SectionOrder           []string `toml:"section_order"`
	ProjectHeadersOrder    []string `toml:"project_header_order"`
	ExperienceHeadersOrder []string `toml:"experience_header_order"`
	EducationHeadersOrder  []string `toml:"education_header_order"`
}

type Header struct {
	Details  [][]string `toml:"details"`
	Name     string     `toml:"name"`
	NameSize int        `toml:"name_size"`
}

type SectionEntry struct {
	SectionType  string      `toml:"section_type"`
	Dates        string      `toml:"dates"`
	Description  string      `toml:"description"`
	Location     string      `toml:"location"`
	Title        string      `toml:"title"`
	Institution  string      `toml:"institution"`
	Bulletpoints []string    `toml:"bulletpoints"`
	Points       [][2]string `toml:"points"`
}

type RawData struct {
	Config  Config                    `toml:"config"`
	Header  Header                    `toml:"header"`
	Skills  map[string]string         // NOTE: what does this to
	Section map[string][]SectionEntry `toml:"section"`
}

func (r *Resume) WriteString(str string) {
	r.Builder.WriteString(str)
}

func (r *Resume) String() string {
	return r.Builder.String()
}
