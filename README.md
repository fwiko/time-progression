# Time Progression API

# Options

### **Format**

The format in which data is returned can be specified with a path parameter.

- `second` - returns the values as seconds
- `percent` - returns the values as a percentage using second precision

### **TimeZone**

The `timezone` query parameter can be used to specify the timezone in which the data is calculated/returned.

e.g. `?timezone=Asia/Singapore` will return the data in New York time.

# Usage

### **Request**

```bash
curl http://localhost:8080/api/percent?timezone=Asia/Singapore -H "Accept: application/json"
```

### **Response**

```json
{
  "timezone": "Asia/Singapore",
  "timestamp": "Wed, 25 May 2022 02:53:23 +08",
  "result": {
    "year": 39.75901509386098,
    "month": 83.73468364197531,
    "week": 30.291501322751323,
    "day": 12.04050925925926,
    "hour": 88.97222222222221,
    "minute": 38.333333333333336
  }
}
```

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
