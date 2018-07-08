package js

type RequestParams struct {
	InputTime  int64               `json:"inputTime"`
	Method     string              `json:"method"`
	Host       string              `json:"host"`
	RemoteAddr string              `json:"remoteAddr"`
	RequestURI string              `json:"requestURI"`
	Body       []byte              `json:"body"`
	UserId     *string             `json:"userId"`
	DeviceId   *string             `json:"deviceId"`
	Headers    map[string][]string `json:"headers"`
}

type Response struct {
	Status int64  `json:"status"`
	Body   []byte `json:"body"`
}
