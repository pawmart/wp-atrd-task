# Not so Secret Service 

There is nothing to write about, it's a secret service after all :)

## To start
* ```docker-compose up``` (default port=7777)


## Example queries

* Add secret
    * `curl -XPOST localhost:7777/v1/secret/ -d "secret=ItsASecretDontTellAnyone&expireAfterViews=5&expireAfter=1"`

* Get Secret
    * JSON format: `curl -H "Accept: application/json" localhost:7777/v1/secret/603f9f9bd33088586f9072d4`
    * XML format : `curl -H "Accept: application/xml" localhost:7777/v1/secret/603f9f9bd33088586f9072d4`

