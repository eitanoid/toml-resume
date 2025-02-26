docker run \
    -v $(pwd):/data:z \
    -v /usr/share/fonts:/usr/share/fonts/sysfonts:ro,z \
    -v $HOME/.local/share/fonts:/usr/share/fonts/userfonts:ro,z \
    toml-resume "$@"
