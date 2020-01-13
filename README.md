# snowflaked

snowflake service in nrpc

## Usage

```shell script
./snowflaked -bind :3000 -cluster-id 1 -worker-id 1
```

## NRPC

* `snowflake,create`

    **Response**

    `{"id":4400868305313792}`
   
* `snowflake,create_s`

    **Response**

    `{"id":"4400868305313792"}`
    
* `snowflake,batch`

    **Request**
    
    `{"size":2}`
    
    **Response**

    `{"id":[4400868305313792,4400868305313792]}`

* `snowflake,batch_s`

    **Request**
    
    `{"size":2}`
    
    **Response**

    `{"id":["4400868305313792","4400868305313792"]}`


## Credits

Guo Y.K. <hi@guoyk.net>
