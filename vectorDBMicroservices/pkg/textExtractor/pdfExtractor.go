package textextractor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	textrank "github.com/DavidBelicza/TextRank"
	"github.com/unidoc/unipdf/v3/extractor"
	pdf "github.com/unidoc/unipdf/v3/model"
)

func PdfExtractor(url string) ([]string, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	buf, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	pdfReader, err := pdf.NewPdfReader(bytes.NewReader(buf))

	if err != nil {
		return nil, err
	}

	enc, _ := pdfReader.IsEncrypted()

	if enc {
		return nil, fmt.Errorf("pdf is encrypted")
	}

	numPages, _ := pdfReader.GetNumPages()
	out := make([]string, 0)
	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			log.Fatal(err)
		}

		ex, err := extractor.New(page)

		if err != nil {
			return nil, err
		}

		text, err := ex.ExtractText()
		if err != nil {
			log.Fatal(err)
		}
		out = append(out, text)
	}
	if len(out) > 45 {
		return getRankedInfo(out, 45), nil
	} else {
		return out, nil
	}
}

func getRankedInfo(docs []string, videoAmt int) []string {
	windowSize := len(docs) / videoAmt
	outVal := make([]string, 0)
	for i := 0; i < len(docs)-windowSize; i += windowSize {
		out := ""
		for j := i; j < i+windowSize; j++ {
			out += strings.ReplaceAll(docs[j], ".", ";;") + "."
		}
		tr := textrank.NewTextRank()
		// Default Rule for parsing.
		rule := textrank.NewDefaultRule()
		// Default Language for filtering stop words.
		language := textrank.NewDefaultLanguage()
		// Default algorithm for ranking text.
		algorithmDef := textrank.NewDefaultAlgorithm()

		tr.Populate(out, language, rule)
		tr.Ranking(algorithmDef)

		outVal = append(outVal, strings.ReplaceAll(textrank.FindSentencesByRelationWeight(tr, 1)[0].Value, ";;", "."))
	}

	return outVal
}
