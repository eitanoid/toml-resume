package main

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"strings"
)

var ( //these describe how each section is treated. section, projectheading, subheading and resumeitem are marcos defined in preamble.
	settingsfile = "preamble.tex"

	section              = "\\section"
	resumeProjectHeading = "\\resumeProjectHeading"
	resumeSubheading     = "\\resumeSubheading"
	resumeItem           = "\\resumeItem"

	resumeSubHeadingListStart = "\\begin{itemize}[leftmargin=0.15in, label={}]"
	resumeSubHeadingListEnd   = "\\end{itemize}"

	resumeItemListStart = "\\begin{itemize}"
	resumeItemListEnd   = "\\end{itemize}\\vspace{-5pt}"

	resumeListSectionStart = "\\begin{itemize}[leftmargin=0.15in, itemsep=-2pt]"
	resumeListSectionEnd   = "\\end{itemize}"

	large_section_seperator = "\n\n\n"
)

func ReadTOML(path string) CV {
	var cv_args = CV{}
	fi, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = toml.Unmarshal(fi, &cv_args)
	if err != nil {
		panic(err)
	}

	return cv_args
}

func main() {
	input, output := ParseInput()
	cv_args := ReadTOML(input)
	cv_builder := strings.Builder{}

	cv_args.ValidateConfig()
	WriteDocSettings(&cv_args, &cv_builder)

	cv_builder.WriteString("\\begin{document}\n")

	WriteHeader(&cv_args, &cv_builder)
	for _, section := range cv_args.Config.Section_order {
		WriteSection(&cv_args, &cv_builder, section)
	}

	cv_builder.WriteString("\\end{document}")

	if output == "" {
		fmt.Println(cv_builder.String())
	} else {
		f, err := os.Create(output)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		_, err = f.WriteString(cv_builder.String())
		if err != nil {
			panic(err)
		}
	}
}

func WriteDocSettings(cv_args *CV, cv_builder *strings.Builder) {
	cv_builder.WriteString(fmt.Sprintf("\\documentclass[%dpt, a4paper]{article}\n", cv_args.Config.Font_size))
	cv_builder.WriteString(fmt.Sprintf("\\input{%s}\n", settingsfile))
	cv_builder.WriteString(fmt.Sprintf("\\usepackage[margin=%fcm]{geometry}\n", cv_args.Config.Page_margin))
	cv_builder.WriteString(fmt.Sprintf("\\setmainfont[Scale=%f]{%s}\n", cv_args.Config.Font_scale, cv_args.Config.Font))
}

func (cv_args *CV) ValidateConfig() {

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
	if cv_args.Config.Font == "" {
		fmt.Printf("font not set, defaulting to: %v.\n", default_font)
		cv_args.Config.Font = default_font
	}

	if cv_args.Config.Font_size < 10 || cv_args.Config.Font_size > 12 {
		fmt.Printf("font size not set or invalid, defaulting to: %v.\n", default_font_size)
		cv_args.Config.Font_size = default_font_size
	}

	if cv_args.Config.Font_scale <= 0.1 {
		fmt.Printf("font scale not set or invalid, defaulting to: %v.\n", default_font_scale)
		cv_args.Config.Font_scale = default_font_scale
	}

	if cv_args.Config.Page_margin <= 0 {
		fmt.Printf("page margin not set or invalid, defaulting to: %v.\n", default_page_margin)
		cv_args.Config.Page_margin = default_page_margin
	}

	//NOTE: validate orders
	if len(cv_args.Config.Section_order) == 0 {
		fmt.Printf("cv_order is not set, defaulting to: %v.\n", default_cv_order)
		cv_args.Config.Section_order = default_cv_order
	}

	//NOTE: process headers, default or pad then trim
	if len(cv_args.Config.Experience_header_order) == 0 {
		fmt.Printf("experience_header_order is not set, defaulting to: %v.\n", default_experience_header)
		cv_args.Config.Experience_header_order = default_experience_header
	} else {
		cv_args.Config.Experience_header_order = append(cv_args.Config.Experience_header_order, "", "", "")[:4] //avoid processing under and over cases
	}

	if len(cv_args.Config.Education_header_order) == 0 {
		fmt.Printf("education_header_order is not set, defaulting to: %v.\n", default_education_header)
		cv_args.Config.Education_header_order = default_education_header
	} else {
		cv_args.Config.Education_header_order = append(cv_args.Config.Education_header_order, "", "", "")[:4] //avoid processing under and over cases
	}

	if len(cv_args.Config.Project_header_order) == 0 {
		fmt.Printf("project_header_order is not set, defaulting to: %v.\n", default_project_header)
		cv_args.Config.Project_header_order = default_project_header
	} else {
		cv_args.Config.Project_header_order = append(cv_args.Config.Project_header_order, "")[:2] //avoid processing under and over cases
	}

}
func WriteBulletpoints(entry SectionEntry, cv_builder *strings.Builder) {
	// process bulletpoints
	if len(entry.Bulletpoints) != 0 {
		cv_builder.WriteString(resumeItemListStart + "\n")
		for _, item := range entry.Bulletpoints {
			cv_builder.WriteString(fmt.Sprintf("	%s{%s}\n", resumeItem, item))
		}
		cv_builder.WriteString(resumeItemListEnd + "\n")
	}
}

