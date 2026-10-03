package profile

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/manifoldco/promptui"
)

func setupProfile(ctx context.Context, profile string) error {
	if ctx.Err() != nil {
		return fmt.Errorf("operation cancelled: %w", ctx.Err())
	}

	configMap, err := openConfigMap(ctx)
	if err != nil {
		slog.Error("unable to open configuration map", slog.String("error", err.Error()))
		return err
	}

	config, ok := configMap[profile]

	if !ok {
		slog.Info("profile does not exist. Creating new profile.", slog.String("name", profile))
		config = Profile{}
	} else {
		slog.Info("profile found. Updating existing profile.", slog.String("name", profile))
	}

	urlPrompt := promptui.Prompt{
		Label: "Server URL",
		Validate: func(s string) error {
			if s == "" {
				return fmt.Errorf("server URL cannot be empty")
			}
			if u, err := url.Parse(s); err == nil {
				if u.Scheme != "http" && u.Scheme != "https" {
					return fmt.Errorf("invalid URL scheme: %s", u.Scheme)
				}
			} else {
				return fmt.Errorf("invalid server URL: %w", err)
			}
			return nil
		},
		Default:   config.URL,
		AllowEdit: true,
	}

	result, err := urlPrompt.Run()
	if err != nil {
		err = fmt.Errorf("prompt failed: %w", err)
		return err
	}
	config.URL = result

	apiPrompt := promptui.Prompt{
		Label:     "API key",
		Default:   config.APIKey,
		AllowEdit: true,
	}

	result, err = apiPrompt.Run()
	if err != nil {
		err = fmt.Errorf("prompt failed: %w", err)
		return err
	}
	config.APIKey = result

	configMap[profile] = config

	err = saveConfigMap(ctx, configMap)
	if err != nil {
		return err
	}

	return nil
}
