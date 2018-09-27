<img align="right" src="https://user-images.githubusercontent.com/20511920/41496418-61a22e4c-715d-11e8-9456-3ef08a2af35d.png" width="50">

# drift 
NSQ Producer/Consumer integration to drift your request smoothly.  
Add/Kill consumer over http request on any topic.  
Publish new request over http on any nsqd.  
[- DOCUMENTATION](https://www.godoc.org/github.com/mayur-tolexo/drift)




### STEPS TO RUN drift
1. install [nsq](https://nsq.io/deployment/installing.html)
1. go get github.com/mayur-tolexo/drift
1. cd $GOPATH/src/github.com/mayur-tolexo/drift
1. install [godep](https://www.github.com/tools/godep)
1. godep restore
1. go run example/drift.go
1. `[in new tab]` cd $GOPATH/src/github.com/mayur-tolexo/drift
1. go run nsqlookup/nsqlookup.go
1. `[in new tab]` cd $GOPATH/src/github.com/mayur-tolexo/drift
1. go run nsqd/nsqd.go --lookupd-tcp-address=127.0.0.1:4160
1. `[in new tab]` cd $GOPATH/src/github.com/mayur-tolexo/drift
1. go run example/producer.go
1. add new consumer as mention below
1. start admin as mention below
1. open http://127.0.0.1:4171/ in browser

### START ADMIN
*POST* localhost:1500/drift/v1/start/admin/
```
{
  "lookup_http_address": ["127.0.0.1:4161"],
  "user": ["drift-user"],
  "acl_http_header": "admin-user"
}
```

### STOP ADMIN
```GET localhost:1500/drift/v1/stop/admin/```

### ADD NEW CONSUMER
*POST* localhost:1500/drift/v1/add/consumer/
```
{
  "lookup_http_address": [
    "127.0.0.1:4161"
  ],
  "topic_detail": [
    {
      "topic": "elastic",
      "channel": "v2.1",
      "count": 1
    },
    {
      "topic": "elastic",
      "channel": "v6.2",
      "count": 2
    }
  ],
  "max_in_flight": 200
}
```

### COUNT CONSUMERS OF A TOPIC ON SPECIFIC CHANNEL
```GET localhost:1500/drift/v1/consumer/?topic=elastic&channel=v2.1```

### COUNT ALL CONSUMERS OF A TOPIC
```GET localhost:1500/drift/v1/consumer/?topic=elastic```

### KILL CONSUMER
*POST* localhost:1500/drift/v1/kill/consumer/
```
{
  "topic": "elastic",
  "channel": "v2.1",
  "count":1
}
```

### PUBLISH REQUEST
*POST* localhost:1500/drift/v1/pub/request/
```
{
  "nsqd_tcp_address": ["127.0.0.1:4150"],
  "topic": "elastic",
  "data": "This is a test over http"
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

