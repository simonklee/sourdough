package main

import (
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/table"
)

func renderTable(w io.Writer, format OutputFormat, tw table.Writer) error {
	var data string
	switch format {
	case OutputFormatHTML:
		data = tw.RenderHTML()
	case OutputFormatMarkdown:
		data = tw.RenderMarkdown()
	default:
		data = tw.Render()
	}
	_, err := fmt.Fprintf(w, "%s\n", data)
	return err
}
