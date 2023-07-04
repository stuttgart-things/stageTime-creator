# stuttgart-things/sweatShop-creator

dynamic rendering and creation of k8s-resources polled from redis streams

## DEPLOY TO CLUSTER

<details><summary><b>REDIS</b></summary>

</details>

<details><summary><b>DEPLOYMENT</b></summary>

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

<details><summary><b>START TEST PRODUCING</b></summary>

```
export REDIS_STREAM=sweatshop:test
export REDIS_PASSWORD=<SETME>
export REDIS_SERVER=redis-pve.labul.sva.de
export REDIS_PORT=6379
task run-test-producer
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
