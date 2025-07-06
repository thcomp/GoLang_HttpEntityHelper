package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"sync"

	root "github.com/thcomp/GoLang_HttpEntityHelper"
	ThcompUtility "github.com/thcomp/GoLang_Utility"
)

var sAutoID int64 = 0
var sAutoIDMutex sync.Mutex

type JSONRPC struct {
	Version string      `json:"jsonrpc"`
	id      interface{} `json:"id"`
}

func newJSONRPC() *JSONRPC {
	return &JSONRPC{Version: "2.0"}
}

func newJSONRPCWithID(id interface{}) *JSONRPC {
	if id == nil {
		sAutoIDMutex.Lock()
		defer sAutoIDMutex.Unlock()
		sAutoID++
		id = sAutoID
	}

	return &JSONRPC{Version: "2.0", id: id}
}

func (rpc *JSONRPC) IsIDNum() bool {
	infHelper := ThcompUtility.NewInterfaceHelper(rpc.id)
	return infHelper.IsNumber()
}

func (rpc *JSONRPC) IsIDString() bool {
	infHelper := ThcompUtility.NewInterfaceHelper(rpc.id)
	return infHelper.IsString()
}

func (rpc *JSONRPC) IDNum() (ret float64, isNum bool) {
	infHelper := ThcompUtility.NewInterfaceHelper(rpc.id)
	ret, isNum = infHelper.GetNumber()

	return
}

func (rpc *JSONRPC) IDString() (ret string, isNum bool) {
	infHelper := ThcompUtility.NewInterfaceHelper(rpc.id)
	ret, isNum = infHelper.GetString()

	return
}

func (rpc *JSONRPC) IDInterface() (ret interface{}) {
	return rpc.id
}

