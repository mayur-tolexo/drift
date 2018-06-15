package lib

//Error Model
type Error struct {
	Code     int    `json:"code"`
	Msg      string `json:"Message"`
	Trace    string `json:"trace,omitempty"`
	DebugMsg string `json:"debug_msg,omitempty"`
}

//Resp : Response Model
type Resp struct {
	Data   interface{} `json:"data,omitempty"`
	Msg    interface{} `json:"message,omitempty"`
	Error  error       `json:"error,omitempty"`
	Status bool        `json:"status"`
}

//MsgResp : Message response for create/update/delete api
type MsgResp struct {
	Data interface{}
	Msg  string
}
