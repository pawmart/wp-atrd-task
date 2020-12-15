# Sample app

# WIP, not finished yet

Simple app to store and fetch secrets saved to Redis.  
Before secret it's stored, it's encrypted server side with AES 128bit

## Dev docs

### Running dev env

```bash
docker-compose up
# REST API http://localhost:3000
# Redis web interface http://admin:admin@localhost:6380
```

### Running tests

```bash
cd cmd/server
godog ../../features
``````
