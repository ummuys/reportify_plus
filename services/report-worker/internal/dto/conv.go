package dto

import "os"

type ConvParams struct {
	Colums []string
	Rows   [][]any
	File   *os.File
	Sep    byte
}
