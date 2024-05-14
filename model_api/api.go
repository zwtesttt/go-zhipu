package model_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"strings"
	"time"

	"github.com/itcwc/go-zhipu/utils"
)

var v4url string = "https://open.bigmodel.cn/api/paas/v4/"

var v3url string = "https://open.bigmodel.cn/api/paas/v4/"

type PostParams struct {
	Model    string     `json:"model"`
	Messages []Messages `json:"messages"`
	Stream   bool       `json:"stream"`
}
type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type StreamResponseData struct {
	ID      string `json:"id"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
	FinishReason string `json:"finish_reason,omitempty"`
	Usage        struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage,omitempty"`
}

func ParseResponse(rawResponse string) (*StreamResponseData, error) {
	// 找到有效数据部分的起始位置
	startIndex := strings.Index(rawResponse, "{")
	if startIndex == -1 {
		return nil, errors.New("未找到有效数据部分")
	}

	// 截取有效数据部分
	validData := rawResponse[startIndex:]

	// 手动解析
	var data *StreamResponseData
	if err := json.Unmarshal([]byte(validData), &data); err != nil {
		return nil, fmt.Errorf("解析数据失败: %v", err)
	}

	return data, nil
}

// 通用模型
func BeCommonModel(expireAtTime int64, postParams PostParams, apiKey string) (map[string]interface{}, error) {

	token, _ := utils.GenerateToken(apiKey, expireAtTime)

	// 示例用法
	apiURL := v4url + "chat/completions"
	timeout := 60 * time.Second

	postResponse, err := utils.Post(apiURL, token, postParams, timeout)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	return postResponse, nil
}

func BeCommonModelStream(expireAtTime int64, postParams PostParams, apiKey string) (io.Reader, error) {
	token, _ := utils.GenerateToken(apiKey, expireAtTime)

	// 示例用法
	apiURL := v4url + "chat/completions"
	timeout := 60 * time.Second

	// 创建管道
	reader, writer := io.Pipe()

	go func() {
		// 在 goroutine 中执行请求，并将结果写入管道
		postResponse, err := utils.Stream(apiURL, token, postParams, timeout)
		if err != nil {
			writer.CloseWithError(fmt.Errorf("创建请求失败: %v", err))
			return
		}

		_, err = io.Copy(writer, postResponse.Body)
		if err != nil {
			writer.CloseWithError(fmt.Errorf("写入管道失败: %v", err))
			return
		}

		defer postResponse.Body.Close()
		defer writer.Close()
	}()

	return reader, nil
}

type PostImageParams struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// 图像大模型
func ImageLargeModel(expireAtTime int64, prompt string, apiKey string, model string) (map[string]interface{}, error) {

	token, _ := utils.GenerateToken(apiKey, expireAtTime)

	// 示例用法
	apiURL := v4url + "images/generations"
	timeout := 60 * time.Second

	// 示例 POST 请求
	postParams := PostImageParams{
		Model:  model,
		Prompt: prompt,
	}

	postResponse, err := utils.Post(apiURL, token, postParams, timeout)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	return postResponse, nil
}

type PostSuperhumanoidParams struct {
	Prompt []Prompt `json:"prompt"`
	Meta   []Meta   `json:"meta"`
}
type Prompt struct {
	Role    string `json:"prompt"`
	Content string `json:"content"`
}
type Meta struct {
	UserInfo string `json:"user_info"`
	BotInfo  string `json:"bot_info"`
	BotName  string `json:"bot_name"`
	UserName string `json:"user_name"`
}

// 超拟人大模型
func SuperhumanoidModel(expireAtTime int64, meta []Meta, prompt []Prompt, apiKey string) (map[string]interface{}, error) {

	token, _ := utils.GenerateToken(apiKey, expireAtTime)

	// 示例用法
	apiURL := v3url + "model-api/charglm-3/sse-invoke"
	timeout := 60 * time.Second

	// 示例 POST 请求
	postParams := PostSuperhumanoidParams{
		Prompt: prompt,
		Meta:   meta,
	}

	postResponse, err := utils.Post(apiURL, token, postParams, timeout)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	return postResponse, nil
}

type PostVectorParams struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

// 向量模型
func VectorModel(expireAtTime int64, input string, apiKey string, model string) (map[string]interface{}, error) {

	token, _ := utils.GenerateToken(apiKey, expireAtTime)

	// 示例用法
	apiURL := v4url + "mbeddings"
	timeout := 60 * time.Second

	// 示例 POST 请求
	postParams := PostVectorParams{
		Input: input,
		Model: model,
	}

	postResponse, err := utils.Post(apiURL, token, postParams, timeout)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	return postResponse, nil
}

type PostFineTuningParams struct {
	Model        string `json:"model"`
	TrainingFile string `json:"training_file"`
}

// 模型微调
func ModelFineTuning(expireAtTime int64, trainingFile string, apiKey string, model string) (map[string]interface{}, error) {

	token, _ := utils.GenerateToken(apiKey, expireAtTime)

	// 示例用法
	apiURL := v4url + "fine_tuning/jobs"
	timeout := 60 * time.Second

	// 示例 POST 请求
	postParams := PostFineTuningParams{
		Model:        model,
		TrainingFile: trainingFile,
	}

	postResponse, err := utils.Post(apiURL, token, postParams, timeout)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	return postResponse, nil
}

type PostFileParams struct {
	File    *FileHeader `json:"file"`
	Purpose string      `json:"purpose"`
}

type FileHeader struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64

	content   []byte
	tmpfile   string
	tmpoff    int64
	tmpshared bool
}

// 文件管理
func FileManagement(expireAtTime int64, purpose string, apiKey string, model string, file *FileHeader) (map[string]interface{}, error) {

	token, _ := utils.GenerateToken(apiKey, expireAtTime)

	// 示例用法
	apiURL := v4url + "files"
	timeout := 60 * time.Second

	// 示例 POST 请求
	postParams := PostFileParams{
		File:    file,
		Purpose: purpose,
	}

	postResponse, err := utils.Post(apiURL, token, postParams, timeout)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	return postResponse, nil
}
