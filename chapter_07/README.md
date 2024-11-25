# Kubernetes Controllers
- In this session, we will talk about the concept of Kubernetes controllers and how to use them to create replicated Kubernetes objects.
- We will talk about the different set of controllers; ReplicaSets, Deployments, DaemonSets, StatefulSets and Jobs - and how to choose a suitable controller for the specific use-case.


## Introduction
- When we deploy our applications into production, it is necessary to waaaaay more than one Pod serving the requests to your application. The chief reason being that - having more than one Pod ensures that the deployed application continues to serve traffic when one of them goes down.
- In addition to handling the failure of Pods, replication ensures load-balancing between Pods replicas so that one Pod will not be the one to handle all the requests (a potential failure point).
- a **controller** can be defined as an object that ensures that your application runs in the desired state for the entirety of its lifetime.
- We will now start exploring the common Controllers that are widely used in Kubernetes:


### ReplicaSets
- We start off with **ReplicaSets**.
- A replicaSet is a Kubernetes controller that keeps a certain number of Pods running at any given time.
- It acts as a supervisor for multiple Pods across the different nodes in a Kubernetes cluster.
- It will terminate or add new Pods to match the declared configuration specified in the ReplicaSet template.
- It is a good idea to use them even if your application will only be running one instance of a Pod.
  - this way, if the one Pod is mistakenly deleted, the ReplicaSet will ensure that a new Pod is created to replace it. This enhances the availability of your application Pod.
- A ReplicaSet can be used to reliably run a single instance of a Pod indefinitely or to run multiple instances of the same Pod.

- How then do we configure a ReplicaSet? We will use this example to show the fields that appear in a ReplicaSet configuration.
```YAML
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: example-replicaset-pod
  labels:
    app: nginx
spec:
  replicas: 2
  selector:
    matchLabels:
      environment: production
  template:
    metadata:
      labels:
        environment: production
    spec:
      containers:
      - name: nginx-container
        image: nginx:latest
        imagePullPolicy: IfNotPresent
```
- As seen, a Replicaset contains the default apiVersion, kind and metadata fields. The juicier information is in the spec field. And they are:
  - **replicas**: this defines the number of Pods the ReplicaSet should be managing and running concurrently. The ReplicaSet controller will delete and add Pods to match the number presented here. The default number of Pods a ReplicaSet runs is 1.
  - **template**: in this field, we write out the template for the Pod to be deployed. This template is essentially a Pod template YAML manifest that we outline here. All the fields of Pod templating can be referenced here.
    - with this template, all Pods to be created as replicas will use this template.
  - **pod selector**: this is an important field. Here, we specify the label selectors that will be used by the ReplicaSet to identify the Pod to manage.
- An interesting thing about the ReplicaSet controller is that; a ReplicaSet will take charge of, and start managing a Pod that was created manually - if the template used in the ReplicaSet is the same as the one used in the standalone Pod template. It matches the Pod on its selector and uses that to manage the standalone Pod.


### Deployments
- A Deployment is a wrapper around a ReplicaSet and makes it easier to use.
- A general rule of thumb is to make Deployments manage replicated services. Deployments manage the ReplicaSets and the ReplicaSets in turn manage the Pods created/stated in the ReplicaSet.
- A major reason to use Deployments is that, it maintains a history of revisions. Everytime a change is made to the ReplicaSet or the underlying Pods, a new Revision of the ReplicaSet is recorded by the Deployment.
- Using a Deployment makes it easy to roll back to a previous state or version. Every rollback will also create a new revision for the Deployment.

