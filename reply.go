package mock

import (
	"encoding/json"
	"encoding/xml"
	"strings"

	log "github.com/cjburchell/uatu-go"
)

type IReply interface {
	Body(body interface{}, logger log.ILog) IReply
	Content(content string) IReply
	Code(code int) IReply
	Header(key, value string) IReply
	BodyString(body string) IReply
	Delay(delayTime int) IReply
}

type reply struct {
	config *endpointConfig
}

func (reply *reply) Body(body interface{}, logger log.ILog) IReply {
	if strings.Contains(reply.config.ContentType, "xml") {
		bodyString, err := xml.Marshal(body)
		if err != nil {
			logger.Error(err)
		}
		return reply.BodyString(string(bodyString))
	}

	if strings.Contains(reply.config.ContentType, "json") {
		bodyString, err := json.Marshal(body)
		if err != nil {
			logger.Error(err)
		}
		return reply.BodyString(string(bodyString))
	}

	return reply
}

func (reply *reply) Content(content string) IReply {
	reply.config.ContentType = content
	return reply
}

func (reply *reply) Code(code int) IReply {
	reply.config.Response = code
	return reply
}

func (reply *reply) Delay(delayTime int) IReply {
	reply.config.ReplyDelay = delayTime
	return reply
}

func (reply *reply) Header(key, value string) IReply {
	if reply.config.Header == nil {
		reply.config.Header = map[string]string{}
	}

	reply.config.Header[key] = value
	return reply
}

func (reply *reply) BodyString(body string) IReply {
	reply.config.StringBody = body
	return reply
}
