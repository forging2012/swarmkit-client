# swarmkit-client

The http client for [swarmkit](https://github.com/docker/swarmkit)

## usage

### run

```
swarmkit-client -s /tmp/manager1/swarm.sock
```

### api

#### node

##### ls node

```
# ls all nodes
curl -X GET http://localhost:8888/nodes

# inspect node and display task
curl -X GET http://localhost:8888/nodes/{nodeid:.*}?all=1

# accept node
curl -X POST http://localhost:8888/nodes/accept

# remove node
// DELETE /nodes/{nodeid:.*}
```

#### service

```
# create service
curl -X POST -d '{"name":"redis", "image":"redis:3.0.5"}' http://localhost:8888/services/create

# ls all running services
curl -X GET http://localhost:8888/services

# inspect service
curl -X GET http://localhost:8888/services/{serviceid:.*}

# update service
curl -X POST -d '{...}' http://localhost:8888/services/{serviceid:.*}/update

# delete service
curl -X DELETE http://localhost:8888/services/7zyp89z8zefrq96jga06vho5f
```