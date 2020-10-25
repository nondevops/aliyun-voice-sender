package alidyvms

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/n9e/alidyvms-sender/config"
	"github.com/toolkits/pkg/logger"
)

// Message 消息主体参数
type TtsParam struct {
	Host     string `json:"host"`
	Alertmsg string `json:"alertmsg"`
	Level    string `json:"level"`
}

var DyvmsConf config.AlidyvmsSection
var DyvmsClient *dyvmsapi.Client

func InitAliDyvms() {
	cfg := config.Get()

	DyvmsConf = cfg.Alidyvms
	accessKey := cfg.Alidyvms.AccessKey
	secret := cfg.Alidyvms.Secret
	regionId := cfg.Alidyvms.RegionId

	client, err := dyvmsapi.NewClientWithAccessKey(regionId, accessKey, secret)
	if err != nil {
		logger.Error("init alidyvms fail: ", err)
		return
	}
	DyvmsClient = client
}
