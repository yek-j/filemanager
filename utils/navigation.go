package utils

import (
	"github.com/yek-j/filemanager/config"
)

// FileInfo: 파일 정보 구조체
type FileInfo struct {
	FilePath  string // 파일 경로
	Name      string // 파일명
	Extension string // 확장자
}

// ScanFiles는 Config에서 가져온 폴더의 유효성을 검증하고 총 파일을 확인한다.
func ScanFiles(cfg *config.Config) ([]FileInfo, error) {
	// config의 sourcePath가 존재하는 확인

	return nil, nil
}
