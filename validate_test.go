package filevalidator

import (
	"testing"
)

func TestValidate(t *testing.T) {
	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}, &testFile{blueJPG, "imageupload", "blue.jpg"}, &testFile{carTARGZ, "fileupload", "car.tar.gz"})

	// this should validate without a problem
	fields := map[string]*Field{
		"imageupload":   NewField(true, false, FileSize(100, 100000), MimeTypes([]string{"image/png", "image/jpeg"}), MinPixelSize(100, 100), MaxPixelSize(500, 500)),
		"fileupload":    NewField(true, true, FileSize(100, 100000), MimeTypes([]string{"application/x-gzip"})),
		"missingupload": NewField(false, true, nil),
	}

	err, validator := New(fields)
	if err != nil {
		t.Errorf("Error making a new file validator type! %s", err.Error())
	}

	_, errors := validator.Validate(req.MultipartForm.File)

	if len(errors) > 0 {
		for field, val := range errors {
			for _, e := range val {
				t.Errorf("TestValidate(): %s: %s", field, e.Error())
			}
		}
	}
}
