package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// mapping for latex macros or formatting
const (
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

func (r *Resume) loadDataFromFile(path string) error {
	fi, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}
	err = toml.Unmarshal(fi, r.Data)
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

	outputPath = strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".tex"

	resume := NewDefaultResume()
	if err := resume.loadDataFromFile(inputPath); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load data from specified file: %v", err)
	}
	if err := resume.ValidateConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid configuration: %v", err)
	}

	resume.CreateLatexDoc()

	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create output file \"%s\": %v", outputPath, err)
			os.Exit(1)
		}
		defer f.Close()
		if _, err := f.WriteString(resume.String()); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to output file \"%s\": %v", outputPath, err)
			os.Exit(1)
		}
	} else {
		fmt.Println(resume.String())
	}

}

func (r *Resume) CreateLatexDoc() {
	r.WriteDocSettings()
	r.WriteString("\\begin{document}\n")
	r.WriteHeader()
	for _, section := range r.Data.Config.SectionOrder {
		r.ProcessSection(section)
	}
	r.WriteString("\\end{document}")
}

func (r *Resume) WriteDocSettings() {
	fmt.Fprintf(&r.Builder, "\\documentclass[%dpt, a4paper]{article}\n", r.Data.Config.FontSize)
	fmt.Fprintf(&r.Builder, "\\input{%s}\n", settingsfile)
	fmt.Fprintf(&r.Builder, "\\usepackage[margin=%fcm]{geometry}\n", r.Data.Config.PageMargin)
	fmt.Fprintf(&r.Builder, "\\setmainfont[Scale=%f]{%s}\n", r.Data.Config.FontScale, r.Data.Config.FontName)
}

func (r *Resume) WriteHeader() {

	r.Builder.WriteString("\\begin{center}\n")
	fmt.Fprintf(&r.Builder, "\\fontsize{%dpt}{12pt}\\selectfont \\textbf{%s}\\\\ \\vspace{1pt}\n", r.Data.Header.NameSize, r.Data.Header.Name)
	r.Builder.WriteString("\\small")

	for i, link_entry := range r.Data.Header.Details {

		switch {
		case len(link_entry) >= 2:
			r.Builder.WriteString(fmt.Sprintf("\\href{%s}{\\underline{%s}}", link_entry[1], link_entry[0]))
		case len(link_entry) == 1:
			r.Builder.WriteString(fmt.Sprintf("\\underline{%s}", link_entry[0]))
		}

		if i != len(r.Data.Header.Details)-1 { // add seperators
			r.Builder.WriteString(" $|$ ")
		}
	}
	r.Builder.WriteString("\n\\end{center}\n")
}

func (r *Resume) ValidateConfig() error {

	if r.Data.Config.FontSize < 10 || r.Data.Config.FontSize > 12 {
		return errors.New("font size must be between 10 and 12")
	}

	// pad to handle exceptions
	r.Data.Config.ExperienceHeadersOrder = append(r.Data.Config.ExperienceHeadersOrder, "", "", "")[:4]
	r.Data.Config.EducationHeadersOrder = append(r.Data.Config.EducationHeadersOrder, "", "", "")[:4]
	r.Data.Config.ProjectHeadersOrder = append(r.Data.Config.ProjectHeadersOrder, "")[:2]

	return nil

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

func (r *Resume) ProcessSection(title string) {

	//validate input
	if len(r.Data.Section[title]) == 0 {
		return
	}

	fmt.Fprint(&r.Builder, section, "{", title, "}\n")

	subheading_count := 0
	section := &strings.Builder{}
	for _, entry := range r.Data.Section[title] {
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
			fmt.Printf("section_type '%s' under section '%s' is not a valid option. Try 'project', 'education', 'experience', 'bulletpoints','list'\n", section_type, title)
		}
	}

	if subheading_count != 0 { // ensure no empty environments in LaTeX
		fmt.Fprint(&r.Builder, resumeSubHeadingListStart+"\n", section.String(), resumeSubHeadingListEnd+"\n")
	} else {
		r.WriteString(section.String())
	}
}

func (r *Resume) WriteExperienceEntryTo(sb *strings.Builder, exp SectionEntry, headerOrder []string) {

	// process subheading, parse order:
	sb.WriteString(resumeSubheading)
	for _, entry := range headerOrder { // only accept the first 4 inputs
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

func (r *Resume) WriteSubExperienceEntryTo(sb *strings.Builder, exp SectionEntry, headerOrder []string) {

	// process subheading, parse order:
	sb.WriteString(resumeSubSubHeading)
	for _, entry := range headerOrder { // only accept the first 4 inputs
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

func (r *Resume) WritePointsSectionTo(sb *strings.Builder, entry SectionEntry) {

	//NOTE: validate input
	if len(entry.Points) == 0 {
		return
	}

	fmt.Fprint(sb,
		"\\begin{itemize}[leftmargin=0.15in, label={}]\n",
		"\\small{\\item{\n",
	)

	for _, entry := range entry.Points {
		fmt.Fprintf(sb, "\\textbf{%s}: %s\\\\ \n", entry[0], entry[1])
	}
	fmt.Fprint(sb,
		"}}\n",
		"\\end{itemize}",
		large_section_seperator,
	)

}

func (r *Resume) WriteListSectionTo(sb *strings.Builder, entry SectionEntry) {
	//NOTE: validate input
	if len(entry.Bulletpoints) == 0 {
		return
	}

	fmt.Fprint(sb,
		"\\begin{itemize}[leftmargin=0.15in, itemsep=-2pt]\n",
		"\\small{\n",
	)

	for _, item := range entry.Bulletpoints {
		fmt.Fprint(sb, "\\item{", item, "}\n")
	}

	fmt.Fprint(sb,
		"}\n",
		"\\end{itemize}",
		large_section_seperator,
	)

}
