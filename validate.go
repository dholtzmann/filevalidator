package filevalidator

import (
	"errors"
	"mime/multipart"
)

/*
	TODO:
	-change errorMsgs map[string]string to map[string]error and add lines in test files that get the default errorMsgs map?
*/

// errors
var ErrNilArguments = errors.New("Arguments must be non-nil!")

type FileValidator struct {
	fields    map[string]*Field
	errorMsgs map[string]string
}

type Field struct {
	required   bool
	singleFile bool
	rules      []Rule
}

type Rule interface {
	Validate(multipart.File, string, map[string]string) (error, []interface{})
}

// Setup all the form validation rules
func NewField(required, singleFile bool, rules ...Rule) *Field {
	var ruleSlice []Rule
	for _, r := range rules {
		ruleSlice = append(ruleSlice, r)
	}

	return &Field{required, singleFile, ruleSlice}
}

// custom error messages option
func (f *FileValidator) SetErrors(msgs map[string]string) error {
	if msgs == nil {
		return ErrNilArguments
	}

	f.errorMsgs = msgs
	return nil
}

func New(fields map[string]*Field) (error, *FileValidator) {

	if fields == nil {
		return ErrNilArguments, nil
	}

	// default error messages for invalid files
	// extra error messages (those that are not used in file validation functions) are here for convenience, grouped for translation.
	var errors = map[string]string{
		"file_failed":    "This file failed to upload. Please retry. (%s)",
		"file_single":    "Only one file upload allowed.",
		"file_too_large": "This file is too large. (%s)",
		"file_too_small": "This file is too small. (%s)",
		"file_type_bad":  "This type of file is not allowed. (%s) [%s]",
		"image_only":     "This file is not an image. (%s)",
		"image_size_min": "This image must be at least %d pixels wide and %d pixels high. (%s)",
		"image_size_max": "This image cannot be more than %d pixels wide and %d pixels high. (%s)",
		"required":       "Please select a file.", // old: "This field is required.",
	}

	return nil, &FileValidator{fields, errors}
}

/*
	Important: This package is case-sensitive! That means a field named "email" is different than "Email"

	Loop through the rules and validate each entry
	returns (bool, if the files are valid, map for error messages)
*/
func (f *FileValidator) Validate(files map[string][]*multipart.FileHeader) (bool, map[string][]FileError) {
	allErrors := make(map[string][]FileError)

	if files == nil {
		return false, nil
	}

	for name, field := range f.fields { // map[string]Field
		var errors []FileError

		headerSlice, ok := files[name]             // this will be an empty slice if 'name' does not exist in the map
		if ok == false && field.required == true { // is the field required?
			errors = appendError(errors, &FileError{f.errorMsgs["required"], nil})
		}

		if field.singleFile == true && len(headerSlice) > 1 {
			errors = appendError(errors, &FileError{f.errorMsgs["file_single"], nil})
			continue // skip to next Field
		}

		for _, header := range headerSlice {
			var originalName = header.Filename
			file, err := header.Open()
			if err != nil {
				errors = appendError(errors, &FileError{f.errorMsgs["file_failed"], []interface{}{originalName}})
				continue // skip to next *multipart.Header
			}
			defer file.Close()

			for _, rule := range field.rules {
				if err, data := rule.Validate(file, originalName, f.errorMsgs); err != nil {
					errors = appendError(errors, &FileError{err.Error(), data}) // format errors for translation (string separate from extra data)
				}
			}
		}

		allErrors[name] = errors // set the errors for the field entry
	}

	if len(allErrors) == 0 {
		return true, allErrors
	}

	return false, allErrors
}
