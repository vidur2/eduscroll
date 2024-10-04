package docdb

import (
	"context"
	"fmt"
	"os"
	"strconv"

	embeddingfunction "github.com/vidur2/vectorMicroservices/pkg/embeddingFunction"
	textextractor "github.com/vidur2/vectorMicroservices/pkg/textExtractor"
	webservertypes "github.com/vidur2/vectorMicroservices/pkg/webserverTypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	embeds, err := EmbeddingFunc.CreateEmbedding(documents)
	if err != nil {
		return err
	}
	return Add(uuid, metadocs, "DEFAULT_COLLECTION_NAME", embeds)
}

func Add(uuid string, documents []webservertypes.Document, collection string, embeds [][]float32) error {
	col := GetOrCreateCol(collection)
	webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
		Status: "Started",
		Info:   webservertypes.StatusMap[uuid].Info,
	}

	bsonFields := make([]interface{}, len(documents))

	for i, doc := range documents {
		bsonFields[i] = bson.D{
			{Key: "subject", Value: doc.Metadata["subject"]},
			{Key: "url", Value: doc.Metadata["url"]},
			{Key: "page", Value: doc.Metadata["page"]},
			{Key: "content", Value: doc.Content},
			{Key: "vectorEmbedding", Value: embeds[i]},
		}
	}

	_, err := col.InsertMany(context.TODO(), bsonFields)

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

func CompareToDb(queryTexts []string, uuid string, subject string, textbook string) (map[string][]bson.M, error) {
	col := GetOrCreateCol(subject)

	if textbook != "-jit" {
		col = GetOrCreateCol(textbook)
	}

	count, err := strconv.Atoi(os.Getenv("MAX_SUBQUERY_RESULTS"))

	if err != nil {
		webservertypes.JitQueryMapDocDB[uuid] = webservertypes.JitQueryResDocDb{
			Status: fmt.Sprintf("Error: %v", err.Error()),
		}
		return nil, err
	}

	val, err := vectorSearchNoSubj(col, queryTexts, count)
	if err != nil {
		return nil, err
	}
	webservertypes.JitQueryMapDocDB[uuid] = webservertypes.JitQueryResDocDb{
		Status: "Finished",
		Body:   val,
	}
	return val, nil
}

func vectorSearchNoSubj(col *mongo.Collection, docs []string, count int) (map[string][]bson.M, error) {
	out := make(map[string][]bson.M)
	embeds, err := EmbeddingFunc.CreateEmbedding(docs)
	if err != nil {
		return nil, err
	}
	for i, embed := range embeds {
		res, _ := col.Aggregate(context.TODO(), bson.D{
			{Key: "$search", Value: bson.D{
				{Key: "vectorSearch", Value: bson.D{
					{Key: "vector", Value: embed},
					{Key: "path", Value: "vectorEmbedding"},
					{Key: "similarity", Value: "euclidian"},
					{Key: "k", Value: count},
					{Key: "probes", Value: 1},
				}},
			}},
		})
		var results []bson.M
		res.All(context.TODO(), &results)
		out[docs[i]] = results
	}

	return out, nil
}

func vectorSearchSubj(col *mongo.Collection, docs []string, subject string, count int) (map[string][]bson.M, error) {
	out := make(map[string][]bson.M)
	embeds, err := EmbeddingFunc.CreateEmbedding(docs)
	if err != nil {
		return nil, err
	}
	for i, embed := range embeds {
		res, _ := col.Aggregate(context.TODO(), bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "subject", Value: subject},
			}},
			{Key: "$search", Value: bson.D{
				{Key: "vectorSearch", Value: bson.D{
					{Key: "vector", Value: embed},
					{Key: "path", Value: "vectorEmbedding"},
					{Key: "similarity", Value: "euclidian"},
					{Key: "k", Value: count},
					{Key: "probes", Value: 1},
				}},
			}},
		})
		var results []bson.M
		res.All(context.TODO(), &results)
		out[docs[i]] = results
	}

	return out, nil
}

func GetOrCreateCol(collection string) *mongo.Collection {
	_ = MongoClient.Database("eduscroll").CreateCollection(context.TODO(), os.Getenv(collection))
	col := MongoClient.Database("eduscroll").Collection(os.Getenv(collection))

	return col
}
