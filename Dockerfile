FROM golang:1.23 AS build

# set working destination
WORKDIR /app

# copy all go files
COPY go.mod go.sum *.go ./

RUN go mod tidy && go mod download

# build the go project
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/toml-resume

# The latex build
FROM debian:stable-slim

# install latex and fontconfig
RUN apt-get update && apt-get install -y fontconfig latexmk texlive-latex-recommended texlive-latex-extra texlive-xetex

WORKDIR /data

COPY --from=build /app/toml-resume /app/toml-resume
COPY preamble.tex ./

VOLUME ["/data"]

ENTRYPOINT ["/bin/sh", "-c", "/app/toml-resume -input=\"$1\" -output=\"$2\" && latexmk -pdf -xelatex \"$2\" && latexmk -c", "--" ]

