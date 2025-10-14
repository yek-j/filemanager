package plugins

import "github.com/yek-j/filemanager/config"

type FileRelocator struct {
	pluginCfg *config.PluginConfig
}

// FileRelocator 플러그인의 설정값 구조체
type FileRelocatorConfig struct {
	// 파일 선택
	FileExtensions []string `json:"file_extensions"`        // 이동할 파일 확장자
	FilePattern    string   `json:"file_pattern,omitempty"` // 파일명 패턴

	// 경로
	SourceLocation  string `json:"source_location"` // 파일 경로
	TargetLocation  string `json:"target_location"` // 이동할 경로

	// 동작 옵션
	CreateFolder   bool     `json:"create_folder"`   // 이동할 폴더가 없을 때 자동 생성 여부
	SearchSubdirs  bool     `json:"search_subdirs"`  // 하위 폴더까지 검색 여부
	OverwriteFiles bool     `json:"overwrite_files"` // 이동할 위치에 이미 파일이 있다면 덮어쓰기 여부
	TargetFolders  []string `json:"target_folders"`  // 이동할 타켓 폴더
	UsePattern     bool     `json:"use_pattern"`	 // depth 사용 시 false, pattern 사용 시 true
}

func (m *FileRelocator) Process(cfg *config.Config) error {
	// TODO
	return nil
}

func (m *FileRelocator) GetName() string {
	return "FILE_RELOCATOR"
}

func (m *FileRelocator) GetDescription() string {
	return "지정된 파일들을 일괄 이동합니다. " +
		"단순 구조는 file_depth 기반, 복잡한 구조는 정규식 패턴을 사용합니다."
}
