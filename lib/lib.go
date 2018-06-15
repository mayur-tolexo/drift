package lib

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/tolexo/aero/conf"
)

//NewResp : Create New Response Objext
func NewResp() *Resp {
	return &Resp{}
}

//GetErrorCode : Get Error Code from *Error
func GetErrorCode(err error) (code int) {
	if err == nil {
		code = NO_ERROR
	} else {
		switch err.(type) {
		case *Error:
			code = err.(*Error).Code
		default:
			code = NO_ERROR
		}
	}
	return
}

//AppendDebugMsg : append debugMsg in error
func AppendDebugMsg(err error, debugMsg ...string) error {
	switch err.(type) {
	case *Error:
		err.(*Error).DebugMsg += " " + strings.Join(debugMsg, " ")
	}
	return err
}

//unsetDebugMsg : Unset debugMsg of *Error
func unsetDebugMsg(err error) error {
	switch err.(type) {
	case *Error:
		err.(*Error).DebugMsg = ""
		err.(*Error).Trace = ""
	}
	return err
}

//GetErrorHTTPCode : Get HTTP code from error code. Please check for httpCode = 0.
func GetErrorHTTPCode(err error) (httpCode int, status bool) {
	errCode := GetErrorCode(err)
	httpCode = ErrorHTTPCode[errCode]
	if err == nil {
		status = true
	}
	return
}

//Error : Implement Error method of error interface
func (e *Error) Error() string {
	return fmt.Sprintf("\nCode:\t\t[%d]\nMessage:\t[%v]\nStackTrace:\t[%v]\nDebugMsg:\t[%v]\n", e.Code, e.Msg, e.Trace, e.DebugMsg)
}

//newError : Create new *Error object
func newError(msg string, err error, code int, debugMsg ...string) error {
	stackDepth := conf.Int("error.stack_depth", 3)
	funcName, fileName, line := StackTrace(stackDepth)
	trace := fileName + " -> " + funcName + ":" + strconv.Itoa(line)
	errStr := strings.Join(debugMsg, " ")

	return &Error{
		Msg:      msg,
		DebugMsg: errStr,
		Trace:    trace,
		Code:     code,
	}
}

//Set Response Values
func (r *Resp) Set(data interface{}, status bool, err error) {
	var (
		respData interface{}
		msg      string
	)
	if val, ok := data.(MsgResp); ok {
		respData = val.Data
		msg = val.Msg
	} else {
		respData = data
	}
	if err == nil {
		r.Data = respData
		r.Msg = msg
	} else {
		r.Error = err
	}
	r.Status = status
	if debugMsg := conf.Bool("error.debug_msg", false); debugMsg == false {
		r.Error = unsetDebugMsg(r.Error)
	}
}

//StackTrace : Get function name, file name and line no of the caller function
//Depth is the value from which it will start searching in the stack
func StackTrace(depth int) (funcName string, file string, line int) {
	var (
		ok bool
		pc uintptr
	)
	for i := depth; ; i++ {
		if pc, file, line, ok = runtime.Caller(i); ok {
			if trackAll := conf.Bool("error.track_all", false); trackAll == false &&
				strings.Contains(file, PACKAGE_NAME) {
				continue
			}
			fileName := strings.Split(file, "github.com")
			if len(fileName) > 1 {
				file = fileName[1]
			}
			_, funcName = packageFuncName(pc)
			break
		} else {
			break
		}
	}
	return
}

//packageFuncName : Package and function name from package counter
func packageFuncName(pc uintptr) (packageName string, funcName string) {
	if f := runtime.FuncForPC(pc); f != nil {
		funcName = f.Name()
		if ind := strings.LastIndex(funcName, "/"); ind > 0 {
			packageName += funcName[:ind+1]
			funcName = funcName[ind+1:]
		}
		if ind := strings.Index(funcName, "."); ind > 0 {
			packageName += funcName[:ind]
			funcName = funcName[ind+1:]
		}
	}
	return
}

//BuildResponse : creates the response of API
func BuildResponse(data interface{}, err error) (int, interface{}) {
	out := NewResp()
	httpCode, status := GetErrorHTTPCode(err)
	out.Set(data, status, err)
	return httpCode, out
}

//Unmarshal : unmarshal of request body and wrap error
func Unmarshal(reqBody string, structModel interface{}) (err error) {
	if err = jsoniter.Unmarshal([]byte(reqBody), structModel); err != nil {
		err = UnmarshalError(err)
	}
	return
}

//GetPriorityValue : get the value by priority
func GetPriorityValue(a ...interface{}) (pVal interface{}) {
	for _, curVal := range a {
		pVal = reflect.Zero(reflect.TypeOf(curVal)).Interface()
		if curVal != pVal {
			pVal = curVal
			break
		}
	}
	return
}
