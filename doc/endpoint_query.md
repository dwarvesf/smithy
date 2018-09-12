#  **POST** | Query 
```
<url>/databases/{database_name}/{table_name}/query
```
### Headers
| | |
|--|--|
| Content-Type | application/json
| Authorization | Bearer {access_token}|
### Body
|Fields| Type | Require | Description |
|--|--|--|--|
| fields | array of string | Yes | Fields's name  you want to query from database. Example: ```["id", "name"]``` |
| filter | object | No | The query's condition. Example: ```{ "operator": "=", "column_name": "name", "value": "Hieu Phan" }``` |
| order | array of string | No | Order the result. Example: ```["id", "desc"]``` |
| offset | int | No | The offset of query |
| limit | int | No | The limit of query |

Example:
```
{
    "fields": ["id", "name"],
    "filter": {
        "column_name": "name",
        "operator": "=",
        "value": "Hieu Phan"
    },
    "order": ["id", "desc"],
    "offset": 40,
    "limit": 20
}
```

### Response
#### Success
|Fields| Type | Description |
|--|--|--|
| status | string | Status of result. Ex: ```success``` |
| columns | array of string | Array of result's columns |
| rows | array of array | Array of query result |
| cols | array of object | Defination of columns in database |

Example:
```
{
    "status": "success",
    "columns": ["id", "name"],
    "rows": [
        [9, "Hieu Phan"],
        [10, "Hieu Phan"]
    ],
    "cols": [
        {
            "name": "id",
            "type": "int",
            "tags": "",
            "is_nullable": false,
            "is_primary": true,
            "default_value": "",
            "foreign_key": {
                "table": "",
                "foreign_column": ""
            }
        },
        {
            "name": "name",
            "type": "string",
            "tags": "",
            "is_nullable": true,
            "is_primary": false,
            "default_value": "",
            "foreign_key": {
                "table": "",
                "foreign_column": ""
            }
        }
    ]
}
```
#### Fail
```
{
    "error": "error message"
}
```