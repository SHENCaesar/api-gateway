package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/SHENCaesar/api-gateway/hz-gateway/kitex_gen/demo"
	"github.com/cloudwego/hertz/pkg/common/json"
)

const gatewayURL = "http://127.0.0.1:8888/gateway/student"

var httpCli = &http.Client{Timeout: 3 * time.Second}

type reqParam struct {
	Method    string `json:"method"`
	BizParams string `json:"biz_params"`
}

func TestStudentService(t *testing.T) {
	for i := 1; i <= 100; i++ {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			newStu := genStudent(i)
			resp, err := request("register", newStu)
			if err != nil {
				t.Errorf("Failed to register student: %v", err)
			}
			if resp["message"] != "ok" {
				t.Errorf("Expected message 'ok', got '%s'", resp["message"])
			}

			jsonData, err := json.Marshal(resp["data"])
			if err != nil {
				t.Errorf("Failed to marshal response data: %v", err)
			}

			var registerResp demo.RegisterResp
			err = json.Unmarshal(jsonData, &registerResp)
			if err != nil {
				t.Errorf("Failed to unmarshal RegisterResp: %v", err)
			}
			if !registerResp.Success {
				t.Errorf("Registration was not successful: %v", registerResp)
			}

			resp, err = request("query", newStu)
			if err != nil {
				t.Errorf("Failed to query student: %v", err)
			}
			if resp["message"] != "ok" {
				t.Errorf("Expected message 'ok', got '%s'", resp["message"])
			}

			jsonData, err = json.Marshal(resp["data"])
			if err != nil {
				t.Errorf("Failed to marshal query response data: %v", err)
			}

			var stu demo.Student
			err = json.Unmarshal(jsonData, &stu)
			if err != nil {
				t.Errorf("Failed to unmarshal Student: %v", err)
			}
			if stu.Id != newStu.Id || stu.Name != newStu.Name || stu.Email[0] != newStu.Email[0] || stu.College.Name != newStu.College.Name {
				t.Errorf("Student data mismatch: expected %+v, got %+v", newStu, stu)
			}
		})
	}
}

func BenchmarkStudentService(b *testing.B) {
	// 基准测试的准备阶段，例如创建测试数据
	prepareData := func(id int) *demo.Student {
		return &demo.Student{
			Id:   int32(id),
			Name: fmt.Sprintf("student-%d", id),
			College: &demo.College{
				Name:    "",
				Address: "",
			},
			Email: []string{fmt.Sprintf("student-%d@pku.com", id)},
		}
	}

	b.ResetTimer() // 重置计时器，忽略准备阶段的时间

	for i := 0; i < b.N; i++ {
		newStu := prepareData(i)
		resp, err := request("register", newStu)
		if err != nil {
			b.Errorf("Failed to register student: %v", err)
			continue
		}
		if resp["message"] != "ok" {
			b.Errorf("Expected message 'ok', got '%s'", resp["message"])
			continue
		}

		jsonData, err := json.Marshal(resp["data"])
		if err != nil {
			b.Errorf("Failed to marshal response data: %v", err)
			continue
		}

		var registerResp demo.RegisterResp
		err = json.Unmarshal(jsonData, &registerResp)
		if err != nil {
			b.Errorf("Failed to unmarshal RegisterResp: %v", err)
			continue
		}
		if !registerResp.Success {
			b.Errorf("Registration was not successful: %v", registerResp)
			continue
		}

		resp, err = request("query", newStu)
		if err != nil {
			b.Errorf("Failed to query student: %v", err)
			continue
		}
		if resp["message"] != "ok" {
			b.Errorf("Expected message 'ok', got '%s'", resp["message"])
			continue
		}

		jsonData, err = json.Marshal(resp["data"])
		if err != nil {
			b.Errorf("Failed to marshal query response data: %v", err)
			continue
		}

		var stu demo.Student
		err = json.Unmarshal(jsonData, &stu)
		if err != nil {
			b.Errorf("Failed to unmarshal Student: %v", err)
			continue
		}
		if stu.Id != newStu.Id || stu.Name != newStu.Name || stu.Email[0] != newStu.Email[0] || stu.College.Name != newStu.College.Name {
			b.Errorf("Student data mismatch: expected %+v, got %+v", newStu, stu)
		}
	}
}

func request(method string, bizParam interface{}) (map[string]interface{}, error) {
	bizParamBody, err := json.Marshal(bizParam)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal business parameters: %v", err)
	}

	reqBody := &reqParam{
		Method:    method,
		BizParams: string(bizParamBody),
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, gatewayURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var rResp map[string]interface{}
	if err := json.Unmarshal(body, &rResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return rResp, nil
}

func genStudent(id int) *demo.Student {
	return &demo.Student{
		Id:   int32(id),
		Name: fmt.Sprintf("student-%d", id),
		College: &demo.College{
			Name:    "",
			Address: "",
		},
		Email: []string{fmt.Sprintf("student-%d@pku.com", id)},
	}
}
