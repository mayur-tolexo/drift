# drift  <img src="https://user-images.githubusercontent.com/20511920/41496418-61a22e4c-715d-11e8-9456-3ef08a2af35d.png" width="30"/>
NSQ Producer/Consumer integration to drift your request smoothly.  
Add/Kill consumer over http request on any topic.  
Publish new request over http on any nsqd.  
[DOC](https://www.godoc.org/github.com/mayur-tolexo/drift)




### STEPS TO RUN drift
1. install [nsq](https://nsq.io/deployment/installing.html)
1. go get github.com/mayur-tolexo/drift
1. cd $GOPATH/src/github.com/mayur-tolexo/drift
1. install [godep](https://www.github.com/tools/godep)
1. godep restore
1. go run example/consumer.go
1. `[in new tab]` go get github.com/nsqio/nsq
1. go get github.com/golang/dep/cmd/dep
1. cd $GOPATH/src/github.com/nsqio/nsq/
1. dep ensure
1. cd $GOPATH/src/github.com/nsqio/nsq/apps/nsqlookupd
1. go run nsqlookupd.go
1. `[in new tab]` cd $GOPATH/src/github.com/nsqio/nsq/apps/nsqd
1. go run nsqd.go --lookupd-tcp-address=127.0.0.1:4160
1. `[in new tab]` cd $GOPATH/src/github.com/nsqio/nsq/apps/nsqadmin
1. go run main.go --lookupd-http-address=127.0.0.1:4161
1. open http://127.0.0.1:4171/ in browser
1. `[in new tab]` cd $GOPATH/src/github.com/mayur-tolexo/drift
1. go run example/producer.go
1. add new consumer as given below


### ADD NEW CONSUMER
*POST* localhost:1500/drift/v1/add/consumer/
```
{
  "lookup_address": [
    "127.0.0.1:4161"
  ],
  "topic_detail": [
    {
      "topic": "elastic",
      "channel": "v2.1"
    },
    {
      "topic": "elastic",
      "channel": "v6.2"
    }
  ],
  "max_in_flight": 200
}
```

### COUNT CONSUMERS OF A TOPIC ON SPECIFIC CHANNEL
```GET localhost:1500/drift/v1//consumer/?topic=elastic&channel=v2.1```

### COUNT ALL CONSUMERS OF A TOPIC
```GET localhost:1500/drift/v1//consumer/?topic=elastic```

### KILL CONSUMER
*POST* localhost:1500/drift/v1/kill/consumer/
```
{
  "topic": "elastic",
  "channel": "v2.1",
  "count":1
}
```


# Example
```
//printIT : function which consumer will call to execute
func printIT(value ...interface{}) error {
  fmt.Println(value)
  return nil
}

func main() {
  d := drift.NewConsumer(printIT)
  d.Start(1500)
}

```


# Handler
```
handler is a function to which the consumer will call.

FUNCTION DEFINATION:
func(value ...interface{}) error

Here PrintIT is a handler function. Define your own handler and pass it in the drift and you are ready to go.
```

