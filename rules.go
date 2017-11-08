package filevalidator

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif" // initialize the package
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"net/http"
	"strings"
)

type fileSize struct { // size in bytes
	min uint64
	max uint64
}

func FileSize(min, max uint64) Rule {
	return &fileSize{min, max}
}

/*
	How big can the file be? Minimum size and maximum size.
*/
func (f *fileSize) Validate(file multipart.File, originalName string, errorMsg map[string]string) (error, []interface{}) {
	var buff bytes.Buffer
	var size int64
	var err error

	if _, err = file.Seek(0, 0); err != nil { // set position back to start of the file
		return errors.New(errorMsg["file_failed"]), []interface{}{originalName}
	}

	if size, err = buff.ReadFrom(file); err != nil { // read the whole file, size is in bytes
		return errors.New(errorMsg["file_failed"]), []interface{}{originalName}
	}

	if uint64(size) < f.min {
		return errors.New(errorMsg["file_too_small"]), []interface{}{originalName}
	}

	if uint64(size) > f.max {
		return errors.New(errorMsg["file_too_large"]), []interface{}{originalName}
	}

	return nil, nil
}

// -----------------------

type mimeTypes struct {
	list []string
}

func MimeTypes(list []string) Rule { // use []string{"*"} to allow all types
	if list == nil {
		list = []string{}
	}

	return &mimeTypes{list}
}

/*
	What mimetypes are allowed: Ex: []string{"image/jpeg", "image/gif", "image/png"}
*/
func (m *mimeTypes) Validate(file multipart.File, originalName string, errorMsg map[string]string) (error, []interface{}) {
	buff := make([]byte, 512)

	if _, err := file.Seek(0, 0); err != nil { // set position back to start of the file
		return errors.New(errorMsg["file_failed"]), []interface{}{originalName}
	}

	if n, err := file.Read(buff); n <= 0 && err != nil { // Copy the first 512 bytes into a buffer to get the mimetype
		return errors.New(errorMsg["file_failed"]), []interface{}{originalName}
	}

	// http.DetectContentType() `always returns a valid MIME type: if it cannot determine a more specific one, it returns "application/octet-stream"`
	mimeType := strings.TrimSpace(strings.ToLower(http.DetectContentType(buff)))

	if len(m.list) == 1 && m.list[0] == "*" { // allow all types
		return nil, nil
	} else if !inSlice(m.list, mimeType) {
		return errors.New(errorMsg["file_type_bad"]), []interface{}{originalName, mimeType}
	}

	return nil, nil
}

// -----------------------

type minPixelSize struct {
	width  uint32
	height uint32
}

func MinPixelSize(width, height uint32) Rule {
	return &minPixelSize{width, height}
}

/*
	How small can an image be in pixels? Minimum width and minimum height. This only works for images.
*/
func (m *minPixelSize) Validate(file multipart.File, originalName string, errorMsg map[string]string) (error, []interface{}) {

	// Some simple verification that the file is an image, this is not fullproof. The user can lie by naming any file "X.jpg".
	ext := getFileExtension(originalName)
	if !(ext == "jpeg" || ext == "jpg" || ext == "gif" || ext == "png") {
		return errors.New(errorMsg["image_only"]), []interface{}{originalName}
	}

	if _, err := file.Seek(0, 0); err != nil { // set position back to start of the file
		return errors.New(errorMsg["file_failed"]), []interface{}{originalName}
	}

	// The upload might be incomplete or maliciously crafted, spoofed mimetime/file extension, different body
	if im, _, err := image.Decode(file); err != nil {
		return errors.New(errorMsg["image_only"]), []interface{}{originalName}
	} else {
		pt := im.Bounds().Size()
		if uint32(pt.X) < m.width {
			return errors.New(errorMsg["image_size_min"]), []interface{}{m.width, m.height, originalName}
		}
		if uint32(pt.Y) < m.height {
			return errors.New(errorMsg["image_size_min"]), []interface{}{m.width, m.height, originalName}
		}
	}

	return nil, nil
}

// -----------------------

type maxPixelSize struct {
	width  uint32
	height uint32
}

func MaxPixelSize(width, height uint32) Rule {
	return &maxPixelSize{width, height}
}

/*
	How large can an image be in pixels? Minimum width and maximum height. This only works for images.
*/
func (m *maxPixelSize) Validate(file multipart.File, originalName string, errorMsg map[string]string) (error, []interface{}) {

	// Some simple verification that the file is an image, this is not fullproof. The user can lie by naming any file "X.jpg".
	ext := getFileExtension(originalName)
	if !(ext == "jpeg" || ext == "jpg" || ext == "gif" || ext == "png") {
		return errors.New(errorMsg["image_only"]), []interface{}{originalName}
	}

	if _, err := file.Seek(0, 0); err != nil { // set position back to start of the file
		return errors.New(errorMsg["file_failed"]), []interface{}{originalName}
	}

	// The upload might be incomplete or maliciously crafted, spoofed mimetime/file extension, different body
	if im, _, err := image.Decode(file); err != nil {
		return errors.New(errorMsg["image_only"]), []interface{}{originalName}
	} else {
		pt := im.Bounds().Size()
		if uint32(pt.X) > m.width {
			return errors.New(errorMsg["image_size_max"]), []interface{}{m.width, m.height, originalName}
		}
		if uint32(pt.Y) > m.height {
			return errors.New(errorMsg["image_size_max"]), []interface{}{m.width, m.height, originalName}
		}
	}

	return nil, nil
}
