package dto

import "io"

type ConvParams struct {
	Writer  io.Writer
	Columns []string
	Rows    [][]any
	Sep     byte
}
