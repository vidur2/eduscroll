package chromadb_test

import (
	"fmt"
	"testing"

	"github.com/joho/godotenv"

	"github.com/vidur2/vectorMicroservices/pkg/chromadb"
	webservertypes "github.com/vidur2/vectorMicroservices/pkg/webserverTypes"
)

func TestAdd(t *testing.T) {
	godotenv.Load("../../.env")
	chromadb.InitVectorDb(false)
	webservertypes.StatusMap = make(map[string]webservertypes.StatusRes, 0)

	err := chromadb.AddTextBooks("https://arxiv.org/pdf/2002.00097.pdf", "test", "", "", "", true)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestList(t *testing.T) {
	godotenv.Load("../../.env")
	chromadb.InitVectorDb(false)
	cols, err := chromadb.ChromaClient.ListCollections()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(cols)
}

func TestListCol(t *testing.T) {
	godotenv.Load("../../.env")
	chromadb.InitVectorDb(false)
	cols, err := chromadb.GetOrCreateCol("ALT_COL_NAME")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	out, err := cols.Get(nil, nil, nil)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(out.CollectionData.Documents)
}

func TestHeartbeat(t *testing.T) {
	godotenv.Load("../../.env")
	chromadb.InitVectorDb(false)
	_, err := chromadb.ChromaClient.Heartbeat()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestVideoGen(t *testing.T) {
	chromadb.MakeRequest("Test", "test", "test")
}
