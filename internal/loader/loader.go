package loader

import (
	"errors"
	"fmt"
	"io"

	"github.com/pelletier/go-toml/v2"
	"sigs.k8s.io/yaml"
)

func LoadFromReader(r io.Reader, target any) error {

	fi, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read content: %w", err)
	}

	// detect filetype by trial and error
	if err := toml.Unmarshal(fi, target); err == nil {
		return nil
	}

	if err := yaml.Unmarshal(fi, target); err == nil {
		return nil
	}

	return errors.New("input is not valid TOML or YAML")
}
