package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func loadGolden(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("failed to load golden file %q: %v", name, err)
	}
	return string(b)
}

func TestGenerateIndex(t *testing.T) {
	testCases := []struct {
		description string
		domain      string
		r           []repository
		goldenFile  string
		expectedErr error
	}{
		{
			description: "basic",
			domain:      "example.com",
			r: []repository{
				{
					Prefix: "pkg1",
					Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}},
					Main:   true,
				},
				{
					Prefix: "pkg2",
					Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2/subsubpkg1"}},
				},
				{
					Prefix: "pkg3",
					Hidden: true,
				},
			},
			goldenFile:  "basic.index.golden.html",
			expectedErr: nil,
		},
		{
			description: "hidden sub-package",
			domain:      "example.com",
			r: []repository{
				{
					Prefix: "pkg1",
					Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}, {Name: "subpkg3", Hidden: true}},
					Main:   true,
				},
				{
					Prefix: "pkg2",
					Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2/subsubpkg1"}, {Name: "subpkg2/subsubpkg2", Hidden: true}},
				},
			},
			goldenFile:  "hidden_sub_package.index.golden.html",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			var out bytes.Buffer
			err := generate_index(&out, tc.domain, tc.r)
			if err != tc.expectedErr {
				t.Errorf("got err=%v, want=%v", err, tc.expectedErr)
			}
			expected := loadGolden(t, tc.goldenFile)
			if out.String() != expected {
				dmp := diffmatchpatch.New()
				diffs := dmp.DiffMain(expected, out.String(), false)
				t.Errorf("output mismatch for %q:\n%s", tc.description, dmp.DiffPrettyText(diffs))
			}
		})
	}
}
