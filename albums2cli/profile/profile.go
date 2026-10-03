package profile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v4"
)

type Network string

type Profile struct {
	URL    string `yaml:"url"`
	APIKey string `yaml:"api_key"`
}

func createConfigPath() (path string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error(
			"failed to get user home directory",
			slog.String("error", err.Error()),
		)
		return
	}
	path = filepath.Join(homeDir, ".albums2cli", "config.yaml")

	return
}

func openConfigMap(ctx context.Context) (profileMap map[string]Profile, err error) {
	if ctx.Err() != nil {
		err = fmt.Errorf("operation cancelled: %w", ctx.Err())
		return
	}
	path, err := createConfigPath()
	if err != nil {
		return
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		slog.Warn(
			"Configuration file does not exist. Creating empty configuration map.",
			slog.String("path", path),
		)
		profileMap = make(map[string]Profile)
		err = nil
		return
	} else if err != nil {
		err = fmt.Errorf("failed to open configuration file: %w", err)
		return
	}
	defer file.Close()

	profileMap = make(map[string]Profile)
	decoder := yaml.NewDecoder(file)

	err = decoder.Decode(&profileMap)
	if errors.Is(err, io.EOF) {
		slog.Warn(
			"Configuration file is empty. Creating empty configuration map.",
			slog.String("path", path),
		)
		profileMap = make(map[string]Profile)
		err = nil
		return
	} else if err != nil {
		err = fmt.Errorf("failed to parse existing configuration file: %w", err)
		return
	}
	return
}

func saveConfigMap(ctx context.Context, configMap map[string]Profile) error {
	if ctx.Err() != nil {
		return fmt.Errorf("operation cancelled: %w", ctx.Err())
	}

	path, err := createConfigPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0777)
	if err != nil {
		return fmt.Errorf("failed to create configuration directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to save configuration file: %w", err)
	}
	defer file.Close()

	yamlEncoder := yaml.NewEncoder(file)
	defer yamlEncoder.Close()

	err = yamlEncoder.Encode(configMap)
	if err != nil {
		return fmt.Errorf("failed to save configuration file: %w", err)
	}

	return nil
}

func LoadProfile(ctx context.Context, name string) (profile Profile, err error) {
	configMap, err := openConfigMap(ctx)
	if err != nil {
		err = fmt.Errorf("unable to open configuration file: %w", err)
		return
	}

	profile, ok := configMap[name]

	if !ok {
		err = fmt.Errorf("profile not found: %s", name)
		return
	}

	return
}
