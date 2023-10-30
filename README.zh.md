# WeTriage

[English](./README.md)

对来自企业微信平台的回调 XML 消息进行分类，将其转换为 JSON 并发布到 MQTT 以供下游使用。

>该项目仍在开发中。

## 功能

- [x] 回调消息签名验证和解密
- [x] 响应回调地址注册测试请求
- [x] 将已识别的 XML 消息转换为 JSON
- [x] 将带有可识别名称的 JSON 消息发布给 MQTT
- [x] 按照文档要求对回调的消息进行正确响应

**非功能**
- [ ] 处理回调消息的任何业务逻辑
- [ ] TLS

## 如何使用

```bash
docker pull ghcr.io/imulab/wetriage:latest

# 对于特定版本，请使用短提交 SHA 作为标记。例如：
# docker pull ghcr.io/imulab/wetriage:117eb11f
#
# 注意这只是一个例子，不是最新的SHA
```

支持以下标志：

| 标志             | 描述                       | 默认          | 环境            |
|----------------|--------------------------|-------------|---------------|
| `--port`       | 要监听的端口                   | `8080`      | `WT_PORT`     |
| `--debug`      | 启用调试模式                   | `false`     | `WT_DEBUG`    |
| `--path`       | 自定义回调端点路径                | `/callback` | `WT_PATH`     |
| `--token`      | 在微信注册的回调令牌               | -           | `WT_TOKEN`    |
| `--aes-key`    | 在微信注册的 Base64 编码的 AES 密钥 | -           | `WT_AES_KEY`  |
| `--topic`，`-t` | 要处理的回调主题。详情请见下文          | -           | -             |
| `--mqtt-url`   | MQTT 经纪商网址。详情见下文         | -           | `WT_MQTT_URL` |

下面显示了使用该镜像的示例。

```bash
docker run -d \
    -p 8080:8080 \
    -e WT_TOKEN=token \
    -e WT_AES_KEY=base64_encoded_aes_key \
    -e WT_MQTT_URL=tcp://localhost:1883 \
    ghcr.io/imulab/wetriage:latest WeTriage server -t suite_ticket_info
```

## 消息话题

当前，支持以下消息主题。

| 话题                          | 描述                                                                                                                        |
|-----------------------------|---------------------------------------------------------------------------------------------------------------------------|
| `suite_ticket_info`         | [推送应用模板凭证](https://developer.work.weixin.qq.com/document/path/97173)                                                      |
| `create_auth_info`          | [授权成功通知](https://developer.work.weixin.qq.com/document/path/97174)                                                        |
| `change_auth_info`          | [変更授权通知](https://developer.work.weixin.qq.com/document/path/97174#%E5%8F%98%E6%9B%B4%E6%8E%88%E6%9D%83%E9%80%9A%E7%9F%A5) |
| `reset_permanent_code_info` | [重置永久授权码通知](https://developer.work.weixin.qq.com/document/path/97175)                                                     |

> 随着项目的进展，将添加更多主题。如果需要申请某个消息话题，请创建一个Issue。

### 消息格式

已识别的主题将转换为一个结构类似的 JSON 对象。例如，`suite_ticket_info`消息 XML：

```xml
<xml>
    <SuiteId><![CDATA[ww4asffe99e54c0fxxxx]]></SuiteId>
    <InfoType><![CDATA[suite_ticket]]></InfoType>
    <TimeStamp>1403610513 </TimeStamp>
    <SuiteTicket><![CDATA[asdfasfdasdfasdf]]></SuiteTicket>
</xml>
```

转换为 JSON 后：

```json
{
  "suite_id": "ww4asffe99e54c0fxxxx",
  "info_type": "suite_ticket",
  "timestamp": 1403610513,
  "suite_ticket": "asdfasfdasdfasdf"
}
```

要确保 JSON 的确切格式，请参阅 [topic package](./topic) 了解详情。

### MQTT

上面的 JSON 消息将封装在结构体中，并以 `T/WeTriage/<topic> `为主题发布到 MQTT 。

例如，上面的 `suite_ticket_info` 消息将在 `T/WeTriage/suite_ticket_info` 下发布：

```json
{
  "id": "a unique id",
  "created_at": 1403610513,
  "topic": "suite_ticket_info",
  "content": {
    "suite_id": "ww4asffe99e54c0fxxxx",
    "info_type": "suite_ticket",
    "timestamp": 1403610513,
    "suite_ticket": "asdfasfdasdfasdf"
  }
}
```
