package manifests

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ClearIgnoredFiles(manifest_name string) error {
	ignore, err := ReadIgnoreSection(manifest_name)
	if err != nil {
		return err
	}
	for _, file := range ignore {
		err = DeleteFile(file)
		if err != nil {
			return err
		}
	}
	fmt.Println("Ignored files cleared")
	return nil
}

func ReadIgnoreSection(manifest_name string) ([]string, error) {
	type Ignore struct {
		Ignore []string `json:"ignore"`
	}

	config := Ignore{}
	optiqueConfig, err := os.ReadFile(manifest_name)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(optiqueConfig, &config)
	if err != nil {
		return nil, err
	}
	return config.Ignore, nil
}

func DeleteFile(glob string) error {
	matches, err := filepath.Glob(glob)
	if err != nil {
		return err
	}
	for _, match := range matches {
		err = os.Remove(match)
		if err != nil {
			return err
		}
	}
	return nil
}
