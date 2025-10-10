package uploadfile

import (
	"autoTest/API/adminApi/login"
	"autoTest/store/logger"
	"autoTest/store/model"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

// Define structs to match the JSON response structure
type Response struct {
	Data          []DataItem  `json:"data"`
	MsgParameters interface{} `json:"msgParameters"`
	Code          int         `json:"code"`
	Msg           string      `json:"msg"`
	MsgCode       int         `json:"msgCode"`
}

type DataItem struct {
	Src   string `json:"src"`
	Title string `json:"title"`
	Size  Size   `json:"size"`
}

type Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

/*
传入文件地址 ../1.png
传入token
*
*/
func UploadFile(filePath, token string) (*model.Response, string, error) {
	// Configuration
	url := "https://sit-tenantadmin-3003.mggametransit.com/api/UploadFile/UploadToOss"
	fileType := "other"
	customPath := "" // Empty as per the request; modify if needed

	// Create a buffer to store the multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file field
	file, err := os.Open(filePath)
	if err != nil {
		errs := fmt.Errorf("failed to open file: %w", err)
		return handlerError(errs), "", errs
	}
	defer file.Close()
	part, err := writer.CreateFormFile("files", filepath.Base(filePath))
	if err != nil {
		errs := fmt.Errorf("failed to create form file: %w", err)
		return handlerError(errs), "", errs
	}
	_, err = io.Copy(part, file)
	if err != nil {
		errs := fmt.Errorf("failed to write file to form: %w", err)
		return handlerError(errs), "", errs
	}

	// Add fileType field
	err = writer.WriteField("fileType", fileType)
	if err != nil {
		errs := fmt.Errorf("failed to write fileType: %w", err)
		return handlerError(errs), "", errs
	}

	// Add customPath field (empty)
	err = writer.WriteField("customPath", customPath)
	if err != nil {
		errs := fmt.Errorf("failed to write customPath: %w", err)
		return handlerError(errs), "", errs
	}

	// Close the writer to finalize the multipart form
	err = writer.Close()
	if err != nil {
		errs := fmt.Errorf("failed to close writer: %w", err)
		return handlerError(errs), "", errs
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		errs := fmt.Errorf("failed to create request: %w", err)
		return handlerError(errs), "", nil
	}

	// Set headers to match the Yakit request
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("ignorecanceltoken", "true")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("referer", "https://sit-tenantadmin-3003.mggametransit.com/")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("origin", "https://sit-tenantadmin-3003.mggametransit.com")
	req.Header.Set("sec-fetch-dest", "empty")
	// Comment out accept-encoding to force uncompressed response for debugging
	// req.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("authorization", "Bearer "+token)
	req.Header.Set("sec-ch-ua", `"Chromium";v="140", "Not=A?Brand";v="24", "Google Chrome";v="140"`)
	req.Header.Set("domainurl", "https://sit-tenantadmin-3003.mggametransit.com")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")

	// Create an HTTP client
	client := &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
		},
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		errs := fmt.Errorf("failed to send request: %w", err)
		return handlerError(errs), "", errs
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		errs := fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(respBody))
		return handlerError(errs), "", errs
	}

	// Check Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		respBody, _ := io.ReadAll(resp.Body)
		errs := fmt.Errorf("unexpected Content-Type: %s, response: %s", contentType, string(respBody))
		return handlerError(errs), "", errs
	}

	// Handle compressed response
	var reader io.Reader = resp.Body
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			errs := fmt.Errorf("failed to create gzip reader: %w", err)
			return handlerError(errs), "", errs
		}
		defer gz.Close()
		reader = gz
	case "br":
		reader = brotli.NewReader(resp.Body)
	case "zstd":
		zstdReader, err := zstd.NewReader(resp.Body)
		if err != nil {
			errs := fmt.Errorf("failed to create zstd reader: %w", err)
			return handlerError(errs), "", errs
		}
		defer zstdReader.Close()
		reader = zstdReader
	case "":
		// No compression
	default:
		errs := fmt.Errorf("unsupported Content-Encoding: %s", resp.Header.Get("Content-Encoding"))
		return handlerError(errs), "", errs
	}

	// Read the response
	respBody, err := io.ReadAll(reader)
	if err != nil {
		errs := fmt.Errorf("failed to read response: %w", err)
		return handlerError(errs), "", errs
	}

	// Check if response is empty
	if len(respBody) == 0 {
		errs := fmt.Errorf("empty response body")
		return handlerError(errs), "", errs
	}

	// Debug: Print response headers and raw response
	//fmt.Printf("Response Headers: %+v\n", resp.Header)
	//fmt.Printf("Raw Response: %s\n", string(respBody))

	// Parse the JSON response
	var response Response
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		errs := fmt.Errorf("failed to parse JSON response: %w, raw response: %s", err, string(respBody))
		return handlerError(errs), "", errs
	}

	// Print the parsed response
	// fmt.Printf("Parsed Response: %+v\n", response.Data[0].Title)
	return &model.Response{
		Code:    response.Code,
		MsgCode: response.MsgCode,
		Msg:     response.Msg,
	}, response.Data[0].Title, nil
}

/*
fileName  上传的文件地址
返回  *model.Response 和 文件地址的地址
*
*/
func RunWorkerOderActiveZx(ctx *context.Context, fileName string, ch chan struct{}) (*model.Response, string) {
	// ctx := context.Background()
	// _, ctxT, err := login.AdminSitLogin(&ctx)
	// if err != nil {
	// 	fmt.Println("Login error:", err)
	// 	return
	// }
	token := (*ctx).Value(login.AuthTokenKey)
	// "./assert/workerOder/1.png"
	if resp, str, err := UploadFile(fileName, token.(string)); err != nil {
		errs := fmt.Errorf("文件上传失败%s", err)
		logger.LogError("文件上传失败", err)
		return handlerError(errs), ""
	} else {
		ch <- struct{}{}
		return resp, str
	}
}

// 专门处理这里的错误信息
func handlerError(err error) *model.Response {
	return &model.Response{
		Code:    -1,
		Msg:     err.Error(),
		MsgCode: -1,
	}
}
