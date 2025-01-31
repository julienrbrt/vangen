package main

import (
	"fmt"
	"html/template"
	"io"
	"strings"
)

func generate_package(w io.Writer, domain, docsDomain, pkg string, r repository) error {
	const html = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Domain}}/{{.Package}}</title>
<meta name="go-import" content="{{.Domain}}{{.Repository.PrefixPath}} {{.Repository.Type}} {{.Repository.URL}}">
<meta name="go-source" content="{{.Domain}}{{.Repository.PrefixPath}} {{.Repository.SourceURLs.Home}} {{.Repository.SourceURLs.Dir}} {{.Repository.SourceURLs.File}}">
<style>
* { font-family: sans-serif; }
body { margin: 16px; background-color: #f4f4f4; }
.content {
  max-width: 600px;
  margin: 0 auto;
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  padding: 16px;
}
code {
  display: block;
  font-family: monospace;
  font-size: 1em;
  background-color: #eee;
  padding: 1em;
  margin-bottom: 16px;
}
ul { margin-top: 16px; margin-bottom: 16px; }
@media (max-width: 600px) {
  .content {
    max-width: 100%;
    margin: 0 8px;
    box-shadow: none;
  }
}
</style>
<body>
<div class="content">
<h2>{{.Domain}}/{{.Package}}</h2>
<code>go get {{.Domain}}/{{.Package}}</code>
<code>import "{{.Domain}}/{{.Package}}"</code>
Home: <a href="{{.HomeURL}}">{{.HomeURL}}</a><br/>
Source: <a href="{{.Repository.URL}}">{{.Repository.URL}}</a><br/>
{{if .Repository.Subs -}}Sub-packages:<ul>{{end -}}
{{range $i, $s := .Repository.Subs -}}{{if not $s.Hidden -}}<li><a href="/{{$.Repository.SubPath $i}}">{{$.Domain}}/{{$.Repository.SubPath $i}}</a></li>{{end -}}{{end -}}
{{if .Repository.Subs -}}</ul>{{end -}}
<a href="http://{{.Domain}}" style="margin-bottom:16px; display:inline-block;">Back</a>
</div>
</body>
</html>`

	tmpl, err := template.New("").Parse(html)
	if err != nil {
		return fmt.Errorf("error loading template: %v", err)
	}

	var homeURL string
	if r.Website.URL != "" {
		homeURL = r.Website.URL
	} else {
		if docsDomain == "" {
			docsDomain = "pkg.go.dev"
		}
		homeURL = fmt.Sprintf("https://%s/%s/%s", docsDomain, domain, pkg)
	}

	if strings.HasPrefix(r.URL, "https://github.com") || strings.HasPrefix(r.URL, "https://gitlab.com") {
		if r.Type == "" {
			r.Type = "git"
		}
		if r.SourceURLs.Home == "" {
			r.SourceURLs.Home = r.URL
		}
		if r.SourceURLs.Dir == "" {
			r.SourceURLs.Dir = r.URL + "/tree/master{/dir}"
		}
		if r.SourceURLs.File == "" {
			r.SourceURLs.File = r.URL + "/blob/master{/dir}/{file}#L{line}"
		}
	}

	if r.SourceURLs.Home == "" {
		r.SourceURLs.Home = "_"
	}
	if r.SourceURLs.Dir == "" {
		r.SourceURLs.Dir = "_"
	}
	if r.SourceURLs.File == "" {
		r.SourceURLs.File = "_"
	}

	data := struct {
		Domain     string
		Package    string
		Repository repository
		HomeURL    string
	}{
		Domain:     domain,
		Package:    pkg,
		Repository: r,
		HomeURL:    homeURL,
	}

	err = tmpl.ExecuteTemplate(w, "", data)
	if err != nil {
		return fmt.Errorf("generating template: %v", err)
	}

	return nil
}
