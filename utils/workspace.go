package utils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/yek-j/filemanager/config"
)

func getPlatform() string {
	return runtime.GOOS // "linux", "windows"..
}

// CopyRootDir: source_path의 폴더를 복사하여 작업 공간을 만든다.
func CopyRootDir(cfg *config.Config) error {
	// work_path 검증
	if !IsWorkPathEmpty(cfg.WorkPath) {
		return fmt.Errorf("work path not empty")
	}

	// 폴더가 없으면 생성
	if err := os.MkdirAll(cfg.WorkPath, 0755); err != nil {
		return err
	}

	if cfg.SelectiveCopy {
		for _, targetFolder := range cfg.TargetFolders {
			sourcePath := filepath.Join(cfg.SourcePath, targetFolder)
			targetPath := filepath.Join(cfg.WorkPath, targetFolder)

			if err := copyByPlatform(sourcePath, targetPath); err != nil {
				return fmt.Errorf("copy failed for %s: %v", targetFolder, err)
			}
		}
	} else {
		if err := copyByPlatform(cfg.SourcePath, cfg.WorkPath); err != nil {
			return fmt.Errorf("copy failed: %v", err)
		}
	}

	if !VerifyWorkspace(cfg) {
		return fmt.Errorf("copy vertification failed")
	}

	return nil
}

// copyWithCommand: 리눅스 명령어를 이용한 copy
func copyWithCommand(sourcePath, workPath string) error {
	cmd := exec.Command("cp", "-r", sourcePath+"/.", workPath)

	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("cp command failed: %v", err)
	}

	return nil
}

// copyWithGo: Go 프로그래밍으로 copy
func copyWithGo(sourcePath, workPath string) error {
	return filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 경로
		rootPath, _ := filepath.Rel(sourcePath, path)
		targetPath := filepath.Join(workPath, rootPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755) // 폴더 생성
		} else {
			return copyFile(path, targetPath)
		}
	})
}

// copyFile: 파일 복사
func copyFile(path, targetPath string) error {
	// 원본 파일 열기
	oriFile, err := os.Open(path)
	if err != nil {
		return err
	}

	defer oriFile.Close()

	// 파일 생성
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}

	defer targetFile.Close()

	// 복사
	_, err = io.Copy(targetFile, oriFile)
	return err
}

// VerifyWorkspace: 복사된 폴더에서 폴더의 존재와 파일 개수를 검증한다.
func VerifyWorkspace(cfg *config.Config) bool {
	// workPath에 폴더가 존재하는지 확인
	if _, err := os.Stat(cfg.WorkPath); os.IsNotExist(err) {
		return false
	}

	// sourcePath의 폴더/파일 개수 세기
	if cfg.SelectiveCopy {
		// selective_copy 모드: target_folders별로 검증
		for _, targetFolder := range cfg.TargetFolders {
			sourceFolderPath := filepath.Join(cfg.SourcePath, targetFolder)
			workFolderPath := filepath.Join(cfg.WorkPath, targetFolder)

			// sourceFolderPath 파일/폴더 개수 세기
			oriFiles, oriDirs := 0, 0
			filepath.WalkDir(sourceFolderPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					oriDirs++
				} else {
					oriFiles++
				}
				return nil
			})

			// workFolderPath 파일/폴더 개수 세기
			workFiles, workDirs := 0, 0
			filepath.WalkDir(workFolderPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					workDirs++
				} else {
					workFiles++
				}
				return nil
			})

			// 개수 비교
			if oriFiles != workFiles || oriDirs != workDirs {
				return false
			}
		}
		return true
	} else {
		oriFiles, oriDirs := 0, 0
		filepath.WalkDir(cfg.SourcePath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				oriDirs++
			} else {
				oriFiles++
			}
			return nil
		})

		workFiles, workDirs := 0, 0
		filepath.WalkDir(cfg.WorkPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				workDirs++
			} else {
				workFiles++
			}
			return nil
		})

		return oriFiles == workFiles && oriDirs == workDirs
	}

}

// copyByPlatform: 플랫폼별 코드 복사 분기
func copyByPlatform(sourcePath, destPath string) error {
	platform := getPlatform()

	switch platform {
	case "linux", "darwin":
		return copyWithCommand(sourcePath, destPath)
	case "windows":
		return copyWithGo(sourcePath, destPath)
	default:
		return copyWithGo(sourcePath, destPath)
	}
}
