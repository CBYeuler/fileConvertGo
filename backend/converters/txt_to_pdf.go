package converters

import (
	"bufio"
	"os"

	"github.com/jung-kurt/gofpdf"
)

func TxtToPdf(inputPath string, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(40, 40, 40)
	pdf.AddPage()

	pdf.SetFont("Arial", "", 12)

	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	lineHeight := 10.0
	pageHeight := pdf.GetPageSizeStr()["h"]
	bottomMargin := pdf.GetMargins().Bottom

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if pdf.GetY()+lineHeight > pageHeight-bottomMargin {
			pdf.AddPage()
		}
		pdf.CellFormat(0, lineHeight, scanner.Text(), "", 1, "", false, 0, "")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return pdf.OutputFileAndClose(outputPath)
}
