package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Person struct {
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
			ctx.AbortWithError(http.StatusInternalServerError, err)
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
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[handlers.handleLinkedInMe] error while executing request, err : %+v\n", err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("[handlers.handleLinkedInMe] error while reading resp body, err : %+v\n", err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var person Person
		err = json.Unmarshal(body, &person)
		if err != nil {
			fmt.Printf("[handlers.handleLinkedInMe] error while unmarshaling person, err : %+v\n", err)
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		fmt.Printf("[handlers.handleLinkedInMe] person : %+v\n", person)

		for key, values := range resp.Header {
			for _, value := range values {
				ctx.Header(key, value)
			}
		}
		ctx.JSON(http.StatusOK, person)

	}
}
