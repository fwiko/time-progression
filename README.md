# Time Progression API
#### A simple API allowing the retreival of percentage data about how far we are through different units of time. All respective data is returned in JSON format.

## Optional parameters

### `format`
- ### `"round"`
    - will round all returned values, using the traditional rounding method.
- ### `"floor"`
    - will round all returned values down, no matter the decimal value.
- ### `"ceil"`
    - will round all returned values up, no matter the decimal value.

### `timezone`
- ### e.g. `"America-Cancun"`
    - All available timezones can be retreived by visiting `/timezones`


## Example Usage

#### Request
```console
curl http://localhost:5000/?timezone=US-Hawaii -H "Accept: application/json"
```

#### Response
```json
{
    "data":{
        "day":56.52314814814815,
        "hour":56.55555555555556,
        "minute":93.33333333333333,
        "month":17.9523596176822,
        "week":36.646164021164026,
        "year":1.5247209538305428
    },
    "datestring":"2022-01-05 13:33:56",
    "timezone":"US/Hawaii"
}
```
### *Percentage values are calculated using second precision.*

