package handler

import (
	"api/constant"
	"api/service"
	"api/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"sort"
	"strings"
	"sync"
)

func DoProxy(c *gin.Context) {

	wg := sync.WaitGroup{}
	success := make(chan map[string]interface{})
	failed := make(chan int)
	request := c.Request

	bodyReadCLoser := c.Request.Body
	bodyByte, err := ioutil.ReadAll(bodyReadCLoser)
	body := make(map[string]interface{})
	if err != nil {
		util.ResponseError(c, 400, constant.PARAMETER_INVALID, "parse body failed")
		return
	}
	if len(bodyByte) != 0 {
		err = json.Unmarshal(bodyByte, &body)
		if err != nil {
			util.ResponseError(c, 400, constant.PARAMETER_INVALID, "parse body failed")
			return
		}
	}

	// search path
	metadataInfo := service.GetPathFromMetadata(request.URL.Path, request.Method)

	if len(metadataInfo.APIs) == 0 {
		util.ResponseError(c, 404, constant.NOT_FOUND, fmt.Sprintf("%s %s not found", request.URL.Path, request.Method))
		return
	} else if len(metadataInfo.APIs) == 1 {
		wg.Add(1)
		go dealWithParallelRequest(&wg, c, metadataInfo.APIs[0], nil, success, failed)
		select {
		case <-failed:
			return
		case subMap := <-success:
			util.ResponseSuccess(c, subMap)
		}
		wg.Wait()
		return
	}

	sort.Sort(&metadataInfo.APIs)
	// do proxy
	metadataInfoList := make(map[int][]service.ApiStruct)
	for i, api := range metadataInfo.APIs {
		metadataInfoList[api.BatchIndex] = append(metadataInfoList[api.BatchIndex], metadataInfo.APIs[i])
	}
	resp := sync.Map{}
	if len(metadataInfoList) == 1 {
		for _, ApiList := range metadataInfoList {
			for _, api := range ApiList {
				wg.Add(1)
				go dealWithParallelRequest(&wg, c, api, body, success, failed)
				select {
				case <-failed:
					return
				case subMap := <-success:
					resp.Store(api.Name, subMap)
				}

			}
			wg.Wait()
		}
	} else {
		for _, ApiList := range metadataInfoList {
			for _, api := range ApiList {
				wg.Add(1)
				go dealWithSerialRequest(&wg, c, api, &resp, body, success, failed)
				select {
				case <-failed:
					return
				case subMap := <-success:
					resp.Store(api.Name, subMap)
				}

			}
			wg.Wait()
		}
	}

	res := make(map[string]interface{})
	resp.Range(func(key, value any) bool {
		keyStr := key.(string)
		res[keyStr] = value
		return true
	})
	util.ResponseSuccess(c, res)
}

func dealWithSerialRequest(wg *sync.WaitGroup, c *gin.Context, api service.ApiStruct, parentRespMap *sync.Map, body map[string]interface{}, success chan map[string]interface{}, failed chan int) {
	defer wg.Done()

	var contentType string
	url := fmt.Sprintf("%s://%s:%s%s", api.Protocol, api.Services.Host, api.Services.Port, api.Path)

	header := make(map[string]string)
	for _, vs := range api.Headers {
		requestHeader := c.GetHeader(vs.Key)
		if requestHeader == "" {
			requestHeader = vs.DefaultValue
		}
		if vs.Key == "Content-Type" {
			contentType = vs.DefaultValue
		}
		header[vs.Key] = requestHeader
	}
	if contentType == "" {
		contentType = c.ContentType()
	}

	newBody := make(map[string]interface{})
	query := make(map[string]string)
	for _, parameter := range api.Parameter {
		if parameter.Type == "query" {
			value := c.Query(parameter.Key)
			if value == "" && parameter.Require {
				value = parameter.DefaultValue
			}
			query[parameter.Key] = value
		} else {
			newBody = parseBody(parameter.Body, body)
		}
	}

	if api.ParentApi != nil {
		for _, parentApi := range api.ParentApi {
			parentResp, exist := parentRespMap.Load(parentApi.ParentName)
			if !exist {
				setDefaultValue(parentApi, query, newBody)
			}
			parent, convertErr := parentResp.(map[string]interface{})
			if !convertErr {
				util.ResponseError(c, 500, constant.REQUSET_FAILED, "parse body failed")
				failed <- 1
				return
			}
			keyList := strings.Split(parentApi.Key, ".")
			v := keyExist(parent, keyList, 0, len(keyList))
			if v == nil {
				setDefaultValue(parentApi, query, newBody)
			} else {
				SetValue(parentApi, query, newBody, v)
			}
		}
	}

	bodyByte, err := json.Marshal(newBody)
	if err != nil {
		util.ResponseError(c, 500, constant.REQUSET_FAILED, "parse body failed")
		failed <- 1
		return
	}

	subResp, code := util.Do(api.Method, url, header, query, bodyByte, contentType)
	subMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(subResp), &subMap)

	if err != nil {
		util.ResponseError(c, 500, constant.REQUSET_FAILED, "parse body failed")
		failed <- 1
		return
	}
	if code != 200 {
		util.ResponseError(c, code, constant.REQUSET_FAILED, subMap)
		failed <- 1
		return
	}
	success <- subMap

}

