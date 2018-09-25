#  **PUT** | Update 
```
<url>/databases/{database_name}/{table_name}/update
```
### Headers
| | |
|--|--|
| Content-Type | application/json
| Authorization | Bearer {access_token}|
### Body
#### Simple
```
{
    "fields": ["id","name", "age"],
    "data":   [ 1,  "Anmoc", 21],
}
```
#### With relationship
```
{
    "fields": [
        "id",
        "name",
        "age",
        {"books": ["id", "name", "description", "author_id"]}
    ],
    "data": [
        1,
        "An Moc",
        21,
        [
            [1, "How to be a handsome man", "How to...", 1],
            [2, "How to be Spiderman", "Spider is a...", 2]
        ]
    ]
    "primary_key": "1"
}
```
### Response
#### Success
```
{
    "status": "success",
    "data": {
        "name": "An Moc",
        "age": 21,
        "books": [
            {
                "author_id": 1,
                "name": "How to be a handsome man",
                "description": "How to..."
            },
            {
                "author_id": 2,
                "name": "How to be Spiderman",
                "description": "Spider is a..."
            }
        ]
    }
}
```
#### Fail
```
{
    "error": "error message"
}
```