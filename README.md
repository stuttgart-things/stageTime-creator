# stuttgart-things/sweatShop-creator

renders & creates k8s resources


## TEST SERVICE LOCALLY

<details><summary><b>START CONSUMER</b></summary>

```
export KUBECONFIG=~/.kube/dev11
export TEMPLATE_PATH=~/projects/go/src/github/sweatShop-creator
export TEMPLATE_NAME=job-template.yaml
export REDIS_STREAM=sweatshop:test
export REDIS_PASSWORD=<SETME>
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

## License

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

Author Information
------------------
Patrick Hermann, stuttgart-things 26/2023