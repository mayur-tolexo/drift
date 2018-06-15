# drift
NSQ integration to drift your request smoothly

### STEPS TO INSTALL drift
1. go get github.com/mayur-tolexo/drift
1. install [godep](https://www.github.com/tools/godep)
1. godep restore
1. go run main.go


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
