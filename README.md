# Sample app

# WIP, not finished yet

Simple app to store and fetch secrets saved to Redis.  
Before writing secret to DB, it's encrypted server side with AES CFB 128bit.  

It supports expiration by number of requests (must be set to more than 0) and minutes entered by user.  
Secret is expired whatever comes first, depleted view counter or time

App is created for being used as single instance. 
For horizontal scaling and strong isolation guarantees it would need to have distributed locks for view counter expiration (can be done by redis, etcd or consul)


## "Prod" env build

```bash
docker build -t wp-atrd-task .
docker run --rm -e "HTTP_DOCS_DIR=/bin/api/swagger/" -p 3000:3000 wp-atrd-task
```

## Dev docs

### Running dev env

```bash
docker-compose up
# REST API http://localhost:3000
# Redis web interface http://admin:admin@localhost:6380
```

### Tests

```bash
cd cmd/server
godog ../../features
``````

#### Manual

```
go run github.com/systemz/wp-atrd-task/cmd/server
```

```bash
curl -XPOST "http://localhost:3000/v1/secret" -d "secret=testtest&expireAfter=1&expireAfterViews=5"
curl http://localhost:3000/v1/secret/216da4-047b-491e-823d-45787d6ea792
```

### Notes

If you don't set custom AES key, it will be automatically generated at server start.  
Warning, this will cause error in decryption after app restart due to AES integrity check.  
Make sure to set your own 16 character key via ENV var `AES_KEY`

Due to time constraints and fact that this isn't production app by design, there are some issues and areas to improvement:
- error and edge cases handling
- smaller DB footprint with using other encoding than json, eg. https://github.com/vmihailenco/msgpack and not using base64 for secret
- less layers of converting secret when setting/getting DB record
