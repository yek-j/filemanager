# FileManager

리눅스/윈도우 시스템을 위한 파일 정리 자동화 도구

## 📋 개요

FileManager는 패턴 기반으로 파일을 정리하는 도구입니다. 안전한 백업과 함께 중복 파일 삭제, 파일 이름 변경을 자동화합니다.

### 주요 기능
- 🔄 **안전한 백업** - 원본 파일을 건드리기 전에 자동 백업
- 📁 **패턴 기반 정리** - 접두사별로 최신 파일만 남기고 나머지 삭제
- 🔧 **유연한 설정** - JSON 설정으로 다양한 정리 규칙 지원
- ⚡ **선택적 복사** - 전체가 아닌 필요한 폴더만 복사하여 효율성 향상
- 🎯 **플러그인별 세부 설정** - 각 플러그인의 고유한 설정 지원

## ⚙️ 설정 파일 구조

### 기본 설정 (`config.json`)

```json
{
  "source_path": "/path/to/root",
  "work_path": "/path/to/work",
  "target_folders": ["paper", "homework"],
  "file_depth": 3,
  "plugin": "underscore_number",
  "selective_copy": true,
  "plugin_config": {
    "allowed_extensions": ["pdf", "txt"]
  }
}
```

### 설정 항목 설명

| 항목 | 설명 | 예시 | 기본값 |
|------|------|------|--------|
| `source_path` | 정리할(복사할) 루트 폴더 경로 | `"/home/user/documents"` | 필수 |
| `work_path` | 작업 경로(복사할 위치) | `"/home/user/work"` | 필수 |
| `target_folders` | 처리할 대상 폴더들(root 하위 폴더 기준) | `["paper", "homework", "assignments"]` | 필수 |
| `file_depth` | target_folder 기준 탐색 깊이 | `3` (paper → 10 → 1001 → 파일) | 필수 |
| `plugin` | 파일 정리 방식을 결정하는 플러그인 | `"underscore_number"` | 필수 |
| `selective_copy` | 선택적 복사 모드 사용 여부 | `true` (target_folders만 복사) | `false` |
| `plugin_config` | 플러그인별 세부 설정 | 플러그인에 따라 다름 | 선택사항 |

### target_depth 설명

`target_depth`는 target_folder를 기준으로 한 상대적 깊이입니다:

- `target_depth: 1` → `paper/` (target_folder 바로 아래)
- `target_depth: 2` → `paper/f1/` (paper 하위 폴더까지)
- `target_depth: 3` → `paper/f1/1111/` (f1 하위  폴더까지, **파일이 있는 위치**)

### 복사 모드 비교

| 모드 | `selective_copy` | 복사 대상 | 장점 | 단점 |
|------|------------------|-----------|------|------|
| **전체 복사** | `false` (기본값) | source_path 전체 | 완전한 백업 | 느림, 공간 많이 사용 |
| **선택적 복사** | `true` | target_folders만 | 빠름, 공간 절약 | 부분 백업만 |


## 🔌 사용 가능한 플러그인

### underscore_number
- **설명**: `문자열_숫자.확장자` 패턴의 파일을 정리
- **동작**: 각 폴더 내에서 같은 접두사 중 가장 큰 숫자 파일만 남기고 `prefix_원하는숫자` 형태로 일괄 변경
- **예시**: `file_123.pdf`, `file_456.pdf` → `file_456.pdf` 외에 같은 접두사 파일은 삭제하고, `file_1.pdf`로 변경
- **로그파일**: 파일 처리 완료 후 WORK_PATH에 삭제된 파일, 이름 변경된 파일, 총 처리된 파일에 대한 로그 파일 생성


#### 플러그인 설정 (`plugin_config`)

```json
{
  "plugin_config": {
    "allowed_extensions": ["pdf", "txt", "docx"]
  }
}
```

| 설정 항목 | 설명 | 예시 | 기본값 |
|-----------|------|------|--------|
| `allowed_extensions` | 처리할 파일 확장자 목록 (점 제외) | `["pdf", "txt"]` | 모든 확장자 |

**확장자 필터링 예시:**
- 설정 없음 → 모든 확장자 처리 (`quiz_1.pdf`, `homework_2.txt`, `image_3.jpg` 모두 처리)
- `["pdf", "txt"]` → PDF와 TXT 파일만 처리 (`quiz_1.pdf`, `homework_2.txt`만 처리, `image_3.jpg` 무시)


## Linux에서 사용

### 1. 빌드
```bash
GOOS=linux GOARCH=amd64 go build -o filemanager-linux
```

### 2. 리눅스에서 실행 권한 설정
```bash
chmod +x filemanger-linux 
```

### 설정 파일 준비 
```bash
# 샘플 설정파일 복사 후 수정해 사용
cp filemanager-config.json my-config.json 
```

### 실행
```bash
./filemanager-linux my-config.json`
```


