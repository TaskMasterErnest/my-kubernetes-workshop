# How To Communicate With The Kubernetes API (API Server)
- In this session, we will look into buulding foundational knowledge and understanding of the API server and the various ways of interacting with it.
- We will see how kubectl and other HTTP clients communicate with the API server. We will see demos of how to trace the communications and see the details of HTTP requests.
- Also, we will see how we can look up the API details so that we can write our own API requests from scratch.


## Introduction
- As a recap, the API server acts as the central hub that communicates with all different components in Kubernetes.
- We are going to look in to the components that make up the API server. We will see how to effectively communicate with the API server and how these API requests are processed.
- We will look at API concepts that will help us understand the HTTP requests that are made to the API server. Concepts such as resources, API groups and API versions.
- We will then interact with the Kubernetes API using multiple REST clients to achieve many of the results we see when we use the kubectl command-line tool.


## The Kubernetes API Server
- All communication and operations in the Kubernetes cluster; between the control plane components and external clients (like kubectl) - are translated into **RESTful API calls** that are handled by the API server.
- So we can say the API server is a RESTful web server that processes RESTful API calls over HTTP. The calls are made, to store and update API object in the etcd datastore.
- The API server can also be seen as a frontend component that acts as a gateway to and from the outside world. A frontend that is accessible by all clients, including the kubectl CLI tool.
- The internal cluster components interact with each other through ONLY through the API server.
  - It is the ONLY component that interacts directly with the **etcd** datastore.
- It is important that the API server be configured correctly since it is the only way for clients to access the cluster.
- The way to see the API server in the components is to use the kubectl command `kubectl get pods -n kube-system`. The API server will have `kube-apiserver` in its name.
- The API server is stateless ie. its behaviour is consistent regardless of the state of the cluster.
- The API server is designed to scale horizontally.


## Kubernetes HTTP Request Flow
- We know now for a fact that, when we run any kubectl command, the command is translated into an HTTP API request in JSON format and is sent to the API server. Then the API server returns a response to the client along with any requested information.
- The following diagram shows the API request life cycle and what happens inside the API server when it receives a request. 
- ![api-request-http-flow](./files/api-server-http-request-flow.png).
- The request goes through the **authentication, authorization and admission** control stages.

### Authentication
- Every API call needs to authenticate with the API server, whether it is coming from outside the cluster (via kubectl) or from a process inside the cluster (such as those made by the kubelet).
- When the request reaches the API server, the API server needs to authenticate the client sending the request.
  - the request should (and usually) contains the username, user ID,and group - all necessary information needed for authentication.
