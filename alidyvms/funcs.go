package alidyvms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
)

func Send(calledShowNumber string, calledNumber string, ttsParam string) (*dyvmsapi.SingleCallByTtsResponse, error) {
	request := dyvmsapi.CreateSingleCallByTtsRequest()
	request.Scheme = "https"
	request.CalledShowNumber = calledShowNumber
	request.CalledNumber = calledNumber
	request.TtsCode = DyvmsConf.TtsCode
	request.TtsParam = ttsParam
	request.PlayTimes = requests.NewInteger(1)
	return DyvmsClient.SingleCallByTts(request)
}
