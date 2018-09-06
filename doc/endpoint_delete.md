# **DELETE** | Delete

```
<url>/tables/{table_name}/delete
```
### Headers
| | |
|--|--|
| Content-Type | application/json
| Authorization | Bearer {jwt_token}|

### Body
```javascript
{
    "filter": {
        "fields": {Tên những fields cần filter},
        "data": {Giá trị tương ứng của những fields cần filter}
    }
}
```

#### Sample
Xóa những trường trong bảng có id = 1
```javascript
{
   "filter": {
        "fields": ["id"],
        "data": ["1"]
   }
}
```

Xóa những trường trong bảng có id = 1 và name = meo con
```javascript
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
{"error": "chi tiết lỗi"}