- The authentication methods will be determined by either the header or the certificate of the request.
- To deal with the different methods for authentication, the API server has different auth plugins, such as ServiceAccount tokens (for authenticating ServiceAccounts) and one other method to authenticate users, usch as the X.509 client certificates.
- **NOTE**: the cluster admin is the one who usually defines the auth plugins during cluster creation. Learn more [here](https://kubernetes.io/docs/reference/access-authn-authz/authentication/).
- The API server will call these plugins one-by-one until one of them authenticates the request. If all of them fail, the auth process fails. If not, the authentication phase completes and the request proceeds to authorization.

### Authorization
- [docs](https://kubernetes.io/docs/reference/access-authn-authz/authorization/)
- The authorization stage is where it is determined whether the user is permitted to perform the requested action.
  - There are various levels of privileges that different users can have; eg. eg listing pods in a namespace, delete deployments etc. These kinds of decisions are made in the authorization phase.
- Provided we have two users, stated below:
![user-privileges](./files/user-privileges.png)
- The above user **ReadOnlyUser** tries to do certain things in the cluster and this is what happens:
![demo-user-privilege-results](./files/demo-user-privilege-results.png)
  - the **Forbidden** error is returned by the authorization plugin.
- Thankfully, kubectl provides a command that can be used to check whether an action is allowed for the **current** user. The command is `kubectl auth can-i`. 
  - a very intuitive command. eg. to check if you can delete deployments: `kubectl auth can-i delete deployments`. 
  - run the `kubectl auth can-i -h` for help.
- Authorization modules are checked in sequence. When multiple authorization modules are configured, if any authorizer approves or denies a request, that decision is immediately returned and no other authorizer will be contacted.

### Admission Control
- [docs](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/)
- After the request is authenticated and authorized, it goes to the admission control modules. 
  - These modules can modify or reject requests.
  - if the request is only doing a READ operation, this stage is bypassed. But if it is trying to create, modify or delete, it is sent to the admission controller plugins.
- A cluster is initiated with a set of predefined admission controllers. Custom admission controllers can be defined as well.
- Like authorization modules, if any admission controller module rejects the request, that request is dropped and not processed further.
- Some examples of how the admission controller objects works is:
  - If we configure a custom rule that every object should have a label, then any request to create an object without a label will be rejected by the admission controllers.
  - When a namespace is deleted, Kubernetes will try to delete all resources in there before deleting. That state is called the **Terminating** state. In this state, we cannot create any new objects in this namespace. The **NamespaceLifecycle** is what prevents that.
  - When there is a request to create a resource in a namespace that does not exist, the **NamespaceExists** admission controller rejects this request.
- Not all of the admission controllers are enabled by default. And the default modules activated chaneg based on the Kubernetes version.
- Cluster admins can decide which modules to enable or disable when initializing the API server.
- To control which of the admission modules to enable or disable other than the default ones, theses flags, `--enable-admission-plugins` and `--disable-admission-plugins` are used.
*see exercise 4.01*

### Validation
- After letting the request pass through all three stages, the API server then validates the object - that is, it checks whether the object specification, carried in the JSON format in the response body, meets the required format and standard.
- After successful validation, the API server stores the object in the etcd datastore and returns a response to the client.


## The Kubernetes API
- The Kubernetes API uses JSON over HTTP for its requests and responses. It follows the REST architectural style.
- The Kubernetes API can be used to read and write Kubernetes resource objects.
- Just like HTTP methods, the Kubernetes API allows clients to create, update, delete or read an object via standard HTTP methods.
  - an example can be seen here: ![http-methods-for-k8s](./files/k8s-api-http-methods.png)
- The API calls carry JSON data and all of them have a JSON schema identified by the **Kind** and **apiVersion** fields.
  - the Kind field is a string the identifies the type odf JSON schema that an object should have. The apiVersion is a string that identifies the version of the JSON schema the object should have.
- The best way to understand how the Kubernetes API works with requests is to trace a kubectl command.

### Tracing kubectl HTTP Requests
- [kubectl verbosity and debugging flags](https://kubernetes.io/docs/reference/kubectl/quick-reference/#kubectl-output-verbosity-and-debugging)
- We can start tracing the HTTP requests that the kubectl sends to the API server.
- Suppose we want to get the pods in the kube-system namespace. We run the command `kubectl get pods -n kube-system`.
  - behind the scenes, what this does is to invoke an HTTP GET request to the API server endpoint and requests information from the `/api/v1/namespaces/kube-system/pods`.
- We can enable a verbose output to our command. If we do this, we get more details in the response. The verbosity ranges from 1 to 10. We enable this by adding the `--v=n` flag to the command, where `n` is the number.
- We get a pretty hefty output: ![verbosity level 8](./files/full-verbosity-output.png)
- Now, let is delve into the output bit by bit to get a better understanding of this.
  - The first part of the output is as follows: ![first part](./files/delve-into-verbosity-1.png)
    - From this, we can see that the kubectl loaded the configuration from the kubeconfig file. That file has the API server endpoint, port, and credentials (certificate or auth token).
  - The next part: ![second part](./files/delve-into-verbosity-2.png)
    - This line contains the action to perform against the API server. The line contains the path to go visit and the limit on the number of resources to return. The limit here is 500; kubectl fetches a large chunk of resources in order to improve latency.
  - The next part of this output is this: ![third part](./files/delve-into-verbosity-3.png)
    - In this part of the output, the Request Headers describes the resource to be fetched or the client requesting the resource. There are two parts of this, and they are used for content negotiation - like this:
      - **Accept**: This is used by HTTP clients to inform the server about the types fo content they (client) will accept. In the output, we can see that the client requests to be handed the content type as an *application/json*. If this content type is not available on the server, it will return the default preconfigured representation type, which is the same *application/json* (as this is what the Kubernetes API server uses as its JSON schema). We can also see that the client is request the output as an APIGroupDiscoveryList as indicated by the **as=** line.
      - **User-Agent**: This header contains information about the client that is requesting the information. In this case, we can see that kubectl is providing information about itself.
  - Moving unto the next part: ![fourth-part](./files/delve-into-verbosity-4.png)
    - Here, we can see that the API server returns a 200 Status Code indicating that the request has been successfully processed. Also,we can see the time taken to process this request.
  - The next part of: ![fifth-part](./files/delve-into-verbosity-5.png)
    - This shows the Response Headers and the main response body sent by the API server.
    - The response body contains the resource data that was requested by the client but in raw JSON format before it can be translated into the response content type the client requested it in. To see the full body response, the verbosity level 10 will do.


### API Resource Type
- In requesting resources from the API Server, we use the HTTP URL which contains the API resource, API Groups and API Version.
- In this section, we will concentrate more on the resource type in the URL such as pods, namespaces, and services. In the JSON format, this is labelled **Kind**.
- The resources can be a collection of resources, or a single resource.
  - For the *collection of resources*, this represents a collection of instances for a resource type. In a URL, this will be how it is represented `GET api/v1/pods`.
  - For a single resource, this represents a single instance of a resource type. The URL to this is more specific; `GET api/v1/namespaces/{namespace}/pods/{pod}`.
- Therefore, in requesting a resource in Kubernetes, you can call for a collection of resources ar a single resource.


### Scope of API Resources
- The resources in Kubernetes can be either cluster-scoped or namespace-scoped. The scope affects the access of that resource and how that resource is managed.

#### Namespace-Scoped Resources
- Kubernetes makes use of Linux namespaces to organize most resources.
- Resources in the same namespace share the same access control policies and authorization checks.
- Let us see what forms the request paths for interacting with namespace-scoped resources:
  - get info about a specific pod in a namespace `GET api/v1/namespaces/{namespace}/pods/{pod}`
  - get info on a collection of Deployments in a namespace `GET apis/apps/v1/namespaces/{namespace}/deployments`.
  - get info on all services across all namespaces `GET ap1/v1/services`
- Now, notice that, when getting information on a resource across all namespaces, it will not have a namespace in its URL.
- To get the full list of all namespace-scoped resources, use this kubectl command: `kubectl api-resources --namespaced=true`.
![namespace-scoped resources](./files/namespace-scoped-resources.png)

#### Cluster-Scoped Resources
- Most Kubernetes resources are namespace-scoped. Any resource that cannot be found under a namespace/within a namespace is cluster-scoped. 
  - an example is a node. You can deploy a pod in a namespace regardless of the namespace you want the pod to be in. A node can also host different pods from different namespaces. Also, you might have guessed it, nodes host the namespace resource ie namespaces are found inside nodes.
- Let us see the request paths for interacting with cluster-scoped resources:
  - get info on a specific node in the cluster: `GET api/v1/nodes/{node}`.
  - get info on all nodes in the cluster: `GET api/v1/nodes`.
- In order to get all the resources that are cluster-scoped, use the kubectl command: `kubectl api-resources --namespaced=false`.
![cluster-scoped resource](./files/cluster-scoped-resources.png)

