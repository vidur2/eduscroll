package pgdb

import (
	"fmt"

	"github.com/pgvector/pgvector-go"
	embeddingfunction "github.com/vidur2/vectorMicroservices/pkg/embeddingFunction"
	textextractor "github.com/vidur2/vectorMicroservices/pkg/textExtractor"
	webservertypes "github.com/vidur2/vectorMicroservices/pkg/webserverTypes"
)

func AddTextBooks(docUrl string, subject string, extension string, uuid string) error {
	webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
		Status: "Pending",
		Info:   webservertypes.StatusMap[uuid].Info,
	}
	docs, err := textextractor.ExtractText(docUrl, uuid, extension)

	if err != nil {
		webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
			Status: fmt.Sprintf("Error: %v", err),
			Info:   webservertypes.StatusMap[uuid].Info,
		}
		return err
	}

	err = AddToDb(docs, subject, uuid, docUrl)

	if err != nil {
		webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
			Status: fmt.Sprintf("Error: %v", err),
			Info:   webservertypes.StatusMap[uuid].Info,
		}
		return err
	}
	return err
}

func AddToDb(documents []string, subject string, uuid string, docUrl string) error {
	metadocs := embeddingfunction.SplitDocuments(documents, docUrl, subject)
	return Add[TextbookEmbedding](uuid, metadocs, TextbookEmbedding{})
}

func Add[T TextbookEmbedding | JitEmbedding](uuid string, metadocs []webservertypes.Document, reflectionType T) error {
	webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
		Status: "Started",
		Info:   webservertypes.StatusMap[uuid].Info,
	}

	documents := make([]string, len(metadocs))
	for i, metadoc := range metadocs {
		documents[i] = metadoc.Content
	}

	function := embeddingfunction.EmbedFunction{}
	embeds, err := function.CreateEmbedding(documents)

	if err != nil {
		webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
			Status: fmt.Sprintf("Error: %v", err),
			Info:   webservertypes.StatusMap[uuid].Info,
		}
		return err
	}

	switch any(reflectionType).(type) {
	case TextbookEmbedding:
		items := make([]TextbookEmbedding, len(metadocs))
		for i := 0; i < len(metadocs); i++ {
			items[i] = TextbookEmbedding{
				Embedding: pgvector.NewVector(embeds[i]),
				Text:      documents[i],
				Subject:   metadocs[i].Metadata["subject"].(string),
				Page:      metadocs[i].Metadata["page"].(int),
			}
		}

		_, err = Db.Model(&items).Insert()
		if err != nil {
			return err
		}
	case JitEmbedding:
		items := make([]JitEmbedding, len(metadocs))
		for i := 0; i < len(metadocs); i++ {
			items[i] = JitEmbedding{
				Embedding: pgvector.NewVector(embeds[i]),
				Text:      documents[i],
				Subject:   metadocs[i].Metadata["subject"].(string),
				Problem:   metadocs[i].Metadata["problem"].(string),
				Images:    metadocs[i].Metadata["images"].([]string),
			}
		}
		_, err := Db.Model(&items).Insert()
		if err != nil {
			return err
		}
	default:
		err = fmt.Errorf("invalid adding type")
		webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
			Status: fmt.Sprintf("Error: %v", err),
			Info:   webservertypes.StatusMap[uuid].Info,
		}
		return err
	}

	if err != nil {
		webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
			Status: fmt.Sprintf("Error while adding %v", err.Error()),
		}
		return err
	}
	webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
		Status: "Finished",
		Info:   webservertypes.StatusMap[uuid].Info,
	}
	return nil
}

func QueryCacheMatch[T JitEmbedding | TextbookEmbedding](docs []string, uuid string) (map[string][]T, error) {
	embeddingFunc := embeddingfunction.EmbedFunction{}
	embeds, err := embeddingFunc.CreateEmbedding(docs)

	if err != nil {
		webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
			Status: fmt.Sprintf("Error while adding %v", err.Error()),
		}
		return nil, err
	}

	outMap := make(map[string][]T)
	for i, embed := range embeds {
		var typedEmbedding []T
		err = Db.Model(&typedEmbedding).OrderExpr("CAST(CAST(embedding AS text) AS vector(384)) <-> CAST(? AS vector(384))", pgvector.NewVector(embed)).ExcludeColumn("embedding").Limit(5).Select()
		if err != nil {
			webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
				Status: fmt.Sprintf("Error while adding %v", err.Error()),
			}
			return nil, err
		}
		outMap[docs[i]] = typedEmbedding
	}
	fmt.Println(outMap)
	return outMap, nil
}
