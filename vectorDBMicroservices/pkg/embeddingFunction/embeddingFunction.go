package embeddingfunction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	webservertypes "github.com/vidur2/vectorMicroservices/pkg/webserverTypes"
)

type EmbedFunction struct {
}

func (e *EmbedFunction) CreateEmbedding(documents []string) ([][]float32, error) {

	client := http.Client{}
	body, err := json.Marshal(webservertypes.EmbedRequest{
		Docs: documents,
	})

	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", os.Getenv("EMBED_URL")+"/embed", bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	r.Header.Add("content-type", "application/json")

	out, err := client.Do(r)

	if err != nil {
		return nil, err
	}

	r.Body.Close()
	body, err = io.ReadAll(out.Body)

	if err != nil {
		return nil, err
	}

	var res webservertypes.EmbedResponse
	err = json.Unmarshal(body, &res)

	if err != nil {
		return nil, err
	}

	return res.Embeds, nil
}

func (e *EmbedFunction) CreateEmbeddingWithModel(documents []string, model string) ([][]float32, error) {
	client := http.Client{}

	body, err := json.Marshal(webservertypes.EmbedRequest{
		Docs: documents,
	})

	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", os.Getenv("EMBED_URL")+"/embed", bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}

	r.Header.Add("content-type", "application/json")

	out, err := client.Do(r)

	if err != nil {
		return nil, err
	}

	r.Body.Close()
	body, err = io.ReadAll(out.Body)

	if err != nil {
		return nil, err
	}

	var res webservertypes.EmbedResponse
	err = json.Unmarshal(body, &res)

	if err != nil {
		return nil, err
	}

	return res.Embeds, nil
}

func SplitDocuments(documents []string, docUrl string, subject string) []webservertypes.Document {
	out := make([]webservertypes.Document, 0)
	i := 0
	for _, doc := range documents {
		splitDocSpaces := strings.Split(doc, " ")
		if len(splitDocSpaces) < 50 {
			metadata := genMetadata(subject, docUrl, i)
			out = append(out, webservertypes.Document{
				Metadata: metadata,
				Content:  doc,
			})
		} else {
			for len(splitDocSpaces) != 0 {
				fmt.Println(len(splitDocSpaces))
				out = append(out, webservertypes.Document{
					Content:  strings.Join(splitDocSpaces[:min(len(splitDocSpaces), 50)], " "),
					Metadata: genMetadata(subject, docUrl, i),
				})
				splitDocSpaces = splitDocSpaces[min(len(splitDocSpaces), 50):]
				i++
			}
		}
		i++
	}

	return out
}

func genMetadata(subject string, docUrl string, i int) map[string]interface{} {
	metadata := make(map[string]interface{})
	metadata["subject"] = subject
	metadata["url"] = docUrl
	metadata["page"] = i
	return metadata
}
