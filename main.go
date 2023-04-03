package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	APIKey    = "9a326c95a58ceb8528c1"
	APISecret = "77ebbf42a68df3fff74d77aeda0781dfe45688bd9a9baa4d264143c1c425172f"
	JWT       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySW5mb3JtYXRpb24iOnsiaWQiOiI5ZjRiMjIxNi04YTA1LTRlNTYtODg2MC03ZmE2OGMzYzU0OTQiLCJlbWFpbCI6ImRvZ2Vmb29kY29pbkBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicGluX3BvbGljeSI6eyJyZWdpb25zIjpbeyJpZCI6IkZSQTEiLCJkZXNpcmVkUmVwbGljYXRpb25Db3VudCI6MX1dLCJ2ZXJzaW9uIjoxfSwibWZhX2VuYWJsZWQiOmZhbHNlLCJzdGF0dXMiOiJBQ1RJVkUifSwiYXV0aGVudGljYXRpb25UeXBlIjoic2NvcGVkS2V5Iiwic2NvcGVkS2V5S2V5IjoiOWEzMjZjOTVhNThjZWI4NTI4YzEiLCJzY29wZWRLZXlTZWNyZXQiOiI3N2ViYmY0MmE2OGRmM2ZmZjc0ZDc3YWVkYTA3ODFkZmU0NTY4OGJkOWE5YmFhNGQyNjQxNDNjMWM0MjUxNzJmIiwiaWF0IjoxNjgwNTA5MjMxfQ.8gkBWIhRL_uvrEhqjjqrJgYwBbHGGDBThklJcnz7RdY"

	URL = "https://api.pinata.cloud/pinning/pinFileToIPFS"
)

func upload(c *gin.Context) {
	h, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	f, err := h.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	part, err := writer.CreateFormFile("file", h.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	_, err = io.Copy(part, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	err = writer.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, URL, buf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	//req.Header.Add("Authorization", "Bearer "+JWT)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("pinata_api_key", APIKey)
	req.Header.Add("pinata_secret_api_key", APISecret)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var mp map[string]interface{}
	if err := json.Unmarshal(data, &mp); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, mp)
}

func main() {
	r := gin.Default()
	r.GET("/upload", upload)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
