# drift
NSQ integration to drift your request smoothly


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