func dealWithParallelRequest(wg *sync.WaitGroup, c *gin.Context, api service.ApiStruct, body map[string]interface{}, success chan map[string]interface{}, failed chan int) {
	defer wg.Done()

	var contentType string
	url := fmt.Sprintf("%s://%s:%s%s", api.Protocol, api.Services.Host, api.Services.Port, api.Path)

	header := make(map[string]string)
	for _, vs := range api.Headers {
		requestHeader := c.GetHeader(vs.Key)
		if requestHeader == "" {
			requestHeader = vs.DefaultValue
		}
		if vs.Key == "Content-Type" {
			contentType = vs.DefaultValue
		}
		header[vs.Key] = requestHeader
	}
	if contentType == "" {
		contentType = c.ContentType()
	}
	newBody := make(map[string]interface{})
	query := make(map[string]string)
	for _, parameter := range api.Parameter {
		if parameter.Type == "query" {
			value := c.Query(parameter.Key)
			if value == "" && parameter.Require {
				value = parameter.DefaultValue
			}
			query[parameter.Key] = value
		} else {
			newBody = parseBody(parameter.Body, body)
		}
	}
	bodyByte, err := json.Marshal(newBody)
	if err != nil {
		util.ResponseError(c, 500, constant.REQUSET_FAILED, "parse body failed")
		failed <- 1
		return
	}
	subResp, code := util.Do(api.Method, url, header, query, bodyByte, contentType)
	subMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(subResp), &subMap)

	if err != nil {
		util.ResponseError(c, 500, constant.REQUSET_FAILED, "parse body failed")
		failed <- 1
		return
	}
	if code != 200 {
		util.ResponseError(c, code, constant.REQUSET_FAILED, subMap)
		failed <- 1
		return
	}
	success <- subMap
}

func parseBody(parameterBody, requestBody map[string]interface{}) map[string]interface{} {
	for key, _ := range parameterBody {
		parameterBody[key] = requestBody[key]
	}
	return parameterBody
}

func SetValue(parentApi service.ParentApiStruct, query map[string]string, body map[string]interface{}, value interface{}) {
	if parentApi.ToType == "query" {
		query[parentApi.ToKey] = util.Strval(value)
	} else {
		keyList := strings.Split(parentApi.ToKey, ".")
		total := len(keyList)
		buildBody(body, keyList, 0, total, value)
	}
}

func setDefaultValue(parentApi service.ParentApiStruct, query map[string]string, body map[string]interface{}) {
	SetValue(parentApi, query, body, parentApi.DefaultValue)
}

func buildBody(body map[string]interface{}, keyList []string, depth, total int, value interface{}) {
	if depth == total-1 {
		body[keyList[depth]] = value
		return
	}
	mapValue := body[keyList[depth]]
	if mapValue == nil {
		body[keyList[depth]] = make(map[string]interface{})
	}
	mapValue = body[keyList[depth]]
	valueType := util.GetValueType(mapValue)
	if valueType == "map[string]interface {}" {
		valueMap, err := mapValue.(map[string]interface{})
		if !err {
			return
		}
		buildBody(valueMap, keyList, depth+1, total, nil)
	}
	return
}

func keyExist(body map[string]interface{}, keyList []string, depth, total int) interface{} {
	if depth == total-1 {
		return body[keyList[depth]]
	}
	key := keyList[depth]
	value := body[key]
	valuemap, err := value.(map[string]interface{})
	if !err {
		return nil
	}
	return keyExist(valuemap, keyList, depth+1, total)
}
