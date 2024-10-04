package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
	"github.com/vidur2/vectorMicroservices/pkg/chromadb"
	webservertypes "github.com/vidur2/vectorMicroservices/pkg/webserverTypes"
)

func handler(ctx *fasthttp.RequestCtx) {
	splitPath := strings.Split(string(ctx.Path()), "/")
	if splitPath[1] == "status" {
		switch splitPath[2] {
		case "get":
			status := webservertypes.JitQueryMap[splitPath[3]]
			res := status
			if status.Status == "Finished" || strings.HasPrefix(status.Status, "Error") {
				delete(webservertypes.JitQueryMap, splitPath[3])
			}
			body, _ := json.Marshal(res)
			ctx.Response.AppendBody(body)
		case "add":
			status := webservertypes.StatusMap[splitPath[3]]
			res := status
			if status.Status == "Finished" || strings.HasPrefix(status.Status, "Error") {
				delete(webservertypes.StatusMap, splitPath[3])
			}
			body, _ := json.Marshal(res)
			ctx.Response.AppendBody(body)
		default:
			ctx.Response.SetStatusCode(404)
		}
	} else {
		switch string(ctx.Path()) {
		case "/add":
			req := ctx.Request.Body()
			var body webservertypes.ContentCacheReq
			err := json.Unmarshal(req, &body)

			body.Textbook = fmt.Sprintf("%v-jit", body.Textbook)
			body.Subject = fmt.Sprintf("%v-jit", body.Subject)

			if err != nil {
				out, _ := json.Marshal(webservertypes.ErrRes{
					Err: err.Error(),
				})
				ctx.Response.SetStatusCode(400)
				ctx.Response.AppendBody(out)
			}
			id, _ := uuid.NewUUID()
			documents := make([]webservertypes.Document, 0)
			for i, query := range body.Queries {
				documents = append(documents, webservertypes.Document{
					Content:  query,
					Metadata: genMetadata(body.Subject, id.String(), i, body.Problems[i], body.S3VideoUri[i]),
				})
			}
			webservertypes.StatusMap[id.String()] = webservertypes.StatusRes{
				Status: "Started adding to vector db",
			}
			if body.Textbook != "-jit" {
				go chromadb.Add(id.String(), documents, body.Textbook)
			}

			go chromadb.Add(id.String(), documents, body.Subject)
			out, _ := json.Marshal(webservertypes.DefaultResponse{
				Res: "Caching in vector db",
				Url: fmt.Sprintf("%v/status/add/%v", os.Getenv("JIT_SERVER_URL"), id.String()),
			})
			ctx.Response.AppendBody(out)

		case "/get":
			var body webservertypes.ModelQueryReq
			err := json.Unmarshal(ctx.Request.Body(), &body)
			body.Textbook = fmt.Sprintf("%v-jit", body.Textbook)
			body.Subject = fmt.Sprintf("%v-jit", body.Subject)
			if err != nil {
				out, _ := json.Marshal(webservertypes.ErrRes{
					Err: err.Error(),
				})
				ctx.Response.SetStatusCode(400)
				ctx.Response.AppendBody(out)
			}

			uuid, err := uuid.NewUUID()
			if err != nil {
				out, _ := json.Marshal(webservertypes.ErrRes{
					Err: err.Error(),
				})
				ctx.Response.SetStatusCode(400)
				ctx.Response.AppendBody(out)
			}
			go chromadb.CompareToDb(body.Queries, uuid.String(), body.Subject, body.Textbook)
			webservertypes.JitQueryMap[uuid.String()] = webservertypes.JitQueryStatusRes{
				Status: "Pending",
			}
			out, _ := json.Marshal(webservertypes.DefaultResponse{
				Res: "Querying from vector db",
				Url: fmt.Sprintf("%v/status/get/%v", os.Getenv("JIT_SERVER_URL"), uuid.String()),
			})
			ctx.Response.AppendBody(out)
		case "/get_raw":
			var body webservertypes.ModelQueryReqRaw
			err := json.Unmarshal(ctx.Request.Body(), &body)
			body.Textbook = fmt.Sprintf("%v-jit", body.Textbook)
			body.Subject = fmt.Sprintf("%v-jit", body.Subject)
			if err != nil {
				out, _ := json.Marshal(webservertypes.ErrRes{
					Err: err.Error(),
				})
				ctx.Response.SetStatusCode(400)
				ctx.Response.AppendBody(out)
			}

			uuid, err := uuid.NewUUID()
			if err != nil {
				out, _ := json.Marshal(webservertypes.ErrRes{
					Err: err.Error(),
				})
				ctx.Response.SetStatusCode(400)
				ctx.Response.AppendBody(out)
			}
			go chromadb.RawCompare(body.Queries, body.Vectors, uuid.String(), body.Subject, body.Textbook)
			webservertypes.JitQueryMap[uuid.String()] = webservertypes.JitQueryStatusRes{
				Status: "Pending",
			}
			out, _ := json.Marshal(webservertypes.DefaultResponse{
				Res: "Querying from vector db",
				Url: fmt.Sprintf("%v/status/get/%v", os.Getenv("JIT_SERVER_URL"), uuid.String()),
			})
			ctx.Response.AppendBody(out)
		default:
			ctx.Response.SetStatusCode(404)
		}
	}
}

func genMetadata(subject string, id string, ord int, problem webservertypes.Problem, videoUrl string) map[string]interface{} {
	fmt.Printf("images: %v", videoUrl)
	out := make(map[string]interface{})
	out["subject"] = subject
	out["id"] = id
	out["ord"] = strconv.Itoa(ord)
	tmp, _ := json.Marshal(problem)
	out["problem"] = string(tmp)
	out["s3VideoUri"] = string(videoUrl)
	return out
}

func main() {
	godotenv.Load()
	chromadb.InitVectorDb(false)
	webservertypes.JitQueryMap = make(map[string]webservertypes.JitQueryStatusRes)
	webservertypes.StatusMap = make(map[string]webservertypes.StatusRes, 0)
	err := fasthttp.ListenAndServe(":8002", handler)
	if err != nil {
		fmt.Println(err)
	}
}