- How do we configure Deployments? The configuration of a Deployment is very similar to that of a ReplicaSet. Here is a configuration:
```YAML
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 1
      maxSurge: 1
  selector:
    matchLabels:
      app: nginx
      environment: production
  template:
    metadata:
      labels:
        app: nginx
        environment: production
    spec:
      containers:
      - name: nginx-container
        image: nginx:latest
        imagePullPolicy: IfNotPresent
```
- There are some new fields in the Deployment configuration:
  - **strategy**: in this field, you specify the strategy the Deployment to use when it replaces old pods with new ones. The types are either *RollingUpdate* and *Recreate*. The default is *RollingUpdate*.
    - the **RollingUpdate** type is one to use to update a Deployment without having a downtime. 
      - In this, the controller updates the Pods one by one; hence, at any given time, there will be some Pods running and these Pods may be running two different versions of the application running.
      - this update strategy should be used for applications where the data stored by the new version can be read and handled by the old version - static data/information is the most preferred in this case.
      - when the *RollingUpdate* is picked, additional configuration is needed as shown:
      ```YAML
      strategy:
        type: RollingUpdate
        rollingUpdate:
          maxUnavailable: 1
          maxSurge: 1
      ```
      These are what they mean and/or represent:
        - **maxUnavailable** is the max number of Pods that can be unavailable during the update. This can be represented as a integer or a string with a percentage. The default value in percentages is **25%**.
        - **maxSurge** is the max number of Pods that can be scheduled/created above the desired number of Pods. This can also be specified as an integer or percentage string. The default value is also **25%**.
        - These two parameters can be tuned for ***availability and speed of scaling up or down***. Eg. **maxUnavailable: 0 && maxSurge: "30%"** - ensures a rapid scale-up while maintaining the desired capacity at all times. **maxUnavailable: "15%" && maxSurge: 0** - ensures that the deployment can be performed without using extra capacity at the cost of having, at worst, 15% fewer Pods running.
      - In the Deployment configuration above, we are specifying that no more than one(1) Pod is ever unavailable during updates and no more than four(4) Pods are ever scheduled.
    - the **Recreate** type strategy ensures that all the Pods are killed before creating new Pods with an updated configuration. In this type, there will be downtime during the update. There will then be fresh Pods with all of them running the same version. This is useful in when working with application Pods that need to have a shared state.
      - the parameters for this strategy type is shown below:
      ```YAML
      strategy:
        type: Recreate
      ```
      - A good use case for the Recreate strategy is if we need to run some data migration or data processing before the new code can be used. In this case, we cannot afford to be new code be running alongside the old code without running the migration or data processing for all Pods first.

- Realistically, there is going to be a chance that you will have to undo a Deployment or rollback a Deployment for a reason.
  - The root of the command used for this is `kubectl rollout`. With this, you can check the revision history and also rollback Deployments.
  ```bash
  kubectl rollout history deployment/<DEPLOYMENT_NAME>
  ```
  - In order to check the history to see which previous Deployments to roll back to, we need to make sure the changes are recorded and then these changes well documented so we know exactly what they are. The command flag to use for this is `--record`, but this flag has been set for deprecation and honestly, it does not do a good job of documenting the revision.
  ```bash
  kubectl apply -f <DEPLOYMENT_YAML_MANIFEST> --record
  ```
  - The agreed upon way to document revisions made to a Deployment manifest is to annotate the deployment with the cause of change.
  ```bash
  kubectl annotate deployment/<DEPLOYMENT_NAME> kubernetes.io/change-cause="<REASON FOR CHANGE>"
  ```
  - **NOTE**: The only way a Deployment can have a revision is when changes are made to the same YAML manifest used to deploy it.
  - Assuming we choose to change the image used in the deployment YAML we have above:
  ```YAML
  template:
    metadata:
      labels:
        app: nginx:alpine
        environment: production
  ```
  - We have to apply these changes to the YAML manifest. For proper documentation sake, we use the command: `kubectl annotate deployment/nginx-deployment kubernetes.io/change-cause="changed image to nginx:alpine"`.
  - We can then check the revisions for the Deployment using the command: `kubectl rollout history deployment/nginx-deployment`.
  - A situation where it is useful to use the `--record` flag will be when a change is made to a Deployment (this is running) using the **kubectl set** command. An example will be to change the image of the Deployment again from nginx to nginx:alpine, like this: `kubectl set image deployment <DEPLOYMENT_NAME> nginx:alpine=nginx --record`.
  - To perform an actual rollback, there are two ways to do this:
    - rollback to the immediate previous revision: `kubectl rollout undo deployment <DEPLOYMENT_NAME>`.
    - rollback to a specific revision in the revision history: `kubectl rollout undo deployment <DEPLOYMENT_NAME> --to-revision=<REVISION_NUMBER>`.
  - You can get the information about a particular revision with the command: `kubectl rollout history deployment/DEPLOYMENT_NAME --revision=<REVISION_NUMBER>`.


