package services

//there is import issues here
import (
	"backend/converters"
	"errors"
	"strings"

	"github.com/CBYeuler/fileConvertGo/backend/converters"
)

func ConvertFile(inputPath string, outputPath string, targetFormat string) error {
	LowerFormat := strings.ToLower(targetFormat)
	switch LowerFormat {
	case "pdf":
		return converters.TxtToPdf()
	case "json":
		return converters.TxtToJson()
	default:
		return errors.New("unsupported target format")
	}
}
