package go_zhipu

import (
	"bytes"
	"fmt"
	"github.com/itcwc/go-zhipu/model_api"
	"io"
	"log"
	"testing"
)

func TestBase(t *testing.T) {
	expireAtTime := int64(1640995200) // token 过期时间
	mssage := []model_api.Messages{
		{
			Role:    "user", // 消息的角色信息 详见文档
			Content: "你好",   // 消息内容
		},
	}
	apiKey := "e2b18c2d32d1ae45ba3e15f9fa21c84a.XbElxPNCc5uy2bpf"
	model := "glm-3-turbo"
	body := model_api.PostParams{
		Stream:   true,
		Model:    model,
		Messages: mssage,
	}

	postResponse, err := model_api.BeCommonModelStream(expireAtTime, body, apiKey)
	if err != nil {
		fmt.Println("响应错误", err)
		return
	}

	// 创建一个缓冲区
	buffer := make([]byte, 1024)

	// 循环读取流中的数据
	for {
		// 从流中读取数据
		n, err := postResponse.Read(buffer)
		if err != nil {
			if err == io.EOF {
				// 如果遇到了流结束，退出循环
				break
			}
			log.Fatalf("读取流数据失败: %v", err)
		}

		// 判断是否已经读取到流的结束标志
		if bytes.Contains(buffer[:n], []byte("data: [DONE]")) {
			fmt.Println("流已经结束")
			break
		}
		fmt.Println("本次数据", string(buffer[:n]))
		m, err := model_api.ParseResponse(string(buffer[:n]))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(m)

		// 输出读取到的数据
		//var data model_api.StreamResponseData
		//err = json.Unmarshal(buffer[:n], &data)
		//if err != nil {
		//	fmt.Println("解析错误", err)
		//	return
		//}
	}

	// 关闭流
	if closer, ok := postResponse.(io.Closer); ok {
		closer.Close()
	}
	//postResponse, err := model_api.BeCommonModelStream(expireAtTime, body, apiKey)
	//if err != nil {
	//	fmt.Println("响应错误", err)
	//	return
	//}
	//
	//// 创建一个 JSON 解码器
	//decoder := json.NewDecoder(postResponse)
	//
	//// 用来存储解析后的数据
	//var responseData *model_api.StreamResponseData
	//
	//// 解码JSON数据
	//if err := decoder.Decode(&responseData); err != nil {
	//	if err == io.EOF {
	//		fmt.Println("完成", err)
	//		// 如果遇到了流结束，退出循环
	//		return
	//	}
	//	fmt.Println("解析错误", err)
	//	return
	//}
	//
	//
	//// 输出解析后的数据
	//fmt.Println("解析后的数据:", responseData)
	//if closer, ok := postResponse.(io.Closer); ok {
	//	closer.Close()
	//}
}

//func TestParse(t *testing.T) {
//	rawResponse := `data: {"id":"8655895084942182103","created":1715653541,"model":"glm-3-turbo","choices":[{"index":0,"delta":{"role":"assistant","content":"你好"}}]}`
//
//	// 解析响应
//	responseData, err := model_api.ParseResponse(rawResponse)
//	if err != nil {
//		fmt.Println("解析响应失败:", err)
//		return
//	}
//
//	// 输出解析后的数据
//	fmt.Println("解析后的数据:", responseData)
//}
