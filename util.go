package drift

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/mayur-tolexo/aqua"
	"github.com/mayur-tolexo/drift/lib"
	nsq "github.com/nsqio/go-nsq"
)

const (
	driftApp          = "-DRIFT-"
	allKey            = "ALL"
	defaultAdminAddrs = "127.0.0.1:4171"
	createAction      = "create"
	emptyAction       = "empty"
	deleteAction      = "delete"
	pauseAction       = "pause"
	unpauseAction     = "unpause"
)

//NewConsumer will create new consumer
func NewConsumer(jobHandler JobHandler) *Drift {
	return &Drift{
		jobHandler:    jobHandler,
		chanelHandler: make(map[string]JobHandler),
		consumers:     make(map[string][]*nsq.Consumer),
	}
}

//StartAdmin will start admin
func StartAdmin(lookupHTTPAddr []string, httpAddrs string) {
	if httpAddrs == "" {
		httpAddrs = defaultAdminAddrs
	}
	d := Drift{}
	d.admin = dAdmin{
		httpAddrs:      strings.TrimPrefix(httpAddrs, "http://"),
		lookupHTTPAddr: lookupHTTPAddr,
	}
	fmt.Println("Starting started at " + httpAddrs)
	go d.sysInterrupt()
	d.admin.startAdmin()
}

//NewPub will create new publisher
func NewPub(nsqDTCPAddrs []string) *Drift {
	for i := range nsqDTCPAddrs {
		nsqDTCPAddrs[i] = strings.TrimPrefix(nsqDTCPAddrs[i], "http://")
	}
	return &Drift{
		jobHandler:    nil,
		chanelHandler: make(map[string]JobHandler),
		consumers:     make(map[string][]*nsq.Consumer),
		pubAddrs:      nsqDTCPAddrs,
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
			case createAction, emptyAction, deleteAction, pauseAction, unpauseAction:
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
		if len(payload.NsqDTCPAddrs) == 0 {
			err = lib.VError("nsqd tcp address required")
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
	config := nsq.NewConfig()
	config.UserAgent = fmt.Sprintf("drift/%s", nsq.VERSION)
	producers := make(map[string]*nsq.Producer)
	for _, addr := range payload.NsqDTCPAddrs {
		var producer *nsq.Producer
		if producer, err = nsq.NewProducer(addr, config); err != nil {
			break
		}
		producers[addr] = producer
	}
	if err == nil {
		for _, producer := range producers {
			var b []byte
			if b, err = jsoniter.Marshal(payload.Data); err == nil {
				if err = producer.Publish(payload.Topic, b); err != nil {
					break
				}
			} else {
				break
			}
		}
	}
	for _, producer := range producers {
		producer.Stop()
	}
	data = "DONE"
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
