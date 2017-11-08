package filevalidator

import (
	"testing"
)

func Test_FileSize(t *testing.T) {
	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}) // 17668 bytes

	for _, fieldSlice := range req.MultipartForm.File { // map[string][]*FileHeader
		for _, header := range fieldSlice {
			file, err := header.Open()
			if err != nil {
				t.Fatalf("Failed to open *multipart.FileHeader: %s", err)
			}
			defer file.Close()

			minRule := FileSize(100000, 1000000)
			if e, _ := minRule.Validate(file, header.Filename, make(map[string]string)); e == nil {
				t.Errorf("Testing minimum size - FileSize(): Should return an error!")
			}

			maxRule := FileSize(2, 1000)
			if e, _ := maxRule.Validate(file, header.Filename, make(map[string]string)); e == nil {
				t.Errorf("Testing maximum size - FileSize(): Should return an error!")
			}

			rightRule := FileSize(17600, 18000)
			if e, _ := rightRule.Validate(file, header.Filename, make(map[string]string)); e != nil {
				t.Errorf("Testing correct size - FileSize(): Should not return an error!")
			}
		}
	}
}

func Test_MimeTypes(t *testing.T) {
	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}, &testFile{blueJPG, "imageupload", "blue.jpg"}) // image/png & image/jpeg

	for _, fieldSlice := range req.MultipartForm.File { // map[string][]*FileHeader
		for _, header := range fieldSlice {
			file, err := header.Open()
			if err != nil {
				t.Fatalf("Failed to open *multipart.FileHeader: %s", err)
			}
			defer file.Close()

			badRule := MimeTypes([]string{"application/zip", "text/html"})
			if e, _ := badRule.Validate(file, header.Filename, make(map[string]string)); e == nil {
				t.Errorf("Testing bad mime type - MimeTypes(): Should return an error!")
			}

			goodRule := MimeTypes([]string{"image/png", "image/jpeg"})
			if e, _ := goodRule.Validate(file, header.Filename, make(map[string]string)); e != nil {
				t.Errorf("Testing good mime type - MimeTypes(): Should not return an error!")
			}
		}
	}
}

func Test_MinPixelSize(t *testing.T) {
	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}) // 250 x 340 pixels

	for _, fieldSlice := range req.MultipartForm.File { // map[string][]*FileHeader
		for _, header := range fieldSlice {
			file, err := header.Open()
			if err != nil {
				t.Fatalf("Failed to open *multipart.FileHeader: %s", err)
			}
			defer file.Close()

			badRule := MinPixelSize(400, 400)
			if e, _ := badRule.Validate(file, header.Filename, make(map[string]string)); e == nil {
				t.Errorf("Testing bad min width and height - MinPixelSize(): Should return an error!")
			}

			goodRule := MinPixelSize(200, 200)
			if e, _ := goodRule.Validate(file, header.Filename, make(map[string]string)); e != nil {
				t.Errorf("Testing good min width and height - MinPixelSize(): Should not return an error!")
			}
		}
	}
}

func Test_MaxPixelSize(t *testing.T) {
	req := setupRequestMultipartForm(&testFile{gopherPNG, "imageupload", "gopher.png"}) // 250 x 340 pixels

	for _, fieldSlice := range req.MultipartForm.File { // map[string][]*FileHeader
		for _, header := range fieldSlice {
			file, err := header.Open()
			if err != nil {
				t.Fatalf("Failed to open *multipart.FileHeader: %s", err)
			}
			defer file.Close()

			badRule := MaxPixelSize(100, 100)
			if e, _ := badRule.Validate(file, header.Filename, make(map[string]string)); e == nil {
				t.Errorf("Testing bad max width and height - MaxPixelSize(): Should return an error!")
			}

			goodRule := MaxPixelSize(400, 400)
			if e, _ := goodRule.Validate(file, header.Filename, make(map[string]string)); e != nil {
				t.Errorf("Testing good max width and height - MaxPixelSize(): Should not return an error!")
			}
		}
	}
}
