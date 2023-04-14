package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LinkedinAccesstokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func HandleLinkedinAccessToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req, err := http.NewRequest("POST", "https://www.linkedin.com/oauth/v2/accessToken?"+ctx.Request.URL.RawQuery, nil)
		if err != nil {
			fmt.Printf("[handlers.HandleLinkedinAccessToken] error while creating request, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}

		req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		respHttp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[handlers.HandleLinkedinAccessToken] error while executing request, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}
		defer respHttp.Body.Close()

		body, err := io.ReadAll(respHttp.Body)
		if err != nil {
			fmt.Printf("[handlers.HandleLinkedinAccessToken] error while reading resp body, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}

		if respHttp.StatusCode != http.StatusOK {
			fmt.Printf("[handlers.HandleLinkedinAccessToken] error from linkedin api, err : %+s\n", string(body))
			ctx.JSON(respHttp.StatusCode, gin.H{
				"code":          respHttp.StatusCode,
				"error_message": string(body),
			})
			return
		}

		var resp LinkedinAccesstokenResp
		err = json.Unmarshal(body, &resp)
		if err != nil {
			fmt.Printf("[handlers.HandleLinkedinAccessToken] error while unmarshaling resp, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}

		// this can be removed
		fmt.Printf("[handlers.HandleLinkedinAccessToken] LinkedinAccesstokenResp : %+v\n", resp)

		for key, values := range respHttp.Header {
			for _, value := range values {
				if ctx.GetHeader(key) != "" {
					ctx.Header(key, value)
				}
			}
		}
		ctx.JSON(http.StatusOK, resp)

	}
}
