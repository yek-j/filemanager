# FileManager

리눅스/윈도우 시스템을 위한 파일 정리 자동화 도구

## 📋 개요

FileManager는 패턴 기반으로 파일을 정리하는 도구입니다. 안전한 백업과 함께 중복 파일 삭제, 파일 이름 변경, 파일 이동을 자동화합니다.

### 주요 기능
- 🔄 **안전한 백업** - 원본 파일을 건드리기 전에 자동 백업
- 📁 **패턴 기반 정리** - 접두사별로 최신 파일만 남기고 나머지 삭제
- 🔧 **유연한 설정** - JSON 설정으로 다양한 정리 규칙 지원
- ⚡ **선택적 복사** - 전체가 아닌 필요한 폴더만 복사하여 효율성 향상
- 🎯 **다중 플러그인** - 여러 플러그인을 조합하여 복잡한 워크플로우 구성

## ⚙️ 설정 파일 구조

### 기본 설정 (`config.json`)

```json
{
  "source_path": "/path/to/root",
  "work_path": "/path/to/work",
  "target_folders": ["paper", "homework"],
  "file_depth": 3,
  "plugin": [
    {
      "name": "underscore_number",
      "config": {
        "allowed_extensions": ["pdf", "txt"],
        "target_folders": ["paper", "homework"]
      }
    },
    {
      "name": "file_relocator",
      "config": {
        "file_extensions": ["mp3"],
        "source_location": "music",
        "target_location": "files",
        "create_folder": true,
        "target_folders": ["paper"]
      }
    }
  ],
  "selective_copy": true
}
```

### 설정 항목 설명

| 항목 | 설명 | 예시 | 기본값 |
|------|------|------|--------|
| `source_path` | 정리할(복사할) 루트 폴더 경로 | `"/home/user/documents"` | 필수 |
| `work_path` | 작업 경로(복사할 위치) | `"/home/user/work"` | 필수 |
| `target_folders` | **복사할 대상 폴더들 (모든 플러그인의 target_folders 포함)** | `["paper", "homework"]` | 필수 |
| `file_depth` | target_folder 기준 탐색 깊이 | `3` | 필수 |
| `plugin` | **실행할 플러그인 목록 (배열, 순서대로 실행)** | `[{name: "...", config: {}}]` | 필수 |
| `selective_copy` | 선택적 복사 모드 사용 여부 | `true` | `false` |

### 📌 target_folders 설정 규칙

**최상위 `target_folders`는 모든 플러그인에서 사용하는 폴더를 포함해야 합니다.**

```
최상위 target_folders = 플러그인1의 target_folders ∪ 플러그인2의 target_folders ∪ ...
```

**예시:**
- underscore_number가 `["paper", "homework"]` 사용
- file_relocator가 `["music"]` 사용
- → 최상위는 `["paper", "homework", "music"]` 필요

> ⚠️ **주의:** `selective_copy: true` 시 최상위 `target_folders`에 없는 폴더는 복사되지 않아 플러그인이 작업할 수 없습니다.

> ⚠️ **플러그인 실행 순서:** `plugin` 배열에 정의된 순서대로 실행됩니다. 순서 변경 시 결과가 달라질 수 있습니다.

### file_depth 설명

`file_depth`는 target_folder를 기준으로 한 상대적 깊이입니다:

- `file_depth: 1` → `paper/` (target_folder 바로 아래)
- `file_depth: 2` → `paper/f1/` (paper 하위 폴더까지)
- `file_depth: 3` → `paper/f1/1111/` (f1 하위 폴더까지, **파일이 있는 위치**)

### 복사 모드 비교

| 모드 | `selective_copy` | 복사 대상 | 장점 | 단점 |
|------|------------------|-----------|------|------|
| **전체 복사** | `false` (기본값) | source_path 전체 | 완전한 백업 | 느림, 공간 많이 사용 |
| **선택적 복사** | `true` | target_folders만 | 빠름, 공간 절약 | 부분 백업만 |

## 🔌 사용 가능한 플러그인

### 개요
- 플러그인은 **배열에 정의된 순서대로** 순차 실행됩니다
- 각 플러그인은 독립적인 설정(`config`)을 가집니다
- 여러 플러그인을 조합하여 복잡한 파일 정리 워크플로우 구성 가능

