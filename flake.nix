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
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };
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
          set -e
          if [[ -z "$1" ]] || [[ ! -f "$1" ]]; then
            echo "Usage: gen-resume <input_file.toml>"
            exit 1
          fi
          INPUT_PATH=$(realpath "$1")
          BASE_NAME=$(basename ''${1%.*})
          CWD=$(pwd)

          TEMP_DIR=$(mktemp -d)

          trap 'rm -rf "$TEMP_DIR"' EXIT

          cd ''${TEMP_DIR}

          toml-resume "''${INPUT_PATH}" -o out_tmp.tex -f

          latexmk -silent -synctex=0 -interaction=nonstopmode -pdf -lualatex out_tmp.tex --jobname="''$BASE_NAME"

          if [[ -f "$BASE_NAME.pdf" ]]; then
            mv "$BASE_NAME.pdf" "$CWD/"
            echo "Created $BASE_NAME.pdf in $CWD"
          else
            echo "Error: failed to generate pdf."
            exit 1
          fi
        '';
        mkTexShell =
          {
            extraFonts ? [ ],
          }:
          let
            fontsConf = pkgs.makeFontsConf {
              fontDirectories = extraFonts;
            };
          in
          pkgs.mkShellNoCC {
            packages = [
              tex
              pkgs.fontconfig
              toml-resume
              gen-resume
            ];
            shellHook = ''
              export FONTCONFIG_FILE="${fontsConf}"
              export OSFONTDIR="${pkgs.lib.makeSearchPath "share/fonts" extraFonts}"
              # This is needed when running in pure mode sometimes
              export TEXMFVAR=$(mktemp -d)
              trap 'rm -rf "$TEXMFVAR"' EXIT
            '';
          };
      in
      {
        # expose builder function in latex shell to --expr
        lib.mkTexShell = mkTexShell;

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
