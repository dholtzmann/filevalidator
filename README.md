filevalidator
========

[Go](http://golang.org) simple package for validating file uploads.

## Installation

```bash
go get -u github.com/dholtzmann/filevalidator
```

## Basic Example

```go
package main

import (
	fv "github.com/dholtzmann/filevalidator"
	"html/template"
	"net/http"
)

func main() {
	// setup HTTP server...
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseGlob("./templates/*.html"))

	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		panic(err.Error())
	}

	fields := map[string]*Field{
		"imageupload":   NewField(true, false, FileSize(100, 100000), MimeTypes([]string{"image/png", "image/jpeg"}), MinPixelSize(100, 100), MaxPixelSize(500, 500)),
		"fileupload":    NewField(true, true, FileSize(100, 100000), MimeTypes([]string{"application/x-gzip"})),
		"missingupload": NewField(false, true, nil),
	}

	err, validator := New(fields)
	if err != nil {
		panic(err.Error())
	}

	isValid, errors := validator.Validate(req.MultipartForm.File)
	if isValid {
		// save file...
	} else {
		// display errors
		tplVars := make(map[string]interface{})
		tplVars["FormErrors"] = errors
		err = tpl.ExecuteTemplate(w, "index.gohtml", tplVars)
	}
}
```

## Validation rules

- NewField(required bool, singleFile bool, rules []Rule)
- FileSize(min, max uint64)
- MimeTypes(list []string)
- MinPixelSize(width, height uint32)
- MaxPixelSize(width, height uint32)

See the file 'validate_test.go' for examples of how to use the validation rules.
