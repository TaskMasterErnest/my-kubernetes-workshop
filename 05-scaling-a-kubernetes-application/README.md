# Deploying and Scaling a Kubernetes-Native Application

## CheckList
- [] Create a Deployment YAML file
- [] Perform some tests to test the cloud-native nature of Kubernetes applications


## Why?
- For business cases, it is possible that your application will have increased usage so it is right to have multiple sets/replicas of your application to serve the increased number of requests.
- When this is done, there are a group of applications waiting to receive and serve requests. These translate to being a groups of Pods having the application being requested in them.
- Kubernetes has a way of abstracting the provisioning of these groups of Pods working together to serve the same purpose. And these can be created/defined in a YAML file.
- By abstracting the way these Pods are provisioned, Kubernetes give the following advantages:
  - Establish redundancy by creating replicas of the Pods.
  - Easy upgrading and rollback of Pods.
- For these, we use the Kubernetes object called Deployments. Deployments make sure that a group of Pods are always online and ready to serve requests. And it has a specification for how these Pods are provisioned, upgraded and rolled-back.


## Creating the Deployment object
- The specification for creating the Deployment YAML file is in the `deploy.yaml` file.
- In the specification, we add a line called Replicas. With this line, we can specify the number of Pods we want in the abstracted group serving requests.
- All these Pods have to be accessed through the same service endpoint so we recycle the `svc.yaml` file from before.


## Test the Cloud-Native nature of the Deployment.
- We will test how the application serves traffic through all the Pods in the Deployment.
  - Simulate a number of requests inside a Bash script; `for i in $(seq 1 50); do curl <External-IP-address>:<Service Port>; done`.
  - Check how the traffic is routed and how the Pods serve these requests. We will count the number of lines each Pod logs in its logs with `kubectl logs <pod-name> | wc -l`. 
- Scale the application up using the command `kubectl scale deploy http-pod --replicas=5`.
- Scale the application down using the command to 2.
- Delete a pod in the Deployment and check how the deploy reacts.
