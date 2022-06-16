package ast_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/commentcov/commentcov-plugin-go/ast"
)

// TestNormalize is the unittest for Normalize.
func TestNormalize(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want string
	}{

		{
			name: "no prefix whitespace",
			str:  "hoge",
			want: "hoge",
		},
		{
			name: "prefix whitespace",
			str:  " hoge",
			want: "hoge",
		},
		{
			name: "multi prefix whitespace",
			str:  "   hoge",
			want: "hoge",
		},
		{
			name: "prefix and suffix whitespace",
			str:  " hoge ",
			want: "hoge ",
		},
		{
			name: "whitespace and return code",
			str:  " hoge \n fuga ",
			want: "hoge \n fuga ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ast.Normalize(tt.str)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("string values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

// TestIsOnlyNoLintAnnotation is the unittest for IsOnlyNoLintAnnotation.
func TestIsOnlyNoLintAnnotation(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{

		{
			name: "nolint annotation without prefix whitespace",
			str:  "nolint:funlen",
			want: true,
		},
		{
			name: "nolint annotation with prefix whitespace",
			str:  " nolint:funlen",
			want: true,
		},
		{
			name: "nolint annotation with prefix and suffix whitespace",
			str:  " nolint:funlen ",
			want: true,
		},
		{
			name: "nolint annotation with return code",
			str:  "nolint:funlen\n",
			want: true,
		},
		{
			name: "any string",
			str:  "This is comment.",
			want: false,
		},
		{
			name: "nolint annotation with multi lines",
			str:  "nolint:funlen\nThis is comment.",
			want: false,
		},
		{
			name: "nolint annotation with multi lines",
			str:  "This is comment.\nnolint:funlen",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ast.IsOnlyNoLintAnnotation(tt.str)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("bool values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}
