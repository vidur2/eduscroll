package chromadb

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	chroma_go "github.com/amikos-tech/chroma-go"
	"github.com/valyala/fasthttp"

	embeddingfunction "github.com/vidur2/vectorMicroservices/pkg/embeddingFunction"
	textextractor "github.com/vidur2/vectorMicroservices/pkg/textExtractor"
	webservertypes "github.com/vidur2/vectorMicroservices/pkg/webserverTypes"
)

func ResetVectorDb() error {
	valid, err := ChromaClient.Reset()
	if !valid {
		return fmt.Errorf("could not reset vector db")
	}

	if err != nil {
		return err
	}

	return nil
}

func AddTextBooks(docUrl string, subject string, extension string, uuid string, textbook string, addTextbook bool) error {
	webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
		Status: "Pending",
		Info:   webservertypes.StatusMap[uuid].Info,
	}
	docs, err := textextractor.ExtractText(docUrl, uuid, extension)

	webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
		Status: "Finished Extraction",
		Info:   webservertypes.StatusMap[uuid].Info,
	}

	if err != nil {
		webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
			Status: fmt.Sprintf("Error: %v", err),
			Info:   webservertypes.StatusMap[uuid].Info,
		}
		return err
	}

	err = AddToDb(docs, subject, uuid, docUrl, textbook, addTextbook)

	if err != nil {
		webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
			Status: fmt.Sprintf("Error: %v", err),
			Info:   webservertypes.StatusMap[uuid].Info,
		}
		return err
	}
	return err
}

func MakeRequest(subject string, textbook string, doc string) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	fmt.Println(fmt.Sprintf("%v/video", os.Getenv("EMBED_URL")))
	req.SetRequestURI(fmt.Sprintf("%v/video", os.Getenv("EMBED_URL")))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	body := webservertypes.VideoGenReq{
		Context:  doc,
		Subject:  subject,
		Textbook: textbook,
	}
	out, err := json.Marshal(body)

	if err != nil {
		return err
	}

	req.SetBodyString(string(out))

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	fasthttp.Do(req, res)
	if res.StatusCode() != 200 {
		return fmt.Errorf("Request failed %v", string(res.Body()))
	}
	return nil
}

