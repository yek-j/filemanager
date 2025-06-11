package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// filemanager-config.json
	configPath := "../filemanager-config.json" 

	config, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// 설정 검증
	if config.SourcePath != "/path/to/root" {
		t.Errorf("Expected SourcePath '/path/to/root', got '%s'", config.SourcePath)
	}

	if config.TargetDepth != 3 {
        t.Errorf("Expected TargetDepth 3, got %d", config.TargetDepth)
    }
}
