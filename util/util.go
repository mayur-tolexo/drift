package util

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mayur-tolexo/drift/lib"
	"github.com/mayur-tolexo/drift/model"
	nsq "github.com/nsqio/go-nsq"
	"github.com/rightjoin/aqua"
)

var consumers []*nsq.Consumer

//TailHandler will implement the nsq handler
type TailHandler struct {
	TopicName string
}

//HandleMessage will define the nsq handler method
func (th *TailHandler) HandleMessage(m *nsq.Message) error {
	_, err := os.Stdout.Write(m.Body)
	if err != nil {
		log.Fatalf("ERROR: failed to write to os.Stdout - %s", err)
	}
	_, err = os.Stdout.WriteString("\n")
	if err != nil {
		log.Fatalf("ERROR: failed to write to os.Stdout - %s", err)
	}
	return nil
}

//VAddConsumer will validate add consumer request
func VAddConsumer(req aqua.Aide) (payload model.AddConstumer, err error) {
	req.LoadVars()
	err = lib.Unmarshal(req.Body, &payload)
	return
}

//PAddConsumer will process add consumer request
func PAddConsumer(payload model.AddConstumer) (data interface{}, err error) {
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
			c.AddHandler(&TailHandler{TopicName: topic})
			if err = c.ConnectToNSQDs(payload.NsqDTCPAddrs); err != nil {
				err = lib.BadReqError(err)
				break
			}
			if err = c.ConnectToNSQLookupds(payload.LookupAddr); err != nil {
				err = lib.BadReqError(err)
				break
			}
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
func SystemInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTSTP)
	fmt.Println("System Exit: ", <-c)
	for _, consumer := range consumers {
		consumer.Stop()
	}
	for _, consumer := range consumers {
		<-consumer.StopChan
	}
	os.Exit(1)
}
