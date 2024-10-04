package docdb

import (
	"context"

	embeddingfunction "github.com/vidur2/vectorMicroservices/pkg/embeddingFunction"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var EmbeddingFunc *embeddingfunction.EmbedFunction
var MongoClient *mongo.Client

func InitVectorDb(uri string) {
	if uri == "" {
		panic("Create a docdb vector db before continuing")
	}
	EmbeddingFunc = &embeddingfunction.EmbedFunction{}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	conn, err := mongo.Connect(context.TODO(), opts)
	MongoClient = conn
	if err != nil {
		panic(err)
	}
}
