# Labels and Annotations
- Metadata ia an extremely useful thing that is used to manage resources in a cluster/enterprise.
- In this session, we will learn to add metadata to Kubernetes resources.
- We will look at labels and annotations and the use-case for each of these.
- We will then see how to organize labels for resources by using label selectors and how to use annotations to add unstructured metadata info to objects.


## Introduction
- We have seen two fields under the *metadata* field in a Pod configuration ie *name* and *namespace*. These are the first metadata on Pods and other Kubernetes resources.
- Now, we can add additional fields to provide extra information as labels and annotations.
- We will see how to metadata to Pod configuration in order to identify them through queries based on metadata and then additional unstructured metadata.
- We will examine both labels and annotations, the difference between them and when to use one or the other.


## Labels
- Labels are the metadata that contain identifiable information pertaining to Kubernetes objects.
- They are key-value pairs that can be attached to objects such as Pods.
- Each key must be unique for an object. These labels contain information that is meaningful to users.
- Labels can be attached to Pods at (1) the time of the Pod creation and also (2) can be added to and/or modified during the runtime.
- Labels appear under the metadata field in the YAML manifest:
```YAML
metadata:
  labels:
    key1: value1
    key2: value2
```

- There are some constraints for labels. These constraints exist because, in this way, querying them will be faster as they have been optimized for faster evaluation by Kubernetes internal data structures and algorithms. The mappings are maintained internally by Kubernetes for faster querying.
- There are two parts of a Label, the key and the value. Here are the constraints for each of these:
  - **Label Keys**: This is an example of a label key, `label_prefix.com/worker-node-1`.
    - There are two parts of the label key. The label prefix and the label name.
    - **Label Prefix**: This prefix is optional BUT it must be a DNS subdomain. It cannot be longer than 254 characters and cannot contain spaces either.
      - the label prefix is always followed by a forward slash (/).
      - if there is not label prefix, the prefix is assumed to be private to the user.
      - some prefixes such as *kubernetes.io/* and *k8s.io/* are reserved for use solely by the Kubernetes core system.
    - **Label Name**: The label name is required. It can only be up tp 63 characters long.
      - the label name can only start and end with alphanumeric characters. It can contain dashes(-), underscores(_), dots(.) but cannot have spaces or line breaks.
  - **Label Values**: Label values take on the same characteristics as the label name part of the label key. 
    - they can only be up to 63 characters long.
    - they can contain dashes, underscores, dots and alphanumeric values but not spaces and line breaks.

- Why then do we need labels? We need labels in order to organize Kubernetes resources in a cluster.
- With these labels in place, we can then filter the resources. Here are two(2) reasons why we need labels:
  - 1. For organizing Pods in the Kubernetes cluster. Let us suppose one cluster is used to serve many teams with different functionalities, labels can be attached to resources only need by specific people on specific teams and for specific projects to ensure that partition. A sample label in this case might be:
  ```YAML
  metadata:
    labels:
      environment: staging
      team: devops
      project: freemind
  ```
  - 2. Labels are useful when we want to selectively run Pods on specific nodes. By slapping on labels to Pods, we can use these labels to ensure the Pods are deployed unto nodes that are specific/suitable for them based on the admin reasons. We can use the **nodeSelector** field in the Pod to assign a Pod to a specific node if it has a specific label. An example is:
  ```YAML
  apiVersion: v1
  kind: Pod
  metadata:
    name: pod-with-node-selector
  spec:
    containers:
    - container: container-with-node-selector-enabled
      image: nginx:latest
    nodeSelector:
      region: us-east
      disktype: ssd
  ```
    - Note that these exact node labels to be used by the nodeSelector are to be made avaialable by the node, either set by the cloud infra or by the admin for a self-managed cluster.

- Labels and labelling in Kubernetes is a curious thing. You can create a Pod withtout a label and add a label when the Pod is running. You can modify the labels on a running Pod and you can also remove the labels from a running Pod.
- Let us use this Pod configuration and we will look at the various ways we can add labels to a Pod.
  - 1. ***Creating a Pod with a label***:
  ```YAML
  apiVersion: v1
  kind: Pod
  metadata:
    name: pod-for-label-experiment
    label:
      env: production
      service: backend
  ```
    - if we run the **`kubectl create -f <pod-config.yaml>`** on this, we will have a Pod with these labels.
  - 2. ***Adding labels to a running Pod***.
    - Using the same Pod configuration above, and using the kubectl CLI tool, we can add more labels to the pod.
    - let us add `team=devops` to the labels. We use the command **`kubectl label pod/<POD_NAME> team=devops`**.
    - You can add multiple labels this way, eg **`kubectl label pod/<POD_NAME> team=devops project=foodie`**. This will have two labels added to the Pod.
    - this will trigger a change in the labelling. To get the new Pod configuration, you can describe the Pod with this command `kubectl describe pods <POD_NAME>`.
  - 3. ***Modifying existing labels on a running Pod***.
    - Using the same Pod configuration, and the kubectl CLI tool, we can overwrite the existing labels of an existing Pod.
    - You need to identify the existing Pod whose label(s) you want to modify. In this case, let's modify the project to be trashie (yeah...real innovative ).
    - You use the command **`kubectl label --overwrite pod/<POD_NAME> project=trashie`**. With this, we have changed the label value from foodie to trashie.
    - The **--overwrite** flag is what is essential for this to work.
  - 4. ***Deleting the labels off a running Pod***.
    - Still using the same Pod configuration and the kubectl CLI tool.
    - You need to identify the Pod from which you will be deleting a label. You also need the label key from the label. Let us assume we want to delete the label `project` from the Pod.
    - We run the command **`kubectl label pod/<POD_NAME> project-`**. Note the hyphen (-) attached to the label key. With this command, we have deleted the label `project` from the Pod.

- Now on to the reason why labels are attached to resources. In order to group various Kubernetes objects you can use their label and group them with a ***label selector***. With label selectors, you can identofy a set of objects matching certain criteria.
- The command we can use to pass the label selector to the kubectl command is: **`kubectl get pods -l {label_selector}`**. The flag is either **-l** or **--label**.
- There are various **{label_selector}** arguments that can be passed to the kubectl command, they are: equality-based selector arguments and set-based selector arguments.
  1. ***Equality-based Selectors***: These selectors allow us to match all objects that have soecific label values for the given label keys. In equality-based selectors, we have inequality matching as well.
  - there are three(3) kinds of operators used in matching: **=**, **==** and **!=**.
  - while using these operators, we can specify more than one condition using any of the operators.
    - assuming we want to get all resources with label keys `environment`, we use `kubectl get pods -l environment=<ENVIRONMENT_NAME>`.
    - let's say we want to match all resources/objects that either do not have a `team` label key OR those for which a `team` label key exists, and the corresponding label value is not `devops`. The command will be `kubectl get pods -l team!=devops`.
    - finally, we can use both selectors together, separated by commas(,). Like this `kubectl get pods -l environment=test,team!=devops`. This will match all the objects that match both criteria specified by the two selectors. The comma acts a logical AND (&&) separator between the two selectors specified.
  2. ***Set-Based Selectors***: These selectors allow us to match all objects that have a given label key with a value in a given set of values.
  - There are three(3) kinds of operators: **in**, **notin** and **exists**.
    - assuming these selector that matches all objects that have an `environment` label key and the value is `production` and `staging`. The command is `kubectl get pods -l 'environment in (production, staging)'`.
    - matching all the objects that have an `team` label key and the value is anything other than `devops`. It also matches those objects that do not have the `team` label key. The command is `kubectl get pods -l team notin (devops)`.
    - also with using the exists operator, an example is `kubectl get pods -l !critical`. This matches on all objects that do not have the `critical` label key.
  - We can combine both equality- and set-based selectors to filter/match objects in the Kubernetes cluster. An example is `kubectl get pods -l 'team=devops,environment in (production,staging)'`.


## Annotations
- Annotations unlike labels have fewer constraints in terms of what kinds of data can be stored in them. However, we cannot filter or select objects by using annotations.
- Annotations are alo key-value pairs that can be used to store the unstructured information pertaining to the Kubernetes objects. Here is how they are defined:
```YAML
metadata:
  annotations:
    key1: value1
    key2: value2
```

- There are constraints for annotations. The rules are more relaxed than the rules for label keys and values the reason is that there is not an optimized way for selecting or filtering objects using annotations.
- Just like labels, annotations have keys and values.
  - 1. ***Annotation Keys***: Similar to label keys, annotation keys have two parts: a prefix and a name. The constraints for both annotation prefixes and names are the same for labels. An annotation key is like this: `annotation_prefix.com/worker-node-identifier`.
  - 2. ***Annotation Values***: There is no restriction in terms of what kinds of data annotation values may contain.

- Annotations have their use-cases. We will delve into these various use cases.
- Annotations are generally used to add metadata that will not be used to filter or select objects. This metadata will be used to get more subjective information regarding the Kubernetes objects.
  - Annotations can be used to add timestamps, commit hashes, issue-tracker limits, or names/information about users who are responsible for specific objects in an organization. The annotations metadata can be written as:
  ```YAML
  metadata:
    annotations:
      timestamp: 1234567890
      commit-SHA: d6s9shb82365yg4ygd782889us28377gf6
      JIRA-issue: "https://your-jira-link.com/issue/ABC-1234"
      owner: "https://internal-link.to.website/username"
  ```
  - Annotations can also be used to add information about client libraries or tools. This information can later be used to for debugging issues in the application. It can be written as:
  ```YAML
  metadata:
    annotations:
      go-version: 1.22.8
      go-documentation: https://go.dev/doc/
  ```
  - Annotations can be used to store the previous Pod configuration deployed. This is useful for figuring out the configuration was deployed before the current revision and what has changed. It can be written as:
  ```YAML
  metadata:
    annotations:
      previous-configuration: "{ some JSON containing the previously deployed configuration of the object }"
  ```


- In working with annotations, they are similar to the way we work with labels. Adding, modifying and deleting annotations can be done to a running Pod. The commands to do so are also similar. We go through them like this:
  - adding annotations to a running Pod can be done using the following command: `kubectl annotate pod <POD_NAME> <ANNOTATION_KEY>=<ANNOTATION_LABEL>`. You can multiple annotations like done with Pods.
  - modifying existing annotations on a running Pod can be done using the following command: `kubectl annotate --overwrite pod <POD_NAME> <ANNOTATION_KEY>=<ANNOTATION_LABEL>`.
  - deleting an annotation from a running Pod can be achieved with the following command: `kubectl annotate pod <POD_NAME> <ANNOTATION_KEY>-`. Note the dash(-) at the end of the annotation key.