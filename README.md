# op4 [resource](https://faun.pub/writing-your-first-kubernetes-operator-8f3df4453234)
**Note**: The document was posted at 2019, the file structure is very different from current operator-sdk version. I have to make some adjustment to make it work. But still this is a very valuable study resource for operators.

## Description
### What is this operator doing? 
* It creates a custom resource replicaset, or in other words, it mocks the logic of built-in replicaset, which determines the amount of running pods by the `replicas` field under spec of the custom resource.

### Steps:
1. operator-sdk init and create API
2. go mod edit to version 1.22, and make sure to get the latest k8s.io libs
```
go 1.22.0

toolchain go1.22.2

require (
	github.com/onsi/ginkgo/v2 v2.17.1
	github.com/onsi/gomega v1.32.0
	k8s.io/api v0.30.0
	k8s.io/apimachinery v0.30.0
	k8s.io/client-go v0.30.0
	sigs.k8s.io/controller-runtime v0.18.0
)
```
3. Update Dockerfile base image to 1.22, and update the manager yaml file that implements the image to look for local image only.
4. Update the CRD file and controller file following the resource.
5. make docker-build, make install
6. make run
7. k apply CR yaml file
8. get pods to see 3 running
9. k edit CR describe directly to adjust replicas number, save
10. check running pods again

## Getting Started

### Prerequisites
- go version v1.20.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/op4:tag
```

**NOTE:** This image ought to be published in the personal registry you specified. 
And it is required to have access to pull the image from the working environment. 
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/op4:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin 
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

