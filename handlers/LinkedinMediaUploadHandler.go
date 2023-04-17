package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func HandleLinkedinMediaUpload() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		file, err := ctx.FormFile("file")
		if err != nil {
			fmt.Printf("[handlers.handleLinkedinMediaUpload] error while fetching file from request, err : %+v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":          http.StatusBadRequest,
				"error_message": err.Error(),
			})
			return
		}

		uploadUrl, err := url.QueryUnescape(ctx.Query("upload_url"))
		if err != nil {
			fmt.Printf("[handlers.handleLinkedinMediaUpload] error while decoding uploadurl, err : %+v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":          http.StatusBadRequest,
				"error_message": err.Error(),
			})
			return
		}

		// Open the file
		src, err := file.Open()
		if err != nil {
			fmt.Printf("[handlers.handleLinkedinMediaUpload] failed to open file, err : %+v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}
		defer src.Close()

		req, err := http.NewRequest("PUT", uploadUrl, src)
		if err != nil {
			fmt.Printf("[handlers.handleLinkedinMediaUpload] error while creating request, err : %+v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}

		req.Header.Set("Authorization", ctx.GetHeader("Authorization")) // pass Authorization header from the client
		req.Header.Set("X-Restli-Protocol-Version", "2.0.0")

		client := &http.Client{}
		respHttp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[handlers.handleLinkedinMediaUpload] error while executing request, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}
		defer respHttp.Body.Close()

		body, err := io.ReadAll(respHttp.Body)
		if err != nil {
			fmt.Printf("[handlers.handleLinkedinMediaUpload] error while reading resp body, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}

		if respHttp.StatusCode != http.StatusCreated && respHttp.StatusCode != http.StatusOK {
			fmt.Printf("[handlers.handleLinkedinMediaUpload] error from linkedin api, err : %+s\n", string(body))
			ctx.JSON(respHttp.StatusCode, gin.H{
				"code":          respHttp.StatusCode,
				"error_message": string(body),
			})
			return
		}

		// no response body in this linkedin api

		for key, values := range respHttp.Header {
			for _, value := range values {
				if ctx.GetHeader(key) != "" && key != "Content-Length" {
					ctx.Header(key, value)
				}
			}
		}
		ctx.Status(http.StatusCreated)
	}
}
