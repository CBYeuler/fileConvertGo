package converters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type JSONToPDFConverter struct{}

func NewJSONToPDFConverter() *JSONToPDFConverter {
	return &JSONToPDFConverter{}
}

func (c *JSONToPDFConverter) Convert(jsonData []byte) ([]byte, error) {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(40, 30, 40)
	pdf.SetAutoPageBreak(true, 30)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	renderStructuredJSON(pdf, data, 0)
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func JsonToPDF(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}
	converter := NewJSONToPDFConverter()
	pdfBytes, err := converter.Convert(data)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, pdfBytes, 0644)
}

func renderStructuredJSON(pdf *gofpdf.Fpdf, data interface{}, indent int) {
	lineHeight := 7.0
	indentX := float64(40 + indent*8)

	switch v := data.(type) {
	case map[string]interface{}:
		// Sort keys for consistent output
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			pdf.SetFont("Arial", "B", 12)
			pdf.SetX(indentX)
			pdf.CellFormat(0, lineHeight, formatKey(key), "", 1, "", false, 0, "")
			renderStructuredJSON(pdf, v[key], indent+1)
			pdf.Ln(2)
		}

	case []interface{}:
		// Only show header for non-empty arrays
		if len(v) > 0 {
			pdf.SetFont("Arial", "B", 12)
			pdf.SetX(indentX)
			pdf.CellFormat(0, lineHeight, fmt.Sprintf("Items (%d):", len(v)), "", 1, "", false, 0, "")
			pdf.Ln(1)
		}

		for _, item := range v {
			pdf.SetFont("Arial", "", 11)
			pdf.SetX(indentX)
			pdf.CellFormat(6, lineHeight, "•", "", 0, "", false, 0, "")
			pdf.SetX(indentX + 6)
			renderStructuredJSON(pdf, item, indent+1)
		}

	default:
		pdf.SetFont("Arial", "", 11)
		pdf.SetX(indentX)
		pdf.MultiCell(0, lineHeight, fmt.Sprintf("%v", v), "", "", false)
	}
}

func formatKey(key string) string {
	key = strings.ReplaceAll(key, "_", " ")
	return cases.Title(language.English).String(key)
}
