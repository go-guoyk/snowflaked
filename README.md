# snowflaked

snowflake service in nrpc

## Environment Variables

* `BIND`, port bind, default to `:3000`
* `CLUSTER_ID`, 5 bits unsigned integer, should not be `0`
* `WORKER_ID`, 5 bits unsigned integer, should not be `0`, automatically load k8s stateful-set sequence id from hostname

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
