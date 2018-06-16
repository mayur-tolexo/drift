package drift

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/mayur-tolexo/drift/lib"
	nsq "github.com/nsqio/go-nsq"
	"github.com/rightjoin/aqua"
)

const (
	driftApp = "-DRIFT-"
	allKey   = "ALL"
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
		pubAddrs:      nsqDHttpAddrs,
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

//vProduceReq will validate produce request
func vPublishReq(req aqua.Aide) (payload Publish, err error) {
	req.LoadVars()
	err = lib.Unmarshal(req.Body, &payload)
	return
}

//pPublishReq will process the producer request
func pPublishReq(payload Publish) (data interface{}, err error) {
	var (
		b    []byte
		req  *http.Request
		resp *http.Response
	)
	if b, err = jsoniter.Marshal(payload.Data); err == nil {
		URL := fmt.Sprintf("%v/pub?topic=%v", payload.NsqDHttpAddrs, payload.Topic)
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

//pAddConsumer will process add consumer request
func (d *Drift) pAddConsumer(payload AddConstumer) (data interface{}, err error) {
	var c *nsq.Consumer
	config := nsq.NewConfig()
	config.MaxInFlight = lib.GetPriorityValue(200, payload.MaxInFlight).(int)
	config.UserAgent = fmt.Sprintf("drift/%s", nsq.VERSION)
	for i := range payload.Topic {

		topic := payload.Topic[i].Topic
		channel := getChannel(payload.Topic[i].Channel)
		handler := d.getHandler(topic, channel)
		if handler == nil {
			continue
		}

		if c, err = nsq.NewConsumer(topic, channel, config); err == nil {
			fmt.Println("Adding consumer for topic:", topic)
			c.AddHandler(&tailHandler{topicName: topic, jobHandler: handler})
			if err = c.ConnectToNSQDs(payload.NsqDTCPAddrs); err != nil {
				err = lib.BadReqError(err)
				break
			}
			if err = c.ConnectToNSQLookupds(payload.LookupAddr); err != nil {
				err = lib.BadReqError(err)
				break
			}
			key := hash(topic, channel)
			d.consumers[key] = append(d.consumers[key], c)
			data = "DONE"
		} else {
			err = lib.BadReqError(err)
			break
		}
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

func (d *Drift) getHandler(topic, channel string) (handler JobHandler) {
	if chHandler, exists := d.chanelHandler[hash(topic, channel)]; exists {
		handler = chHandler
	} else if chHandler, exists := d.chanelHandler[hash(topic, allKey)]; exists {
		handler = chHandler
	}
	return
}

func hash(a, b string) string {
	return a + driftApp + b
}

func decrypt(h string) (topic string, channel string) {
	val := strings.Split(h, driftApp)
	vLen := len(val)
	if vLen > 0 {
		topic = val[0]
	}
	if vLen > 1 {
		channel = val[1]
	}
	return
}

//sysInterrupt will handle system interrupt
func (d *Drift) sysInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTSTP)
	fmt.Println("System Exit: ", <-c)
	for _, topic := range d.consumers {
		for _, consumer := range topic {
			consumer.Stop()
		}
	}
	for _, topic := range d.consumers {
		for _, consumer := range topic {
			<-consumer.StopChan
		}
	}
	os.Exit(1)
}
