package textextractor

import (
	"fmt"
	"strings"

	webservertypes "github.com/vidur2/vectorMicroservices/pkg/webserverTypes"
)

func ExtractText(url string, uuid string, extension string) ([]string, error) {
	webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
		Status: "Started Parsing",
		Info:   webservertypes.StatusMap[uuid].Info,
	}
	out := strings.Split(url, ".")
	if len(out) == 0 {
		return nil, fmt.Errorf("bad url")
	}
	if extension == "" {
		extension = out[len(out)-1]
	}
	switch extension {
	case "pdf":
		webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
			Status: "Started Extraction",
			Info:   webservertypes.StatusMap[uuid].Info,
		}
		return PdfExtractor(url)
	default:
		return nil, fmt.Errorf("file format not specified")
	}
}
