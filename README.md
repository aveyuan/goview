# goview

forked https://github.com/foolin/goview

Goview is a lightweight, minimalist and idiomatic template library based on golang [html/template](https://golang.org/pkg/html/template/) for building Go web application.

## Contents

- [Install](#install)
- [Features](#features)
- [Docs](#docs)
- [Supports](#supports)
- [Usage](#usage)
    - [Overview](#overview)
    - [Config](#config)
    - [Include syntax](#include-syntax)
    - [Render name](#render-name)
	- [Custom template functions](#custom-template-functions)
	- [Built-in template functions (goview)](#built-in-template-functions-goview)
	- [Go `html/template` built-ins (quick reference)](#go-htmltemplate-built-ins-quick-reference)
- [Examples](#examples)
    - [Basic example](#basic-example)
    - [Gin example](#gin-example)
	- [more examples](#more-examples)


## Install
```bash
go get github.com/aveyuan/goview
```


## Features

* **Lightweight** - use golang html/template syntax.
* **Easy** - easy use for your web application.
* **Fast** - Support configure cache template.
* **Include syntax** - Support include file.
* **Master layout** - Support configure master layout file.
* **Extension** - Support configure template file extension.
* **Easy** - Support configure templates directory.
* **Auto reload** - Support dynamic reload template(disable cache mode).
* **Multiple Engine** - Support multiple templates for frontend and backend.
* **No external dependencies** - plain ol' Go html/template.
* **Gin** - Provide Gin adapter via `supports/ginview`.


## Docs
See <https://pkg.go.dev/github.com/aveyuan/goview>


## Supports
- **ginview** goview adapter for Gin framework (see `supports/ginview`)


## Usage

### Overview

Project structure:

```go
|-- app/views/
    |--- index.html          
    |--- page.html
    |-- layouts/
        |--- footer.html
        |--- master.html
   
```

Use default instance:

```go
    //write http.ResponseWriter
    //"index" -> index.html
    goview.Render(writer, http.StatusOK, "index", goview.M{})
```

Use new instance with config:

```go

    gv := goview.New(&goview.Config{
        Root:      "views",
        Extension: ".tpl",
        Master:    "layouts/master",
        Partials:  []string{"partials/ad"},
        Funcs: template.FuncMap{
            "sub": func(a, b int) int {
                return a - b
            },
            "copy": func() string {
                return time.Now().Format("2006")
            },
        },
        DisableCache: true,
	Delims:    Delims{Left: "{{", Right: "}}"},
    })
    
    //Set new instance
    goview.Use(gv)
    
    //write http.ResponseWriter
    goview.Render(writer, http.StatusOK, "index", goview.M{})

```


Use multiple instance with config:

```go
    //============== Frontend ============== //
    gvFrontend := goview.New(&goview.Config{
        Root:      "views/frontend",
        Extension: ".tpl",
        Master:    "layouts/master",
        Partials:  []string{"partials/ad"},
        Funcs: template.FuncMap{
            "sub": func(a int, b int) int {
                return a + b
            },
            "copy": func() string {
                return time.Now().Format("2006")
            },
        },
        DisableCache: true,
	Delims:       Delims{Left: "{{", Right: "}}"},
    })
    
    //write http.ResponseWriter
    gvFrontend.Render(writer, http.StatusOK, "index", goview.M{})
    
    //============== Backend ============== //
    gvBackend := goview.New(&goview.Config{
        Root:      "views/backend",
        Extension: ".tpl",
        Master:    "layouts/master",
        Partials:  []string{"partials/ad"},
        Funcs: template.FuncMap{
            "sub": func(a int, b int) int {
                return a - b
            },
            "copy": func() string {
                return time.Now().Format("2006")
            },
        },
        DisableCache: true,
	Delims:       Delims{Left: "{{", Right: "}}"},
    })
    
    //write http.ResponseWriter
    gvBackend.Render(writer, http.StatusOK, "index", goview.M{})

```

### Config

```go
goview.Config{
    Root:      "views", //template root path
    Extension: ".tpl", //file extension
    Master:    "layouts/master", //master layout file
    Partials:  []string{"partials/head"}, //partial files
    Funcs: template.FuncMap{
        "sub": func(a, b int) int {
            return a - b
        },
        // more funcs
    },
    DisableCache: false, //if disable cache, auto reload template file for debug.
    Delims:       Delims{Left: "{{", Right: "}}"},
}
```

### Include syntax

```go
//template file
{{include "layouts/footer"}}
```

### Render name: 

Render name use `index` without `.html` extension, that will render with master layout.

- **"index"** - Render with master layout.
- **"index.html"** - Not render with master layout.

```
Notice: `.html` is default template extension, you can change with config
```


Render with master

```go
//use name without extension `.html`
goview.Render(w, http.StatusOK, "index", goview.M{})
```

The `w` is instance of  `http.ResponseWriter`

Render only file(not use master layout)
```go
//use full name with extension `.html`
goview.Render(w, http.StatusOK, "page.html", goview.M{})
```

### Custom template functions

We have two type of functions `global functions`, and `temporary functions`.

`Global functions` are set within the `config`.

```go
goview.Config{
	Funcs: template.FuncMap{
		"reverse": e.Reverse,
	},
}
```

```go
//template file
{{ reverse "route-name" }}
```

`Temporary functions` are set inside the handler.

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	err := goview.Render(w, http.StatusOK, "index", goview.M{
		"reverse": e.Reverse,
	})
	if err != nil {
		fmt.Fprintf(w, "Render index error: %v!", err)
	}
})
```
```go
//template file
{{ call $.reverse "route-name" }}
```

### Built-in template functions (goview)

goview provides a set of built-in template functions globally.

- These functions are available in all templates without extra configuration.
- `Config.Funcs` can override any built-in function with the same name.

#### default

`default(d, v) any`

Returns `d` when `v` is considered empty.

Empty means:

- `nil`
- `false`
- numeric `0`
- empty string
- empty slice/map/array
- zero-value struct

Example:

```go
{{ default "N/A" .Name }}
```

#### list

`list(v1, v2, ...) []any`

Example:

```go
{{ index (list 1 2 3) 1 }}
```

#### dict

`dict(k1, v1, k2, v2, ...) (map[string]any, error)`

Keys must be strings and the argument count must be even.

Example:

```go
{{ $m := dict "a" 1 "b" 2 }}
{{ index $m "a" }}
```

#### urlquery

`urlquery(v) template.URL`

URL-escapes the value for use in query strings.

Example:

```go
<a href="/search?q={{ urlquery .Keyword }}">search</a>
```

#### date

`date(t, layout) (string, error)`

Supported input types: `time.Time`, `*time.Time`, `int64`, `int`.

Example:

```go
{{ date .CreatedAt "2006-01-02 15:04:05" }}
```

#### json

`json(v) (template.JS, error)`

Marshals to JSON using `encoding/json`.

Example:

```go
<script>window.__DATA__ = {{ json . }}</script>
```

#### escape

`escape(v) template.HTML`

HTML-escapes the value (uses `template.HTMLEscapeString`).

Example:

```go
{{ escape .UnsafeText }}
```

#### safeHTML

`safeHTML(v) template.HTML`

Marks the value as trusted HTML.

Use this only when you know the content is safe. Do NOT use it directly on user input.

Example:

```go
{{ safeHTML .TrustedHTML }}
```

#### include

`include(name) (template.HTML, error)`

Renders another template and returns the result.

Example:

```go
{{ include "layouts/footer" }}
```

### Go `html/template` built-ins (quick reference)

goview is based on Go `html/template`, so the standard template syntax and the built-in functions are available.

#### Syntax essentials

- **Action**: `{{ ... }}`
- **Pipeline**: `{{ .Title | printf "%q" }}`
- **If/Else**:

```go
{{ if .OK }}yes{{ else }}no{{ end }}
```

- **Range**:

```go
{{ range .Items }}{{ . }}{{ end }}
```

- **With**:

```go
{{ with .User }}{{ .Name }}{{ end }}
```

- **Template / define**:

```go
{{ define "content" }}...{{ end }}
{{ template "content" . }}
```

#### Common built-in functions

- **Logic/compare**: `and`, `or`, `not`, `eq`, `ne`, `lt`, `le`, `gt`, `ge`
- **Length**: `len`
- **Index**: `index`
- **Print**: `print`, `printf`, `println`
- **Call function**: `call`

See the official docs for the full list and details:

- https://pkg.go.dev/html/template
- https://pkg.go.dev/text/template



## Examples

See `_examples/` in this repository.


### Basic example

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/aveyuan/goview"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := goview.Render(w, http.StatusOK, "index", goview.M{"title": "Index title!"})
		if err != nil {
			fmt.Fprintf(w, "Render index error: %v!", err)
		}
	})

	fmt.Println("Listening and serving HTTP on :9090")
	_ = http.ListenAndServe(":9090", nil)
}
```

### Gin example

```go
package main

import (
	"net/http"

	"github.com/aveyuan/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.HTMLRender = ginview.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index", gin.H{"title": "Index title!"})
	})

	_ = router.Run(":9090")
}
```

### More examples

See `_examples/` in this repository.
 
