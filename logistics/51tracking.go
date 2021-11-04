/**
 * @Author: Lee
 * @Description:
 * @File:  51tracking
 * @Version: 1.0.0
 * @Date: 2021/10/23 5:28 下午
 */

package logistics

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// https://www.51tracking.com/v3/api-index?language=Golang#api%E7%89%88%E6%9C%AC
const (
	baseUrl    = "https://api.51tracking.com"
	apiKey     = ""
	apiVersion = "v3"
)

// 物流状态
const (
	statusPending     = "pending"     // 查询中：新增包裹正在查询中，请等待
	statusNotfound    = "notfound"    // 查询不到：包裹信息目前查询不到
	statusTransit     = "transit"     // 运输途中：物流商已揽件，包裹正被发往目的地
	statusPickup      = "pickup"      // 到达待取：包裹正在派送中，或到达当地收发点
	statusDelivered   = "delivered"   // 成功签收：包裹已被成功投递
	statusExpired     = "expired"     // 运输过久：包裹在很长时间内都未投递成功。快递包裹超过30天、邮政包裹超过60天未投递成功，该查询会被识别为此状态
	statusUndelivered = "undelivered" // 投递失败：快递员投递失败（通常会留有通知并再次尝试投递）
	statusException   = "exception"   // 可能异常：包裹退回、包裹丢失、清关失败等异常情况
)

// 物流子状态
const (
	// 查询不到子状态
	subStatusNotfound001 = "notfound001" // 包裹正在等待被揽件
	subStatusNotfound002 = "notfound002" // 无任何物流信息

	// 运输途中子状态
	subStatusTransit001 = "transit001" // 包裹正在运输途中
	subStatusTransit002 = "transit002" // 包裹到达分拣中心
	subStatusTransit003 = "transit003" // 包裹到达目的网点
	subStatusTransit004 = "transit004" // 包裹达到目的国家
	subStatusTransit005 = "transit005" // 包裹清关完成
	subStatusTransit006 = "transit006" // 包裹已经封装好，即将送出
	subStatusTransit007 = "transit007" // 包裹已经交给航空公司

	// 到达待取子状态
	subStatusPickup001 = "pickup001" // 包裹正在派送中或即将派送
	subStatusPickup002 = "pickup002" // 包裹到达代收点等待自取
	subStatusPickup003 = "pickup003" // 派送中，物流商已经与收件人联系过至少一次

	// 成功签收子状态
	subStatusDelivered001 = "delivered001" // 包裹投递成功
	subStatusDelivered002 = "delivered002" // 客户自提成功
	subStatusDelivered003 = "delivered003" // 包裹由客户签字签收
	subStatusDelivered004 = "delivered004" // 包裹放至快递柜，或由物业、门卫代收

	// 投递失败子状态
	subStatusUndelivered001 = "undelivered001" // 包裹由于地址问题投递失败
	subStatusUndelivered002 = "undelivered002" // 派送时无人在家
	subStatusUndelivered003 = "undelivered003" // 无法联系到收件人
	subStatusUndelivered004 = "undelivered004" // 由于其他原因导致派送失败

	// 可能异常子状态
	subStatusException004 = "exception004" // 包裹无人认领
	subStatusException005 = "exception005" // 其他异常情况
	subStatusException006 = "exception006" // 包裹被海关扣留
	subStatusException007 = "exception007" // 包裹破损、遗失或被丢弃
	subStatusException008 = "exception008" // 订单在揽件前被取消
	subStatusException009 = "exception009" // 收件人拒收包裹
	subStatusException010 = "exception010" // 包裹已退回发件人处
	subStatusException011 = "exception011" // 包裹正在被退回发件人的途中
)

