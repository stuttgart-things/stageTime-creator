# stuttgart-things/sweatShop-creator

dynamic rendering and creation of k8s-resources polled from redis streams

<details><summary><b>DEPLOYMENT TOOLS</b></summary>

* golang
* helm
* nerdctl
* kubectl
* redis-cli
* taskfile

</details>

## DEPLOY DEV CODE TO CLUSTER

<details><summary><b>DEPLOYMENT INCLUDING REDIS + TO DIFFERENT NAMESPACE</b></summary>

```
helm pull oci://eu.gcr.io/stuttgart-things/sweatshop-creator --version v0.1.44

cat <<EOF > creator.yaml
---
namespace: sweatshop-creator

tektonResources:
  enabled: false
pipelineRuns:
  enableRuns: false

redis:
  enabled: true
  sentinel:
    enabled: true
  master:
    service:
      type: ClusterIP
    persistence:
      enabled: false
      medium: ""
  replica:
    replicaCount: 1
    persistence:
      enabled: false
      medium: ""
  auth:
    password: <path:apps/data/sweatshop#redisPassword>

configmaps:
  creator:
    TEMPLATE_PATH: /templates
    REDIS_STREAM: sweatshop:manifests
    REDIS_SERVER: sweatshop-creator-redis-headless.sweatshop-creator.svc.cluster.local
    REDIS_PORT: "6379"

secrets:
  redis:
    name: redis
    secretKVs:
      REDIS_PASSWORD: <path:apps/data/sweatshop#redisPassword>

clusterRoleBindings:
  sweatshop-creator:
    subjects:
      - kind: ServiceAccount
        name: sweatshop-creator
        namespace: sweatshop-creator
roleBindings:
  sweatshop-creator:
    subjects:
      - kind: ServiceAccount
        name: sweatshop-creator
        namespace: sweatshop-creator
EOF

helm upgrade --install sweatshop-creator oci://eu.gcr.io/stuttgart-things/sweatshop-creator --version v0.1.44 --values ankit.yaml -n sweatshop-creator --create-namespace
```

</details>

<details><summary><b>CHECK REDIS DATA w/ CLI</b></summary>

```
# Install redis-cli #
sudo apt-get update
sudo apt-get install redis

kubectl -n sweatshop port-forward creator-redis-node-0 28015:6379
redis-cli -h 127.0.0.1 -p 28015 -a ankit
# CHECK ALL REDIS KEYS
KEYS *
# READ STREAM
XREAD COUNT 2 STREAMS sweatshop:manifests writers 0-0 0-0
# DELETE STREAM
DEL sweatshop:manifests
```

</details>


## TEST SERVICE LOCALLY (OUTSIDE CLUSTER)

<details><summary><b>START CONSUMER</b></summary>

```
export KUBECONFIG=~/.kube/dev11
export TEMPLATE_PATH=~/projects/go/src/github/sweatShop-creator/tests
export TEMPLATE_NAME=job-template.yaml
export REDIS_STREAM=sweatshop:test
export REDIS_PASSWORD=<SET-ME>
export REDIS_SERVER=redis-pve.labul.sva.de
export REDIS_PORT=6379
task run
```

</details>

<details><summary><b>START TEST PRODUCING (EXTERNAL REDIS)</b></summary>


```
# kubectl -n sweatshop-redis port-forward redis-sweatshop-deployment-node-0 28015:6379
task run-test
```

</details>

<details><summary><b>START TEST PRODUCING (REDIS INSIDE CLUSTER)</b></summary>

```
kubectl -n <REDIS-NS> port-forward redis-sweatshop-deployment-node-0 <HOST-PORT>:<CONTAINER-PORT>

# kubectl -n sweatshop-redis port-forward redis-sweatshop-deployment-node-0 28015:6379

export REDIS_STREAM=sweatshop:manifests
export REDIS_PASSWORD=<SETME>
export REDIS_SERVER=127.0.0.1
export REDIS_PORT=28015 # HOST-PORT
task run-test-producer
```

</details>

<details><summary><b>VERIFY REDIS</b></summary>

```
redis-cli -h <REDIS_SERVER>-p <HOST-PORT> -a <SETME>

# redis-cli -h 127.0.0.1 -p 28015 -a test

KEYS *
# GET VALUE
GET <KEYNAME>
# GET STREAM
XREAD COUNT 2 STREAMS <STREAM-NAME> writers 0-0 0-0
```

</details>


## LICENSE

<details><summary><b>APACHE 2.0</b></summary>

Copyright 2023 patrick hermann.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

</details>

Author Information
------------------
Patrick Hermann, stuttgart-things 06/2023
