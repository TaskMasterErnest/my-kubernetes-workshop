# Using kubectl - The Command Center of Kubernetes
In this session, we will look at the end-to-end process of using kubectl to communicate with a Kubernetes cluster.


## What Is kubectl
**kubectl** is the command-line utility for communicating with Kubernetes clusters and performing various operations in them.
kubectl uses two ways to manage the Kubernetes cluster:
- imperative management: this is by using commands rather than YAML manifests to achieve the desired state of a deployment & 
- declarative management: which focuses on creating and updating YAML manifest files to achieve a desired state.
kubectl uses these to manage the Kubernetes API objects (also called API primitives).
kubectl allows the wielder to send commands to Kubernetes clusters; deploy apps, inspect, manage objects, troubleshoot and view logs.
kubectl although used to manage Kubernetes, has to be downloaded separately as a command-line tool.


## How kubectl Communicates With A Kubernetes Cluster
The simple and uncomplicated way kubectl communicates with the cluster is this: 
- the kubectl command is translated into an API call, which is then sent to the API Server
- the API Server authenticates and validates the requests
- After authentication and validation, the API server retrieves and updates data in the \*etcd\* and responds with the requested information
The API Server is the one that manages communications between the user and Kubernetes and it also acts as an API Gateway to the cluster.
- to do this, it implements a RESTful API over HTTP and HTTPS protocols to perform CRUD operations to populate and update Kubernetes objects based on the instructions sent by the user via kubectl.


## The kubeconfig File Configuration
The kubeconfig file is used to store information about the users, clusters, namespaces and authentication mechanisms about Kubernetes clusters.
In an enterprise environment, an administrator can be managing more that one Kubernetes cluster. This user will have to interact with all these cluster and swrich between them if needed to perform different operations. The user relies on the kubeconfig files to do so.
The kubeconfig file contains information about each cluster and the info needed to communicate with its API Server.
This file, by default is stored in the `$HOME/.kube/` directory under the `config` filename. In other cases, a `KUBECONFIG` env-var or the `--kubeconfig` flag may be used to specify the location in which it is placed.
To view the contents of the kubeconfig file, do any of these:
- If a cluster is vavailable and running, run `kubectl config view` to get the contents
- Or run `cat $HOME/.kube/config` to get the contents.

By now, you must have guessed what a context in a Kubernetes cluster is.
A **context** is the information needed to access a cluster. It contains the name of the cluster, the user and the namespace.
The **current-context** field shows which particular context is being worked with currently.
To switch contexts, use the `kubectl config use-context <name-of-cluster-to-switch-to>`.


## Common kubectl Commands
The kubectl command-line tool has a lot of useful commands to work in Kubernetes with.
Here are come common ones used to create, manage and delete Kubernetes objects:
- **get < object >**: this is used to get the list of the objects specified. Use the **all** to get the list of all kinds of objects. Eg. `kubectl get pods`
- **describe < object-type > < object-name >**: this command is used to check very relevant information concerning a specific object. Eg. `kubectl describe pod http-pod`
- **logs < object-name >**: this is used to check all the relevant logs of a specific object. Eg. `kubectl logs http-pod`
- **edit < object-type > < object-name >**: this command is used to edit a specific object. Eg. `kubectl edit pod http-pod`
- **delete < object-type > < object-name >**: this command is used to delete a specific object. Eg. `kubectl delete pod http-pod`
- **create < filename.yaml >**: this is used to create Kubernetes objects that have been specified in the YAML manifest. Eg. `kubectl create -f http-pod.yaml`
- **apply < filename.yaml >**: this is used to either create or update the Kubernetes objects that have been defined in the YAML manifest. Eg. `kubectl apply -f http-pod.yaml`


### Some Useful Flags For The get Command 
Ths getc command is a pretty standard command we use to get the list of object present in the Kubernetes cluster.
- the **--all-namespaces** flag is used to get a particular type of resource/object in all the available namespaces present in the cluster. Eg. `kubectl get deployments --all-namespaces` lists all the Deployment objects from all the namespaces present in the cluster.
- the **-n** flag is used to specify a particular namespace to list an object(s) from. Eg. `kubectl get deployments -n kube-system` will list all Deployments in the kube-system namespace.
- the **--show-labels** flag is used to add the labels an object/resource has attached to it in the Kubernetes cluster. Eg. `kubectl get pods --show-labels` will list all the Pods in the default namespace and the labels attached to them.
- the **-o wide** flag is used to display more information about objects. The '-o' stands for output so it direcly translates to showing more output from the brief information shown when object info is queried. Eg. `kubectl get nodes -o wide` will show all nodes with additional information.


