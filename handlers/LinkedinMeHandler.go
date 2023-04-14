package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LinkedinMeResp struct {
	LocalizedLastName string `json:"localizedLastName"`
	FirstName         struct {
		Localized       map[string]string `json:"localized"`
		PreferredLocale struct {
			Country  string `json:"country"`
			Language string `json:"language"`
		} `json:"preferredLocale"`
	} `json:"firstName"`
	LastName struct {
		Localized       map[string]string `json:"localized"`
		PreferredLocale struct {
			Country  string `json:"country"`
			Language string `json:"language"`
		} `json:"preferredLocale"`
	} `json:"lastName"`
	ID                 string `json:"id"`
	LocalizedFirstName string `json:"localizedFirstName"`
}

func HandleLinkedInMe() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req, err := http.NewRequest("GET", "https://api.linkedin.com/v2/me", nil)
		if err != nil {
			fmt.Printf("[handlers.handleLinkedInMe] error while creating request, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}

		req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
		// req.Header.Set("Authorization", "Bearer AQUL6Ei5u42FjzNj6416-hX5_1MdIMY67qq1N3zK_znmN4sAYgEOl_-xgvJmqy5MMVUnbqOX9_SDLDuZY5kWaGfi-3Qmvkkb7w8Jv0bZgYDhGjIcUGJbQGZOiMVxqFFN_KRHMHd43W3KVmq3Ij1W45LwGREr2-PSQm3_oWz6zzP-aMdf2dQLiSUPjKfsYVYimqPim8B9YFA6XUtXe50I_XvF1vCcdZuqlHuuRZd9E5Qzk_rQNWCinZlHP4Tef8aG3_qadPZRiEb8lAwmfivSrgQFUh6iWKB8u1YpwL1szDILP4Fu6ki5XAlgu6gz-h70SW83gzb_nriX8BhnXRmHNXonkF62kA")
		req.Header.Set("Authorization", ctx.GetHeader("Authorization")) // pass Authorization header from the client
		req.Header.Set("Access-Control-Allow-Credentials", "true")
		req.Header.Set("Access-Control-Allow-Headers", "Accept, X-Access-Token, X-Application-Name, X-Request-Sent-Time, Authorization")
		req.Header.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		req.Header.Set("Access-Control-Allow-Origin", "*")
		req.Header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Encoding")

		client := &http.Client{}
		respHttp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[handlers.handleLinkedInMe] error while executing request, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}
		defer respHttp.Body.Close()

		body, err := io.ReadAll(respHttp.Body)
		if err != nil {
			fmt.Printf("[handlers.handleLinkedInMe] error while reading resp body, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}

		if respHttp.StatusCode != http.StatusOK {
			fmt.Printf("[handlers.handleLinkedInMe] error from linkedin api, err : %+s\n", string(body))
			ctx.JSON(respHttp.StatusCode, gin.H{
				"code":          respHttp.StatusCode,
				"error_message": string(body),
			})
			return
		}

		var resp LinkedinMeResp
		err = json.Unmarshal(body, &resp)
		if err != nil {
			fmt.Printf("[handlers.handleLinkedInMe] error while unmarshaling resp, err : %+v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":          http.StatusInternalServerError,
				"error_message": err.Error(),
			})
			return
		}

		// this can be removed
		fmt.Printf("[handlers.handleLinkedInMe] LinkedinMeResp : %+v\n", resp)

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
