package jsonrpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/thcomp/GoLang_HttpEntityHelper/entity"
)

type JSONRPCParser struct {
	executorMap map[string](methodInfo)
}

const CondMapKeyMethod = "method"

func NewJSONRPCParser() *JSONRPCParser {
	return &JSONRPCParser{}
}

func (parser *JSONRPCParser) Parse(obj interface{}) (ret entity.HttpEntity, retErr error) {
	body := (io.Reader)(nil)
	header := http.Header(nil)
	jsonrpcIntf := interface{}(nil)

	switch v := obj.(type) {
	case *http.Request:
		body = v.Body
		header = v.Header
		jsonrpcIntf = JSONRPCRequest{}
	case *http.Response:
		body = v.Body
		header = v.Header
		jsonrpcIntf = JSONRPCResponse{}
	default:
		retErr = fmt.Errorf("can not parse on JSONRPC Request, unsupported type: %v", v)
	}

	if retErr == nil {
		if parser.IsJSON(header) {
			if parseErr := json.NewDecoder(body).Decode(&jsonrpcIntf); parseErr == nil {
				switch jsonrpcIns := jsonrpcIntf.(type) {
				case JSONRPCRequest:
					if jsonrpcIns.Version != "2.0" || jsonrpcIns.Method == "" {
						retErr = fmt.Errorf("can not parse on JSONRPC Request")
					} else {
						ret = &jsonrpcIns
					}
				case JSONRPCResponse:
					if jsonrpcIns.Version != "2.0" {
						retErr = fmt.Errorf("can not parse on JSONRPC Response")
					} else {
						ret = &jsonrpcIns
					}
				}
			} else {
				retErr = parseErr
			}
		} else {
			retErr = entity.ErrUnsupportEntity
		}
	}

	return
}

func (parser *JSONRPCParser) IsJSON(headers http.Header) (ret bool) {
	mimetype := headers.Get("Content-type")
	lowerMimetype := strings.ToLower(mimetype)
	if strings.HasPrefix(lowerMimetype, "application/json") {
		ret = true
	}

	return
}
