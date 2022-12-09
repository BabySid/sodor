package alert

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BabySid/proto/sodor"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var _ alert = (*dingDing)(nil)

type dingDing struct {
	webhook string
	sign    string
}

func NewDingDing() *dingDing {
	return &dingDing{}
}

func (d *dingDing) GetName() string {
	return sodor.AlertPluginName_APN_DingDing.String()
}

const (
	DingDingText      = "text"
	DingDingAtMobiles = "atMobiles"
)

func (d *dingDing) GetParams() []*sodor.AlertPluginParam {
	return []*sodor.AlertPluginParam{
		{
			Name:   DingDingText,
			Type:   sodor.AlertPluginParamType_APPT_String,
			DescEn: "content that sent to dingDing",
		},
		{
			Name:   DingDingAtMobiles,
			Type:   sodor.AlertPluginParamType_APPT_String,
			DescEn: "mobiles that want to be @. format is 1;2",
		},
	}
}

func (d *dingDing) GiveAlarm(param map[string]interface{}) error {
	// https://open.dingtalk.com/document/robots/custom-robot-access
	url := d.webhook
	if d.sign != "" {
		timestamp := time.Now().UnixMilli()
		stringToSign := fmt.Sprintf("%d\n%s", timestamp, d.sign)
		hash := hmac.New(sha256.New, []byte(d.sign))
		hash.Write([]byte(stringToSign))
		sign := base64.StdEncoding.EncodeToString(hash.Sum(nil))
		url = fmt.Sprintf("%s&timestamp=%d&sign=%s", d.webhook, timestamp, sign)
	}

	msg := "empty message"
	if v, ok := param[DingDingText]; ok {
		msg = fmt.Sprintf("%v", v)
	}
	msg += "\n"

	/*
		{
		     "msgtype": "markdown",
		     "markdown": {
		         "title":"杭州天气",
		         "text": "#### 杭州天气 @150XXXXXXXX \n > 9度，西北风1级，空气良89，相对温度73%\n"
		     },
		      "at": {
		          "atMobiles": [
		              "150XXXXXXXX"
		          ],
		          "atUserIds": [
		              "user123"
		          ],
		          "isAtAll": false
		      }
		 }
	*/
	request := make(map[string]interface{})

	if v, ok := param[DingDingAtMobiles]; ok {
		atMap := make(map[string][]string)
		vs := strings.Split(fmt.Sprintf("%v", v), ";")
		for _, m := range vs {
			msg += "@" + m
		}
		atMap["atMobiles"] = vs
		request["at"] = atMap
	}

	request["msgtype"] = "markdown"

	request["markdown"] = map[string]string{
		"title": "Alert",
		"text":  msg,
	}

	bs, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bs, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(bs))
	}

	return nil
}
