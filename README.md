# Time Progression API

https://progression.didntlaugh.com/

# Usage

### Request

```bash
curl http://localhost:8080/?timezone=US-Hawaii -H "Accept: application/json"
```

### Response

```json
{
  "timezone": "US/Hawaii",
  "timestamp": "Tue, 24 May 2022 03:01:10 HST",
  "result": {
    "year": 39.4865233384069,
    "month": 2412.5810185185182,
    "week": 1.8052248677248677,
    "day": 12.581018518518519,
    "hour": 1.9444444444444444,
    "minute": 16.666666666666664
  }
}
```

_All values are calculated to the second._

# Containerisation with Docker

This will use the included Dockerfile to create an install dependencies and build the application.

### Build Image

```bash
docker build -t time-progression .
```

### Start Container

```bash
docker run -d -p 80:80 \
-v $(pwd)/config.json:/usr/src/app/config.json \
--name time-progression time-progression:latest
```
