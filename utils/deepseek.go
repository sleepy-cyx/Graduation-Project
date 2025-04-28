package utils

import (
	"Graduation-Project/log"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type ApiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type ScheduleJson struct {
	Title     string `gorm:"not null" json:"title"`
	Location  string `gorm:"not null" json:"location"`
	Comment   string `gorm:"not null" json:"comment"`
	StartTime string `gorm:"not null" json:"start_time"`
	EndTime   string `gorm:"not null" json:"end_time"`
}

func TextToJson(text string) (ScheduleJson, error) {
	url := "https://api.deepseek.com/chat/completions"
	method := "POST"
	schedule := ScheduleJson{}
	todayStr := time.Now().Format("2006-01-02")
	// 修复后的请求体结构
	requestBody := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role": "user",
				"content": "将文本中的日程信息提取并且转化成格式化的json，日程结构体为type Schedule struct { " +
					"title string; location string; comment string; start_time string; end_time string " +
					"} 返回的json的key值跟结构体的变量需要一一对应，如果提取不到对应信息，对应字段可以填空字符串，开始时间和结束时间形如yyyy-mm-ddT14:00:00+08:00格式。" +
					"文本内容是：" + text + "今天的日期为:" + todayStr,
			},
		},
		"model":           "deepseek-chat",
		"temperature":     0.3,
		"response_format": map[string]string{"type": "json_object"},
	}

	jsonBody, _ := json.Marshal(requestBody)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Logger.Error(err)
		return schedule, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer sk-a7fc2b76e58e4bcca774f30181e589d1")

	res, err := client.Do(req)
	if err != nil {
		log.Logger.Error(err)
		return schedule, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Logger.Error(err)
		return schedule, err
	}
	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Logger.Errorf("解析API响应失败:%v", err)
		return schedule, err
	}

	// 检查是否存在有效响应
	if len(apiResponse.Choices) == 0 || apiResponse.Choices[0].Message.Content == "" {
		log.Logger.Error("API返回空响应")
		return schedule, err
	}

	// 提取并解析内层JSON
	content := apiResponse.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &schedule); err != nil {
		log.Logger.Errorf("解析日程信息失败:%v", err)
		return schedule, err
	}
	return schedule, nil
}
