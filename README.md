# WeTriage

[中文](./README.zh.md)

Triage callback XML messages from the WeCom (Enterprise WeChat) platform and convert them to JSON format before
publishing them to a MQTT message broker for downstream consumption.

> This project is still under active development.

## Features

- [x] Callback message signature verification and decryption
- [x] Callback echo test
- [x] Convert identified XML messages to JSON
- [x] Publish converted JSON messages with an identifiable name to MQTT broker for consumption
- [x] Properly respond to WeCom server with a response required by the documentation for the incoming message

**This project does NOT try to**
- [ ] Handle any business logic of the callback messages
- [ ] Enable TLS for the HTTP server

## Getting Started

```bash
docker pull ghcr.io/imulab/wetriage:latest

# For a specific version, use the short commit SHA as the tag. For example:
#   docker pull ghcr.io/imulab/wetriage:117eb11f
#
# Note this is just an example, that's not the latest commit hash
```

The following flags are supported:

| Flag            | Description                                  | Default     | Env           |
|-----------------|----------------------------------------------|-------------|---------------|
| `--port`        | Port to listen on                            | `8080`      | `WT_PORT`     |
| `--debug`       | Enable debug mode                            | `false`     | `WT_DEBUG`    |
| `--path`        | Customize the callback endpoint path         | `/callback` | `WT_PATH`     |
| `--token`       | Callback token registered with WeCom         | -           | `WT_TOKEN`    |
| `--aes-key`     | Base64 encoded AES key registered with WeCom | -           | `WT_AES_KEY`  |
| `--topic`, `-t` | Callback topic to process. See details below | -           | -             |
| `--mqtt-url`    | MQTT broker URL. See details below           | -           | `WT_MQTT_URL` |

Below shows an example of using the image.

```bash
docker run -d \
    -p 8080:8080 \
    -e WT_TOKEN=token \
    -e WT_AES_KEY=base64_encoded_aes_key \
    -e WT_MQTT_URL=tcp://localhost:1883 \
    ghcr.io/imulab/wetriage:latest WeTriage server -t suite_ticket_info
```

## Topics

Currently, the following topics are supported.

| Topic                       | Description                                                                                                               |
|-----------------------------|---------------------------------------------------------------------------------------------------------------------------|
| `suite_ticket_info`         | [推送应用模板凭证](https://developer.work.weixin.qq.com/document/path/97173)                                                      |
| `create_auth_info`          | [授权成功通知](https://developer.work.weixin.qq.com/document/path/97174)                                                        |
| `change_auth_info`          | [変更授权通知](https://developer.work.weixin.qq.com/document/path/97174#%E5%8F%98%E6%9B%B4%E6%8E%88%E6%9D%83%E9%80%9A%E7%9F%A5) |
| `reset_permanent_code_info` | [重置永久授权码通知](https://developer.work.weixin.qq.com/document/path/97175)                                                     |

> More topics will be added as the project progresses. To request a topic, please open an issue.

### Message Format

An identified topic is converted to a similar object in JSON. For example, an original `suite_ticket_info` topic
message could be the following XML:

```xml
<xml>
    <SuiteId><![CDATA[ww4asffe99e54c0fxxxx]]></SuiteId>
    <InfoType><![CDATA[suite_ticket]]></InfoType>
    <TimeStamp>1403610513</TimeStamp>
    <SuiteTicket><![CDATA[asdfasfdasdfasdf]]></SuiteTicket>
</xml>
```

After converting to JSON, it could look like:

```json
{
  "suite_id": "ww4asffe99e54c0fxxxx",
  "info_type": "suite_ticket",
  "timestamp": 1403610513,
  "suite_ticket": "asdfasfdasdfasdf"
}
```

To make sure the exact format of the JSON, please refer to types in the [topic package](./topic) for details.

### MQTT

The above JSON message is to be wrapped in a message envelope and published to the MQTT broker under a topic of `T/WeTriage/<topic>`.

For example, the above `suite_ticket_info` message will be published to the broker under `T/WeTriage/suite_ticket_info`, with a payload of:

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
