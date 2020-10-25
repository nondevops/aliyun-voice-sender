# alidyvms-sender

Nightingale的理念，是将告警事件扔到redis里就不管了，接下来由各种sender来读取redis里的事件并发送，毕竟发送报警的方式太多了，适配起来比较费劲，希望社区同仁能够共建。

这里提供一个阿里语音通知的sender，参考了[https://github.com/yanjunhui/chat](https://github.com/yanjunhui/chat)，具体如何获取企业微信信息，也可以参看yanjunhui这个repo

##参考
阿里接口文档:
https://help.aliyun.com/document_detail/114035.html?spm=a2c4g.11186623.2.25.539c9152QXKfRY#doc-api-Dyvmsapi-SingleCallByTts
https://help.aliyun.com/document_detail/150016.html?spm=a2c4g.11186623.6.557.59aa3c2dtBIypg

## compile

```bash
cd $GOPATH/src
mkdir -p github.com/n9e
cd github.com/n9e
git clone https://github.com/orcswang-lang/alidyvms-sender.git
cd alidyvms-sender
go build
```

如上编译完就可以拿到二进制了。

## configuration

直接修改etc/alidyvms-sender.yml即可。另外n9e-monapi这个模块默认的发送通道只打开了mail，如果要同时使用voice，需要在notify这里打开相关配置：

```yaml
notify:
  p1: ["mail", "voice"]
  p2: ["mail", "voice"]
  p3: ["mail", "voice"]
```

## 阿里语音模板变量:
```
host,alertmsg,level  
```
如果需要重新定义变量,需自行修改 cron/sender.go和alidyvms/alidyvms.go后重新编译
## pack

编译完成之后可以打个包扔到线上去跑，将二进制和配置文件打包即可：

```bash
tar zcvf alidyvms-sender.tar.gz 
alidyvms-sender -f etc/alidyvms-sender.yml
```

## run

使用systemd或者supervisor之类的托管起来，systemd的配置实例：


```
$ cat alidyvms-sender.service
[Unit]
Description=Nightingale wechat sender
After=network-online.target
Wants=network-online.target

[Service]
User=root
Group=root

Type=simple
ExecStart=/home/n9e/alidyvms-sender
WorkingDirectory=/home/n9e

Restart=always
RestartSec=1
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
```
