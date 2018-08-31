## Hook scripting language 

### What is hook scripting language?

- Hook scripting language is a scripting language embedded in model configuration,  It used to enable some action was happen in before or after user interact with model

- Currently hook language is [anko](https://github.com/mattn/anko), it syntax is very same with `golang`.

### How to use

- Hook script can be added by editing configuration file (in model section) like the example below

```
table_name: "users"
...
...
hooks:
  before_create:
    enable: true
    content: "println(\"this is a script\")"
```
- Hook can also can registered in runtime by calling POST endpoint `{tableName}/hooks`

- Currently only support hook for: `before_create`, `after_create`, `before_update`, `after_update`, `before_delete`, `after_delete`

### Script support function

- `println`: is same with `fmt.Println`
- `ctx`: get current data value in a object (map key value of current model), change value of `ctx` will also change data of column in current context, ex: `this()["name"]` will return value of `name` column of current context model
- `callAPI`: is a utility to call API to a URL, have 3 arguments: method, URL, data (currently only support json string format), return `false`, `true` if it correct
