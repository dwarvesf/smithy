#  **POST** | Create 
```
<url>/tables/{table_name}/create
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
    "fields": ["name", "age"],
    "data": ["Hieu Phan", 21]
}
```
#### With relationship
```
{
    "fields": [
        "name",
        "age",
        {"books": ["name", "description"]}
    ],
    "data": [
        "Hieu Phan",
        21,
        [
            ["How to be a handsome man", "How to..."],
            ["How to be Spiderman", "Spider is a..."]
        ]
    ]
}
```
### Response
#### Success
```
{
    "status": "success",
    "data": {
        "id": 1,
        "name": "Hieu Phan",
        "age": 21,
        "books": [
            {
                "id": 1,
                "name": "How to be a handsome man",
                "description": "How to..."
            },
            {
                "id": 2,
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