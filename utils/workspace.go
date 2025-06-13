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

	// 플랫폼별 복사 진행
	platform := getPlatform()

	var err error

	switch platform {
	case "linux", "darwin":
		err = copyWithCommand(cfg.SourcePath, cfg.WorkPath)
	case "windows":
		err = copyWithGo(cfg.SourcePath, cfg.WorkPath)
	default:
		err = copyWithGo(cfg.SourcePath, cfg.WorkPath)
	}

	if err != nil {
		return fmt.Errorf("copy failed: %v", err)
	}

	if !VerifyWorkspace(cfg.SourcePath, cfg.WorkPath) {
		return fmt.Errorf("copy vertification failed")
	}

	return nil
}

// copyWithCommand: 리눅스 명령어를 이용한 copy
func copyWithCommand(sourcePath, workPath string) error {
	cmd := exec.Command("cp", "-r", sourcePath, workPath)

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
func VerifyWorkspace(sourcePath, workPath string) bool {
	// workPath에 폴더가 존재하는지 확인
	if _, err := os.Stat(workPath); os.IsNotExist(err) {
		return false
	}

	// sourcePath의 폴더/파일 개수 세기
	oriFiles, oriDirs := 0, 0
	filepath.WalkDir(sourcePath, func(path string, d fs.DirEntry, err error) error {
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
	filepath.WalkDir(workPath, func(path string, d fs.DirEntry, err error) error {
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
