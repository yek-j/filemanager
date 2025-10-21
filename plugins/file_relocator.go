package plugins

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/yek-j/filemanager/config"
	"github.com/yek-j/filemanager/utils"
)

type FileRelocator struct {
	pluginCfg *config.PluginConfig
}

// FileRelocator 플러그인의 설정값 구조체
type FileRelocatorConfig struct {
	// 파일 선택
	FileExtensions []string `json:"file_extensions,omitempty"` // 이동할 파일 확장자
	FilePattern    string   `json:"file_pattern,omitempty"`    // 파일명 패턴

	// 경로
	SourceLocation string `json:"source_location"` // 파일 경로
	TargetLocation string `json:"target_location"` // 이동할 경로

	// 동작 옵션
	CreateFolder   bool     `json:"create_folder"`   // 이동할 폴더가 없을 때 자동 생성 여부
	SearchSubdirs  bool     `json:"search_subdirs"`  // 하위 폴더까지 검색 여부
	OverwriteFiles bool     `json:"overwrite_files"` // 이동할 위치에 이미 파일이 있다면 덮어쓰기 여부
	TargetFolders  []string `json:"target_folders"`  // 이동할 타켓 폴더
	UsePattern     bool     `json:"use_pattern"`     // depth 사용 시 false, pattern 사용 시 true
}

type FileRelocatorLog struct {
	MovedFiles  map[string]string // 원본경로 -> 대상경로
	FailedMoves []string          // 실패한 파일 (전체 경로)
	TotalFiles  int
}

func (m *FileRelocator) Process(cfg *config.Config) error {
	totalProcessed := 0
	log := &FileRelocatorLog{
		MovedFiles: make(map[string]string),
	}

	// 설정 구조체
	var pluginConfig FileRelocatorConfig

	// Config 파싱
	if m.pluginCfg != nil && len(m.pluginCfg.Config) > 0 {
		err := json.Unmarshal(m.pluginCfg.Config, &pluginConfig)
		if err != nil {
			return fmt.Errorf("failed to parse plugin config: %v", err)
		}
	}

	// UsePatter에 따라 작업 방식 분기
	if pluginConfig.UsePattern {
		// TODO: usePattern에 따라 작업할 폴더 찾기
	} else {
		//
		for _, targetDir := range pluginConfig.TargetFolders {

			basePath := filepath.Join(cfg.WorkPath, targetDir)

			// 작업할 경로
			workDirs := utils.GetTargetDirs(basePath, cfg.TargetDepth)
			for _, dir := range workDirs {
				count, err := processMoveFiles(dir, pluginConfig, log)
				totalProcessed += count
				if err != nil {
					return err
				}
			}
		}
	}

	log.TotalFiles = totalProcessed

	logFileName := fmt.Sprintf("file_relocator_log_%s.txt",
		time.Now().Format("20060102_150405"))
	logPath := filepath.Join(cfg.WorkPath, logFileName)

	if err := writeFileRelocatorLogFile(log, logPath); err != nil {
		fmt.Printf("Warning: Failed to write log file: %v\n", err)
	} else {
		fmt.Printf("📝 Log file created: %s\n", logPath)
	}

	return nil
}

func processMoveFiles(workDir string, pluginConfig FileRelocatorConfig, log *FileRelocatorLog) (int, error) {
	finalDir := filepath.Join(workDir, pluginConfig.SourceLocation)
	processFileCount := 0 // 총 파일 개수

	// SearchSubdirs(하위 폴더까지 검색하는 옵션)에 따라 분기
	if pluginConfig.SearchSubdirs {
		err := filepath.WalkDir(finalDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}

			if err := processFile(path, d.Name(), workDir, pluginConfig, log); err != nil {
				return err
			}
			processFileCount++
			return nil
		})

		if err != nil {
			return processFileCount, err
		}
	} else {
		entires, err := os.ReadDir(finalDir)

		if err != nil {
			return processFileCount, err
		}

		for _, entry := range entires {
			if entry.IsDir() {
				continue
			}

			sourcePath := filepath.Join(finalDir, entry.Name())
			if err := processFile(sourcePath, entry.Name(), workDir, pluginConfig, log); err != nil {
				return processFileCount, err
			}
			processFileCount++

		}
	}

	return 0, nil
}

// processFile: 단일 파일을 target_location으로 이동
func processFile(sourcePath string, fileName string, baseDir string,
	pluginConfig FileRelocatorConfig, log *FileRelocatorLog) error {
	// 확장자 체크
	ext := strings.TrimPrefix(filepath.Ext(fileName), ".")

	if len(pluginConfig.FileExtensions) > 0 &&
		!slices.Contains(pluginConfig.FileExtensions, ext) {
		return nil // 건너뛰기
	}

	// target_location 경로
	targetDirPath := filepath.Join(baseDir, pluginConfig.TargetLocation)

	// target_location 디렉터리 확인
	err := ensureDir(targetDirPath, pluginConfig.CreateFolder)
	if err != nil {
		return err
	}

	// 최종 경로
	targetPath := filepath.Join(targetDirPath, fileName)

	// 덮어쓰기 체크
	if _, err := os.Stat(targetPath); err == nil {
		if !pluginConfig.OverwriteFiles {
			return nil // 건너뛰기
		}
	}

	os.Rename(sourcePath, targetPath)
	log.MovedFiles[sourcePath] = targetPath
	return nil
}

// 디렉터리 존재 확인 + 필요시 생성
func ensureDir(path string, create bool) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if create {
			return os.MkdirAll(path, 0755)
		}
		return fmt.Errorf("directory does not exist: %s", path)
	}
	return nil
}

func writeFileRelocatorLogFile(log *FileRelocatorLog, logPath string) error {
	file, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "FileManager FileRelocator Processing Log\n")
	fmt.Fprintf(file, "Total files processed: %d\n\n", log.TotalFiles)

	fmt.Fprintf(file, "=== MOVED FILES ===\n")
	for original, moved := range log.MovedFiles {
		fmt.Fprintf(file, "MOVED: %s -> %s\n", original, moved)
	}

	return nil
}

func (m *FileRelocator) GetName() string {
	return "FILE_RELOCATOR"
}

func (m *FileRelocator) GetDescription() string {
	return "지정된 파일들을 일괄 이동합니다. " +
		"단순 구조는 file_depth 기반, 복잡한 구조는 정규식 패턴을 사용합니다."
}
