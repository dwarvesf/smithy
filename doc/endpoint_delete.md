# **DELETE** | Delete

```
<url>/tables/{table_name}/delete
```
### Headers
| Key|Value |
|--|--|
| Content-Type | application/json
| Authorization | Bearer {jwt_token}|

### Body
```
{
    "filter": {
        "fields": {fields name need to fill},
        "data": {data of fields need to fill}
    }
}
```

#### Sample
Delete record have id = 1
```
{
   "filter": {
        "fields": ["id"],
        "data": ["1"]
   }
}
```

Delete record have id = 1 and name = "Mèo con"
```
{
    "filter": {
        "fields": ["id", "name"],
        "data": ["1", "meo con"]
    }
} 
```

### Response
#### Success
```javascript
{"status": "success"}
```

#### Fail
```javascript
{"error": "error detail"}

