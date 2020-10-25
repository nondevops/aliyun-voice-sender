package cron

import (
	"encoding/json"
	"fmt"
	"github.com/n9e/alidyvms-sender/alidyvms"
	"math/rand"
	"strings"
	"time"

	"github.com/n9e/alidyvms-sender/config"
	"github.com/n9e/alidyvms-sender/dataobj"
	"github.com/n9e/alidyvms-sender/redisc"
	"github.com/toolkits/pkg/logger"
)

var semaphore chan int

func SendDyvms() {
	c := config.Get()

	semaphore = make(chan int, c.Consumer.Worker)
	for {
		messages := redisc.Pop(1, c.Consumer.Queue)
		if len(messages) == 0 {
			time.Sleep(time.Duration(300) * time.Millisecond)
			continue
		}

		sendDyvms(messages, c.MaxDelayTime)
	}
}

func sendDyvms(messages []*dataobj.Message, maxDelayTime int) {
	for _, message := range messages {
		semaphore <- 1
		// 检查当前消息事件是否超过延迟阈值
		intervalSeconds := time.Now().Unix() - message.Event.Etime
		if intervalSeconds > int64(maxDelayTime) {
			continue
		}
		go sendVoice(message)
	}
}

func sendVoice(message *dataobj.Message) {
	defer func() {
		<-semaphore
	}()

	logger.Info("<-- hashid: %v -->", message.Event.HashId)
	logger.Infof("hashid: %d: endpoint: %s, metric: %s, tags: %s", message.Event.HashId, message.ReadableEndpoint, strings.Join(message.Metrics, ","), message.ReadableTags)

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

		ttsParamStr, err := json.Marshal(alidyvms.TtsParam{
			Host:     strings.Join(hosts, ","),
			Alertmsg: message.Event.Value,
			Level:    fmt.Sprintf("P%d%s", message.Event.Priority, ET[message.Event.EventType]),
		})

		if err != nil {
			logger.Errorf("json.Marshal TtsParam failed, error:%v", err)
			continue
		}
		perm := rand.Perm(len(alidyvms.DyvmsConf.CalledShowNumber))
		for i := range perm {
			response, err := alidyvms.Send(alidyvms.DyvmsConf.CalledShowNumber[perm[i]], message.Tos[i], string(ttsParamStr))

			if err != nil {
				logger.Errorf("Send voice failed, error:%v", err)
				continue
			}
			if response.Code != "OK" {
				logger.Errorf("send to %s fail,RequestId:%v, CallId:%v,Message:%v", message.Tos[i], response.CallId, response.RequestId, response.Message)
				continue
			}
			logger.Infof("send to %s succ", message.Tos[i])
			break
		}

	}

	logger.Info("<-- hashid: %v -->", message.Event.HashId)
}

var ET = map[string]string{
	"alert":    "告警",
	"recovery": "恢复",
}
