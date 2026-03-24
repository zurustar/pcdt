// Package model contains the data models for the countdown timer.
package model

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// AppConfig はアプリケーション設定を保持する
type AppConfig struct {
	AlwaysOnTop  bool    `json:"alwaysOnTop"`  // 常に最前面表示
	WindowWidth  float32 `json:"windowWidth"`  // ウィンドウ幅
	WindowHeight float32 `json:"windowHeight"` // ウィンドウ高さ
	LastMinutes  int     `json:"lastMinutes"`  // 最後に入力した分
	LastSeconds  int     `json:"lastSeconds"`  // 最後に入力した秒
}

// DefaultConfig はデフォルト設定を返す
func DefaultConfig() *AppConfig {
	return &AppConfig{
		AlwaysOnTop:  false,
		WindowWidth:  800,
		WindowHeight: 500,
		LastMinutes:  5,
		LastSeconds:  0,
	}
}

// GetConfigDir はOS別の設定ディレクトリパスを返す
func GetConfigDir() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		// macOS: ~/Library/Application Support/CountdownTimer/
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, "Library", "Application Support", "CountdownTimer"), nil
	case "windows":
		// Windows: %APPDATA%\CountdownTimer\
		appData := os.Getenv("APPDATA")
		if appData == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		return filepath.Join(appData, "CountdownTimer"), nil
	default:
		// Linux/その他: ~/.config/CountdownTimer/
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, ".config", "CountdownTimer"), nil
	}
}

// GetConfigFilePath は設定ファイルのフルパスを返す
func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

// ToJSON は設定をJSON形式にシリアライズする
func (c *AppConfig) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "    ")
}

// FromJSON はJSON形式から設定をデシリアライズする
func FromJSON(data []byte) (*AppConfig, error) {
	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Save は設定をファイルに保存する
func (c *AppConfig) Save() error {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	// 設定ディレクトリを作成
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// JSONにシリアライズ
	data, err := c.ToJSON()
	if err != nil {
		return err
	}

	// ファイルに書き込み
	return os.WriteFile(configPath, data, 0644)
}

// Load は設定ファイルから設定を読み込む
// ファイルが存在しない場合やエラーの場合はデフォルト設定を返す
func Load() (*AppConfig, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return DefaultConfig(), err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return DefaultConfig(), err
	}

	config, err := FromJSON(data)
	if err != nil {
		return DefaultConfig(), err
	}

	return config, nil
}
