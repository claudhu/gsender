package pkg

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type GSender struct{}

type Payload struct {
	To        string   `json:"to"`
	Subject   string   `json:"subject"`
	Body      string   `json:"body"`
	Cc        string   `json:"cc"`
	Bcc       string   `json:"bcc"`
	FilePaths []string `json:"file_paths"`
}

func (g *GSender) Send(p Payload) {

	url := os.Getenv("MAIL_API")
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("to", p.To)
	_ = writer.WriteField("subject", p.Subject)
	_ = writer.WriteField("body", p.Body)
	if p.Cc != "" {
		_ = writer.WriteField("cc", p.Cc)
	}
	if p.Bcc != "" {
		_ = writer.WriteField("bcc", p.Bcc)
	}
	for _, filePath := range p.FilePaths {
		file, errFile := os.Open(filePath)
		if errFile != nil {
			fmt.Println(errFile)
			return
		}
		part,
			errFile := writer.CreateFormFile("attachment", filepath.Base(filePath))
		if errFile != nil {
			fmt.Println(errFile)
			return
		}
		_, errFile = io.Copy(part, file)
		if errFile != nil {
			fmt.Println(errFile)
			return
		}
		errFile = file.Close()
		if errFile != nil {
			fmt.Println(errFile)
			return
		}
	}
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func NewGSender() *GSender {
	// check env
	if os.Getenv("MAIL_API") == "" {
		fmt.Println("MAIL_API is not set")
		return nil
	}
	return &GSender{}
}
