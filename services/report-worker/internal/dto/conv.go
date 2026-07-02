package dto

import "io"

type ConvParams struct {
	Writer io.Writer
	Colums []string
	Rows   [][]any
	Sep    byte
}
