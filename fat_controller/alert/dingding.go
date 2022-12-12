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
	"time"
)

var _ Alert = (*DingDing)(nil)

type DingDing struct {
	webhook   string
	sign      string
	atMobiles []string
}

func NewDingDing(webhook string, sign string, atMobiles []string) *DingDing {
	d := &DingDing{
		webhook:   webhook,
		sign:      sign,
		atMobiles: atMobiles,
	}

	return d
}

func (d *DingDing) GetName() string {
	return sodor.AlertPluginName_APN_DingDing.String()
}

func (d *DingDing) GiveAlarm(content string) error {
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

	msg := content + "\n"

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

	if len(d.atMobiles) > 0 {
		atMap := make(map[string][]string)
		for _, m := range d.atMobiles {
			msg += "@" + m
		}
		atMap["atMobiles"] = d.atMobiles
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
