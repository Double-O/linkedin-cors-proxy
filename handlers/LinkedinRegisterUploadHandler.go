package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LinkedinRegisterUploadResp struct {
	Value Value `json:"value"`
}

type Headers struct {
	MediaTypeFamily string `json:"media-type-family"`
}

type ComLinkedinDigitalmediaUploadingMediaUploadHTTPRequest struct {
	UploadURL string  `json:"uploadUrl"`
	Headers   Headers `json:"headers"`
}

type UploadMechanism struct {
	ComLinkedinDigitalmediaUploadingMediaUploadHTTPRequest ComLinkedinDigitalmediaUploadingMediaUploadHTTPRequest `json:"com.linkedin.digitalmedia.uploading.MediaUploadHttpRequest"`
}

type Value struct {
	MediaArtifact      string          `json:"mediaArtifact"`
	UploadMechanism    UploadMechanism `json:"uploadMechanism"`
	Asset              string          `json:"asset"`
	AssetRealTimeTopic string          `json:"assetRealTimeTopic"`
}

func HandleLinkedinRegisterUpload() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var postBody map[string]interface{}
		if err := ctx.BindJSON(&postBody); err != nil {
			fmt.Printf("[handlers.HandleLinkedinRegisterUpload] error while fetching postbody from request, err : %+v\n", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":          http.StatusBadRequest,
				"error_message": err.Error(),
			})
			return
		}

		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(postBody)
		req, err := http.NewRequest("POST", "https://api.linkedin.com/v2/assets?action=registerUpload", buf)
		if err != nil {
			fmt.Printf("[handlers.HandleLinkedinRegisterUpload] error while creating request, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}

		req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", ctx.GetHeader("Authorization")) // pass Authorization header from the client

		client := &http.Client{}
		respHttp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[handlers.HandleLinkedinRegisterUpload] error while executing request, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}
		defer respHttp.Body.Close()

		body, err := io.ReadAll(respHttp.Body)
		if err != nil {
			fmt.Printf("[handlers.HandleLinkedinRegisterUpload] error while reading resp body, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}

		if respHttp.StatusCode != http.StatusCreated && respHttp.StatusCode != http.StatusOK {
			fmt.Printf("[handlers.HandleLinkedinRegisterUpload] error from linkedin api, err : %+s\n", string(body))
			ctx.JSON(respHttp.StatusCode, gin.H{
				"code":          respHttp.StatusCode,
				"error_message": string(body),
			})
			return
		}

		var resp LinkedinRegisterUploadResp
		err = json.Unmarshal(body, &resp)
		if err != nil {
			fmt.Printf("[handlers.HandleLinkedinRegisterUpload] error while unmarshaling resp, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}

		// this can be removed
		fmt.Printf("[handlers.HandleLinkedinRegisterUpload] LinkedinAccesstokenResp : %+v\n", resp)

		for key, values := range respHttp.Header {
			for _, value := range values {
				if ctx.GetHeader(key) != "" && key != "Content-Length" {
					ctx.Header(key, value)
				}
			}
		}
		ctx.JSON(http.StatusCreated, resp)
	}
}
