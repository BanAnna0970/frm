package frm

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/buger/jsonparser"
)

type HCaptcha struct {
	ClienKey string `json:"clientKey,omitempty"`
	TaskType struct {
		Type        string `json:"type,omitempty"`
		WebsiteURL  string `json:"websiteURL,omitempty"`
		WebsiteKey  string `json:"websiteKey,omitempty"`
		IsInvisible bool   `json:"isInvisible,omitempty"`
		Data        string `json:"data,omitempty"`
		UserAgent   string `json:"userAgent,omitempty"`
		Cookies     string `json:"cookies,omitempty"`
	} `json:"task,omitempty"`
}

type V2 struct {
	ClienKey string `json:"clientKey,omitempty"`
	TaskType struct {
		Type         string `json:"type,omitempty"`
		WebsiteURL   string `json:"websiteURL,omitempty"`
		WebsiteKey   string `json:"websiteKey,omitempty"`
		ReDataSValue string `json:"recaptchaDataSValue,omitempty"`
		UserAgent    string `json:"userAgent,omitempty"`
		Cookies      string `json:"cookies,omitempty"`
	} `json:"task,omitempty"`
}

type V3 struct {
	ClienKey string `json:"clientKey,omitempty"`
	TaskType struct {
		Type       string  `json:"type,omitempty"`
		WebsiteURL string  `json:"websiteURL,omitempty"`
		WebsiteKey string  `json:"websiteKey,omitempty"`
		MinScore   float32 `json:"minScore,omitempty"`
		PageAction string  `json:"pageAction,omitempty"`
	} `json:"task,omitempty"`
}

func (c *HCaptcha) SolveCaptcha() (solution string, err error) {
	for i := 0; i < 3; i++ {
		taskID, err := c.createTask()
		if err != nil {
			continue
		}
		solution, err := getTaskResult(c.ClienKey, taskID)
		if err != nil {
			return "", err
		}

		return solution, err
	}
	return "", errors.New("error solving captcha")
}
func (c *HCaptcha) createTask() (taskID int64, err error) {
	c.TaskType.Type = "HCaptchaTaskProxyless"
	json, _ := json.Marshal(c)

	data := FastData{
		URL:     "https://api.capmonster.cloud/createTask",
		Method:  "POST",
		Payload: string(json),
	}

	fh := data.Build()

	err = fh.DoRequest()
	if err != nil {
		Logger.Error().Err(err)
		return 0, err
	}

	body := fh.Request.Body()

	errorID, err := jsonparser.GetInt(body, "errorID")
	if err != nil || errorID != 0 {
		errorDescription, _ := jsonparser.GetString(body, "errorDescription")
		err = fmt.Errorf("err: %v\nerrorId: %v", err, errorDescription)
		Logger.Error().Err(err)
		return 0, err
	}

	taskID, err = jsonparser.GetInt(body, "taskID")
	if err != nil {
		Logger.Error().Err(err)
		return 0, err
	}
	return taskID, nil
}

func getTaskResult(clientKey string, taskID int64) (solution string, err error) {
	data := FastData{
		URL:     "https://api.capmonster.cloud/getTaskResult",
		Method:  "POST",
		Payload: fmt.Sprintf(`{"clientKey":"%v","taskID":%v}`, clientKey, taskID),
	}

	fh := data.Build()

	for i := 0; i < 20; i++ {
		time.Sleep(2 * time.Second)
		err := fh.DoRequest()
		if err != nil {
			Logger.Error().Err(err)
			return "", err
		}

		body := fh.Request.Body()
		solution, err = jsonparser.GetString(body, "solution", "gRecaptchaResponse")
		if err != nil {
			continue
		}
		return solution, err
	}

	return "", errors.New("failed to get captcha response")
}

// func (c *V2) CreateTask() {
// 	json, _ := json.Marshal(c)

// }

// func (c *V3) CreateTask() {
// 	json, _ := json.Marshal(c)

// }