func AddToDb(documents []string, subject string, uuid string, docUrl string, textbook string, addTextbook bool) error {
	metadocs := embeddingfunction.SplitDocuments(documents, docUrl, subject)
	webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
		Status: "Started Uploading Video",
		Info:   webservertypes.StatusMap[uuid].Info,
	}
	for i, doc := range metadocs {
		// Call API route and return to JIT here
		if i == 28 || i == 40 || i == 39 || i == 25 || i == 43 {
			err := MakeRequest(subject, textbook, doc.Content)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	if addTextbook {
		if textbook != "" {
			return Add(uuid, metadocs, textbook)
		} else {
			return Add(uuid, metadocs, subject)
		}
	}

	return nil
}

func Add(uuid string, documents []webservertypes.Document, collection string) error {
	col, err := GetOrCreateCol(collection)
	webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
		Status: "Started Adding to DB",
		Info:   webservertypes.StatusMap[uuid].Info,
	}
	if err != nil {
		webservertypes.StatusMap[uuid] = webservertypes.StatusRes{
			Status: fmt.Sprintf("Error while creating collection %v", err.Error()),
		}
		return err
	}
	ids := make([]string, len(documents))
	metadatas := make([]map[string]interface{}, len(documents))
	content := make([]string, len(documents))

	for i := 0; i < len(documents); i++ {
		ids[i] = fmt.Sprintf("%v%v", uuid, i)
		metadatas[i] = documents[i].Metadata
		content[i] = documents[i].Content
	}

	_, err = col.Add(nil, metadatas, content, ids)

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

func RawCompare(queryTexts []string, queryVectors [][]float32, uuid string, subject string, textbook string) (map[string][]webservertypes.DistanceResponse, error) {
	col, err := getNilChroma(), getNilError()
	if textbook != "-jit" {
		col, err = GetColWithEmbedMap(textbook, queryTexts, queryVectors)
	} else {
		col, err = GetColWithEmbedMap(subject, queryTexts, queryVectors)
	}
	if err != nil {
		webservertypes.JitQueryMap[uuid] = webservertypes.JitQueryStatusRes{
			Status: fmt.Sprintf("Error: %v", err.Error()),
		}
		return nil, err
	}

	count, err := strconv.Atoi(os.Getenv("MAX_SUBQUERY_RESULTS"))

	if err != nil {
		webservertypes.JitQueryMap[uuid] = webservertypes.JitQueryStatusRes{
			Status: fmt.Sprintf("Error: %v", err.Error()),
		}
		return nil, err
	}
	var res *chroma_go.QueryResults

	if subject != "" {
		subjectMap := make(map[string]interface{})
		subjectMap["subject"] = subject
		res, err = col.Query(queryTexts, int32(count), subjectMap, nil, nil)
	} else {
		res, err = col.Query(queryTexts, int32(count), nil, nil, nil)
	}

	if err != nil {
		webservertypes.JitQueryMap[uuid] = webservertypes.JitQueryStatusRes{
			Status: fmt.Sprintf("Error: %v", err.Error()),
		}
		return nil, err
	}

	out := make(map[string][]webservertypes.DistanceResponse)
	for i := 0; i < len(queryTexts); i++ {
		distRes := make([]webservertypes.DistanceResponse, len(res.Documents[i]))
		for j := 0; j < len(res.Documents[i]); j++ {
			outMeta := make(map[string]interface{})
			for k, v := range res.Metadatas[i][j] {
				outMeta[k] = string(v.([]byte))
			}
			distRes[j] = webservertypes.DistanceResponse{
				Response: res.Documents[i][j],
				Distance: res.Distances[i][j],
				Metadata: outMeta,
			}
		}
		out[queryTexts[i]] = distRes
	}

	webservertypes.JitQueryMap[uuid] = webservertypes.JitQueryStatusRes{
		Status: "Finished",
		Body:   out,
	}

	return out, nil
}

func CompareToDb(queryTexts []string, uuid string, subject string, textbook string) (map[string][]webservertypes.DistanceResponse, error) {
	col := getNilChroma()
	err := getNilError()

	if textbook != "-jit" {
		col, err = GetOrCreateCol(textbook)
	} else {
		col, err = GetOrCreateCol(subject)
	}
	if err != nil {
		webservertypes.JitQueryMap[uuid] = webservertypes.JitQueryStatusRes{
			Status: fmt.Sprintf("Error: %v", err.Error()),
		}
		return nil, err
	}

	count, err := strconv.Atoi(os.Getenv("MAX_SUBQUERY_RESULTS"))

	if err != nil {
		webservertypes.JitQueryMap[uuid] = webservertypes.JitQueryStatusRes{
			Status: fmt.Sprintf("Error: %v", err.Error()),
		}
		return nil, err
	}
	var res *chroma_go.QueryResults

	res, err = col.Query(queryTexts, int32(count), nil, nil, nil)

	if err != nil {
		webservertypes.JitQueryMap[uuid] = webservertypes.JitQueryStatusRes{
			Status: fmt.Sprintf("Error: %v", err.Error()),
		}
		return nil, err
	}

	out := make(map[string][]webservertypes.DistanceResponse)
	for i := 0; i < len(queryTexts); i++ {
		distRes := make([]webservertypes.DistanceResponse, len(res.Documents[i]))
		for j := 0; j < len(res.Documents[i]); j++ {
			outMeta := make(map[string]interface{})
			for k, v := range res.Metadatas[i][j] {
				outMeta[k] = string(v.([]byte))
			}
			distRes[j] = webservertypes.DistanceResponse{
				Response: res.Documents[i][j],
				Distance: res.Distances[i][j],
				Metadata: outMeta,
			}
		}
		out[queryTexts[i]] = distRes
	}

	webservertypes.JitQueryMap[uuid] = webservertypes.JitQueryStatusRes{
		Status: "Finished",
		Body:   out,
	}

	return out, nil
}

func getNilError() error {
	return nil
}

func getNilChroma() *chroma_go.Collection {
	return nil
}

func getFullString(col *chroma_go.Collection, metadata map[string]interface{}) (string, error) {
	where := make(map[string]interface{}, 1)
	where["id"] = metadata["id"]
	count, err := strconv.Atoi(os.Getenv("MAX_SUBQUERY_RESULTS"))
	if err != nil {
		return "", err
	}
	out, err := col.Query(nil, int32(count), where, nil, nil)

	if err != nil {
		return "", err
	}

	outMap := make([]string, len(out.Documents))
	for i, str := range out.Documents[0] {
		ord := out.Metadatas[0][i]["ord"].(int)
		if err != nil {
			return "", err
		}
		outMap[ord] = str
	}

	outStr := ""
	for _, str := range outMap {
		outStr += str
	}
	return outStr, nil
}

func GetOrCreateCol(colName string) (*chroma_go.Collection, error) {
	metadata := make(map[string]interface{})
	col, err := ChromaClient.GetCollection(colName, embeddingFunc)
	if err != nil {
		col, err = ChromaClient.CreateCollection(colName, metadata, true, embeddingFunc, chroma_go.L2)
	}
	if err != nil {
		return nil, err
	}
	return col, nil
}

func GetColWithEmbedMap(colName string, documents []string, embed [][]float32) (*chroma_go.Collection, error) {
	metadata := make(map[string]interface{})
	embedMap := make(map[string][]float32)
	for i, doc := range documents {
		embedMap[doc] = embed[i]
	}
	embedfunction := embeddingfunction.PreEmbedFunction{
		EmbedMap: embedMap,
	}
	col, err := ChromaClient.GetCollection(colName, &embedfunction)

	if err != nil {
		col, err = ChromaClient.CreateCollection(colName, metadata, true, embeddingFunc, chroma_go.L2)
	}
	if err != nil {
		return nil, err
	}
	return col, nil
}
