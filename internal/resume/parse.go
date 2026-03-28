package resume

import (
	"errors"
	"fmt"
	"os"

	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"sigs.k8s.io/yaml"
)

var (
	ErrNoInput = errors.New("missing path to input file")
)

func (r *Resume) LoadDataFromFile(path string) error {
	ext := filepath.Ext(path)
	fi, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	switch ext {

	case ".toml":
		if err := toml.Unmarshal(fi, r.Data); err != nil {
			return fmt.Errorf("failed to unmarshal toml file: %w", err)
		}

	case ".yaml", ".yml":
		if err := yaml.Unmarshal(fi, r.Data); err != nil {
			return fmt.Errorf("failed to unmarshal yaml file: %w", err)
		}

	default:
		return errors.New("failed to determine file type")
	}
	return nil
}
