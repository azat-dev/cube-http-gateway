package cube_http_gateway

import (
	"encoding/json"
	"fmt"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/akaumov/cube"
	"github.com/akaumov/cube-http-gateway/js"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const Version = "1"

type Handler struct {
	httpServer   *http.Server
	cubeInstance cube.Cube
	timeoutMs    uint64
	jwtSecret    string
}

func (h *Handler) OnInitInstance() []cube.InputChannel {
	return []cube.InputChannel{}
}

func (h *Handler) OnStart(cubeInstance cube.Cube) {
	fmt.Println("Starting http gateway...")

	h.cubeInstance = cubeInstance
	h.jwtSecret = cubeInstance.GetParam("jwtString")

	h.timeoutMs = 30000
	timeoutString := cubeInstance.GetParam("timeoutMs")

	if timeoutString != "" {
		timeoutMs, err := strconv.ParseUint(timeoutString, 10, 64)
		if err != nil {
			cubeInstance.LogError("Wrong ")
		}

		h.timeoutMs = timeoutMs
	}

	go h.startHttpServer(cubeInstance)
}

func (h *Handler) OnStop(c cube.Cube) {
}

func (h *Handler) OnReceiveMessage(instance cube.Cube, channel cube.Channel, message cube.Message) {
	fmt.Println("OnReceiveMessage: is not implemented")
	instance.LogError("OnReceiveMessage: is not implemented")
}

//From bus
func (h *Handler) OnReceiveRequest(instance cube.Cube, channel cube.Channel, request cube.Request) (*cube.Response, error) {
	fmt.Println("OnReceiveRequest: is not implemented")
	instance.LogError("OnReceiveRequest: is not implemented")
	return &cube.Response{
		Version: Version,
		Result:  nil,
		Errors: &[]cube.Error{
			{
				Code:        "400",
				Name:        "NotImplemented",
				Description: "OnReceiveRequest: is not implemented",
			},
		},
	}, nil
}

func (h *Handler) startHttpServer(cubeInstance cube.Cube) {

	srv := http.Server{
		Addr:    ":80",
		Handler: h,
	}

	h.httpServer = &srv

	fmt.Println("Start http listening")
	cubeInstance.LogInfo("Start http listening")
	err := srv.ListenAndServe()

	fmt.Println("Stop http listenning", err)
	cubeInstance.LogFatal(err.Error())
}

func (h *Handler) getAuthData(tokenString string) (*string, *string, error) {

	if tokenString == "" {
		return nil, nil, fmt.Errorf("empty token")
	}

	newToken, err := jws.ParseJWT([]byte(tokenString))
	if err != nil {
		return nil, nil, err
	}

	err = newToken.Validate([]byte(h.jwtSecret), crypto.SigningMethodHS256)
	if err != nil {
		return nil, nil, err
	}

	claims := newToken.Claims()
	userId := claims.Get("userId").(string)
	deviceId := claims.Get("deviceId").(string)

	return &userId, &deviceId, nil
}

func (h *Handler) packRequest(userId *string, deviceId *string, request *http.Request) (*cube.Request, error) {
	var err error
	var body []byte

	if request.Body != nil {
		body, err = ioutil.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		request.Body.Close()

		if body != nil && len(body) == 0 {
			body = nil
		}
	}

	headers := map[string][]string{}

	for key, value := range request.Header {
		headers[key] = value
	}

	params := js.RequestParams{
		DeviceId:   deviceId,
		UserId:     userId,
		Method:     request.Method,
		InputTime:  time.Now().UnixNano(),
		Host:       request.Host,
		RequestURI: request.RequestURI,
		Body:       body,
		RemoteAddr: request.RemoteAddr,
		Headers:    headers,
	}

	packedParams, err := json.Marshal(params)

	requestData := &cube.Request{
		Version: "1",
		Method:  request.Method,
		Params:  (*json.RawMessage)(&packedParams),
	}

	return requestData, nil
}

func (h *Handler) handleResponse(responseMessage *cube.Response, writer http.ResponseWriter) error {

	var response js.Response

	err := json.Unmarshal(*responseMessage.Result, &response)
	if err != nil {
		return err
	}

	writer.WriteHeader(int(response.Status))
	if response.Body != nil && len(response.Body) > 0 {
		writer.Write(response.Body)
	}

	return nil
}

//Request from gateway
func (h *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	fmt.Println("onReceiveRequest", request)

	var userId, deviceId *string
	var err error
	token := ""

	if h.jwtSecret != "" {
		token = request.Header.Get("X-Auth-Token")
		userId, deviceId, err = h.getAuthData(token)

		if err != nil {
			http.Error(writer,
				http.StatusText(http.StatusUnauthorized),
				http.StatusUnauthorized)
			return
		}
	}

	requestData, err := h.packRequest(userId, deviceId, request)
	if err != nil {
		http.Error(writer,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	timeout := time.Duration(h.timeoutMs) * time.Millisecond
	cubeChannel := cube.Channel(request.Method)

	response, err := h.cubeInstance.CallMethod(cubeChannel, *requestData, timeout)

	if err != nil {
		if err == cube.ErrorTimeout {
			http.Error(writer,
				http.StatusText(http.StatusGatewayTimeout),
				http.StatusGatewayTimeout)
			return
		}

		http.Error(writer,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	h.handleResponse(response, writer)
	if err != nil {
		http.Error(writer,
			http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

var _ cube.HandlerInterface = (*Handler)(nil)
