# FileManager

리눅스/윈도우 시스템을 위한 파일 정리 자동화 도구

## 📋 개요

FileManager는 패턴 기반으로 파일을 정리하는 도구입니다. 안전한 백업과 함께 중복 파일 삭제, 파일 이름 변경을 자동화합니다.

### 주요 기능
- 🔄 **안전한 백업** - 원본 파일을 건드리기 전에 자동 백업
- 📁 **패턴 기반 정리** - 접두사별로 최신 파일만 남기고 나머지 삭제
- 🔧 **유연한 설정** - JSON 설정으로 다양한 정리 규칙 지원
- ✅ **미리보기 모드** - dry_run으로 안전하게 테스트


## ⚙️ 설정 파일 구조

### 기본 설정 (`config.json`)

```json
{
  "source_path": "/path/to/root",
  "backup_path": "/path/to/backup",
  "target_folders": ["paper", "homework"],
  "target_depth": 3,
  "dry_run": true
}
```

### 설정 항목 설명

| 항목 | 설명 | 예시 |
|------|------|------|
| `source_path` | 정리할(복사할) 루트 폴더 경로 | `"/home/user/documents"` |
| `work_path` | 작업 경로(복사할 위치) | `"/home/user/backup"` |
| `target_folders` | 처리할 대상 폴더들(root 하위 폴더 기준) | `["paper", "homework", "assignments"]` |
| `target_depth` | target_folder 기준 탐색 깊이 | `3` (paper → 10 → 1001 → 파일) |
| `dry_run` | 미리보기 모드 (실제 실행 안함) | `true` / `false` |

### target_depth 설명

`target_depth`는 target_folder를 기준으로 한 상대적 깊이입니다:

- `target_depth: 1` → `paper/` (target_folder 바로 아래)
- `target_depth: 2` → `paper/f1/` (paper 하위 폴더까지)
- `target_depth: 3` → `paper/f1/1111/` (f1 하위  폴더까지, **파일이 있는 위치**)