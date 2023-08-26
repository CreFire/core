package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func ClientApi() {
	proxyUrl, _ := url.Parse("http://127.0.0.1:7890")
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	url := "https://api.openai.com/v1/models"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer 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")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)

	responseJson, _ := json.MarshalIndent(response, "", "    ")
	os.WriteFile("response.txt", responseJson, 0644)

}
