package main

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"strings"
)

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

	title_index := 0
	for _, section := range cv_args.Config.Cv_order {
		section_title := cv_args.Config.Section_titles[title_index]

		switch strings.ToLower(section) {

		case "experience":
			WriteExperience(&cv_args, &cv_builder, section_title)
			title_index++
		case "education":
			WriteEducation(&cv_args, &cv_builder, section_title)
			title_index++
		case "projects":
			WriteProjects(&cv_args, &cv_builder, section_title)
			title_index++
		case "skills":
			WriteSkills(&cv_args, &cv_builder, section_title)
			title_index++
		case "header":
			WriteHeader(&cv_args, &cv_builder)
		default:
			fmt.Printf("%s is not an accepted order, ignoring. Try 'header', 'experience', 'education', 'projects', skills'.\n", section)
		}

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

	if len(cv_args.Config.Cv_order) == 0 {
		fmt.Printf("cv_order is not set, defaulting to: %v.\n", default_cv_order)
		cv_args.Config.Cv_order = default_cv_order
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

func WriteExperience(cv_args *CV, cv_builder *strings.Builder, section_title string) {

	//NOTE: validate input
	if len(cv_args.Experience) == 0 {
		return
	}

	WriteExperienceEntry := func(exp Experience, cv_builder *strings.Builder) {

		// process subheading, parse order:
		cv_builder.WriteString(resumeSubheading)
		for _, entry := range cv_args.Config.Experience_header_order { // only accept the first 4 inputs
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
		cv_builder.WriteString(resumeItemListStart)
		cv_builder.WriteString("\n")
		for _, item := range exp.Bulletpoints {
			cv_builder.WriteString(fmt.Sprintf("	%s{%s}\n", resumeItem, item))
		}
		cv_builder.WriteString(resumeItemListEnd)
		cv_builder.WriteString("\n\n")
	}

	cv_builder.WriteString(fmt.Sprintf("%s{%s}\n", section, section_title)) //TODO: experience is controlled by the config

	cv_builder.WriteString(resumeSubHeadingListStart)
	cv_builder.WriteString("\n")

	for _, exp := range cv_args.Experience {
		WriteExperienceEntry(exp, cv_builder)
	}

	cv_builder.WriteString(resumeSubHeadingListEnd)
	cv_builder.WriteString(large_section_seperator)
}

func WriteEducation(cv_args *CV, cv_builder *strings.Builder, section_title string) {

	if len(cv_args.Education) == 0 {
		return
	}

	WriteEducationEntry := func(edu Education, cv_builder *strings.Builder) {

		//Process subheading
		cv_builder.WriteString(resumeSubheading)
		for _, entry := range cv_args.Config.Education_header_order {
			cv_builder.WriteString("{")
			switch strings.ToLower(entry) {
			case "title":
				cv_builder.WriteString(edu.Title)
			case "dates":
				cv_builder.WriteString(edu.Dates)
			case "institution":
				cv_builder.WriteString(edu.Institution)
			case "location":
				cv_builder.WriteString(edu.Location)
			}
			cv_builder.WriteString("}")
		}
		cv_builder.WriteString("\n")

		cv_builder.WriteString(resumeItemListStart)
		cv_builder.WriteString("\n")

		for _, item := range edu.Bulletpoints {
			cv_builder.WriteString(fmt.Sprintf("	%s{%s}\n", resumeItem, item))
		}
		cv_builder.WriteString(resumeItemListEnd)
		cv_builder.WriteString("\n\n")
	}

	cv_builder.WriteString(fmt.Sprintf("%s{%s}\n", section, section_title)) //TODO: experience is controlled by the config

	cv_builder.WriteString(resumeSubHeadingListStart)
	cv_builder.WriteString("\n")

	for _, edu := range cv_args.Education {
		WriteEducationEntry(edu, cv_builder)
	}
	cv_builder.WriteString(resumeSubHeadingListEnd)
	cv_builder.WriteString(large_section_seperator)
}

func WriteProjects(cv_args *CV, cv_builder *strings.Builder, section_title string) {

	if len(cv_args.Project) == 0 {
		return
	}

	WriteProjectEntry := func(project Project, cv_builder *strings.Builder) {
		if project.Description != "" {
			cv_builder.WriteString(fmt.Sprintf("%s{\\textbf{%s} $|$ \\textit{%s}}{%s}\n", resumeProjectHeading, project.Title, project.Description, project.Dates))
		} else {
			cv_builder.WriteString(fmt.Sprintf("%s{\\textbf{%s}}{%s}\n", resumeProjectHeading, project.Title, project.Dates)) //TODO: should redo resumeProjectHeading, dont like using \textbf and \textit here
		}

		cv_builder.WriteString(resumeItemListStart)
		cv_builder.WriteString("\n")

		for _, item := range project.Bulletpoints {
			cv_builder.WriteString(fmt.Sprintf("	%s{%s}\n", resumeItem, item))
		}
		cv_builder.WriteString(resumeItemListEnd)
		cv_builder.WriteString("\n\n")
	}

	cv_builder.WriteString(fmt.Sprintf("%s{%s}\n", section, section_title)) //TODO: experience is controlled by the config

	cv_builder.WriteString(resumeSubHeadingListStart)
	cv_builder.WriteString("\n")

	for _, project := range cv_args.Project {
		WriteProjectEntry(project, cv_builder)
	}
	cv_builder.WriteString(resumeSubHeadingListEnd)
	cv_builder.WriteString(large_section_seperator)
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

func WriteSkills(cv_args *CV, cv_builder *strings.Builder, section_title string) {

	//NOTE: validate input
	if len(cv_args.Skills) == 0 {
		return
	}

	cv_builder.WriteString(fmt.Sprintf("\\section{%s}\n", section_title)) //TODO: key skills to be determined by config file

	cv_builder.WriteString("\\begin{itemize}[leftmargin=0.15in, label={}]\n")
	cv_builder.WriteString("\\small{\\item{\n")

	for title, entry := range cv_args.Skills {
		cv_builder.WriteString(fmt.Sprintf("\\textbf{%s}: %s\\\\ \n", title, entry))
	}
	cv_builder.WriteString("}}\n")
	cv_builder.WriteString("\\end{itemize}")
	cv_builder.WriteString(large_section_seperator)
}

// func WriteHobbies(cv_args *CV, cv_builder *strings.Builder) {
//
// 	//NOTE: validate input
// 	if len(cv_args.Skills) == 0 {
// 		return
// 	}
//
// 	cv_builder.WriteString(fmt.Sprintf("\\section{%s}\n", "Key Skills")) //TODO: key skills to be determined by config file
//
// 	cv_builder.WriteString("\\begin{itemize}[leftmargin=0.15in, label={}]\n")
// 	cv_builder.WriteString("\\small{\\item{\n")
//
// 	for title, entry := range cv_args.Skills {
// 		cv_builder.WriteString(fmt.Sprintf("\\textbf{%s}: %s\\\\ \n", title, entry))
// 	}
// 	cv_builder.WriteString("}}\n")
// 	cv_builder.WriteString("\\end{itemize}")
// 	cv_builder.WriteString(large_section_seperator)
// }
