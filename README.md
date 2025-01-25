# Description
A generate easy to write `latex` typset resumes using a `toml` configuration file, parser is written in `Go`. 
The resume template is taken from [This Repo](https://github.com/jakegut/resume/), being minorly changed to allow for any font use.

# Usage Guide

Cofigure the `.toml` file as desired, then run the parser, resulting in a `.tex` file.
To compile the `.tex` file, ensure LaTeX is installed, the chosen font is installed as part of your system and compile using `XeLaTeX` or `LuaLaTeX`.


![example](https://github.com/eitanoid/toml-resume/blob/main/examples/examplecv.png)

## Running the program:

- Git clone this repo into a directory.

- Install all external dependancies:
  `go get github.com/pelletier/go-toml/v2`

- Build the Go project by running `go build`

- Run the binary with the flags: `./tomlresume -input="input.toml" -output="output.tex"`

- Compile the resulting `.tex` file using `xelatex` or `lualatex` for example using latexmk: `latexmk -xelatex out.tex`. Alternatively this step can be done on Overleaf.
<details>
  <summary>Compiling the output file with Overleaf</summary>
    
### Guide:
1. After running the toml interpreter, create a new Overleaf project and upload `preamble.tex` and your `output.tex` file.
2. Upload your desired font files into the Overleaf document (eg. calibri-xyz.tff).
3. Inside the `output.tex` (or whichever name you gave it) document, changed the line:
```tex
\setmainfont[
	...
]{Calibri} % Where Calibri can be any system font name.
```
To
```tex
\setmainfont[
    ...
    BoldFont=calibri-bold.ttf,
    ItalicFont=calibri-italic.ttf,
    BoldItalicFont=calibri-bold-italic.ttf]{calibri-regular.ttf} % each being the path to the corresponding font file in the Overleaf project.
```

4. In the settings menu, change the rendering engine from 'PDFLaTeX' to 'XeLaTeX' or 'LuaLaTeX'.

5. Compile the document to get an output pdf.

</details>

## Config guide:

The format and settings of the output structure can be configured:

```toml
[config]
font_size=10 # 10, 11, 12 as per latex's document class sizes
font_scale=1 # multiplicative font scaling
page_margin=1.5 #in centimeters
font="Calibri" #your system font of choice

# Case sensitive. Specify each section exactly as you defined it. Repeats are allowed.
section_order=["Technical Skills","Education","Work Experience","Projects","Hobbies and Interests"] 

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
section_type="experience" #other options are: "Education", "Project", "List", "Points". Not case sensitive.
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
section_type="experience"
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
section_type="education" # education entries accept the same values as experience ones, only differing on the header order set in the config section.
title="Your Degree"
institution="University of Universities"
location="Place, Locaiton"
dates="Date - Another Date"
bulletpoints=[] # bulletpoints can be empty.


[[section."Projects"]] 
section_type="project" # project headers are shorter and only accept by a title and a date.
title="This is my first project"
dates="Date - Another Date"
bulletpoints=[
			"Really cool project, does really cool things.",
			]

[[section."Projects"]] 
section_type="project" 
title="This is my first project"
description="Tools used in this project" # An inline description or summary may be added.
dates="Date - Another Date"
bulletpoints=[
			"Really cool project, does really cool things.",
			]


[[section."Technical Skills"]] # this section is composed of title-description pairs with each point's 'title' being displayed in bold.
section_type="points" 
Points=[
    ["These points","Are stored as key-value pairs."],
    [ "This test is bold","this text isn't." ],
]

[[section."Hobbies and Interests"]] # this section is composed of bulletpoints only.
section_type="list"
bulletpoints=[
	"This is a normal list.",
	"each entry is displayed in a new line."
]
```
