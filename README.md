# stuttgart-things/stageTime-creator

dynamic rendering and creation of k8s manifests/resources polled from redis streams/json

## DEV-TASKS

```bash
task --list: Available tasks for this project:
* build:               Build the app
* build-image:         Build image
* git-push:            Commit & push the module
* lint:                Lint code
* package:             Update Chart.yaml and package archive
* push:                Push to registry
* release:             Build amd release to github w/ goreleaser
* run:                 Run app
* run-container:       Run container
* run-test:            Run test-producer
* tag:                 Commit, push & tag the module
* test:                Test code
```

## HELMFILE-DEPLOYMENTS

<details><summary>SET VAULT CONNECTION</summary>

```bash
export VAULT_ADDR=https://${VAULT_FQDN}}
export VAULT_NAMESPACE=root

# APPROLE AUTH
export VAULT_AUTH_METHOD=approle
export VAULT_ROLE_ID=${VAULT_ROLE_ID}
export VAULT_SECRET_ID=${VAULT_SECRET_ID}

# TOKEN AUTH
export VAULT_AUTH_METHOD=token #default
export VAULT_TOKEN=${VAULT_TOKEN}
```

</details>

<details><summary>RENDER/APPLY</summary>

```bash
helmfile template --environment labul-pve-dev
helmfile sync --environment labul-pve-dev
```

</details>

## DEPLOY DEV CODE TO CLUSTER

<details><summary><b>DEPLOYMENT</b></summary>

```bash
helm pull oci://eu.gcr.io/stuttgart-things/stagetime-creator --version v0.1.44
```

```yaml
cat <<EOF > stageTime-creator.yaml
---
secrets:
  redis-connection:
    name: redis-connection
    labels:
      app: stagetime-server
    dataType: stringData
    secretKVs:
      REDIS_SERVER: redis-stack-deployment-headless.redis-stack.svc.cluster.local
      REDIS_PORT: 6379
      REDIS_PASSWORD: <PASSWORD>
EOF
```

```bash
helm upgrade --install stagetime-creator oci://eu.gcr.io/stuttgart-things/stagetime-creator --version v0.1.44 --values stageTime-creator.yaml -n stagetime-creator --create-namespace
```

</details>

<details><summary><b>CHECK REDIS DATA w/ CLI</b></summary>

```
# Install redis-cli #
sudo apt-get update
sudo apt-get install redis

kubectl -n sweatshop port-forward redis-sweatshop-deployment-node-0 28015:6379 -n sweatshop-redis
redis-cli -h 127.0.0.1 -p 28015 -a ${PASSWORD}
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
export TEMPLATE_PATH=~/projects/go/src/github/stageTime-creator/tests
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

export REDIS_STREAM=stagetime-revisionruns
export REDIS_PASSWORD=<SETME>
export REDIS_SERVER=127.0.0.1
export REDIS_PORT=28015 # HOST-PORT
go run tests/test-json-set.go
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
