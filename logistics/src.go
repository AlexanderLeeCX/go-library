/**
 * @Author: Lee
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/10/16 6:21 下午
 */

package logistics

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// https://market.aliyun.com/products/57126001/cmapi010996.html?spm=5176.2020520132.101.2.6d0c7218sIw5be#sku=yuncode499600008
const (
	// 查询物流信息
	expApi = "https://ali-deliver.showapi.com/showapi_expInfo"
	// 查询物流公司
	expressApi = "https://ali-deliver.showapi.com/showapi_expressList?maxSize=9999999"
)

const (
	appCode = "0bba53f3bc734d3baf4a7f3260673437"
)

const (
	SuccessCode = 0
)

type Logistics struct {
	appCode string
}

func NewLogistics(appCode string) (logistics *Logistics, err error) {
	if len(appCode) == 0 {
		err = errors.New("params error")
		return
	}
	return &Logistics{
		appCode: appCode,
	}, nil
}

// Express 物流公司
type Express struct {
	ImgUrl     string `json:"imgUrl"`     //快递公司图标
	SimpleName string `json:"simpleName"` //拼音或英文简称
	Phone      string `json:"phone"`      //官方电话
	ExpName    string `json:"expName"`    //快递名称
	Url        string `json:"url"`        //官方网址
	Note       string `json:"note"`       //备注信息
}

// ExpInfo 物流信息
type ExpInfo struct {
	Context string `json:"context"`
	Time    string `json:"time"`
}

type ExpressListResponse struct {
	ShowapiResCode  int    `json:"showapi_res_code"`
	ShowapiResError string `json:"showapi_res_error"`
	ShowapiResBody  struct {
		RetCode     int       `json:"ret_code"`
		Flag        bool      `json:"flag"`
		ExpressList []Express `json:"expressList"`
	} `json:"showapi_res_body"`
}

type ExpInfoResponse struct {
	ShowapiResCode  int    `json:"showapi_res_code"`
	ShowapiResError string `json:"showapi_res_error"`
	ShowapiResBody  struct {
		Update          int64     `json:"update"`          // 更新时间戳
		UpgradeInfo     string    `json:"upgrade_info"`    // 提示信息，用于提醒用户可能出现的情况
		Status          int       `json:"status"`          // 快递状态 1 暂无记录 2 在途中 3 派送中 4 已签收 (完结状态) 5 用户拒签 6 疑难件 7 无效单 (完结状态) 8 超时单 9 签收失败 10 退回
		ExpSpellName    string    `json:"expSpellName"`    // 快递编码
		Msg             string    `json:"msg"`             // 返回提示信息
		UpdateStr       string    `json:"updateStr"`       // 更新时间戳字符串
		PossibleExpList []string  `json:"possibleExpList"` // 自动识别结果
		Tel             string    `json:"tel"`             // 快递公司联系方式
		Logo            string    `json:"logo"`            // 快递公司logo
		ExpTextName     string    `json:"expTextName"`     // 快递简称
		MailNo          string    `json:"mailNo"`          // 快递单号
		RetCode         int       `json:"ret_code"`
		Flag            bool      `json:"flag"`
		ExpInfoList     []ExpInfo `json:"data"` // 在途跟踪数据Ò
	} `json:"showapi_res_body"`
}

// GetExpressList 获取物流公司列表
func (l *Logistics) GetExpressList() (expressList []Express, reqBody string, err error) {
	var (
		client          = &http.Client{}
		request         *http.Request
		body            []byte
		expressListResp ExpressListResponse
	)
	request, err = http.NewRequest("GET", expressApi, nil)
	if err != nil {
		return
	}
	request.Header.Add("Authorization", fmt.Sprintf("APPCODE %s", l.appCode))
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	reqBody = string(body)
	err = json.Unmarshal(body, &expressListResp)
	if err != nil {
		return
	}

	if response.StatusCode == http.StatusOK && expressListResp.ShowapiResCode == SuccessCode {
		expressList = expressListResp.ShowapiResBody.ExpressList
		return
	}
	// 接口报错异常判断
	errMsg := response.Header.Get("X-Ca-Error-Message")
	if response.StatusCode == http.StatusInternalServerError {
		err = errors.New("ali-deliver.showapi.com 查询物流公司接口网关错误")
		return
	} else if len(errMsg) > 0 {
		err = errors.New(fmt.Sprintf("ali-deliver.showapi.com 查询物流公司接口错误 %s", errMsg))
		return
	} else {
		err = errors.New("ali-deliver.showapi.com 查询物流公司接口其它错误")
		return
	}
}

// GetExpInfo 获取物流信息
func (l *Logistics) GetExpInfo(nu string, com string, phone string) (expInfoList []ExpInfo, reqBody string, err error) {
	var (
		client          = &http.Client{}
		request         *http.Request
		body            []byte
		expInfoListResp ExpInfoResponse
	)
	if len(nu) == 0 {
		err = errors.New("params error")
		return
	}
	// 若不指定物流公司，则自动获取
	if len(com) == 0 {
		com = "auto"
	}
	// 顺丰快速，手机号后四位为必填
	if (com == "shunfengen" || com == "shunfeng" || com == "nsf") && len(phone) == 0 {
		err = errors.New("params error")
		return
	}
	uri, err := url.Parse(expApi)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 填充参数
	params := make(url.Values)
	params["com"] = []string{com}
	params["nu"] = []string{nu}
	params["phone"] = []string{phone}
	uri.RawQuery = params.Encode()
	fmt.Println(uri.String())
	request, err = http.NewRequest("GET", uri.String(), nil)
	if err != nil {
		return
	}
	request.Header.Add("Authorization", fmt.Sprintf("APPCODE %s", l.appCode))
	response, err := client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	reqBody = string(body)
	err = json.Unmarshal(body, &expInfoListResp)
	if err != nil {
		return
	}

	if response.StatusCode == http.StatusOK && expInfoListResp.ShowapiResCode == SuccessCode {
		expInfoList = expInfoListResp.ShowapiResBody.ExpInfoList
		return
	}
	// 接口报错异常判断
	errMsg := response.Header.Get("X-Ca-Error-Message")
	if response.StatusCode == http.StatusInternalServerError {
		err = errors.New("ali-deliver.showapi.com 查询物流公司接口网关错误")
		return
	} else if len(errMsg) > 0 {
		err = errors.New(fmt.Sprintf("ali-deliver.showapi.com 查询物流公司接口错误 %s", errMsg))
		return
	} else {
		err = errors.New("ali-deliver.showapi.com 查询物流公司接口其它错误")
		return
	}
}
