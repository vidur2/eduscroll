package textextractor_test

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	textextractor "github.com/vidur2/vectorMicroservices/pkg/textExtractor"
)

func TestPdfExtractor(t *testing.T) {
	out, err := textextractor.PdfExtractor("https://www.people.vcu.edu/~rhammack/BookOfProof/BookOfProof.pdf")
	if err != nil {
		log.Fatal(err)
		t.Fatal()
	}
	assert.NotEqual(t, len(out), 0)
}

func TestExtractor(t *testing.T) {
	out, err := textextractor.ExtractText("https://arxiv.org/pdf/2002.00097.pdf", "", "")
	if err != nil {
		log.Fatal(err)
		t.Fatal()
	}

	assert.NotEqual(t, len(out), 0)
}
