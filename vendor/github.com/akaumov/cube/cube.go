package cube

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrorTimeout = errors.New("cube: request timeout")
)

type Message struct {
	Id     string           `json:"id"`
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

type Request struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params"`
}

type Response struct {
	Id     string           `json:"id"`
	Result *json.RawMessage `json:"result"`
	Error  *Error           `json:"error"`
}

type Error struct {
	Name    string `json:"name"`
	Message string `json:"description"`
}

func NewResultResponse(requestId string, result *json.RawMessage) Response {
	return Response{
		Id:     requestId,
		Result: result,
		Error:  nil,
	}
}

func NewErrorResponse(requestId string, name string, message string) Response {
	return Response{
		Id:     requestId,
		Result: nil,
		Error: &Error{
			Name:    name,
			Message: message,
		},
	}
}

type Channel string
type InputChannel Channel

type Cube interface {
	GetParam(param string) string
	GetClass() string
	GetInstanceId() string

	PublishMessage(channel Channel, message Message) error
	CallMethod(channel Channel, request Request, timeout time.Duration) (*Response, error)

	Stop()

	LogDebug(text string) error
	LogError(text string) error
	LogFatal(text string) error
	LogInfo(text string) error
	LogWarning(text string) error
	LogTrace(text string) error
}

type HandlerInterface interface {
	OnInitInstance() []InputChannel
	OnStart(instance Cube) error
	OnStop(instance Cube)
	OnReceiveMessage(instance Cube, channel Channel, message Message)
	OnReceiveRequest(instance Cube, channel Channel, request Request) Response
}
