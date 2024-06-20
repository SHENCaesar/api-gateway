Github 地址
https://github.com/SHENCaesar/api-gateway
项目代码运行方式：
1. 启动etcd
2. 启动数据库（在 github.com/SHENCaesar/api-gateway/kitex/dal/mysql/init.go 里设置数据库端口)
3. 启动两个终端分别进入api-gateway/kitex和api-gateway/hz-gateway目录，分别用两个终端运行：
go run .                       
API 文档
注册学生信息
请求方式
- POST
访问地址
http://localhost:8888/gateway/student
输入参数
- method: register 
- biz_params: JSON格式的字符串，包含查询条件。
示例
curl -X POST http://localhost:8888/gateway/student \
     -H "Content-Type: application/json" \
     -d '{
          "method": "register",
          "biz_params": "{\"id\": 1, \"name\":\"Emma\", \"college\": {\"name\": \"software college\", \"address\": \"逸夫\"}, \"email\": [\"emma@pku.com\"]}"
        }'
返回值
- 成功时返回学生注册成功的响应
{"code":0,"data":{"message":"","success":true},"message":"ok"}
[图片]
- 失败时返回错误信息。

查询学生信息
请求方式
- POST
访问地址
http://localhost:8888/gateway/student
输入参数
- method: query 
- biz_params: JSON格式的字符串，包含查询条件。
示例
curl -X POST http://localhost:8888/gateway/student \
     -H "Content-Type: application/json" \
     -d '{
          "method": "query",
          "biz_params": "{\"id\": 1}"
        }'
返回值
- 成功时返回查询到的学生信息。
{"code":0,"data":{"college":{"address":"","name":""},"email":["student-1@pku.com"],"id":1,"name":"student-1"},"message":"ok"}
[图片]
- 失败时返回错误信息。
数据库示例
原本数据库为空：
[图片]
注册学生后：
[图片]

性能测试和优化
测试方案
main_test.go  单元测试和基准测试
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



