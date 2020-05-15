package mock

import (
	"encoding/json"
)

type endpointConfig struct {
	Path         string            `json:"path"`
	Method       string            `json:"method"`
	ResponseBody json.RawMessage   `json:"response_body"`
	StringBody   string            `json:"string_body"`
	ContentType  string            `json:"content_type"`
	Response     int               `json:"response"`
	Header       map[string]string `json:"header"`
	Name         string            `json:"-"`
	ReplyDelay   int               `json:"reply_delay"`
}

type endpoint struct {
	config               endpointConfig
	customRequestHandler func(string, string, string)
}

//IEndpoint interface
type IEndpoint interface {
	RequestHandler(handler func(string, string, string)) IEndpoint
	Reply(handler func(IReply)) IEndpoint
}

func (endpoint *endpoint) RequestHandler(handler func(string, string, string)) IEndpoint {
	endpoint.customRequestHandler = handler
	return endpoint
}

func (endpoint *endpoint) Reply(handler func(IReply)) IEndpoint {
	handler(&reply{config: &endpoint.config})
	return endpoint
}
