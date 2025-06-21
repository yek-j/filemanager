package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/yek-j/filemanager/config"
)

// ScanReport: 설정을 검증하여 작업이 가능한 상태임을 확인할 수 있는 구조체
type ScanReport struct {
	RootExists     bool
	TargetFolders  map[string]bool  // 폴더명: 존재여부
	FoldersByDepth map[int][]string // 깊이별 폴더 목록
	FilesByExt     map[string]int   // 확장자별 개수
	TotalFiles     int
	ReadyToProcess bool
}

// ScanFiles는 Config에서 가져온 폴더의 유효성을 검증하고 총 작업 파일 수 확인
// 작업이 가능한지 확인한다.
func ScanFiles(cfg *config.Config) (*ScanReport, error) {
	scanReport := &ScanReport{
		TargetFolders:  make(map[string]bool),
		FoldersByDepth: make(map[int][]string),
		FilesByExt:     make(map[string]int),
	}

	// config의 sourcePath가 존재하는 확인
	sourceInfo, err := os.Stat(cfg.SourcePath)

	if err != nil {
		scanReport.RootExists = false
		scanReport.ReadyToProcess = false
		return scanReport, fmt.Errorf("source path not found: %v", err)
	}

	if !sourceInfo.IsDir() {
		scanReport.RootExists = false
		scanReport.ReadyToProcess = false
		return scanReport, fmt.Errorf("source path is not a directory")
	}

	scanReport.RootExists = true // root 파일 존재 확인

	// TargetFolders 존재 여부 확인
	for _, targetFolder := range cfg.TargetFolders {
		// 전체 경로 생성
		targetPath := filepath.Join(cfg.SourcePath, targetFolder)

		// 존재 확인
		if _, err := os.Stat(targetPath); err == nil {
			scanReport.TargetFolders[targetFolder] = true
		} else {
			scanReport.TargetFolders[targetFolder] = false
		}
	}

	// FoldersByDepth 깊이별 폴더 목록 확인
	// FilesByExt 최종 TargetDepth에서 확장자별 파일 수 확인
	// TotalFiles 총 파일 수 확인
	err = filepath.WalkDir(cfg.SourcePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		var depth int
		relativePath := strings.TrimPrefix(path, cfg.SourcePath)
		relativePath = strings.TrimPrefix(relativePath, string(os.PathSeparator)) // '/' 제거

		if relativePath == "" {
			depth = 0 // root 자체
		} else {
			depth = strings.Count(relativePath, string(os.PathSeparator)) + 1
		}

		if d.IsDir() {
			// 깊이 제한 체크 + root 제외
			if depth > 0 && depth <= cfg.TargetDepth {
				scanReport.FoldersByDepth[depth] = append(scanReport.FoldersByDepth[depth], path)
			}
		} else {
			ext := filepath.Ext(d.Name())
			if ext != "" {
				scanReport.FilesByExt[ext]++
				scanReport.TotalFiles++
			}
		}
		return nil
	})

	if err != nil {
		return scanReport, fmt.Errorf("failed to scan directory structure: %v", err)
	}

	// ReadyToProcess 위에 ROOT 폴더, Target 폴더 모두 존재한다면 작업 준비 완료
	existingCount := 0
	for _, exists := range scanReport.TargetFolders {
		if exists {
			existingCount++
		}
	}

	if scanReport.RootExists && existingCount == len(cfg.TargetFolders) {
		scanReport.ReadyToProcess = true
	} else {
		scanReport.ReadyToProcess = false
	}

	return scanReport, nil
}

// IsWorkPathEmpty: 폴더가 없으면 true, 폴더가 있지만 비어있으면 true,
// 폴더에 파일/폴더 있으면 false 반환
func IsWorkPathEmpty(workPath string) bool {
	workDir, err := os.ReadDir(workPath)

	if os.IsNotExist(err) {
		return true // 폴더 없음
	}

	if err != nil {
		return false // error
	}

	return len(workDir) == 0 // 빈 폴더면 true
}
