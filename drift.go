package drift

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	jsoniter "github.com/json-iterator/go"
	nsq "github.com/nsqio/go-nsq"
	"github.com/nsqio/nsq/nsqadmin"
	"github.com/rightjoin/aqua"
	"github.com/tolexo/tachyon/lib"
)

//AddTopicHandler will add a new handler with the given topic
func (d *Drift) AddTopicHandler(topic string, jobHandler JobHandler) {
	d.chanelHandler[hash(topic, allKey)] = jobHandler
}

//AddChanelHandler will add a new handler with the channel of given topic
func (d *Drift) AddChanelHandler(topic, channel string, jobHandler JobHandler) {
	d.chanelHandler[hash(topic, channel)] = jobHandler
}

//Start will start the drift server
func (d *Drift) Start(port int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in drift", r)
			os.Exit(1)
		}
	}()
	runtime.GOMAXPROCS(runtime.NumCPU())
	d.Server = aqua.NewRestServer()
	d.Server.Port = port
	d.Server.AddService(&ds{drift: d})
	go d.sysInterrupt()
	d.Server.Run()
}

//Publish will broadcast the data to the nsqd
func (d *Drift) Publish(topic string, data interface{}) (resp interface{}, err error) {
	payload := Publish{
		NsqDHTTPAddrs: d.pubAddrs,
		Topic:         topic,
		Data:          data,
	}
	resp, err = pPublishReq(payload)
	return
}

//addConsumer will process add consumer request
func (d *Drift) addConsumer(payload AddConstumer) (data interface{}, err error) {
	var c *nsq.Consumer
	config := nsq.NewConfig()
	config.MaxInFlight = lib.GetPriorityValue(200, payload.MaxInFlight).(int)
	config.UserAgent = fmt.Sprintf("drift/%s", nsq.VERSION)
	for i := range payload.Topic {

		topic := payload.Topic[i].Topic
		channel := getChannel(payload.Topic[i].Channel)
		n := lib.GetPriorityValue(payload.Topic[i].Count, 1).(int)
		handler := d.getHandler(topic, channel)
		if handler == nil {
			continue
		}
		for j := 0; j < n; j++ {
			if c, err = nsq.NewConsumer(topic, channel, config); err == nil {
				fmt.Println("Adding consumer for topic:", topic)
				c.AddHandler(&tailHandler{topicName: topic, jobHandler: handler})
				if err = c.ConnectToNSQDs(payload.NsqDTCPAddrs); err != nil {
					err = lib.BadReqError(err)
					break
				}
				if err = c.ConnectToNSQLookupds(payload.LookupHTTPAddr); err != nil {
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
		if err != nil {
			break
		}
	}
	if payload.StartAdmin && !d.admin.adminRunning {
		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		d.admin.httpAddrs = defaultAdminAddrs
		d.admin.lookupHTTPAddr = payload.LookupHTTPAddr
		d.admin.nsqDTCPAddrs = payload.NsqDTCPAddrs
		d.admin.adminUser = []string{"mayur-drift-" + hostname}
		d.admin.aclHTTPHeader = "X-Drift"
		go d.admin.startAdmin()
	}
	return
}

func (d *Drift) getHandler(topic, channel string) (handler JobHandler) {
	if chHandler, exists := d.chanelHandler[hash(topic, channel)]; exists {
		handler = chHandler
	} else if chHandler, exists := d.chanelHandler[hash(topic, allKey)]; exists {
		handler = chHandler
	} else {
		handler = d.jobHandler
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

//killConsumer will process kill consumer request
func (d *Drift) killConsumer(payload KillConsumer) (data interface{}, err error) {
	c := 0
	key := hash(payload.Topic, payload.Channel)
	if _, exists := d.consumers[key]; exists {
		for _, consumer := range d.consumers[key] {
			consumer.Stop()
			c++
			if c == payload.Count {
				break
			}
		}
		s := 0
		for _, consumer := range d.consumers[key] {
			<-consumer.StopChan
			s++
			if s == c {
				d.consumers[key] = d.consumers[key][c:]
				break
			}
		}
		if c < payload.Count {
			d.consumers[key] = nil
		}
	}
	data = "DONE"
	return
}

//vStartAdmin will validate start admin request
func (d *dAdmin) vStartAdmin(req aqua.Aide) (err error) {
	req.LoadVars()
	var payload AddAdmin
	if err = lib.Unmarshal(req.Body, &payload); err == nil {
		if payload.HTTPAddrs == "" {
			payload.HTTPAddrs = defaultAdminAddrs
		}
		d.httpAddrs = payload.HTTPAddrs
		d.adminUser = payload.AdminUser
		d.lookupHTTPAddr = payload.LookupHTTPAddr
		d.nsqDTCPAddrs = payload.NsqDTCPAddrs
		d.aclHTTPHeader = payload.ACLHTTPHeader
		d.notificationHTTPEndpoint = payload.NotificationHTTPEndpoint
	}
	return
}

//startAdmin will add new admin
func (d *dAdmin) startAdmin() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in admin", r)
			os.Exit(1)
		}
	}()
	signalChan := make(chan os.Signal, 1)
	d.exitAdmin = make(chan int)
	go func() {
		<-signalChan
		d.exitAdmin <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	opts := nsqadmin.NewOptions()
	opts.HTTPAddress = d.httpAddrs
	opts.AdminUsers = d.adminUser
	opts.NSQLookupdHTTPAddresses = d.lookupHTTPAddr
	opts.NSQDHTTPAddresses = d.nsqDTCPAddrs
	opts.AclHttpHeader = d.aclHTTPHeader
	opts.NotificationHTTPEndpoint = d.notificationHTTPEndpoint

	nsqadmin := nsqadmin.New(opts)
	nsqadmin.Main()
	d.adminRunning = true
	<-d.exitAdmin
	d.adminRunning = false
	nsqadmin.Exit()
	d.exitAdmin <- 1
}

//doAction will validate start admin request
func (d *dAdmin) doAction(payload Admin, aclValue string) (data interface{}, err error) {
	var (
		b    []byte
		req  *http.Request
		resp *http.Response
	)
	reqBody := map[string]string{"action": payload.Action}
	if b, err = jsoniter.Marshal(reqBody); err == nil {
		method := "POST"
		if payload.Action == "delete" {
			method = "DELETE"
		}
		URL := fmt.Sprintf("http://%v/api/topics/%v", d.httpAddrs, payload.Topic)
		if payload.Channel != "" {
			URL += "/" + payload.Channel
		}
		if req, err = http.NewRequest(method,
			URL, bytes.NewBuffer(b)); err == nil {
			HTTPClient := &http.Client{}
			if d.aclHTTPHeader == "" {
				d.aclHTTPHeader = "X-Forwarded-User"
			}
			req.Header.Set(d.aclHTTPHeader, aclValue)
			if resp, err = HTTPClient.Do(req); err == nil {
				defer resp.Body.Close()
				// if resp.StatusCode == http.StatusOK {
				bodyBytes, _ := ioutil.ReadAll(resp.Body)
				if err = jsoniter.Unmarshal(bodyBytes, &data); err != nil {
					err = lib.UnmarshalError(err)
				}
				// }
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
