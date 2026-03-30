package resume

import "strings"

type Resume struct {
	Builder strings.Builder
	Data    *RawData
}

func NewDefaultResume() *Resume {

	defaultConfig := Config{
		FontSize:               12,
		FontScale:              1,
		PageMargin:             1.5,
		FontName:               "Calibri",
		SectionOrder:           []string{"header", "skills", "experience", "education", "projects"},
		ProjectHeadersOrder:    []string{"title", "dates"},
		ExperienceHeadersOrder: []string{"title", "dates", "institution", "location"},
		EducationHeadersOrder:  []string{"title", "dates", "institution", "location"},
	}

	return &Resume{
		Builder: strings.Builder{},
		Data: &RawData{
			Config: defaultConfig,
		},
	}
}

type Config struct {
	FontSize               int      `toml:"font_size" json:"font_size"`
	FontScale              float64  `toml:"font_scale" json:"font_scale"`
	PageMargin             float64  `toml:"page_margin" json:"page_margin"`
	FontName               string   `toml:"font" json:"font"`
	SectionOrder           []string `toml:"section_order" json:"section_order"`
	ProjectHeadersOrder    []string `toml:"project_header_order" json:"project_header_order"`
	ExperienceHeadersOrder []string `toml:"experience_header_order" json:"experience_header_order"`
	EducationHeadersOrder  []string `toml:"education_header_order" json:"education_header_order"`
}

type Header struct {
	Details  [][]string `toml:"details" json:"details"`
	Name     string     `toml:"name" json:"name"`
	NameSize int        `toml:"name_size" json:"name_size"`
}

type SectionEntry struct {
	SectionType  string      `toml:"section_type" json:"section_type"`
	Dates        string      `toml:"dates" json:"dates"`
	Description  string      `toml:"description" json:"description"`
	Location     string      `toml:"location" json:"location"`
	Title        string      `toml:"title" json:"title"`
	Institution  string      `toml:"institution" json:"institution"`
	Bulletpoints []string    `toml:"bulletpoints" json:"bulletpoints"`
	Points       [][2]string `toml:"points" json:"points"`
}

type RawData struct {
	ApiVersion string                    `toml:"apiVersion" json:"apiVersion"`
	Config     Config                    `toml:"config" json:"config"`
	Header     Header                    `toml:"header" json:"header"`
	Section    map[string][]SectionEntry `toml:"section" json:"section"`
}

func (r *Resume) WriteString(str string) {
	r.Builder.WriteString(str)
}

func (r *Resume) String() string {
	return r.Builder.String()
}