type JSONRPCRequest struct {
	JSONRPC
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func NewJSONRPCNotificationRequest(method string, params interface{}) (*JSONRPCRequest, error) {
	ret := &JSONRPCRequest{
		JSONRPC: *newJSONRPC(),
		Method:  method,
	}
	retErr := error(nil)

	if reader, assertionOK := params.(io.Reader); assertionOK {
		if paramsBytes, readErr := ioutil.ReadAll(reader); readErr == nil {
			ret.Params = paramsBytes
		} else {
			retErr = readErr
		}
	} else {
		ret.Params = params
	}

	return ret, retErr
}

func NewJSONRPCRequest(id interface{}, method string, params interface{}) (*JSONRPCRequest, error) {
	ret := &JSONRPCRequest{
		JSONRPC: *newJSONRPCWithID(id),
		Method:  method,
	}
	retErr := error(nil)

	if reader, assertionOK := params.(io.Reader); assertionOK {
		if paramsBytes, readErr := ioutil.ReadAll(reader); readErr == nil {
			ret.Params = paramsBytes
		} else {
			retErr = readErr
		}
	} else {
		ret.Params = params
	}

	return ret, retErr
}

func ParseJSONRequest(reader io.Reader) (*JSONRPCRequest, error) {
	ret := (*JSONRPCRequest)(nil)
	tempRet := map[string]interface{}{}
	retErr := json.NewDecoder(reader).Decode(&tempRet)

	if retErr == nil {
		ret = &JSONRPCRequest{}
		if valueInf, exist := tempRet["jsonrpc"]; exist {
			valueInfHelper := ThcompUtility.NewInterfaceHelper(valueInf)
			if valueInfHelper.IsString() {
				ret.JSONRPC.Version, _ = valueInfHelper.GetString()
			}
		}
		if valueInf, exist := tempRet["id"]; exist {
			ret.JSONRPC.id = valueInf
		}
		if valueInf, exist := tempRet["method"]; exist {
			valueInfHelper := ThcompUtility.NewInterfaceHelper(valueInf)
			if valueInfHelper.IsString() {
				ret.Method, _ = valueInfHelper.GetString()
			}
		}
		if valueInf, exist := tempRet["params"]; exist {
			ret.Params = valueInf
		}
	}

	return ret, retErr
}

func (req *JSONRPCRequest) EntityType() root.HttpEntityType {
	return root.JSONRPC_Request
}

func (req *JSONRPCRequest) EncodeByJSON() ([]byte, error) {
	tempMap := map[string]interface{}{
		"jsonrpc": req.JSONRPC.Version,
		"method":  req.Method,
	}

	if req.JSONRPC.id != nil {
		tempMap["id"] = req.JSONRPC.id
	}
	if req.Params != nil {
		tempMap["params"] = req.Params
	}

	return json.Marshal(tempMap)
}

func (req *JSONRPCRequest) ParseParams(toPtr interface{}) (retErr error) {
	if reflect.TypeOf(toPtr).Kind() == reflect.Pointer {
		if tempJsonBytes, marshalErr := json.Marshal(req.Params); marshalErr == nil {
			retErr = json.Unmarshal(tempJsonBytes, toPtr)
		} else {
			retErr = marshalErr
		}
	} else {
		retErr = fmt.Errorf("toPtr is not pointer")
	}

	return retErr
}

type JSONRPCResponse struct {
	JSONRPC
	Result interface{}   `json:"result"`
	Error  *JSONRPCError `json:"error"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const JSONRPCParseError = -32700
const JSONRPCInvalidRequest = -32600
const JSONRPCMethodNotFound = -32601
const JSONRPCInvalidParams = -32602
const JSONRPCInternalError = -32603
const JSONRPCServerErrorMax = -32000
const JSONRPCServerErrorMin = -32099

func NewJSONRPCResponse(id interface{}) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: *newJSONRPCWithID(id),
	}
}

func NewJSONRPCResponseFromRequest(request *JSONRPCRequest) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: request.JSONRPC,
	}
}

func ParseJSONResponse(reader io.Reader) (*JSONRPCResponse, error) {
	ret := (*JSONRPCResponse)(nil)
	tempRet := map[string]interface{}{}
	retErr := json.NewDecoder(reader).Decode(&tempRet)

	if retErr == nil {
		ret = &JSONRPCResponse{}
		if valueInf, exist := tempRet["jsonrpc"]; exist {
			valueInfHelper := ThcompUtility.NewInterfaceHelper(valueInf)
			if valueInfHelper.IsString() {
				ret.JSONRPC.Version, _ = valueInfHelper.GetString()
			}
		}
		if valueInf, exist := tempRet["id"]; exist {
			ret.JSONRPC.id = valueInf
		}
		if valueInf, exist := tempRet["result"]; exist {
			ret.Result = valueInf
		}
		if valueInf, exist := tempRet["error"]; exist {
			if errorMap, assertionOK := valueInf.(map[string]interface{}); assertionOK {
				ret.Error = &JSONRPCError{}
				if valueInf, exist := errorMap["code"]; exist {
					valueInfHelper := ThcompUtility.NewInterfaceHelper(valueInf)
					if valueInfHelper.IsNumber() {
						tempValue, _ := valueInfHelper.GetNumber()
						ret.Error.Code = int(tempValue)
					}
				}
				if valueInf, exist := errorMap["message"]; exist {
					valueInfHelper := ThcompUtility.NewInterfaceHelper(valueInf)
					if valueInfHelper.IsString() {
						ret.Error.Message, _ = valueInfHelper.GetString()
					}
				}
				if valueInf, exist := errorMap["data"]; exist {
					ret.Error.Data = valueInf
				}
			}
		}
	}

	return ret, retErr
}

func (res *JSONRPCResponse) EntityType() root.HttpEntityType {
	return root.JSONRPC_Response
}

func (res *JSONRPCResponse) EncodeByJSON() ([]byte, error) {
	tempMap := map[string]interface{}{
		"jsonrpc": res.JSONRPC.Version,
	}

	if res.JSONRPC.id != nil {
		tempMap["id"] = res.JSONRPC.id
	}
	if res.Result != nil {
		tempMap["result"] = res.Result
	}
	if res.Error != nil {
		tempMap["error"] = res.Error
	}

	return json.Marshal(tempMap)
}

func (res *JSONRPCResponse) ParseResult(toPtr interface{}) (retErr error) {
	if reflect.TypeOf(toPtr).Kind() == reflect.Pointer {
		if tempJsonBytes, marshalErr := json.Marshal(res.Result); marshalErr == nil {
			retErr = json.Unmarshal(tempJsonBytes, toPtr)
		} else {
			retErr = marshalErr
		}
	} else {
		retErr = fmt.Errorf("toPtr is not pointer")
	}

	return retErr
}

func (res *JSONRPCResponse) Reader() (ret io.Reader, retErr error) {
	if jsonBytes, encodeErr := res.EncodeByJSON(); encodeErr == nil {
		ret = bytes.NewReader(jsonBytes)
	} else {
		retErr = encodeErr
	}
	return
}

func NewJSONRPCError(code int, message string, data interface{}) *JSONRPCError {
	return &JSONRPCError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func (err *JSONRPCError) EntityType() root.HttpEntityType {
	return root.JSONRPC_Error
}

func (err *JSONRPCError) ParseData(toPtr interface{}) (retErr error) {
	if reflect.TypeOf(toPtr).Kind() == reflect.Pointer {
		if tempJsonBytes, marshalErr := json.Marshal(err.Data); marshalErr == nil {
			retErr = json.Unmarshal(tempJsonBytes, toPtr)
		} else {
			retErr = marshalErr
		}
	} else {
		retErr = fmt.Errorf("toPtr is not pointer")
	}

	return retErr
}
