package docdb

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	webservertypes "github.com/vidur2/vectorMicroservices/pkg/webserverTypes"
)

func HandlerTextbook(ctx *fasthttp.RequestCtx) {
	splitPath := strings.Split(string(ctx.Path()), "/")
	if splitPath[1] == "status" {
		status := webservertypes.StatusMap[splitPath[2]]
		res := status
		if status.Status == "Finished" || strings.HasPrefix(status.Status, "Error") {
			delete(webservertypes.StatusMap, splitPath[2])
		}
		body, _ := json.Marshal(res)
		ctx.Response.AppendBody(body)
	} else {
		switch string(ctx.Path()) {
		case "/":
			ctx.Response.AppendBodyString("Please send a POST request /add with a body of { url: string, subject: string }")
		case "/add":
			if string(ctx.Method()) != "POST" {
				ctx.Response.SetStatusCode(405)
			} else {
				var parsed webservertypes.VectorDbAddRequest
				body := ctx.Request.Body()
				err := json.Unmarshal(body, &parsed)
				if err != nil {
					ctx.Response.SetStatusCode(400)
				} else {
					id := uuid.New()
					webservertypes.StatusMap[id.String()] = webservertypes.StatusRes{
						Status: "Queued",
						Info:   parsed,
					}
					go AddTextBooks(parsed.DocUrl, parsed.DocSubject, parsed.Ext, id.String())

					out, _ := json.Marshal(webservertypes.DefaultResponse{
						Res: "Adding to vector db",
						Url: fmt.Sprintf("%v/status/%v", os.Getenv("DOMAIN_URL"), id.String()),
					})

					ctx.Response.AppendBody(out)
				}
			}
		default:
			ctx.Response.SetStatusCode(404)
		}
	}
}
