package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var ( //these describe how each section is treated. section, projectheading, subheading and resumeitem are marcos defined in preamble.
	settingsfile = "preamble.tex"

	section              = "\\section"
	resumeProjectHeading = "\\resumeProjectHeading"
	resumeSubheading     = "\\resumeSubheading"
	resumeSubSubHeading  = "\\resumeSubSubheading"
	resumeItem           = "\\resumeItem"

	resumeSubHeadingListStart = "\\begin{itemize}[leftmargin=0.15in, label={}]"
	resumeSubHeadingListEnd   = "\\end{itemize}"

	resumeItemListStart = "\\begin{itemize}"
	resumeItemListEnd   = "\\end{itemize}\\vspace{-5pt}"

	resumeListSectionStart = "\\begin{itemize}[leftmargin=0.15in, itemsep=-2pt]"
	resumeListSectionEnd   = "\\end{itemize}"

	large_section_seperator = "\n\n\n"
)

var (
	ErrNoInput = errors.New("missing path to input file")
)

func getDataFromFile(res *Resume, path string) error {
	fi, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}
	err = toml.Unmarshal(fi, res.Data)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal toml file: %w", err)
	}
	return nil
}

func main() {
	var inputPath string
	var outputPath string

	if len(os.Args) > 1 {
		inputPath = os.Args[1]
	}

	if inputPath == "" {
		fmt.Fprint(os.Stderr, ErrNoInput)
		os.Exit(1)
	}
	if i := strings.LastIndex(".", inputPath); i >= 0 {
		outputPath = inputPath[:i] + "tex"
	}

	resume := NewResume()
	err := getDataFromFile(resume, inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load data from specified file: %v", err)
	}

	resume.ValidateConfig()
	resume.WriteDocSettings()

	resume.WriteString("\\begin{document}\n")

	resume.WriteHeader()
	for _, section := range resume.Data.Config.SectionOrder {
		resume.WriteSection(section)
	}

	resume.WriteString("\\end{document}")

	if outputPath == "" {
		fmt.Println(resume.String())
	} else {
		f, err := os.Create(outputPath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = f.WriteString(resume.String())
		if err != nil {
			panic(err)
		}
	}
}

func (r *Resume) WriteDocSettings() {
	fmt.Fprintf(&r.Builder, "\\documentclass[%dpt, a4paper]{article}\n", r.Data.Config.FontSize)
	fmt.Fprintf(&r.Builder, "\\input{%s}\n", settingsfile)
	fmt.Fprintf(&r.Builder, "\\usepackage[margin=%fcm]{geometry}\n", r.Data.Config.PageMargin)
	fmt.Fprintf(&r.Builder, "\\setmainfont[Scale=%f]{%s}\n", r.Data.Config.FontScale, r.Data.Config.FontName)
}

func (r *Resume) ValidateConfig() {

	var ( // default values
		default_font              = "Calibri"
		default_font_size         = 12
		default_font_scale        = 1.0
		default_page_margin       = 1.5
		default_cv_order          = []string{"header", "skills", "experience", "education", "projects"}
		default_experience_header = []string{"title", "dates", "institution", "locaiton"}
		default_education_header  = []string{"title", "dates", "institution", "locaiton"}
		default_project_header    = []string{"title", "dates"}
	)

	//NOTE: validate page settings:
	if r.Data.Config.FontName == "" {
		fmt.Printf("font not set, defaulting to: %v.\n", default_font)
		r.Data.Config.FontName = default_font
	}

	if r.Data.Config.FontSize < 10 || r.Data.Config.FontSize > 12 {
		fmt.Printf("font size not set or invalid, defaulting to: %v.\n", default_font_size)
		r.Data.Config.FontSize = default_font_size
	}

	if r.Data.Config.FontScale <= 0.1 {
		fmt.Printf("font scale not set or invalid, defaulting to: %v.\n", default_font_scale)
		r.Data.Config.FontScale = default_font_scale
	}

	if r.Data.Config.PageMargin <= 0 {
		fmt.Printf("page margin not set or invalid, defaulting to: %v.\n", default_page_margin)
		r.Data.Config.PageMargin = default_page_margin
	}

	//NOTE: validate orders
	if len(r.Data.Config.SectionOrder) == 0 {
		fmt.Printf("cv_order is not set, defaulting to: %v.\n", default_cv_order)
		r.Data.Config.SectionOrder = default_cv_order
	}

	//NOTE: process headers, default or pad then trim
	if len(r.Data.Config.ExperienceHeadersOrder) == 0 {
		fmt.Printf("experience_header_order is not set, defaulting to: %v.\n", default_experience_header)
		r.Data.Config.ExperienceHeadersOrder = default_experience_header
	} else {
		r.Data.Config.ExperienceHeadersOrder = append(r.Data.Config.ExperienceHeadersOrder, "", "", "")[:4] //avoid processing under and over cases
	}

	if len(r.Data.Config.EducationHeadersOrder) == 0 {
		fmt.Printf("education_header_order is not set, defaulting to: %v.\n", default_education_header)
		r.Data.Config.EducationHeadersOrder = default_education_header
	} else {
		r.Data.Config.EducationHeadersOrder = append(r.Data.Config.EducationHeadersOrder, "", "", "")[:4] //avoid processing under and over cases
	}

	if len(r.Data.Config.ProjectHeadersOrder) == 0 {
		fmt.Printf("project_header_order is not set, defaulting to: %v.\n", default_project_header)
		r.Data.Config.ProjectHeadersOrder = default_project_header
	} else {
		r.Data.Config.ProjectHeadersOrder = append(r.Data.Config.ProjectHeadersOrder, "")[:2] //avoid processing under and over cases
	}

}
func (r *Resume) WriteBulletpointsTo(sb *strings.Builder, entry SectionEntry) {
	if len(entry.Bulletpoints) != 0 {
		sb.WriteString(resumeItemListStart + "\n")
		for _, item := range entry.Bulletpoints {
			fmt.Fprintf(sb, "	%s{%s}\n", resumeItem, item)
		}
		sb.WriteString(resumeItemListEnd + "\n")
	}
}

func (r *Resume) WriteSection(section_title string) {

	//validate input
	if len(r.Data.Section[section_title]) == 0 {
		return
	}

	r.WriteString(fmt.Sprintf("\\section{%s}\n", section_title))

	subheading_count := 0
	section := &strings.Builder{}

	for _, entry := range r.Data.Section[section_title] {
		switch section_type := strings.ToLower(entry.SectionType); section_type {
		case "project":
			r.WriteProjectEntryTo(section, entry, r.Data.Config.ProjectHeadersOrder)
			subheading_count++
		case "education":
			r.WriteExperienceEntryTo(section, entry, r.Data.Config.EducationHeadersOrder)
			subheading_count++
		case "experience":
			r.WriteExperienceEntryTo(section, entry, r.Data.Config.ExperienceHeadersOrder)
			subheading_count++
		case "subexperience":
			r.WriteSubExperienceEntryTo(section, entry, r.Data.Config.ExperienceHeadersOrder)
			subheading_count++
		case "list": // these 2 do not have headings
			r.WriteListSectionTo(section, entry)
		case "points":
			r.WritePointsSectionTo(section, entry)
		default: // is default
			fmt.Printf("section_type '%s' under section '%s' is not a valid option. Try 'project', 'education', 'experience', 'bulletpoints','list'\n", section_type, section_title)
		}
	}

	if subheading_count != 0 { // ensure no empty environments in LaTeX
		r.WriteString(resumeSubHeadingListStart + "\n")
		r.WriteString(section.String())
		r.WriteString(resumeSubHeadingListEnd + "\n")
	} else {
		r.WriteString(section.String())
	}
}

func (r *Resume) WriteExperienceEntryTo(sb *strings.Builder, exp SectionEntry, header_format []string) {

	// process subheading, parse order:
	sb.WriteString(resumeSubheading)
	for _, entry := range header_format { // only accept the first 4 inputs
		sb.WriteString("{")
		switch strings.ToLower(entry) {
		case "title":
			sb.WriteString(exp.Title)
		case "dates":
			sb.WriteString(exp.Dates)
		case "institution":
			sb.WriteString(exp.Institution)
		case "location":
			sb.WriteString(exp.Location)
		}
		sb.WriteString("}")
	}
	sb.WriteString("\n")
	r.WriteBulletpointsTo(sb, exp)
}

func (r *Resume) WriteSubExperienceEntryTo(sb *strings.Builder, exp SectionEntry, header_format []string) {

	// process subheading, parse order:
	sb.WriteString(resumeSubSubHeading)
	for _, entry := range header_format { // only accept the first 4 inputs
		sb.WriteString("{")
		switch strings.ToLower(entry) {
		case "title":
			sb.WriteString(exp.Title)
		case "dates":
			sb.WriteString(exp.Dates)
		}
		sb.WriteString("}")
	}
	sb.WriteString("\n")
	r.WriteBulletpointsTo(sb, exp)
}

func (r *Resume) WriteProjectEntryTo(sb *strings.Builder, project SectionEntry, headerOrder []string) {

	//NOTE: validate input:
	if len(project.Title) == 0 {
		return
	}

	if project.Description != "" {
		fmt.Fprintf(sb, "%s{\\textbf{%s} $|$ \\textit{%s}}{%s}\n", resumeProjectHeading, project.Title, project.Description, project.Dates)
	} else { // add "| <desc>" or don't if empty.
		fmt.Fprintf(sb, "%s{\\textbf{%s}}{%s}\n", resumeProjectHeading, project.Title, project.Dates)
	}

	r.WriteBulletpointsTo(sb, project)
}

func (r *Resume) WriteHeader() {
	r.Builder.WriteString("\\begin{center}\n")

	r.Builder.WriteString(fmt.Sprintf("\\fontsize{%dpt}{12pt}\\selectfont \\textbf{%s}\\\\ \\vspace{1pt}\n", r.Data.Header.NameSize, r.Data.Header.Name))
	r.Builder.WriteString("\\small")

	for i, link_entry := range r.Data.Header.Details {

		switch {
		case len(link_entry) >= 2:
			r.Builder.WriteString(fmt.Sprintf("\\href{%s}{\\underline{%s}}", link_entry[1], link_entry[0]))
		case len(link_entry) == 1:
			r.Builder.WriteString(fmt.Sprintf("\\underline{%s}", link_entry[0]))
		}

		if i != len(r.Data.Header.Details)-1 { // add seperators while not final entry
			r.Builder.WriteString(" $|$ ")
		}
	}
	r.Builder.WriteString("\n\\end{center}\n")
}

func (r *Resume) WritePointsSectionTo(sb *strings.Builder, entry SectionEntry) {

	//NOTE: validate input
	if len(entry.Points) == 0 {
		return
	}

	sb.WriteString("\\begin{itemize}[leftmargin=0.15in, label={}]\n")
	sb.WriteString("\\small{\\item{\n")

	for _, entry := range entry.Points {
		fmt.Fprintf(sb, "\\textbf{%s}: %s\\\\ \n", entry[0], entry[1])
	}
	sb.WriteString("}}\n")
	sb.WriteString("\\end{itemize}")
	sb.WriteString(large_section_seperator)

}

func (r *Resume) WriteListSectionTo(sb *strings.Builder, entry SectionEntry) {
	//NOTE: validate input
	if len(entry.Bulletpoints) == 0 {
		return
	}

	sb.WriteString("\\begin{itemize}[leftmargin=0.15in, itemsep=-2pt]\n")
	sb.WriteString("\\small{\n")

	for _, item := range entry.Bulletpoints {
		sb.WriteString("\\item{")
		sb.WriteString(item)
		sb.WriteString("}\n")
	}

	sb.WriteString("}\n")
	sb.WriteString("\\end{itemize}")
	sb.WriteString(large_section_seperator)

}
