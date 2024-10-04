package pgdb_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/vidur2/vectorMicroservices/pkg/pgdb"
	webservertypes "github.com/vidur2/vectorMicroservices/pkg/webserverTypes"
)

func TestAddTextbook(t *testing.T) {
	godotenv.Load("../../.env")
	pgdb.InitVectorDb()
	defer pgdb.Db.Close()
	webservertypes.StatusMap = make(map[string]webservertypes.StatusRes)
	err := pgdb.AddTextBooks("https://arxiv.org/pdf/2002.00097.pdf", "test", "", "")

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestQueryTextbook(t *testing.T) {
	godotenv.Load("../../.env")
	pgdb.InitVectorDb()
	webservertypes.StatusMap = make(map[string]webservertypes.StatusRes)
	id := uuid.NewString()
	out, err := pgdb.QueryCacheMatch[pgdb.TextbookEmbedding]([]string{"testing thing this is cool"}, id)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(out)
}

func TestAddJit(t *testing.T) {
	godotenv.Load("../../.env")
	pgdb.InitVectorDb()
	webservertypes.StatusMap = make(map[string]webservertypes.StatusRes)
	docs := make([]webservertypes.Document, 0)
	meta := make(map[string]interface{})
	meta["subject"] = "test"
	meta["problem"] = "test2"
	meta["images"] = []string{"test3"}

	docs = append(docs, webservertypes.Document{
		Content:  "thing is cool",
		Metadata: meta,
	})

	err := pgdb.Add[pgdb.JitEmbedding](uuid.NewString(), docs, pgdb.JitEmbedding{})

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestJitQuery(t *testing.T) {
	godotenv.Load("../../.env")
	pgdb.InitVectorDb()
	webservertypes.StatusMap = make(map[string]webservertypes.StatusRes)
	id := uuid.NewString()
	out, err := pgdb.QueryCacheMatch[pgdb.JitEmbedding]([]string{"testing thing this is cool"}, id)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	fmt.Println(out)
}

func TestInit(t *testing.T) {
	err := pgdb.InitVectorDb()
	defer pgdb.Db.Close()

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
