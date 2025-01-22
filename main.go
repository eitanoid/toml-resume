package main

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"strings"
)

//TODO:
//output to a file
// change font settings in toml
//TODO: redo how I do sections, make sections maps, that way [[section]] is enoguh to add a section
//TODO: Decide how to do latex mathmode and special chars.

var ( //TODO: redo the template Im using, keep these for now because they are convenient
	section                   = "\\section"
	subheading                = "\\resumeSubheading"
	resumeItem                = "\\resumeItem"
	resumeSubheading          = "\\resumeSubheading"
	resumeItemListStart       = "\\resumeItemListStart"
	resumeItemListEnd         = "\\resumeItemListEnd"
	resumeSubHeadingListStart = "\\resumeSubHeadingListStart"
	resumeSubHeadingListEnd   = "\\resumeSubHeadingListEnd"
	resumeProjectHeading      = "\\resumeProjectHeading"
	large_section_seperator   = "\n\n\n"
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
	cv_args := ReadTOML("test.toml")
	cv_builder := strings.Builder{}

	cv_builder.WriteString("\\input{preamble.tex}\n")
	cv_builder.WriteString("\\begin{document}\n")

	WriteHeader(&cv_args, &cv_builder)
	for _, section := range cv_args.Config.Section_order {
		WriteSection(&cv_args, &cv_builder, section)
	}

	cv_builder.WriteString("\\end{document}")
	fmt.Println(cv_builder.String())
}

func ValidateConfig(cv_args *CV) {

	var ( // default values
		// default_font              = "Calibri"
		// default_font_size         = 12
		// default_font_scale        = 1.0
		// default_page_margin		 = 1.5
		default_cv_order          = []string{"header", "skills", "experience", "education", "projects"}
		default_experience_header = []string{"title", "dates", "institution", "locaiton"}
		default_education_header  = []string{"title", "dates", "institution", "locaiton"}
		default_project_header    = []string{"title", "dates"}
	)

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

func WriteSection(cv_args *CV, cv_builder *strings.Builder, section_title string) {

	//NOTE: Validate Input
	//--------------------
	if len(cv_args.Section[section_title]) == 0 {
		return
	}

	cv_builder.WriteString(fmt.Sprintf("\\section{%s}\n", section_title))

	if len(cv_args.Section[section_title][0].Points) == 0 { // if points are empty
		cv_builder.WriteString(resumeSubHeadingListStart + "\n")
		for _, entry := range cv_args.Section[section_title] {
			switch header_type := strings.ToLower(entry.Header_style); header_type {

			case "project":
				WriteProjectEntry(entry, cv_builder, cv_args.Config.Project_header_order)
			case "education":
				WriteExperienceEntry(entry, cv_builder, cv_args.Config.Education_header_order)
			default: // "experience" is default
				WriteExperienceEntry(entry, cv_builder, cv_args.Config.Experience_header_order)
			}
		}
		cv_builder.WriteString(resumeSubHeadingListEnd + "\n\n")
	} else {
		for _, entry := range cv_args.Section[section_title] {
			WriteItemList(entry.Points, cv_builder)
		}
	}

}

func WriteExperienceEntry(exp SectionEntery, cv_builder *strings.Builder, header_format []string) {

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

	// process bulletpoints
	cv_builder.WriteString(resumeItemListStart + "\n")

	for _, item := range exp.Bulletpoints {
		cv_builder.WriteString(fmt.Sprintf("	%s{%s}\n", resumeItem, item))
	}
	cv_builder.WriteString(resumeItemListEnd + "\n\n")
}

func WriteProjectEntry(project SectionEntery, cv_builder *strings.Builder, header_format []string) {

	//NOTE: validate input:
	if len(project.Title) == 0 {
		return
	}

	if project.Description != "" {
		cv_builder.WriteString(fmt.Sprintf("%s{\\textbf{%s} $|$ \\textit{%s}}{%s}\n", resumeProjectHeading, project.Title, project.Description, project.Dates))
	} else {
		cv_builder.WriteString(fmt.Sprintf("%s{\\textbf{%s}}{%s}\n", resumeProjectHeading, project.Title, project.Dates)) //TODO: should redo resumeProjectHeading, dont like using \textbf and \textit here
	}

	cv_builder.WriteString(resumeItemListStart + "\n")

	for _, item := range project.Bulletpoints {
		cv_builder.WriteString(fmt.Sprintf("	%s{%s}\n", resumeItem, item))
	}
	cv_builder.WriteString(resumeItemListEnd + "\n\n")
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
	cv_builder.WriteString("\n\\end{center}\n\n")
}

func WriteItemList(items map[string]string, cv_builder *strings.Builder) {

	//NOTE: validate input
	if len(items) == 0 {
		return
	}

	cv_builder.WriteString("\\begin{itemize}[leftmargin=0.15in, label={}]\n")
	cv_builder.WriteString("\\small{\\item{\n")

	for title, entry := range items {
		cv_builder.WriteString(fmt.Sprintf("\\textbf{%s}: %s\\\\ \n", title, entry))
	}
	cv_builder.WriteString("}}\n")
	cv_builder.WriteString("\\end{itemize}")
	cv_builder.WriteString(large_section_seperator)
}
