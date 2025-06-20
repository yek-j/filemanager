package plugins

import (
	"testing"
)

type parseFileNameTest struct {
	input          string // 입력 파일명
	expectedPrefix string
	expectedNumber int
	expectedExt    string
	expectedValid  bool // 예상 결과
}

func TestParseFileName(t *testing.T) {
	tests := []parseFileNameTest{
		// 정상 케이스
		{"quiz_123.pdf", "quiz", 123, ".pdf", true},
		{"file_456.txt", "file", 456, ".txt", true},
		{"homework_1.docx", "homework", 1, ".docx", true},
		{"test_999.jpg", "test", 999, ".jpg", true},
		{"a_0.pdf", "a", 0, ".pdf", true},

		// 실패 케이스 - 언더스코어 없음
		{"quiz123.pdf", "", 0, "", false},
		{"file.txt", "", 0, "", false},

		// 실패 케이스 - 숫자 아님
		{"quiz_abc.pdf", "", 0, "", false},
		{"file_text.txt", "", 0, "", false},
		{"homework_.pdf", "", 0, "", false},

		// 실패 케이스 - 특수 상황
		{"_123.pdf", "", 0, "", false}, // 접두사 없음
		{"", "", 0, "", false},         // 빈 문자열
		{"quiz_", "", 0, "", false},    // 확장자도 숫자도 없음
		{"_", "", 0, "", false},        // 언더스코어만

		// 엣지 케이스
		{"quiz_123", "quiz", 123, "", true}, // 확장자 없음 (허용?)
	}

	for _, tc := range tests {
		prefix, number, ext, valid := parseFileName(tc.input)

		if prefix != tc.expectedPrefix {
			t.Fatalf("input: %q, expected prefix: %q, got: %q",
				tc.input, tc.expectedPrefix, prefix)
		}

		if number != tc.expectedNumber {
			t.Fatalf("input: %q, expected number: %q, got: %q",
				tc.input, tc.expectedNumber, number)
		}

		if ext != tc.expectedExt {
			t.Fatalf("input: %q, expected ext: %q, got: %q",
				tc.input, tc.expectedExt, ext)
		}

		if valid != tc.expectedValid {
			t.Fatalf("input: %q, expected valid: %t, got: %t",
				tc.input, tc.expectedValid, valid)
		}
	}
}