### 플러그인 목록
1. [underscore_number](#1-underscore_number) - 패턴 기반 파일 정리
2. [file_relocator](#2-file_relocator) - 파일 일괄 이동

---

### 1. underscore_number

**설명:** `문자열_숫자.확장자` 패턴의 파일을 정리

**동작:**
- 각 폴더 내에서 같은 접두사 중 가장 큰 숫자 파일만 남기고 삭제 남은 파일은 `prefix_1` 형태로 일괄 변경

**예시:**
- `quiz_10.pdf`, `quiz_25.pdf`, `quiz_5.pdf` 존재
- → `quiz_25.pdf`만 남기고 나머지 삭제
- → `quiz_25.pdf`를 `quiz_1.pdf`로 변경

**로그 파일:**
- 작업 완료 후 `work_path`에 자동 생성
- 삭제된 파일, 이름 변경된 파일 목록 포함

#### 플러그인 설정 (`config`)

```json
{
  "allowed_extensions": ["pdf", "txt", "docx"],
  "target_folders": ["paper", "homework"]
}
```

| 설정 항목 | 설명 | 예시 | 필수 | 기본값 |
|-----------|------|------|------|--------|
| `allowed_extensions` | 처리할 파일 확장자 목록 (점 제외) | `["pdf", "txt"]` | ❌ | 모든 확장자 |
| `target_folders` | 작업할 대상 폴더 목록 | `["paper", "homework"]` | ✅ | - |

**확장자 필터링 예시:**
- 설정 없음 → 모든 확장자 처리 (`quiz_1.pdf`, `homework_2.txt` 모두 처리)
- `["pdf", "txt"]` → PDF와 TXT만 처리 (`image_3.jpg` 무시)

#### 사용 예시

```json
{
  "name": "underscore_number",
  "config": {
    "allowed_extensions": ["pdf"],
    "target_folders": ["paper", "assignments"]
  }
}
```

---

### 2. file_relocator

**설명:** 지정된 파일들을 일괄 이동합니다.
- 단순 구조: `file_depth` 기반 탐색
- 복잡한 구조: 정규식 패턴 사용

**주요 기능:**
- ✅ 확장자별 파일 필터링
- ✅ 하위 폴더 검색 옵션
- ✅ 대상 폴더 자동 생성
- ⚠️ **이동만 지원** (복사 X - 원본 파일 삭제됨)

#### 플러그인 설정 (`config`)

```json
{
  "file_extensions": ["mp3", "wav"],
  "file_pattern": "[0-9]+\\.pdf",
  "source_location": "downloads",
  "target_location": "music",
  "target_folders": ["paper"],
  "create_folder": true,
  "search_subdirs": true,
  "overwrite_files": false,
  "use_pattern": false
}
```

| 설정 항목 | 설명 | 타입 | 필수 | 기본값 |
|-----------|------|------|------|--------|
| `file_extensions` | 이동할 파일 확장자 목록 (점 제외) | `string[]` | ✅ | - |
| `file_pattern` | 파일명 필터링 정규식 (선택적) | `string` | ❌ | - |
| `source_location` | 파일이 위치한 경로 (상대 경로) | `string` | ✅ | - |
| `target_location` | 이동할 대상 경로 (상대 경로) | `string` | ✅ | - |
| `target_folders` | 작업 대상 폴더 목록 | `string[]` | ✅ | - |
| `create_folder` | 대상 폴더 자동 생성 여부 | `boolean` | ❌ | `false` |
| `search_subdirs` | 하위 폴더까지 검색 여부 | `boolean` | ❌ | `false` |
| `overwrite_files` | 파일 덮어쓰기 허용 여부 | `boolean` | ❌ | `false` |
| `use_pattern` | depth 기반(false) / pattern 기반(true) | `boolean` | ❌ | `false` |
| `file_pattern` | 경로 탐색 방식 - depth 기반(false) / pattern 기반(true) | `boolean` | ❌ | `false` |

#### 사용 예시

**시나리오 1: 음악 파일을 music 폴더로 이동**
```json
{
  "name": "file_relocator",
  "config": {
    "file_extensions": ["mp3", "wav", "flac"],
    "source_location": "downloads",
    "target_location": "music",
    "target_folders": ["paper"],
    "create_folder": true,
    "search_subdirs": true
  }
}
```

**시나리오 2: 이미지 파일 정리**
```json
{
  "name": "file_relocator",
  "config": {
    "file_extensions": ["jpg", "png"],
    "source_location": "temp",
    "target_location": "images",
    "target_folders": ["homework"],
    "create_folder": true,
    "overwrite_files": false
  }
}
```


## 🚀 사용 방법

### Linux에서 사용

#### 1. 빌드
```bash
GOOS=linux GOARCH=amd64 go build -o filemanager-linux
```

#### 2. 실행 권한 설정
```bash
chmod +x filemanager-linux 
```

#### 3. 설정 파일 준비
```bash
# 샘플 설정파일 복사 후 수정
cp filemanager-config.json my-config.json
# 에디터로 경로 수정
nano my-config.json
```

#### 4. 실행
```bash
./filemanager-linux my-config.json
```

---

## ⚠️ 주의사항

### 플러그인 사용 시
- **실행 순서 중요:** 플러그인은 배열 순서대로 실행되므로 순서 변경 시 결과가 달라질 수 있습니다
- **file_relocator 사용 시:** 이동(move) 기능만 지원하므로 원본 파일이 삭제됩니다
- **백업 권장:** 중요한 데이터는 별도 백업 후 작업하세요

### 설정 파일 작성 시
- **target_folders 일치:** 최상위 `target_folders`는 모든 플러그인의 `target_folders`를 포함해야 합니다
- **경로 확인:** `source_path`와 `work_path`가 올바른지 확인하세요
- **file_depth 검증:** 실제 폴더 구조와 맞는지 확인하세요

### 첫 실행 시
- 작은 테스트 폴더로 먼저 실행해보세요
- 로그 파일을 확인하여 의도한 대로 작동하는지 검증하세요