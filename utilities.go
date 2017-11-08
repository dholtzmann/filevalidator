package filevalidator

import (
	"fmt"
	"strings"
)

/*
	Custom error type

	Str can be accessed to do translation, it has the raw data so Sprintf flags are included (%d,%s,etc.)

	Example:

	var e = FormError{"Integer must be between %d and %d", []interface{}{5, 10}}
	i18n(e.Str, e.Data) -> "Escriba Usted un valor entre 5 y 10"
*/
type FileError struct {
	Str  string
	Data []interface{}
}

func (e *FileError) Error() string {
	if len(e.Data) > 0 { // if 'Data' is not empty format the string
		return fmt.Sprintf(e.Str, e.Data...)
	}
	return e.Str
}

// appendError(err, errors...) do not forget the dots for slices
func appendError(err []FileError, arg ...*FileError) []FileError {
	for _, e := range arg {
		err = append(err, *e)
	}
	return err
}

// Is some value in the slice?
func inSlice(slice []string, val string) bool {
	for _, j := range slice {
		if j == val {
			return true
		}
	}
	return false
}

func getFileExtension(name string) string {
	if pos := strings.LastIndex(name, "."); pos != -1 {
		return strings.ToLower(name[pos+1:])
	}
	return ""
}
