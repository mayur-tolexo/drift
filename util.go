package drift

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mayur-tolexo/drift/lib"
	nsq "github.com/nsqio/go-nsq"
	"github.com/rightjoin/aqua"
)

//HandleMessage will define the nsq handler method
func (th *TailHandler) HandleMessage(m *nsq.Message) error {
	return th.jobHandler(string(m.Body))
}

//vAddConsumer will validate add consumer request
func vAddConsumer(req aqua.Aide) (payload AddConstumer, err error) {
	req.LoadVars()
	err = lib.Unmarshal(req.Body, &payload)
	return
}

//pAddConsumer will process add consumer request
func (d *Drift) pAddConsumer(payload AddConstumer) (data interface{}, err error) {
	var c *nsq.Consumer
	config := nsq.NewConfig()
	maxInFlight := lib.GetPriorityValue(200, payload.MaxInFlight).(int)
	config.MaxInFlight = maxInFlight
	config.UserAgent = fmt.Sprintf("drift/%s", nsq.VERSION)
	for i := range payload.Topic {

		topic := payload.Topic[i].Topic
		channel := payload.Topic[i].Channel
		if channel == "" {
			rand.Seed(time.Now().UnixNano())
			channel = fmt.Sprintf("drift%06d#ephemeral", rand.Int()%999999)
		}

		if c, err = nsq.NewConsumer(topic, channel, config); err == nil {
			fmt.Println("Adding consumer for topic:", topic)
			c.AddHandler(&TailHandler{topicName: topic, jobHandler: d.jobHandler})
			if err = c.ConnectToNSQDs(payload.NsqDTCPAddrs); err != nil {
				err = lib.BadReqError(err)
				break
			}
			if err = c.ConnectToNSQLookupds(payload.LookupAddr); err != nil {
				err = lib.BadReqError(err)
				break
			}
			d.consumers = append(d.consumers, c)
		} else {
			err = lib.BadReqError(err)
			break
		}
	}
	if err == nil {
		data = "DONE"
	}
	return
}

//SystemInterrupt will handle system interrupt
func (d *Drift) SystemInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTSTP)
	fmt.Println("System Exit: ", <-c)
	for _, consumer := range d.consumers {
		consumer.Stop()
	}
	for _, consumer := range d.consumers {
		<-consumer.StopChan
	}
	os.Exit(1)
}