// 接口响应状态码
const (
	responseCode_200 = 200 // 请求成功
	responseCode_203 = 203 // API服务只提供给付费账户,请付费购买单号以解锁API服务
	respondeCode_204 = 204 // 请求成功，但未获取到数据,可能是该单号、所查询目标数据不存在
	responseCode_400 = 400 // 请求类型错误,请查询API文档，确定该接口的调用方法（POST、GET等）
	responseCode_401 = 401 // 授权失败或没有权限,请检查并确保你API Key正确无误
	responseCode_403 = 403 // 该页面不存在,请检查并确保你的访问链接正确无误
	responseCode_404 = 404 // 该页面不存在,请检查并确保你的访问链接正确无误
	responseCode_408 = 408 // 请求超时,官网没有返回数据，请稍后再试
	responseCode_411 = 411 // 请求参数长度超过限制,	请检查并确保请求参数长度符合要求
	responseCode_412 = 412 // 请求参数格式不合要求, 请检查并确保请求参数格式符合要求
	responseCode_413 = 413 // 请求参数数量超过限制, 请查看API文档以获取该接口请求数量限制
	responseCode_417 = 417 // 缺少请求参数或者请求参数无法解析,请检查并确保请求参数完整、格式正确
	responseCode_421 = 421 // 部分必填参数为空, 部分物流商需上传特殊参数才可查询物流信息（特殊物流商）
	responseCode_422 = 422 // 物流商简码无法识别或者不支持该物流商,请检查并确保物流商简码正确（物流商简码）
	responseCode_423 = 423 // 跟踪单号已存在，无需再次创建
	responseCode_424 = 424 // 跟踪单号不存在,	请先调用「添加物流单号」接口创建单号
	responseCode_429 = 429 // API请求频率次限制，请稍后再试,	请查看API文档以获取该接口请求频率限制
	responseCode_511 = 511 // 系统错误, 请联系我们：service@51tracking.org
	responseCode_512 = 512 // 系统错误, 请联系我们：service@51tracking.org
	responseCode_513 = 513 // 系统错误, 请联系我们：service@51tracking.org
)

type BaseResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func doRequest(url string, data interface{}, method string) (content string, err error) {
	jsonStr, err := json.Marshal(data)
	method = strings.ToUpper(method)
	requestUrl := baseUrl + "/" + apiVersion + "/trackings/" + url
	req, err := http.NewRequest(method, requestUrl, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Tracking-Api-Key", apiKey)
	defer req.Body.Close()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err2 := client.Do(req)
	if err2 != nil {
		panic(err2)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return content, nil
}

// Carrier 物流商
type Carrier struct {
	CourierName string `json:"courier_name"` // 物流商的名称
	CourierCode string `json:"courier_code"` // 物流商对应的唯一简码
	CountryCode string `json:"country_code"` // 该偏远地区所对应的国家名称
	CourierUrl  string `json:"courier_url"`  // 物流商官网的链接
	CourierType string `json:"courier_type"` // 运输类型（快递或邮政）
	CourierLogo string `json:"courier_logo"` // 物流商 logo 的链接
}

// CarriersResponse 物流商接口响应
type CarriersResponse struct {
	BaseResponse
	Data []Carrier `json:"data"`
}

// GetCarriers 获取物流商列表
func GetCarriers(lang string) (carrierList []Carrier, err error) {
	var response CarriersResponse
	result, err := doRequest(fmt.Sprintf("carriers?lang=%s", lang), nil, "GET")
	if err != nil {
		return
	}
	dec := json.NewDecoder(bytes.NewBuffer([]byte(result)))
	dec.UseNumber() //关键步骤
	err = dec.Decode(&response)
	if err != nil {
		return
	}
	if response.Code != responseCode_200 {
		err = errors.New(response.Message)
		return
	}
	carrierList = response.Data
	return
}

type ModifyCourierResponse struct {
	BaseResponse
	Data struct {
		TrackingNumber string `json:"tracking_number"`
		CourierCode    string `json:"courier_code"`
	} `json:"data"`
}

// UpdateCourierCode 修改指定订单的物流简码
func UpdateCourierCode(trackingNumber string, courierCode string, newCourierCode string) error {
	var response ModifyCourierResponse
	params := map[string]string{
		"tracking_number":  trackingNumber,
		"courier_code":     courierCode,
		"new_courier_code": newCourierCode,
	}
	result, err := doRequest("modifycourier", params, "PUT")
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewBuffer([]byte(result)))
	dec.UseNumber() //关键步骤
	err = dec.Decode(&response)
	if err != nil {
		return err
	}
	if response.Code != responseCode_200 {
		err = errors.New(response.Message)
		return err
	}
	return nil
}

type TrackingInfo struct {
	TrackingNumber          string `json:"tracking_number"`           // 包裹物流单号
	CourierCode             string `json:"courier_code"`              // 物流商对应的唯一简码
	DestinationCode         string `json:"destination_code"`          // 目的国的二字简码
	Title                   string `json:"title"`                     // 包裹名称
	LogisticsChannel        string `json:"logistics_channel"`         // 自定义字段，用于填写物流渠道（比如某货代商）
	CustomerName            string `json:"customer_name"`             // 客户姓名
	CustomerEmail           string `json:"customer_email"`            // 客户邮箱
	CustomerPhone           string `json:"customer_phone"`            // 顾客接收短信的手机号码。手机号码的格式应该为：“+区号手机号码”（例子：+8612345678910）
	OrderNumber             string `json:"order_number"`              // 包裹的订单号，由商家/平台所产生的订单编号
	ShippingDate            string `json:"shipping_date"`             // 包裹发货时间（例子：2020-09-17 16:51）
	TrackingShippingDate    string `json:"tracking_shipping_date"`    // 包裹的发货时间，其格式为：YYYYMMDD，有部分的物流商（如 deutsch-post）需要这个参数（例子：20200102）
	TrackingPostalCode      string `json:"tracking_postal_code"`      // 收件人所在地邮编，仅有部分的物流商（如 postnl-3s）需要这个参数
	TrackingDestinationCode string `json:"tracking_destination_code"` // 目的国对应的二字简码，部分物流商（如postnl-3s）需要这个参数
	TrackingCourierAccount  string `json:"tracking_courier_account"`  // 物流商的官方账号，仅有部分的物流商（如 dynamic-logistics）需要这个参数
	Lang                    string `json:"lang"`                      // 查询结果的语言（例子：cn, en），若未指定该参数，结果会以英文或中文呈现。注意：只有物流商支持多语言查询结果时，该指定才会生效
	Note                    string `json:"note"`
}

type TrackingInfoResponse struct {
	BaseResponse
	Data struct {
		Success []struct {
			TrackingNumber string `json:"tracking_number"`
			CourierCode    string `json:"courier_code"`
			OrderNumber    string `json:"order_number"`
		} `json:"success"`
		Error []struct{} `json:"error"`
	} `json:"data"`
}

// CreateTracking 添加物流单号
func CreateTracking(trackingList []TrackingInfo) (trackingNumberList []string, err error) {
	var response TrackingInfoResponse
	if len(trackingList) > 40 {
		err = errors.New("tracking info length beyond 40")
		return
	}
	result, err := doRequest("create", trackingNumberList, "POST")
	if err != nil {
		return
	}
	dec := json.NewDecoder(bytes.NewBuffer([]byte(result)))
	dec.UseNumber() //关键步骤
	err = dec.Decode(&response)
	if response.Code != responseCode_200 {
		err = errors.New(response.Message)
		return
	}
	for _, tracking := range trackingList {
		trackingNumberList = append(trackingNumberList, tracking.TrackingNumber)
	}
	return
}

type TrackingResult struct {
	TrackingNumber         string     `json:"tracking_number"`   // 包裹物流单号
	CourierCode            string     `json:"courier_code"`      // 物流商对应的唯一简码
	OrderNumber            string     `json:"order_number"`      // 包裹的订单号，由商家/平台所产生的订单编号
	DeliveryStatus         string     `json:"delivery_status"`   // 被人为指定的快递状态，该参数的值可被设置为4-已送达、7-异常、8-原始状态
	Archived               bool       `json:"archived"`          // “true”表示该单号已被归档，“false”表示该单号处于未归档状态
	Updating               bool       `json:"updating"`          // “true”表示该单号会被继续更新，“false”表示该单号已停止更新
	CreatedAt              string     `json:"created_at"`        // 创建查询的时间
	UpdateDate             string     `json:"update_date"`       // 系统最后更新查询的时间
	ShippingDate           string     `json:"shipping_date"`     // 包裹发货时间（例子：2020-09-17 16:51）
	CustomerName           string     `json:"customer_name"`     // 客户姓名
	CustomerEmail          string     `json:"customer_email"`    // 客户邮箱
	CustomerPhone          string     `json:"customer_phone"`    // 顾客接收短信的手机号码
	Title                  string     `json:"title"`             // 包裹名称
	LogisticsChannel       string     `json:"logistics_channel"` // 自定义字段，用于填写物流渠道（比如某货代商）
	Note                   string     `json:"note"`              // 备注，可自定义
	Destination            string     `json:"destination"`       // 目的国的二字简码
	Original               string     `json:"original"`          // 发件国的名称
	ServiceCode            string     `json:"service_code"`      // 快递服务类型，比如次日达（部分物流商返回）
	Weight                 string     `json:"weight"`            // 该货物的重量（多个包裹会被打包成一个“货物”）
	Substatus              string     `json:"substatus"`
	StatusInfo             string     `json:"status_info"` // 最新的一条物流信息
	Previously             string     `json:"previously"`
	DestinationTrackNumber string     `json:"destination_track_number"` // 该包裹在目的国的物流单号
	ExchangeNumber         string     `json:"exchangeNumber"`           // 该包裹在中转站的物流商单号
	Consignee              string     `json:"consignee"`                // 签收人
	ScheduledDeliveryDate  string     `json:"scheduled_delivery_date"`
	ScheduledAddress       string     `json:"Scheduled_Address"`
	LatestEvent            string     `json:"latest_event"`            // 最新物流信息的梗概，包括以下信息：状态、地址、时间
	LastestCheckpointTime  string     `json:"lastest_checkpoint_time"` // 最新物流信息的更新时间
	TransitTime            int        `json:"transit_time"`            // 包裹的从被揽收至被送达的时长（天）
	StayTime               int        `json:"stay_time"`               // 物流信息未更新的时长（单位：天），由当前时间减去物流信息最近更新时间得到
	OriginInfo             OriginInfo `json:"origin_info"`             // 发件国的物流信息
	DestinationInfo        OriginInfo `json:"destination_info"`        // 目的国的物流信息
}

type OriginInfo struct {
	CourierCode            string    `json:"courier_code"`
	CourierPhone           string    `json:"courier_phone"`            // 物流商官网上的电话
	Weblink                string    `json:"weblink"`                  // 物流商的官网的链接
	ReferenceNumber        string    `json:"reference_number"`         // 包裹对应的另一个单号，作用与当前单号相同（仅有少部分物流商提供）
	ReceivedDate           string    `json:"received_date"`            // 物流商接收包裹的时间（也称为上网时间）
	DispatchedDate         string    `json:"dispatched_date"`          // 包裹封发时间，封发指将多个小包裹打包成一个货物（方便运输）
	DepartedAirportDate    string    `json:"departed_airport_date"`    // 包裹离开此出发机场的时间
	ArrivedAbroadDate      string    `json:"arrived_abroad_date"`      // 包裹达到目的国的时间
	CustomsReceivedDate    string    `json:"customs_received_date"`    // 包裹移交给海关的时间
	ArrivedDestinationDate string    `json:"arrived_destination_date"` // 包裹达到目的国、目的城市的时间
	Trackinfo              TrackInfo `json:"trackinfo"`                // 详细物流信息
}

type TrackInfo struct {
	CheckpointDate              string `json:"checkpoint_date"`               // 本条物流信息的更新时间，由物流商提供（包裹被扫描时，物流信息会被更新）
	TrackingDetail              string `json:"tracking_detail"`               // 具体的物流情况
	Location                    string `json:"location"`                      // 物流信息更新的地址（该包裹被扫描时，所在的地址）
	CheckpointDeliveryStatus    string `json:"checkpoint_delivery_status"`    // 根据具体物流情况所识别出来的物流状态
	CheckpointDeliverySubstatus string `json:"checkpoint_delivery_substatus"` // 物流状态的子状态
}

type TrackingResultResponse struct {
	BaseResponse
	Data []TrackingResult `json:"data"`
}

// GetTrackingResult 获取物流查询结果
func GetTrackingResult(trackingNumberList []string) (err error) {
	if len(trackingNumberList) > 40 {
		err = errors.New("tracking number length beyond 40")
		return
	}
	var response TrackingResultResponse
	params := strings.Join(trackingNumberList, ",")
	result, err := doRequest("get?tracking_numbers="+params, nil, "GET")
	dec := json.NewDecoder(bytes.NewBuffer([]byte(result)))
	dec.UseNumber() //关键步骤
	err = dec.Decode(&response)
	if response.Code != responseCode_200 {
		err = errors.New(response.Message)
		return
	}
	return
}

// UpdateTrackingInfo 修改单号信息
func UpdateTrackingInfo(trackingList []TrackingInfo) (err error) {
	if len(trackingList) > 40 {
		err = errors.New("tracking info length beyond 40")
		return
	}
	var response TrackingInfoResponse
	result, err := doRequest("modifyinfo", trackingList, "PUT")
	if err != nil {
		return
	}
	dec := json.NewDecoder(bytes.NewBuffer([]byte(result)))
	dec.UseNumber() //关键步骤
	err = dec.Decode(&response)
	if response.Code != responseCode_200 {
		err = errors.New(response.Message)
		return
	}
	return
}

type DeleteTrackingInfoRequest struct {
	TrackingNumber string `json:"tracking_number"`
	CourierCode    string `json:"courier_code"`
}

type DeleteTrackingInfoResponse struct {
	BaseResponse
	Data struct {
		Success []struct {
			TrackingNumber string `json:"tracking_number"`
			CourierCode    string `json:"courier_code"`
		} `json:"success"`
		Error []struct{} `json:"error"`
	} `json:"data"`
}

// DeleteTrackingInfo 删除查询单号
func DeleteTrackingInfo(trackingList []DeleteTrackingInfoRequest) (err error) {
	if len(trackingList) > 40 {
		err = errors.New("tracking info length beyond 40")
		return
	}
	var response DeleteTrackingInfoResponse
	result, err := doRequest("delete", trackingList, "DELETE")
	dec := json.NewDecoder(bytes.NewBuffer([]byte(result)))
	dec.UseNumber() //关键步骤
	err = dec.Decode(&response)
	if response.Code != responseCode_200 {
		err = errors.New(response.Message)
		return
	}
	return
}

type NotUpdateTrackingInfoRequest struct {
	TrackingNumber string `json:"tracking_number"`
	CourierCode    string `json:"courier_code"`
}

type NotUpdateTrackingInfoResponse struct {
	BaseResponse
	Data struct {
		Success []struct {
			TrackingNumber string `json:"tracking_number"`
			CourierCode    string `json:"courier_code"`
		} `json:"success"`
		Error []struct{} `json:"error"`
	} `json:"data"`
}

// NotUpdateTrackingInfo 停止单号更新
func NotUpdateTrackingInfo(trackingList []NotUpdateTrackingInfoRequest) (err error) {
	if len(trackingList) > 40 {
		err = errors.New("tracking info length beyond 40")
		return
	}
	var response NotUpdateTrackingInfoResponse
	result, err := doRequest("notupdate", trackingList, "POST")
	dec := json.NewDecoder(bytes.NewBuffer([]byte(result)))
	dec.UseNumber() //关键步骤
	err = dec.Decode(&response)
	if response.Code != responseCode_200 {
		err = errors.New(response.Message)
		return
	}
	return
}

type CheckRemoteRegionRequest struct {
	PostalCode  string `json:"postal_code"`  // 该偏远地区所对应的邮编
	Country     string `json:"country"`      // 城市名称或者国家二字简码（例子：CN）
	CourierCode string `json:"courier_code"` // 物流商的唯一简码，仅可填：dhl, ups, fedex, tnt
}

type CheckRemoteRegionResponse struct {
	BaseResponse
	Data struct {
		CountryCode       string `json:"country_code"`        // 该偏远地区所对应的国家名称
		PostalCode        string `json:"postal_code"`         // 该偏远地区所对应的邮编
		RemoteCourierCode string `json:"remote_courier_code"` // 支持该偏远地区的物流商。51tracking 目前仅支持DHL, UPS, Fedex, TNT四家物流商的偏远地区查询
	} `json:"data"`
}

// CheckRemoteRegion 检测偏远地区
func CheckRemoteRegion(params CheckRemoteRegionRequest) (err error) {
	var response CheckRemoteRegionResponse
	result, err := doRequest("remote", params, "POST")
	dec := json.NewDecoder(bytes.NewBuffer([]byte(result)))
	dec.UseNumber() //关键步骤
	err = dec.Decode(&response)
	if response.Code != responseCode_200 {
		err = errors.New(response.Message)
		return
	}
	return
}

type TotalTrackingStatusRequest struct {
	CourierCode     string `json:"courier_code"`      // 物流商对应的唯一简码
	CreatedDateMin  int64  `json:"created_date_min"`  // 创建查询的起始时间，时间戳格式（例子：1076599161）
	CreatedDateMax  int64  `json:"created_date_max"`  // 创建查询的结束时间，时间戳格式（例子：1076599161），与上一参数联用，可筛选出创建时间为created_date_min～created_date_max范围内的查询结果
	ShippingDateMin int64  `json:"shipping_date_min"` // 发货的起始时间，时间戳格式（例子：1076599161）
	ShippingDateMax int64  `json:"shipping_date_max"` // 发货的结束时间，时间戳格式（例子：1076599161），与上一参数联用，可以筛选出更新时间为shipping_date_min～shipping_date_max范围内的查询结果
}

type TotalTrackingStatusResponse struct {
	BaseResponse
	Data struct {
		Pending     int64 `json:"pending"`
		Notfound    int64 `json:"notfound"`
		Transit     int64 `json:"transit"`
		Pickup      int64 `json:"pickup"`
		Delivered   int64 `json:"delivered"`
		Expired     int64 `json:"expired"`
		Undelivered int64 `json:"undelivered"`
		Exception   int64 `json:"exception"`
	} `json:"data"`
}

// TotalTrackingStatus 统计包裹状态 单号量超过 200万 的用户不允许调用。
func TotalTrackingStatus(params TotalTrackingStatusRequest) (err error) {
	var response TotalTrackingStatusResponse
	url := fmt.Sprintf("status?courier_code=%s&created_date_min=%d&created_date_max=%d&shipping_date_min=%d&shipping_date_max=%d",
		params.CourierCode, params.CreatedDateMin, params.CreatedDateMax, params.ShippingDateMin, params.ShippingDateMax)
	result, err := doRequest(url, nil, "GET")
	dec := json.NewDecoder(bytes.NewBuffer([]byte(result)))
	dec.UseNumber() //关键步骤
	err = dec.Decode(&response)
	if response.Code != responseCode_200 {
		err = errors.New(response.Message)
		return
	}
	return
}

type GetTransitTimeRequest struct {
	CourierCode     string `json:"courier_code"`     // 物流商对应的唯一简码
	OriginalCode    string `json:"original_code"`    // 发件国二字简码（例子：CN）
	DestinationCode string `json:"destination_code"` // 目的国的二字简码
}

type GetTransitTimeResponse struct {
	BaseResponse
	Data struct {
		CourierCode         string  `json:"courier_code"`          // 物流商对应的唯一简码
		OriginalCode        string  `json:"original_code"`         // 发件国二字简码（例子：CN）
		DestinationCode     string  `json:"destination_code"`      // 目的国的二字简码
		Total               int64   `json:"total"`                 // 未签收的总单号数量
		Delivered           int64   `json:"delivered"`             // 已签收的总单号数量
		Range_1_7           float64 `json:"range_1_7"`             // 送达时间为0～7天的单号的占比
		Range_8_15          float64 `json:"range_8_15"`            // 送达时间为7～15天的单号的占比
		Range_16_30         float64 `json:"range_16_30"`           // 送达时间为16～30天的单号的占比
		Range_31_60         float64 `json:"range_31_60"`           // 送达时间为31～60天的单号的占比
		Range_60_up         float64 `json:"range_60_up"`           // 送达时间大于60天的单号的占比
		AverageDeliveryTime float64 `json:"average_delivery_time"` // 平均送达时间（单位：天）
	} `json:"data"`
}

// GetTransitTime 获取物流时效
func GetTransitTime(params GetTransitTimeRequest) (err error) {
	var response GetTransitTimeResponse
	result, err := doRequest("transittime", params, "POST")
	dec := json.NewDecoder(bytes.NewBuffer([]byte(result)))
	dec.UseNumber() //关键步骤
	err = dec.Decode(&response)
	if response.Code != responseCode_200 {
		err = errors.New(response.Message)
		return
	}
	return
}

type GetUserInfoResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"meta"`
	Data struct {
		Email       string  `json:"email"`        // 登录邮箱
		Regtime     string  `json:"regtime"`      // 账户注册时间
		Phone       string  `json:"phone"`        // 账户绑定的手机号码
		Money       float64 `json:"money"`        // 账户剩余的充值金额
		TrackNumber int64   `json:"track_number"` // 账户剩余的单号额度
	} `json:"data"`
}

// GetUserInfo 获取账号信息
func GetUserInfo() (err error) {
	var response GetUserInfoResponse
	result, err := doRequest("getuserinfo", nil, "GET")
	dec := json.NewDecoder(bytes.NewBuffer([]byte(result)))
	dec.UseNumber() //关键步骤
	err = dec.Decode(&response)
	if response.Meta.Code != responseCode_200 {
		err = errors.New(response.Meta.Message)
		return
	}
	return
}
