package dingding

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const DingBaseUrl = "https://oapi.dingtalk.com/robot/send?access_token="

type Robot struct {
	Webhook string
	Secret  string
}

func NewRobot(webhook string, secret string) *Robot {
	if len(webhook) == 64 {
		webhook = DingBaseUrl + webhook
	} 
	return &Robot{Webhook: webhook, Secret: secret}
}

func (robot *Robot) toSign() (paras string) {
	timestamp := strconv.FormatInt(time.Now().Unix()*1000, 10)
	toSign := timestamp + "\n" + robot.Secret
	sha256 := hmac.New(sha256.New, []byte(robot.Secret))
	sha256.Write([]byte(toSign))
	base := base64.StdEncoding.EncodeToString(sha256.Sum(nil))
	v := url.Values{}
	v.Add("sign", base)
	v.Add("timestamp", timestamp)
	paras = v.Encode()
	return
}

func (robot *Robot) SendMessage(message interface{}) error {
	body := bytes.NewBuffer(nil)
	err := json.NewEncoder(body).Encode(message)
	if err != nil {
		return fmt.Errorf("message json failed, message: %v, err: %v", message, err.Error())
	}

	paras := robot.toSign()
	urlAddr := robot.Webhook + "&" + paras

	request, err := http.NewRequest(http.MethodPost, urlAddr, body)
	if err != nil {
		return fmt.Errorf("error request: %v", err.Error())
	}
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	res, err := (&http.Client{}).Do(request)
	if err != nil {
		return fmt.Errorf("send failed, error: %v", err.Error())
	}
	defer func() { _ = res.Body.Close() }()
	result, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return fmt.Errorf("send failed, %s", httpError(request, res, result, "http code is not 200"))
	}
	if err != nil {
		return fmt.Errorf("send failed, %s", httpError(request, res, result, err.Error()))
	}

	type response struct {
		ErrCode int `json:"errcode"`
	}
	var ret response

	if err := json.Unmarshal(result, &ret); err != nil {
		return fmt.Errorf("send failed, %s", httpError(request, res, result, err.Error()))
	}

	if ret.ErrCode != 0 {
		return fmt.Errorf("send failed, %s", httpError(request, res, result, "errcode is not 0"))
	}

	return nil
}

func httpError(request *http.Request, response *http.Response, body []byte, error string) string {
	return fmt.Sprintf(
		"http request failure, error: %s, status code: %d, %s %s, body:\n%s",
		error,
		response.StatusCode,
		request.Method,
		request.URL.String(),
		string(body),
	)
}

func (robot *Robot) SendTextMessage(content string, atMobiles []string, isAtAll bool) error {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
		"at": map[string]interface{}{
			"atMobiles": atMobiles,
			"isAtAll":   isAtAll,
		},
	}

	return robot.SendMessage(msg)
}

func (robot *Robot) SendMarkdownMessage(title string, text string, atMobiles []string, isAtAll bool) error {
	msg := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
		"at": map[string]interface{}{
			"atMobiles": atMobiles,
			"isAtAll":   isAtAll,
		},
	}

	return robot.SendMessage(msg)
}

func (robot *Robot) SendLinkMessage(title string, text string, messageUrl string, picUrl string) error {
	msg := map[string]interface{}{
		"msgtype": "link",
		"link": map[string]string{
			"title":      title,
			"text":       text,
			"messageUrl": messageUrl,
			"picUrl":     picUrl,
		},
	}

	return robot.SendMessage(msg)
}

func (robot *Robot) SendActionCard(title, text, singleTitle, singleURL, btnOrientation, hideAvatar string) error {
	msg := map[string]interface{}{
		"msgtype": "actionCard",
		"actionCard": map[string]string{
			"title":          title,
			"text":           text,
			"hideAvatar":     hideAvatar,
			"btnOrientation": btnOrientation,
			"singleTitle":    singleTitle,
			"singleURL":      singleURL,
		},
	}
	return robot.SendMessage(msg)
}