func WriteSection(cv_args *CV, cv_builder *strings.Builder, section_title string) {

	//NOTE: Validate Input
	if len(cv_args.Section[section_title]) == 0 {
		return
	}

	cv_builder.WriteString(fmt.Sprintf("\\section{%s}\n", section_title))

	subheading_count := 0
	section_builder := strings.Builder{}

	for _, entry := range cv_args.Section[section_title] { // write section
		switch section_type := strings.ToLower(entry.Section_type); section_type {

		case "project":
			WriteProjectEntry(entry, &section_builder, cv_args.Config.Project_header_order)
			subheading_count++
		case "education":
			WriteExperienceEntry(entry, &section_builder, cv_args.Config.Education_header_order)
			subheading_count++
		case "experience":
			WriteExperienceEntry(entry, &section_builder, cv_args.Config.Experience_header_order)
			subheading_count++
		case "list": // these 2 do not have headings
			WriteListSection(entry, &section_builder)
		case "points":
			WritePointsSection(entry, &section_builder)
		default: // is default
			fmt.Printf("section_type '%s' under section '%s' is not a valid option. Try 'project', 'education', 'experience', 'bulletpoints','list'\n", section_type, section_title)
		}
	}

	if subheading_count != 0 { // verify no empty environments in LaTeX, causing an error
		cv_builder.WriteString(resumeSubHeadingListStart + "\n")
		cv_builder.WriteString(section_builder.String())
		cv_builder.WriteString(resumeSubHeadingListEnd + "\n")
	} else {
		cv_builder.WriteString(section_builder.String())
	}
}

func WriteExperienceEntry(exp SectionEntry, cv_builder *strings.Builder, header_format []string) {

	// process subheading, parse order:
	cv_builder.WriteString(resumeSubheading)
	for _, entry := range header_format { // only accept the first 4 inputs
		cv_builder.WriteString("{")
		switch strings.ToLower(entry) {
		case "title":
			cv_builder.WriteString(exp.Title)
		case "dates":
			cv_builder.WriteString(exp.Dates)
		case "institution":
			cv_builder.WriteString(exp.Institution)
		case "location":
			cv_builder.WriteString(exp.Location)
		}
		cv_builder.WriteString("}")
	}
	cv_builder.WriteString("\n")
	WriteBulletpoints(exp, cv_builder)
}

func WriteProjectEntry(project SectionEntry, cv_builder *strings.Builder, header_format []string) {

	//NOTE: validate input:
	if len(project.Title) == 0 {
		return
	}

	if project.Description != "" {
		cv_builder.WriteString(fmt.Sprintf("%s{\\textbf{%s} $|$ \\textit{%s}}{%s}\n", resumeProjectHeading, project.Title, project.Description, project.Dates))
	} else { // add "| <desc>" or don't if empty.
		cv_builder.WriteString(fmt.Sprintf("%s{\\textbf{%s}}{%s}\n", resumeProjectHeading, project.Title, project.Dates))
	}

	WriteBulletpoints(project, cv_builder)
}

func WriteHeader(cv_args *CV, cv_builder *strings.Builder) { //TODO: decide how I want to do this, want to be able to align to any direction in toml
	cv_builder.WriteString("\\begin{center}\n")

	cv_builder.WriteString(fmt.Sprintf("\\fontsize{%dpt}{12pt}\\selectfont \\textbf{%s}\\\\ \\vspace{1pt}\n", cv_args.Header.Name_size, cv_args.Header.Name))
	cv_builder.WriteString("\\small")

	for i, entry := range cv_args.Header.Header_format {

		switch entry {

		case "email":
			cv_builder.WriteString(fmt.Sprintf("\\href{mailto:%s}{\\underline{%s}}", cv_args.Header.Email, cv_args.Header.Email))

		case "linkedin":
			cv_builder.WriteString(fmt.Sprintf("\\href{https://%s}{\\underline{%s}}", cv_args.Header.Linkedin, cv_args.Header.Linkedin))

		case "github":
			cv_builder.WriteString(fmt.Sprintf("\\href{https://%s}{\\underline{%s}}", cv_args.Header.Github, cv_args.Header.Github))

		case "phone":
			cv_builder.WriteString(fmt.Sprintf("%s", cv_args.Header.Phone))
		}

		if i != len(cv_args.Header.Header_format)-1 { // add seperators while not final entry
			cv_builder.WriteString(" $|$ ")
		}
	}
	cv_builder.WriteString("\n\\end{center}\n")
}

func WritePointsSection(entry SectionEntry, cv_builder *strings.Builder) {

	//NOTE: validate input
	if len(entry.Points) == 0 {
		return
	}

	cv_builder.WriteString("\\begin{itemize}[leftmargin=0.15in, label={}]\n")
	cv_builder.WriteString("\\small{\\item{\n")

	for _, entry := range entry.Points {
		cv_builder.WriteString(fmt.Sprintf("\\textbf{%s}: %s\\\\ \n", entry[0], entry[1]))
	}
	cv_builder.WriteString("}}\n")
	cv_builder.WriteString("\\end{itemize}")
	cv_builder.WriteString(large_section_seperator)

}

func WriteListSection(entry SectionEntry, cv_builder *strings.Builder) {
	//NOTE: validate input
	if len(entry.Bulletpoints) == 0 {
		return
	}

	cv_builder.WriteString("\\begin{itemize}[leftmargin=0.15in, itemsep=-2pt]\n")
	cv_builder.WriteString("\\small{\n")

	for _, item := range entry.Bulletpoints {
		cv_builder.WriteString("\\item{")
		cv_builder.WriteString(item)
		cv_builder.WriteString("}\n")
	}

	cv_builder.WriteString("}\n")
	cv_builder.WriteString("\\end{itemize}")
	cv_builder.WriteString(large_section_seperator)

}
