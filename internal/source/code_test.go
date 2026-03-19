package source

import (
	"reflect"
	"testing"
)

func TestCleanedComments(t *testing.T) {
	tests := []struct {
		name     string
		rawSQL   string
		cs       CommentSyntax
		expected []string
	}{
		{
			name: "single line dash",
			rawSQL: "-- comment\nSELECT 1;",
			cs: CommentSyntax{Dash: true},
			expected: []string{" comment"},
		},
		{
			name: "multi line slash star",
			rawSQL: "/*\nline 1\nline 2\n*/\nSELECT 1;",
			cs: CommentSyntax{SlashStar: true},
			expected: []string{"", "line 1", "line 2", ""},
		},
		{
			name: "mixed comments",
			rawSQL: "-- dash\n/* star */\n# hash",
			cs: CommentSyntax{Dash: true, SlashStar: true, Hash: true},
			expected: []string{" dash", " star ", " hash"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := CleanedComments(tt.rawSQL, tt.cs)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
