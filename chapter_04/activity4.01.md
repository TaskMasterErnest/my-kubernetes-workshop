# Create a Deployment Using a ServiceAccount Identity
- We will be using various operations on our cluster and using different methods to access the API server.

1. create a new namespace called `activity-example`:
  - `kubectl create ns activity-example`

2. Create a new ServiceAccount called `activity-sa`:
  - `kubectl create sa activity-sa -n activity-example`

3. Create a new RoleBinding called `activity-sa-clusteradmin` to attach the `activity-sa` ServiceAccount to the `cluster-admin` ClusterRole (which exists by default). This step is to ensure that our ServiceAccount has the necessary permissions to interact with the API server as a cluster admin.
  - `kubectl create rolebinding activity-sa-clusteradmin --clusterrole=cluster-admin --serviceaccount=activity-example:activity-sa --namespace=activity-example`

4. Create a new NGINX Deployment with the identity of the `activity-sa` ServiceAccount
  - use the kubectl proxy comamnd to create a proxy so we can use cURL to create resources.
  - use this command to create the deployment in the `activity-example` namespace.
  ```bash
  curl -X POST http://127.0.0.1:8080/apis/apps/v1/namespaces/activity-example/deployments -H 'Content-Type: application/yaml' --data "
  apiVersion: apps/v1
  kind: Deployment
  metadata:
    creationTimestamp: null
    labels:
      app: activity-nginx
    name: activity-nginx
    namespace: activity-example
  spec:
    replicas: 1
    selector:
      matchLabels:
        app: activity-nginx
    strategy: {}
    template:
      metadata:
        creationTimestamp: null
        labels:
          app: activity-nginx
      spec:
        serviceAccountName: activity-sa
        containers:
        - image: nginx:latest
          name: nginx
          resources: {}
  status: {}
  "
  ```
  - check if the serviceAccount has been used correctly using the command `kubectl get pod <pod_name> -n activity-example -o yaml | grep serviceAccount`