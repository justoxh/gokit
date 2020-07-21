package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func HttpGetRequest(strUrl string, mapParams map[string]string) (string, error) {
	httpClient := &http.Client{}
	var strRequestUrl string
	if len(mapParams) == 0 {
		strRequestUrl = strUrl
	} else {
		strParams := Map2UrlQuery(mapParams)
		strRequestUrl = strUrl + "?" + strParams
	}

	// 构建Request
	request, err := http.NewRequest("GET", strRequestUrl, nil)
	if nil != err {
		return err.Error(), err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	request.Close = true
	// 发出请求
	response, err := httpClient.Do(request)

	if err != nil {
		fmt.Printf(" err :%v\n response:%v\n", err, response)
		return err.Error(), err
	}
	if response != nil {
		defer response.Body.Close()
	}
	// 解析响应内容
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error(), err
	}
	return string(body), nil
}

func HttpPostRequest(strUrl string, mapParams map[string]string) (string, error) {
	httpClient := &http.Client{}

	jsonParams := ""
	if nil != mapParams {
		bytesParams, err := json.Marshal(mapParams)
		if err != nil {
			fmt.Printf("spider post request json marshal  error: %v\n", err)
			return "", nil
		}
		jsonParams = string(bytesParams)
	}

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error(), err
	}

	request.Header.Add("Content-Type", "application/json")
	response, err := httpClient.Do(request)

	if nil != err {
		return err.Error(), err
	}
	if response != nil {
		defer response.Body.Close()
	}
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error(), err
	}

	return string(body), nil
}

func Map2UrlQuery(mapParams map[string]string) string {
	var strParams string
	for key, value := range mapParams {
		param := key + "=" + value + "&"
		strParams += param
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len([]rune(strParams))-1])
	}
	return strParams
}
