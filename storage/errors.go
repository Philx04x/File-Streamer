package storage

import "fmt"

var ErrFileWrite error = fmt.Errorf("failed to write file data to file")
var ErrFileOpen error = fmt.Errorf("failed to open file")
var ErrFileReader error = fmt.Errorf("failed to read file data")
