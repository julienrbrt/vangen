package main

import (
	"bytes"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestGenerate(t *testing.T) {
	testCases := []struct {
		description string
		domain      string
		docsDomain  string
		pkg         string
		r           repository
		goldenFile  string
		expectedErr error
	}{
		{
			description: "simple",
			domain:      "example.com",
			docsDomain:  "godoc.org",
			pkg:         "pkg1",
			r: repository{
				Prefix: "pkg1",
				Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}},
				Type:   "git",
				URL:    "https://repositoryhost.com/example/go-pkg1",
			},
			goldenFile:  "simple.pkgs.golden.html",
			expectedErr: nil,
		},
		{
			description: "hidden",
			domain:      "example.com",
			docsDomain:  "godoc.org",
			pkg:         "pkg1",
			r: repository{
				Prefix: "pkg1",
				Hidden: true,
				Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}},
				Type:   "git",
				URL:    "https://repositoryhost.com/example/go-pkg1",
			},
			goldenFile:  "hidden.pkgs.golden.html",
			expectedErr: nil,
		},
		{
			description: "custom source urls",
			domain:      "example.com",
			docsDomain:  "pkg.go.dev",
			pkg:         "pkg1",
			r: repository{
				Prefix: "pkg1",
				Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}},
				Type:   "git",
				URL:    "https://repositoryhost.com/example/go-pkg1",
				SourceURLs: sourceURLs{
					Home: "https://repositoryhost.com/example/go-pkg1/home",
					Dir:  "https://repositoryhost.com/example/go-pkg1/browser{/dir}",
					File: "https://repositoryhost.com/example/go-pkg1/view{/dir}{/file}",
				},
				Website: website{
					URL: "https://www.example.com",
				},
			},
			goldenFile:  "custom-source-urls.pkgs.golden.html",
			expectedErr: nil,
		},
		{
			description: "sub-package",
			domain:      "example.com",
			pkg:         "pkg1/subpkg1",
			r: repository{
				Prefix: "pkg1",
				Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}},
				Type:   "git",
				URL:    "https://repositoryhost.com/example/go-pkg1",
				SourceURLs: sourceURLs{
					Home: "https://repositoryhost.com/example/go-pkg1/home",
					Dir:  "https://repositoryhost.com/example/go-pkg1/browser{/dir}",
					File: "https://repositoryhost.com/example/go-pkg1/view{/dir}{/file}",
				},
				Website: website{
					URL: "https://www.example.com",
				},
			},
			goldenFile:  "sub-package.pkgs.golden.html",
			expectedErr: nil,
		},
		{
			description: "sub-package hidden",
			domain:      "example.com",
			pkg:         "pkg1/subpkg1",
			r: repository{
				Prefix: "pkg1",
				Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}, {Name: "subpkg3", Hidden: true}},
				Type:   "git",
				URL:    "https://repositoryhost.com/example/go-pkg1",
				SourceURLs: sourceURLs{
					Home: "https://repositoryhost.com/example/go-pkg1/home",
					Dir:  "https://repositoryhost.com/example/go-pkg1/browser{/dir}",
					File: "https://repositoryhost.com/example/go-pkg1/view{/dir}{/file}",
				},
				Website: website{
					URL: "https://www.example.com",
				},
			},
			goldenFile:  "sub-package-hidden.pkgs.golden.html",
			expectedErr: nil,
		},
		{
			description: "github defaults",
			domain:      "example.com",
			docsDomain:  "pkg.go.dev",
			pkg:         "pkg1",
			r: repository{
				Prefix: "pkg1",
				Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}},
				URL:    "https://github.com/example/go-pkg1",
			},
			goldenFile:  "github-defaults.pkgs.golden.html",
			expectedErr: nil,
		},
		{
			description: "sub-package github defaults",
			domain:      "example.com",
			docsDomain:  "pkg.go.dev",
			pkg:         "pkg1/subpkg1",
			r: repository{
				Prefix: "pkg1",
				Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}},
				URL:    "https://github.com/example/go-pkg1",
			},
			goldenFile:  "sub-package-github-defaults.pkgs.golden.html",
			expectedErr: nil,
		},
		{
			description: "gitlab defaults",
			domain:      "example.com",
			docsDomain:  "",
			pkg:         "pkg1",
			r: repository{
				Prefix: "pkg1",
				Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}},
				URL:    "https://gitlab.com/example/go-pkg1",
			},
			goldenFile:  "gitlab-defaults.pkgs.golden.html",
			expectedErr: nil,
		},
		{
			description: "sub-package gitlab defaults",
			domain:      "example.com",
			docsDomain:  "pkg.go.dev",
			pkg:         "pkg1/subpkg1",
			r: repository{
				Prefix: "pkg1",
				Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}},
				URL:    "https://gitlab.com/example/go-pkg1",
			},
			goldenFile:  "sub-package-gitlab-defaults.pkgs.golden.html",
			expectedErr: nil,
		},
		{
			description: "github defaults with custom source",
			domain:      "example.com",
			docsDomain:  "",
			pkg:         "pkg1",
			r: repository{
				Prefix: "pkg1",
				Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}},
				Type:   "git",
				URL:    "https://github.com/example/go-pkg1",
				SourceURLs: sourceURLs{
					Home: "https://github.com/example/go-pkg1",
					Dir:  "https://github.com/example/go-pkg1/tree/branch{/dir}",
					File: "https://github.com/example/go-pkg1/blob/branch{/dir}/{file}#L{line}",
				},
				Website: website{
					URL: "https://www.example.com",
				},
			},
			goldenFile:  "github-defaults-custom-source.pkgs.golden.html",
			expectedErr: nil,
		},
		{
			description: "single module deployment that has no 'prefix'",
			domain:      "example.com",
			docsDomain:  "",
			pkg:         "",
			r: repository{
				Prefix: "",
				Subs:   []sub{{Name: "subpkg1"}, {Name: "subpkg2"}},
				Type:   "git",
				URL:    "https://github.com/example/go-pkg1",
				SourceURLs: sourceURLs{
					Home: "https://github.com/example/go-pkg1",
					Dir:  "https://github.com/example/go-pkg1/tree/branch{/dir}",
					File: "https://github.com/example/go-pkg1/blob/branch{/dir}/{file}#L{line}",
				},
				Website: website{
					URL: "https://www.example.com",
				},
			},
			goldenFile:  "single-module-deployment-no-prefix.pkgs.golden.html",
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			var out bytes.Buffer
			err := generate_package(&out, tc.domain, tc.docsDomain, tc.pkg, tc.r)
			if err != tc.expectedErr {
				t.Errorf("Test case %q got err %#v, want %#v", tc.description, err, tc.expectedErr)
			} else {
				expected := loadGolden(t, tc.goldenFile)
				if out.String() != expected {
					dmp := diffmatchpatch.New()
					diffs := dmp.DiffMain(expected, out.String(), false)
					t.Errorf("Test case %q mismatch:\n%s", tc.description, dmp.DiffPrettyText(diffs))
				}
			}
		})
	}
}