## Creating a Deployment in Kubernetes
Deployments are a convenient way to manage and update Pods. 
When we create a Deployment in Kubernetes, we have created a way to effectively and efficiently provide declarative updates to the application we are introducing into the Kubernetes cluster.

We can create the Deployment in the Kubernetes cluster using either imperative commands or using the declatrative format using YAML manifest files.

The imperative commands are of the form; `kubectl run nginx --image=nginx:alpine --replicas=3`, the declarative commands are of the form; `kubectl apply/create -f sample-deployment.yaml` where the 'sample-deployment.yaml' is a YAML manifest file.

To get a list of all the Deployment objects in a particular namespace, we use the command `kubectl get deployments -n <namespace>`.

A more encouraged way to create Deployments is by using a declarative YAML manifest file. This way, the changes can be tracked with a version control tool like Git.


## Updating a Deployment
Updating a Deployment is in two ways: moving from the current version to a recently-published (newer) version **and** rolling back from a current version to the previous version.

Let us assume for the current Deployment, I want to raise the image version of NGINX to 1.9.1. There are two ways to do this - using the **kubectl set image** command and/or using the **kubectl apply** command.

We can update the Deployment of the NGINX image using the following imperative command; `kubectl set image deployment/nginx-deployment nginx=nginx:1.9.1 --record.` The `deployment/nginx-deploy` is the name of the Deployment when you look it up in the Kubernetes cluster with the `kubectl get deploy` command (the **< object-type >/< object-name >** format). 

The **--record** flag is to save the updates that have been made by the kubectl tool to the current resource/object.

In using the declarative command format, we get the YAML manifest file that was used to create the NGINX Deployment, manually change the version of NGINX to 1.9.1 and apply the changes with the command `kubectl apply -f nginx-deployment.yaml --record`.

### Performing a Rollback on Deployments
It comes without saying that, in order to perform a rollback, we must have already performed an update to the resource/object.

Using the previous update as our lauching pad, we can directly perform a rollback using the command, `kubectl  rollout undo deployments nginx-deployment`.

This is no doubt fascinating but there is more to this. We can view the entire rollout history (mostly updates) done to our specific Deployment by running the command `kubectl rollout history deployment nginx-deployment`. It shows us a list of all the revisions made to this paticular Deployment. If you used the record flag on any update, the change that happens is stated in the rollout history for you to peruse.

We can check out the details of a specific Deployment by running the `kubectl rollout history deployment nginx-deployment --revison=<number>` command to select the specific revision to review.

We can also rollback a Deployment to a specific revision by using the **--to-revison** flag, like this, `kubectl rollout undo deployments nginx-deployment --to-revision=<number>`.


## EMERGENCY ALERT / ACTIVITY
You are a SysOps Engineer managing a Kubernetes cluster that has a web application deployed on it and made available to the public.
You have been monitoring the application and have detected that it experiences throttling issues at peak times.
Based on the monitoring, the solution you have concluded upon is to increase the memory and CPU allocated to the application.
And in this case, you have to do a Live Edit of the application as it is running in the Kubernetes cluster.

The steps to do this are simple (if you have beedn following every project in this workhop so far.), they are:
1. deploy the application `sample-application.yaml` on the Kubernetes cluster.
2. get the IP address of the service that exposes the application.
3. edit the live deployment of the `melonvote-front` Deployment. Change the following:
  - resources.limits.cpu
  - resources.limits.memory
  - resources.requests.cpu
  - resources.requests.memory
  
  basically double these values and youy should be able to see the UI of the application.

A lot of things are deprecated, so you will need to fix the application, rebuild it and deploy it to Dockerhub and pull it into you new Deployment. The link to the resources are [here](https://github.com/Azure-Samples/azure-voting-app-redis/blob/master/azure-vote/Dockerfile-for-app-service) and [here](https://learn.microsoft.com/en-us/azure/aks/tutorial-kubernetes-prepare-app).