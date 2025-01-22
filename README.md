# Description
A generate easy to write `latex` typset resumes using a `toml` configuration file, parser is written in `Go`. 
The resume template is taken from [This Repo](https://github.com/jakegut/resume/), being minorly changed to allow for any font use.

# Usage Guide

Cofigure the `.toml` file as desired, then run the parser, resulting in a `.tex` file.
To compile the `.tex` file, ensure LaTeX is installed, the chosen font is installed as part of your system and compile using `XeLaTeX` or `LuaLaTeX`.
Example command to render: `latexmk -xelatex cv.tex`

![example](https://github.com/eitanoid/toml-resume/blob/main/examples/examplecv.png)

## TOML guide:


The format and settings of the output structure can be configured:
```toml
[config]
font_size=10 # 10, 11, 12 as per latex's document class sizes
font_scale=1 # multiplicative font scaling
page_margin=1.5 #in centimeters
font="Calibri" #your system font of choice

# Case sensitive. Specify each section exactly as you defined it. Repeats are allowed.
section_order=["Technical Skills","Education","Work Experience","Projects"] 

# Format of each entry types's headers, play around with it until you're happy!
education_header_order=["institution", "dates", "title", "location"]
experience_header_order=["institution","dates","title","location"]
project_header_order=["title, dates"]
```

As part of the `.toml` file, the resume structure is as follows:

```toml
[header]
#settings you can play around with:
header_format=["email", "linkedin","github", "phone"]
name_size=14 #font size in pt.

#your personal details:
name="Your Full Name"
location="Location, Earth"
phone="+xxx xxxxx xxxxxx"
email="your\\_email@email.com" 
# as the document is interpreted into tex, special characters like '_' '&' '%' in must be escaped by adding a '\' before them.
# \ beings an escape sequence in the toml specification. To add a '\' into the tex code, we must escape the '\'. That is '\\' interpretes into '\'.
# and '\\_' interpretes in to '\_' in tex, and to '_' in the final document.
linkedin="linkedin.com/your-linkedin"
github="github.com/your-github"

[[section."Work Experience"]] # creates the section "Work Experience" and adds an entry into it.
header_style="experience" #other options are: "Education", "Project". not setting this defaults to "experience".
title="Your Job Title"
institution="Company Name"
location="Place, Locaiton"
dates="November 2022 - December 2024"
bulletpoints=[
			"Accomplishment 1.",
			"Accomplishment 2.",
			"Accomplishment 3.",
			]


[[section."Work Experience"]] # this is another addition to "Work Experience".
header_style="experience"
title="Your Job Title"
institution="2nd Company Name"
location="Place, Locaiton"
dates="Date - Another Date"
bulletpoints=[
			"Accomplishment 1 at 2nd company.",
			"Accomplishment 2 at 2nd company.",
			"Accomplishment 3 at 2nd company.",
			]

[[section."Education"]] 
header_style="education" # education entries are identical to experience ones, only differing on the header order config.
title="Your Degree"
institution="University of Universities"
location="Place, Locaiton"
dates="Date - Another Date"
bulletpoints=[] # bulletpoints can be empty.


[[section."Projects"]] 
header_style="project" # project headers are shorter and only determined by a title and a date.
title="This is my first project"
dates="Date - Another Date"
bulletpoints=[
			"Really cool project, does really cool things.",
			]

[[section."Projects"]] 
header_style="project" 
title="This is my first project"
description="Tools used in this project" # Inline description or summary.
dates="Date - Another Date"
bulletpoints=[
			"Really cool project, does really cool things.",
			]


[[section."Technical Skills"]] # this section is composed of bulletpoints only. Any other entries will be ignored when Points are present.
Points."These points"="Are stored as key-value pairs."
Points."This test is bold"="this text isn't."
```

# To do list:
- Parse font, font size etc into the preamble document.
- 
