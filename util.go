package drift

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/mayur-tolexo/drift/lib"
	nsq "github.com/nsqio/go-nsq"
	"github.com/rightjoin/aqua"
)

const (
	driftApp          = "-DRIFT-"
	allKey            = "ALL"
	defaultAdminAddrs = "127.0.0.1:4171"
)

//NewConsumer will create new consumer
func NewConsumer(jobHandler JobHandler) *Drift {
	return &Drift{
		jobHandler:    jobHandler,
		chanelHandler: make(map[string]JobHandler),
		consumers:     make(map[string][]*nsq.Consumer),
	}
}

//NewPub will create new publisher
func NewPub(nsqDHttpAddrs string) *Drift {
	return &Drift{
		jobHandler:    nil,
		chanelHandler: make(map[string]JobHandler),
		consumers:     make(map[string][]*nsq.Consumer),
		pubAddrs:      strings.TrimPrefix(nsqDHttpAddrs, "http://"),
	}
}

//NewAdmin will create new admin model
func NewAdmin(nsqDHttpAddrs string) *Drift {
	return &Drift{
		jobHandler:    nil,
		chanelHandler: make(map[string]JobHandler),
		consumers:     make(map[string][]*nsq.Consumer),
		pubAddrs:      strings.TrimPrefix(nsqDHttpAddrs, "http://"),
	}
}

//HandleMessage will define the nsq handler method
func (th *tailHandler) HandleMessage(m *nsq.Message) error {
	return th.jobHandler(string(m.Body))
}

//vAddConsumer will validate add consumer request
func vAddConsumer(req aqua.Aide) (payload AddConstumer, err error) {
	req.LoadVars()
	err = lib.Unmarshal(req.Body, &payload)
	return
}

//vAdmin will validate admin action request
func vAdmin(req aqua.Aide) (payload Admin, err error) {
	req.LoadVars()
	if err = lib.Unmarshal(req.Body, &payload); err == nil {
		if payload.Topic == "" {
			err = lib.VError("Empty Topic")
		} else {
			switch payload.Action {
			case "empty", "delete", "pause", "unpause":
			default:
				err = lib.VError("Invalid action")
			}
		}
	}
	return
}

//vProduceReq will validate produce request
func vPublishReq(req aqua.Aide) (payload Publish, err error) {
	req.LoadVars()
	if err = lib.Unmarshal(req.Body, &payload); err == nil {
		if payload.NsqDHTTPAddrs == "" {
			err = lib.VError("Invalid nsqd address")
		}
	}
	return
}

//vKillConsumer will validate kill consumer request
func vKillConsumer(req aqua.Aide) (payload KillConsumer, err error) {
	req.LoadVars()
	if err = lib.Unmarshal(req.Body, &payload); err == nil {
		if payload.Count <= 0 {
			err = lib.BadReqError(err, "Invalid count")
		}
	}
	return
}

//pPublishReq will process the publish request
func pPublishReq(payload Publish) (data interface{}, err error) {
	var (
		b    []byte
		req  *http.Request
		resp *http.Response
	)
	if b, err = jsoniter.Marshal(payload.Data); err == nil {
		payload.NsqDHTTPAddrs = strings.TrimPrefix(payload.NsqDHTTPAddrs, "http://")
		URL := fmt.Sprintf("http://%v/pub?topic=%v", payload.NsqDHTTPAddrs, payload.Topic)
		if req, err = http.NewRequest("POST",
			URL, bytes.NewBuffer(b)); err == nil {
			HTTPClient := &http.Client{}
			if resp, err = HTTPClient.Do(req); err == nil {
				defer resp.Body.Close()
				data = resp.StatusCode
			} else {
				err = lib.BadReqError(err)
			}
		} else {
			err = lib.BadReqError(err)
		}
	} else {
		err = lib.BadReqError(err)
	}
	return
}

func getChannel(c string) (channel string) {
	channel = c
	if channel == "" {
		rand.Seed(time.Now().UnixNano())
		channel = fmt.Sprintf("drift%06d#ephemeral", rand.Int()%999999)
	}
	return
}

func hash(a, b string) string {
	return a + driftApp + b
}
