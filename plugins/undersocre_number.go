package plugins

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/yek-j/filemanager/config"
)

type UnderscoreNumber struct{}

type FileInfo struct {
	FileName string
	Number   int
	FullPath string
}

func (u *UnderscoreNumber) Process(cfg *config.Config) error {
	// 작업할 폴더들 찾기
	// cfg.WorkPath + cfg.TargetFolders + cfg.TargetDepth 조합
	// 원하는 위치에서 파일 수집
	for _, targetFolder := range cfg.TargetFolders {
		// workPath/targetFolder
		basePath := filepath.Join(cfg.WorkPath, targetFolder)

		// target_depth로 들어가기
		err := processTargetDepth(basePath, cfg.TargetDepth)
		if err != nil {
			return err
		}
	}

	return nil
}

func processTargetDepth(basePath string, depth int) error {
	currentDirs := []string{basePath}

	// 'depth - 1' 반복으로 최종 폴더 찾기
	for i := 1; i < depth; i++ {
		nextDirs := []string{}

		for _, dir := range currentDirs {
			// dir의 하위 폴더들 읽기
			entries, err := os.ReadDir(dir)

			if err != nil {
				continue // 읽을 수 없다면 스킵
			}

			for _, entry := range entries {
				if entry.IsDir() {
					nextDirs = append(nextDirs, filepath.Join(dir, entry.Name()))
				}
			}
		}

		currentDirs = nextDirs
	}

	// 최종 폴더들에서 파일을 처리
	for _, finalDir := range currentDirs {
		if err := processFilesInDirectory(finalDir); err != nil {
			return err
		}
	}

	return nil
}

func processFilesInDirectory(finalDir string) error {
	// 폴더 안의 파일들만 읽기(하위폴더 제외)
	entires, err := os.ReadDir(finalDir)

	if err != nil {
		return err
	}

	//prefix_숫자.확장자 패턴인 파일만 읽기
	// prefix별로 그룹핑
	groups := make(map[string][]FileInfo)

	for _, entry := range entires {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		prefix, number, ext, valid := parseFileName(filename)
		if !valid {
			continue
		}

		// prefix.pdf 형식으로 키 생성
		key := prefix + ext
		groups[key] = append(groups[key], FileInfo{
			FileName: filename,
			Number:   number,
			FullPath: filepath.Join(finalDir, filename),
		})
	}

	// 각 그룹에서 최대 숫자 파일만 남기고 삭제
	// 남은 파일을 prefix_1.확장자로 변경
	for _, files := range groups {
		// 최대 숫자 찾기
		maxFile := files[0]
		for _, file := range files {
			if file.Number > maxFile.Number {
				maxFile = file
			}
		}

		for _, file := range files {
			// 나머지 파일 삭제
			if file.FullPath != maxFile.FullPath {
				os.Remove(file.FullPath)
			} else {
				// 파일 이름 변경
				prefix, _, ext, _ := parseFileName(file.FileName)
				newName := prefix + "_1" + ext
				newPath := filepath.Join(finalDir, newName)

				os.Rename(file.FullPath, newPath)
			}
		}

	}

	return nil
}

func parseFileName(filename string) (prefix string, number int, ext string, valid bool) {
	// 확장자 분리
	ext = filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// _ 찾기
	lastUnderscoreIndex := strings.LastIndex(nameWithoutExt, "_")
	if lastUnderscoreIndex == -1 {
		return "", 0, "", false
	}

	// prefix와 숫자부분 분리
	prefix = nameWithoutExt[:lastUnderscoreIndex]
	numberStr := nameWithoutExt[lastUnderscoreIndex+1:]

	// 숫자 변환
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return "", 0, "", false // 숫자가 아닌 경우
	}

	return prefix, number, ext, true
}

func (u *UnderscoreNumber) GetName() string {
	return "UNDERSCORE_NUMBER"
}

func (u *UnderscoreNumber) GetDescription() string {
	return "prefix_1.txt 형식의 파일을 폴더별로 찾아서 같은 prefix별로 숫자가 가장 큰 수를 제외하고 삭제한다. 남은 파일은 prefix_1로 일괄로 변경한다."
}
