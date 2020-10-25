package alidyvms

import (
	"encoding/json"
	"fmt"
	"github.com/n9e/alidyvms-sender/config"
	"github.com/n9e/alidyvms-sender/dataobj"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var ET = map[string]string{
	"alert":    "告警",
	"recovery": "恢复",
}

func TestSend(t *testing.T) {
	confPath := "E:\\gopath\\src\\github.com\\n9e\\alidyvms-sender\\etc\\alidyvms-sender.yml"
	if err := config.ParseConfig(confPath); err != nil {
		fmt.Println("cannot parse configuration file:", err)
	} else {
		fmt.Println("parse configuration file:", confPath)
	}
	InitAliDyvms()

	message := dataobj.Message{
		Tos: []string{"150xxxxxx"},
		ClaimLink: "",
		StraLink:  "http://n9e.xx.com/#/monitor/strategy/5",
		EventLink: "http://n9e.xx.com/#/monitor/history/his/91185",
		Bindings: []string{
			"127.0.0.1 - n9e-svc - ops.offline",
		},
		NotifyType: "voice",
		Metrics: []string{
			"ProcAgentAlive",
			"ProcAgentAlive1",
		},
		ReadableEndpoint: "127.0.0.1(n9e-svc)",
		ReadableTags:     "",
		IsUpgrade:        false,
		Event: &dataobj.Event{
			Id:            83584,
			Sid:           5,
			Sname:         "监控agent失联",
			NodePath:      "ops.offline",
			Endpoint:      "127.0.0.1",
			EndpointAlias: "n9e-svc",
			Priority:      1,
			EventType:     "recovery",
			Category:      1,
			Status:        0,
			HashId:        140423094364084604,
			Etime:         1590645382,
			Value:         "proc.agent.alive1",
			Info:          "proc.agent.alive (nodata,60s)",
			Created:       time.Now(),
			Detail:        "[{\"metric\":\"proc.agent.alive\",\"points\":[{\"timestamp\":1590645340,\"value\":null},{\"timestamp\":1590645360,\"value\":1.000000},{\"timestamp\":1590645380,\"value\":null}],\"extra\":\"\"}]",
			Users:         "[]",
			Groups:        "[4]",
			Nid:           7,
			NeedUpgrade:   0,
			AlertUpgrade:  "{\"users\":\"[]\",\"groups\":\"[]\",\"duration\":60,\"level\":1}",
		},
	}
	count := len(message.Tos)
	for i := 0; i < count; i++ {
		var hosts []string
		for _, v := range message.Metrics {
			var metricStr []string
			vals := strings.Split(v, ".")
			for _, val := range vals {
				metricStr = append(metricStr, strings.Title(val))
			}
			hosts = append(metricStr, strings.Join(metricStr, ""))
		}
		ttsParamStr, err := json.Marshal(TtsParam{
			Host:     strings.Join(hosts, ","),
			Alertmsg: message.Event.Value,
			Level:    fmt.Sprintf("P%d%s", message.Event.Priority, ET[message.Event.EventType]),
		})
		if err != nil {
			fmt.Printf("json.Marshal TtsParam failed, error:%v", err)
			continue
		}
		perm := rand.Perm(len(DyvmsConf.CalledShowNumber))
		for i := range perm {
			response, err := Send(DyvmsConf.CalledShowNumber[perm[i]], message.Tos[i], string(ttsParamStr))

			if err != nil {
				fmt.Printf("Send voice failed, error:%v", err)
				continue
			}
			if response.Code != "OK" {
				fmt.Printf("send to %s fail,RequestId:%v, CallId:%v,Message:%v", message.Tos[i], response.CallId, response.RequestId, response.Message)
				continue
			}
			fmt.Printf("send to %s succ", message.Tos[i])
			break
		}

	}
}
