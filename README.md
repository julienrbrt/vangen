# vangen

Vangen is a tool for generating static HTML for Go vanity import paths.

Go vanity import paths work by serving a HTML file that tells the `go get` tool where to download the source from. You can still host the source code at Github, BitBucket, but the vanity URL gives you portability and other benefits.

## Why

* Maintain Go vanity import paths with a simple definition file `vangen.json`.
* Host Go vanity import paths using static hosting. No need for Google AppEngine, Heroku, etc. Host the files on Github Pages, AWS S3, Google Cloud Storage, etc.

## Install

```sh
go install github.com/julienrbrt/vangen@latest
```

## Usage

1. Create a `vangen.json` (see examples below)
2. Run `vangen`
3. Host the files outputted in `vangen/` at your domain
4. Try it out with `go get [domain]/[package]`

```sh
$ vangen -help
Vangen is a tool for generating static HTML for hosting Go repositories at a vanity import path.

Usage:

  vangen [-config=vangen.json] [-out=vangen/]

Flags:

  -config filename
        vangen json configuration filename (default "vangen.json")
  -help
        print this help list
  -out directory
        output directory that static files will be written to (default "vangen/")
  -verbose
        print verbose output when run
  -version
        print program version
```

## Examples

### Minimal

The repository `type` and `source` properties will be set automatically when `url` begins with `https://github.com` or `https://gitlab.com`. Below is a minimal config for a project hosted on GitHub.

```json
{
  "domain": "4d63.com",
  "repositories": [
    {
      "prefix": "optional",
      "subs": [
        "template"
      ],
      "url": "https://github.com/leighmcculloch/go-optional"
    }
  ]
}
```

### All fields

```json
{
  "domain": "4d63.com",
  "docsDomain": "pkg.go.dev",
  "repositories": [
    {
      "prefix": "optional",
      "subs": [
        "template"
      ],
      "type": "git",
      "hidden": false,
      "url": "https://github.com/leighmcculloch/go-optional",
      "source": {
        "home": "https://github.com/leighmcculloch/go-optional",
        "dir": "https://github.com/leighmcculloch/go-optional/tree/master{/dir}",
        "file": "https://github.com/leighmcculloch/go-optional/blob/master{/dir}/{file}#L{line}"
      },
      "website": {
        "url": "https://github.com/leighmcculoch/go-optional"
      }
    }
  ]
}
```

### Thanks

This project is a fork of [leighmcculloch/vangen](https://leighmcculloch/vangen) with the following changes:

* Changed go.mod name to `github.com/julienrbrt/vangen` for easier installation of this fork
* UI improvements