### StatefulSets
- StatefulSets are used to manage stateful replicas. Stateful repilcas are Pods of the same nature that manage/work with persisted data.
- StatefulSets maintain a unique identity for each of their Pods. 
- The Pods are of identical specification but they cannot be interchanged. This is because, StatefulSets attach a sticky identity to each Pod. This sticky identity is then used by the application to manage its data and/or state.
- For a StatefulSet with *n* replicas, the Pods are assigned a unique identity integer ordinals from 0 to n-1.
- When a StatefulSet is brought up, the Pods can created in the order of the unique interger ordinal, starting from 0 to n-1.
- A great thing about StatefulSets is that, when a Pod crashes and it is being brought up again, the unique sticky identity of the previously crashed Pod is assigned to the new Pod.

- A StatefuSet has the following configuration:
```YAML
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: sample-statefulset
spec:
  replicas: 3
  selector:
    matchLabels:
      environment: production
  template:
    metadata:
      labels:
        environment: production
    spec:
      containers:
      - name: nginx-stateful-container
        image: nginx:alpine
```

- What are some of the use cases for StatefulSets?
  - StatefulSets are useful when we want to work with persisted data / have persistent storage. With StatefulSets, the data can be partitioned in stored in various Pods. In case a Pod is being brought up to replace a crashed Pod, the new Pod will inherit the unique ID and the data partition of the crashed Pod.
  - A statefulSet is also useful if you want to create or update Pods in the order of the unique IDs assigned to them.


### DaemonSets
- DaemonSet are Kubernetes controller objects used to manage the creation of a particular Pod on all nodes or selected nodes in a cluster.
- When configured, they create Pods on nodes in the cluster. If new nodes are added, DaemonSets will make sure to run the particular Pods on these. If nodes are removed from the cluster, the daemonSet Pods are destroyed on the cluster.

- What are some use cases for DaemonSets?
  - DaemonSets can be configured to deploy log collection Pods on all nodes. These Pods collect logs and forward them to a log aggregator/processor.
  - DaemonSets can be used to create caching Pods on all nodes. These Pods can be used by the application Pods for temporary storage of cached data locally (on the node).
  - DaemonSets can be used to create Pods that collect system- to application-level metrics on all nodes.

- The configuration of a DaemonSet is as follows:
```YAML
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: sample-daemonset
  labels: 
    app: daemonset-example
spec:
  selector:
    matchLabels:
      app: daemonset-example
  template:
    metadata:
      labels:
        app: daemonset-example
    spec:
      containers:
      - name: log-container
        image: fluentbit:latest
```

**NOTE**: All the controllers that have been talked about so far are for workloads that are to be run continually. There are controllers that manage workloads that are brought to a graceful conclusion after the tasks in them have been performed. For the latter kind of workload, Kubernetes uses a Job controller for it.

### Job
- A Job is a coontroller that is used to manage Pods that are supposed to run a determined task and thn terminate gracefully afterwards.
- A Job creates the designated number of Pods and ensures that they complete their task successfully.
- When a Job is created, it creates and tracks the Pods that are specified in it configuration. 
- A Job is only complete when a specified number of Job Pods complete successfully.
- If underlying node issues cause a Job to fail, the Job controller will create a new Pod and start the Job again.
- The Job Pods created are not removed/deleted by default when they complete successfully. They stay in the cluster with a *Completed* status.
- btw Jobs have a nice cozy niche in machine learning and data analysis workloads.

- What are some of the use cases for a Job then?
  - The simplest use case is to create a Job that runs only one Pod to completion. If failures arise, additional Pods will be created to run the job to completion. This can be applied to a one-off or recurring data analysis task or for the training of an ML model.
  - Jobs can be used for parallel processing. We specify more than one successful Pod completion to ensure the Job will complete only when a certain number of Pods terminate successfully.

- A sample configuration for a Job is:
```YAML
apiVersion: batch/v1
kind: Job
metadata:
  name: one-time-job
spec:
  template:
    spec:
      containers:
      - name: get-date
        image: busybox
        args:
        - /bin/sh
        - c
        - date
      restartPolicy: OnFailure
```

- Let us look at a use case for Jobs in Machine Learning.
  - Jobs are perfect for batch processing tasks. These are tasks that run for a certain amount of time before exiting.
  - This makes Jobs ideal for production machine learning tasks like feature engineering, cross-validation, model training and batch inference.
  - You can create a Job that trains an ML model and persist the model and training metadata to an external storage. You can create Jobs for batch inference - ie. tasks that fetch a pre-trained model from storage, loads the model and data into memory, performs inference and stores the predictions.