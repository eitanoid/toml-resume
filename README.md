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
```nix
# impure, if your system has your specified font
nix develop github:eitanoid/toml-resume --command gen-resume <file>.toml

# pure, use flake with bundled dependancies
nix develop --ignore-environment --impure --expr '
  let
    flake = builtins.getFlake github:eitanoid/toml-resume;
    pkgs = flake.inputs.nixpkgs.legacyPackages.''${builtins.currentSystem};
  in
    flake.lib.''${builtins.currentSystem}.mkTexShell {
      extraFonts = [
        # nixpkgs package with any font you need
        pkgs.source-serif
        pkgs.vista-fonts
      ];
    }
' --command gen-resume <file>.toml
```

> [!NOTE]
> `nix-command` and `flakes` need to be enabled as experimental features either in a config file or with the `--extra-experimental-features` flag.

## Config guide:

Check out the [examples](https://github.com/eitanoid/toml-resume/blob/main/examples) directory!
