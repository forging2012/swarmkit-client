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
curl -X GET http://localhost:8888/nodes
[
    {
        "id":"711dquxq93v1hwi721ysbr75n",
        "meta":{
            "version":{"index":23},
            "created_at":{"seconds":1465989057,"nanos":709750008},
            "updated_at":{"seconds":1466041314,"nanos":829778302}
        },
        "spec":{
            "annotations":{},
            "role":1,
            "membership":1
        },
        "description":{
            "hostname":"node-1",
            "platform":{"architecture":"x86_64","os":"linux"},
            "resources":{"nano_cpus":2000000000,"memory_bytes":2098659328},
            "engine":{
                "engine_version":"1.11.2",
                "plugins":[
                    {"type":"Volume","name":"local"},
                    {"type":"Network","name":"overlay"},
                    {"type":"Network","name":"bridge"},
                    {"type":"Network","name":"null"},
                    {"type":"Network","name":"host"}
                ]
            }
        },
        "status":{"state":2},
        "manager_status":{"raft_id":5128967476676048958,"addr":"10.10.16.200:4242","leader":true,"reachability":2},
        "attachment":{
            "network":{
                "id":"7ngh108lpqn51iysue8d7eaa7",
                "meta":{
                    "version":{"index":7},
                    "created_at":{"seconds":1465989057,"nanos":787336869},
                    "updated_at":{"seconds":1465989057,"nanos":788146300}
                },
                "spec":{
                    "annotations":{
                        "name":"ingress",
                        "labels":{"com.docker.swarm.internal":"true"}
                    },
                    "driver_config":{},
                    "ipam":{
                        "driver":{},
                        "configs":[{"subnet":"10.255.0.0/16","gateway":"10.255.0.1"}]
                    }
                },
                "driver_state":{
                    "name":"overlay",
                    "options":{"com.docker.network.driver.overlay.vxlanid_list":"256"}
                },
                "ipam":{
                    "driver":{"name":"default"},
                    "configs":[{"subnet":"10.255.0.0/16","gateway":"10.255.0.1"}]
                }
            },
            "addresses":["10.255.0.3/16"]
        },
        "certificate":{
            "role":1,
            "status":{"state":3},
            "cn":"711dquxq93v1hwi721ysbr75n"
        }
    }
]

```

#### service

##### create service

```
curl -X POST -d '{"name":"redis", "image":"redis:3.0.5"}' http://localhost:8888/services/create
{
    "id":"7zyp89z8zefrq96jga06vho5f",
    "meta":{
        "version":{"index":31},
        "created_at":{"seconds":1466053992,"nanos":848560070},
        "updated_at":{"seconds":1466053992,"nanos":848560070}
    },
    "spec":{
        "annotations":{"name":"redis"},
        "task":{
            "Runtime":{
                "Container":{"image":"redis:3.0.5","mounts":null}
            }
        },
        "Mode":{
            "Replicated":{"replicas":1}
        },
        "endpoint":{}
    }
}
```