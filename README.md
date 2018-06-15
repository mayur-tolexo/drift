# drift
NSQ integration to drift your request smoothly

### STEPS TO INSTALL drift
1. install [nsq](https://nsq.io/deployment/installing.html)
1. go get github.com/mayur-tolexo/drift
1. install [godep](https://www.github.com/tools/godep)
1. godep restore
1. go run main.go
1. go get github.com/nsqio/nsq
1. cd $GOPATH/src/github.com/nsqio/nsq/apps/nsqlookupd
1. go run nsqlookupd.go
1. cd $GOPATH/src/github.com/nsqio/nsq/apps/nsqd
1. go run nsqd.go --lookupd-tcp-address=127.0.0.1:4160
1. cd $GOPATH/src/github.com/nsqio/nsq/apps/nsqadmin
1. go run main.go --lookupd-http-address=127.0.0.1:4161
1. open http://127.0.0.1:4171/ in browser


### ADD NEW CONSUMER
*POST* localhost:1500/drift/v1/add/consumer/
```
{
  "loopup_address": [
    "127.0.0.1:4161"
  ],
  "topic_detail": [
    {
      "topic": "elastic",
      "channel": "v2.1"
    },
    {
      "topic": "elastic",
      "channel": "v2.1"
    }
  ],
  "max_in_flight": 200
}
```
