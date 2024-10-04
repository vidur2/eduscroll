package chromadb

import (
	"os"

	chroma "github.com/amikos-tech/chroma-go"
	embeddingfunction "github.com/vidur2/vectorMicroservices/pkg/embeddingFunction"
)

var embeddingFunc *embeddingfunction.EmbedFunction
var ChromaClient *chroma.Client

func InitVectorDb(isTesting bool) {
	embeddingFunc = &embeddingfunction.EmbedFunction{}
	if !isTesting {
		ChromaClient = chroma.NewClient(os.Getenv("CHROMA_URI"))
	} else {
		ChromaClient = chroma.NewClient(os.Getenv("CHROMA_TEST_URI"))
	}
}
