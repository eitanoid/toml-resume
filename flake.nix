{
  description = "A Nix-flake for generating resumes from TOML/YAML";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
    flake-utils.url = "github:numtide/flake-utils";
    toml-resume-src = {
      url = "github:eitanoid/toml-resume";
      flake = false;
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      toml-resume-src,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
        tex = pkgs.texlive.combine {
          inherit (pkgs.texlive)
            latexmk
            scheme-basic
            fontspec
            enumitem
            titlesec
            hyperref
            anyfontsize
            ;
        };
        toml-resume = pkgs.buildGoModule {
          pname = "toml-resume";
          version = "latest";
          src = toml-resume-src;
          vendorHash = "sha256-y3eTrOW1Y3THGaiMUKr4igAfoJgQF9mwLAGYzLcKA0g=";
          env = {
            CGO_ENABLED = 0;
          };
        };
        gen-resume = pkgs.writeShellScriptBin "gen-resume" ''
          if [ -z $1 ]; then
              echo "input is required"
              exit 1
          fi
          # I'm not sure if these are needed. when I was testing earlier, this fixed a lualatex null-font issue
          # export TEXMFVAR=$(mktemp -d)
          # trap 'rm -rf "$TEXMFVAR"' EXIT

          # BASENAME=$(basename $1)
          JOBNAME=''${1%.*}
          toml-resume ''$1 -o out.tex -f | latexmk -interaction=errorstopmode -pdf -lualatex out.tex --jobname=''$JOBNAME
          latexmk -c ''$JOBNAME.pdf  && rm out.tex
        '';
      in
      {
        devShells = {
          default = pkgs.mkShellNoCC {
            packages = [
              toml-resume
              tex
              gen-resume
            ];
          };

        };
      }
    );
}
