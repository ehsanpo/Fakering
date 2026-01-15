package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// App represents the main application service
type App struct{}

// AppConfig represents the application configuration
type AppConfig struct {
	Theme          string                 `json:"theme"`
	Language       string                 `json:"language"`
	WindowWidth    int                    `json:"windowWidth"`
	WindowHeight   int                    `json:"windowHeight"`
	CustomSettings map[string]interface{} `json:"customSettings"`
}


// GetConfigPath returns the path to the config file
func (a *App) GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".fakeRing")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig loads the application configuration
func (a *App) LoadConfig() (*AppConfig, error) {
	configPath, err := a.GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Return default config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return a.GetDefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// SaveConfig saves the application configuration
func (a *App) SaveConfig(config *AppConfig) error {
	configPath, err := a.GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetDefaultConfig returns the default configuration
func (a *App) GetDefaultConfig() *AppConfig {
	return &AppConfig{
		Theme:        "light",
		Language:     "en",
		WindowWidth:  1024,
		WindowHeight: 768,
		CustomSettings: make(map[string]interface{}),
	}
}

// GetSetting gets a specific setting value
func (a *App) GetSetting(key string) (interface{}, error) {
	config, err := a.LoadConfig()
	if err != nil {
		return nil, err
	}

	if value, exists := config.CustomSettings[key]; exists {
		return value, nil
	}

	return nil, nil
}

// SetSetting sets a specific setting value
func (a *App) SetSetting(key string, value interface{}) error {
	config, err := a.LoadConfig()
	if err != nil {
		return err
	}

	if config.CustomSettings == nil {
		config.CustomSettings = make(map[string]interface{})
	}

	config.CustomSettings[key] = value
	return a.SaveConfig(config)
}
