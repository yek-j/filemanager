package config

import (
	"encoding/json"
	"os"
)

// Config 구조체 - JSON 설정과 파일 매핑
type Config struct {
	SourcePath    string   `json:"source_path"`
	BackupPath    string   `json:"backup_path"`
	TargetFolders []string `json:"target_folders"`
	TargetDepth   int      `json:"file_depth"`
	DryRun        bool     `json:"dry_run"`
}

// LoadConfig JSON 파일에서 설정을 읽어 온다.
// configPath: 설정 파일 경로
// return: Config 구조체 포인터와 에러
func LoadConfig(configPath string) (*Config, error) {
	// 파일 읽기
	configData, err := os.ReadFile(configPath)

	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(configData, &config) // json->Config

	if err != nil {
		return nil, err
	}

	return &config, nil
}
