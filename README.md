# Description
Easily generate beautiful `LaTeX` typset resumes using a `toml` or `yaml` configuration file.

The resume template is taken from [this repo](https://github.com/jakegut/resume/), and updated to use arbitrary fonts.

![example](https://github.com/eitanoid/toml-resume/blob/main/examples/example.jpg)

# Usage Guide

Install the program:
`go install github.com/eitanoid/toml-resume@main`

Run the program on a file:
```bash
# standard use
$ toml-resume example.toml
> Successfully generated: example.tex

# output to stdout
$ toml-resume example.toml -o -

# read from stdin or redirection (output to stdout by default)
$ cat file.toml | toml-resume
$ toml-resume < example.toml

# specify output path
$ toml-resume example.toml -o custompath
> Successfully generated: custompath
```

## Run and Compile LaTeX with Nix

Create a pdf document from the config file by running:
```bash
nix develop github:eitanoid/toml-resume --command gen-resume <file>.toml
```

## Config guide:

Check out the [examples](https://github.com/eitanoid/toml-resume/blob/main/examples) directory!
