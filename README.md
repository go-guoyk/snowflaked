# snowflaked

snowflake service in nrpc

## Environment Variables

* `BIND`, port bind, default to `:3000`
* `CLUSTER_ID`, 5 bits unsigned integer, should not be `0`
* `WORKER_ID`, 5 bits unsigned integer, should not be `0`, automatically load k8s stateful-set sequence id from hostname

## HTTP

* Single

    ```
    GET /snowflake/next_id
  
    QUERY
    format = ("str_oct" | "str_dec" | "str_hex" | null)
  
    BODY
    {"id":"1234567890"}
    ```
    
* Batch

    ```
    GET /snowflake/next_ids
  
    QUERY
    size = number of ids to return
    format = ("str_oct" | "str_dec" | "str_hex" | null)
  
    BODY
    {"ids":["1234567890"]}
    ```
     

## Credits

Guo Y.K. <hi@guoyk.net>
