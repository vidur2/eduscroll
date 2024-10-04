package webservertypes

import "go.mongodb.org/mongo-driver/bson"

var StatusMap map[string]StatusRes
var JitQueryMap map[string]JitQueryStatusRes
var JitQueryMapDocDB map[string]JitQueryResDocDb

type VideoGenReq struct {
	Context  string `json:"textbook_context"`
	Subject  string `json:"subject"`
	Textbook string `json:"textbook"`
}

type JitQueryStatusRes struct {
	Status string                        `json:"status"`
	Body   map[string][]DistanceResponse `json:"body"`
}

type DistanceResponse struct {
	Response string      `json:"response"`
	Distance float32     `json:"distance"`
	Metadata interface{} `json:"metadata"`
}

type DefaultResponse struct {
	Res string `json:"res"`
	Url string `json:"url"`
}

type ContentCacheReq struct {
	Subject    string    `json:"subject"`
	Queries    []string  `json:"queries"`
	Problems   []Problem `json:"problems"`
	S3VideoUri []string  `json:"s3VideoUri"`
	Textbook   string    `json:"textbook"`
}

type Problem struct {
	Q string `json:"question"`
	A string `json:"answer"`
}

type ModelQueryReq struct {
	Queries  []string `json:"queries"`
	Subject  string   `json:"subject"`
	Textbook string   `json:"textbook"`
}

type ErrRes struct {
	Err string `json:"res"`
}

type StatusRes struct {
	Status string             `json:"status"`
	Info   VectorDbAddRequest `json:"info"`
}

type VectorDbAddRequest struct {
	DocUrl              string `json:"url"`
	DocSubject          string `json:"subject"`
	Ext                 string `json:"extension"`
	Textbook            string `json:"textbook"`
	AddToBothCollection bool   `json:"add_both"`
}

type EmbedResponse struct {
	Embeds [][]float32 `json:"embeddings"`
}

type EmbedRequest struct {
	Docs []string `json:"docs"`
}

type Document struct {
	Content  string
	Metadata map[string]interface{}
}

type JitQueryResDocDb struct {
	Body   map[string][]bson.M
	Status string
}

type ModelQueryReqRaw struct {
	Queries  []string    `json:"queries"`
	Vectors  [][]float32 `json:"vectors"`
	Subject  string      `json:"subject"`
	Textbook string      `json:"textbook"`
}
