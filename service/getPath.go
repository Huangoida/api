package service

import (
	"api/util"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"
)

type MetadataRequestStruct struct {
	APIs ApiSlice
}

type ApiSlice []ApiStruct

type ApiStruct struct {
	Services   ServicesStruct
	Path       string
	Protocol   string
	Method     string
	Parameter  []ParametersStruct
	Headers    []HeaderStruct
	BatchIndex int
	ParentApi  []ParentApiStruct
	Name       string
}

type ParentApiStruct struct {
	ParentName   string
	Key          string
	DefaultValue interface{}
	ToKey        string
	ToType       string
}

type ServicesStruct struct {
	Host string
	Port string
}

type HeaderStruct struct {
	Key          string
	DefaultValue string
}

type ParametersStruct struct {
	Type         string
	Key          string
	DefaultValue string
	Require      bool
	Body         map[string]interface{}
}

type DslInfoStruct struct {
	Id        int64          `gorm:"column:id" bson:"_id"`
	Name      string         `gorm:"column:name" bson:"name"`
	Path      string         `gorm:"column:path" bson:"path"`
	Content   string         `gorm:"column:content" bson:"content"`
	Method    string         `gorm:"column:method" bson:"method"`
	CreatedAt time.Time      `gorm:"created_at;<-:create" bson:"created_at"`
	UpdatedAt time.Time      `gorm:"updated_at;<-:update" bson:"updated_at"`
	Deleted   gorm.DeletedAt `gorm:"deleted" bson:"deleted"`
}

func (api *ApiSlice) Len() int {
	return len(*api)
}

func (api *ApiSlice) Less(i, j int) bool {
	return (*api)[i].BatchIndex < (*api)[j].BatchIndex
}

func (api *ApiSlice) Swap(i, j int) {
	(*api)[i], (*api)[j] = (*api)[j], (*api)[i]
}

// 一个朴素的想法就是直接调用metadata的dsl/list()接口
func GetPathFromMetadata(path string, method string) MetadataRequestStruct {
	//根据传入的聚合请求的path和method构造对dsl/list的查询请求
	newBody := make(map[string]interface{})
	query := make(map[string]string)
	header := make(map[string]string)
	var contentType string
	failed := false

	Method := "Get"
	//注意端口映射
	url := "http://100.100.30.74:32346/v1/dsl/list"
	//有一个大问题就是这里的token。这玩意要登录的.为了防止麻烦。直接写死一个。
	header["ApiToken"] = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxNjA3MjkwNDY1ODQzOTQxMzc2LCJleHAiOjE2NzQ1NDQ5NjUsImlzcyI6Imp3dCJ9.n4aOKtwM-hvSDHYlvLqWIZN1udWNhbNT90c860ZFy6w"
	query["Path"] = path
	query["Method"] = method

	bodyByte, err := json.Marshal(newBody)
	if err != nil {
		failed = true
	}

	subResp, code := util.Do(Method, url, header, query, bodyByte, contentType)
	subMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(subResp), &subMap)
	if code != 200 {
		failed = true
	}
	if err != nil {
		failed = true
	}

	//这里有个问题就是，content=>ApiStruct在这里写还是在metadata里写
	//在这里写需要把DslInfoStruct在这里重新定义
	//在metadata里写需要添加接口或者修改接口
	dslContentStr := strings.Split(subMap["res"].(DslInfoStruct).Content, "\"APIs\":")[1]
	var dslInfoStructs []ApiStruct
	err = json.Unmarshal([]byte(dslContentStr), &dslInfoStructs)
	if err != nil {
		failed = true
	}

	request := make(map[string]MetadataRequestStruct)
	request["testOne"] = MetadataRequestStruct{APIs: dslInfoStructs}

	if failed == false {
		return request["testOne"]
	} else {
		//如果前面发生错误，这里就返回一个模拟的APi数组
		api := ApiStruct{
			Services: ServicesStruct{
				Host: "127.0.0.1",
				Port: "8081",
			},
			Path:     "/ping",
			Protocol: "http",
			Method:   "GET",
			Parameter: []ParametersStruct{{
				Key:          "id",
				DefaultValue: "string",
				Require:      true,
				Type:         "query",
				Body:         nil,
			}, {
				Key:          "id",
				DefaultValue: "string",
				Require:      true,
				Type:         "body",
				Body: map[string]interface{}{
					"test": "12345",
				},
			}},
			BatchIndex: 0,
			Name:       "ping",
		}
		api2 := ApiStruct{
			Services: ServicesStruct{
				Host: "127.0.0.1",
				Port: "8081",
			},
			Path:     "/ping1",
			Protocol: "http",
			Method:   "GET",
			Parameter: []ParametersStruct{{
				Key:          "id",
				Type:         "query",
				DefaultValue: "string",
				Require:      true,
				Body:         nil,
			}},
			Headers: []HeaderStruct{{
				Key:          "jwt",
				DefaultValue: "token",
			}, {
				Key:          "Content-Type",
				DefaultValue: "application/pdf",
			}},
			BatchIndex: 1,
			ParentApi: []ParentApiStruct{{
				ParentName:   "ping",
				Key:          "rewq",
				DefaultValue: 0,
				ToType:       "query",
				ToKey:        "ewqe",
			}, {
				ParentName:   "ping",
				Key:          "message",
				DefaultValue: 0,
				ToType:       "body",
				ToKey:        "test1",
			}, {
				ParentName:   "ping",
				Key:          "message1",
				DefaultValue: 0,
				ToType:       "body",
				ToKey:        "test2",
			}, {
				ParentName:   "ping",
				Key:          "message2.test.try",
				DefaultValue: 0,
				ToType:       "body",
				ToKey:        "test3.test.try",
			}},
			Name: "ping1",
		}
		api3 := ApiStruct{
			Services: ServicesStruct{
				Host: "127.0.0.1",
				Port: "8081",
			},
			Path:     "/postPing",
			Protocol: "http",
			Method:   "POST",
			Parameter: []ParametersStruct{{
				Key:          "postIds",
				DefaultValue: "12324 ",
				Require:      true,
				Type:         "query",
				Body:         nil,
			}, {
				Key:          "",
				DefaultValue: "",
				Require:      true,
				Type:         "json",
				Body: map[string]interface{}{
					"rerqwe": "123",
				},
			}},
			BatchIndex: 0,
			Name:       "postPing",
		}
		request["ErrorOne"] = MetadataRequestStruct{APIs: []ApiStruct{api, api2, api3}}
		return request["ErrorOne"]
	}
}
