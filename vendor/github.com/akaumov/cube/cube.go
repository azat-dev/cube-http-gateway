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
	Version string           `json:"version"`
	Id      *string          `json:"id"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params"`
}

type Error struct {
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Data        *json.RawMessage `json:"data"`
}

type Request struct {
	Version string           `json:"version"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params"`
}

type Response struct {
	Version string           `json:"version"`
	Result  *json.RawMessage `json:"result"`
	Errors  *[]Error         `json:"errors"`
}

type Cube interface {
	GetParam(param string) string
	GetClass() string
	GetInstanceId() string

	PublishMessage(channel string, message Message) error
	CallMethod(channel string, request Request, timeout time.Duration) (*Response, error)

	LogDebug(text string) error
	LogError(text string) error
	LogFatal(text string) error
	LogInfo(text string) error
	LogWarning(text string) error
	LogTrace(text string) error
}

type HandlerInterface interface {
	OnStart(instance Cube)
	OnStop(instance Cube)
	OnReceiveMessage(instance Cube, channel string, message Message)
	OnReceiveRequest(instance Cube, channel string, request Request) (*Response, error)
}
