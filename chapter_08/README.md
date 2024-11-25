# Service Discovery
- In this session, we will loo at how to route traffic between the various kinds of objects in the Kubernetes cluster.
- We will look at how to make our applications discoverable internally (in the cluster) and externally from the cluster.
- We will also look at the different types of Services and how to use these Services to enable interaction with and between Pods.


## Introduction
- Networking is important in making sure the application deployed in Pods is accessible by the end-users.
- Networking involves the use of IP addresses. IP addresses are assigned to Pods when they are created. These IP addresses are how we can discover our application in the cluster.
- Assume we have replicated Pods from a Deployment serving up the application, how do we ensure that all these Pods with different IP addresses are accessible by the end users?
- This is why Kubernetes Service objects exist. Services group logical Pods together - logical here means Pods that are meant to do the same thing - and makes sure they are discoverable and/or accessible by the end-users.
- Applications inside Pods can either be accessed by another application or by the real users. Eg. A frontend application needs to be accessible to users, and it also needs to access the backend service (if there is any).


### What is a Service?
- A Service is an object that defines the policies by which a set of logical Pods can be accessed.
- Services make sure there is communication between the user and the application or between components of the application.
- The abstraction of the network in the form of Services, allows applications to be decoupled in the cluster and still be able to effectively communicate with each other.
- **NOTE**: The native way Kubernetes links various different resources together is through the use of `labels` and `label selectors`.
- They way a Deployment manages scaling and availability of its Pods is by utilizing labels and label selectors to identify and manage its Pods.
  - Services provide a static layer between these Pod IP addresses and the selector-based mechanism of reaching these Pods. In this way, they know how to reach these Pods and can route traffic to these Pods.
- A Service spans all the nodes in the cluster. It provides a flat interface across nodes in the cluster where all the selected Pods can plug into.
- You can declare a Service by writing a Service YAML manifest or by using the ***kubectl eexpose*** command.


### How do you configure a Service?
- Configure a Service object with a YAML file similar to how Kubernetes controllers are configured.
- A sample configuration is as follows:
```YAML
apiVersion: v1
kind: Service
metadata:
  name: sample-service
  label:
    ## optional
  annotations:
    ## optional
spec:
  ports:
  - port: 80
    targetPort: 80
  selector:
    key1: value1
```
- The spec field is where we truly define the type of Kubernetes Service to use.


### What types of Kubernetes Services are available?
- There are 4 types of Kubernetes Service available; ClusterIP. NodePort, LoadBalancer and ExternalName

1. ClusterIP: the default type of Service. This is used to expose the Service object on a certain IP address only in the cluster.
2. NodePort: this Service object makes internal Pods accesseble externally via a port on the node on which the Pod is running.
3. LoadBalancer: this Service exposed the Pod externally via a LoadBalancer object; whether bare-metal or cloud.
4. ExternalName: this Service points to a DNS rather than a set of Pods.

- All these Service types use labels and label-selectors to match Pods, except for the ExternalName service.

- Let us delve into these Services more.

#### NodePort Service
- The same port on all nodes is exposed for access to the selected Pods. 
- The application can be accessed via the IP-Port combination: ***<NodeIP>:<NodePort>***.
- The sample configuration for a NodePort service is as follows:
```YAML
apiVersion: v1
kind: Service
metadata:
  name: example-nodeport
spec:
  type: NodePort ## spotlight here
  ports:
  - port: 80
    targetPort: 80
    nodePort: 32023
  selector:
    app: nginx
    environment: production
```
- In the ports field, we have:
  - ***targetPort***: this is the port the application in the Pod is running on. This is the port the Service forwards requests to. Usually, this and port are the same value.
  - ***port***: this is the Port of the Service itself.
  - ***NodePort***: this is the port on the node where the Service is exposed.
- In the selector field, these are the labels a Pod needs to have in order to be selected by a Service. Once a Service is deployed, it looks for the Pods that have these labels and adds *endpoints* for them.

- Assuming there is a Deployment with Pods that do not have Services attached to them. You can modify the Deployment manifest file to contain the Service YAML; you can write the Service YAML for the Deployment separately; or you can use the kubectl expose command.
  - when you use the kubectl expose command, you can conveniently leave out the NodePort and have that port assigned automatically. An format of the command is: `kubectl expose deployment <DEPLOYMENT_NAME> --name=<SERVICE_NAME> --port=<PORT_VALUE> --target-port=<TARGET_PORT_VALUE> --type=<SERVICE_TYPE>`.

#### ClusterIP Service
- ClusterIP is used to expose Pods internally in the cluster.
- ClusterIP is ideal for situations where we have decoupled application components communicating with each other. Eg. A frontend application Pod accessing the backend application Pods internally.
- The configuration of the ClusterIP Service is:
```YAML
apiVersion: v1
kind: Service
metadata:
  name: example-clusterIP
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 80
  selector:
    app: nginx
    environment: production
```
- In using the ClusterIP Service object, you can specify the IP address to use instead of letting it be randomly generated from the IP address pool available. The only requirement is that the custom IP should be in the range of the IP pool.
  - the IP pool is the CIDR on the default kubernetes Service in the cluster.
- With a ClusterIP specified, we can SSH into a node, and access Services from inside the cluster at that IP address.

#### LoadBalancer Service
- The LoadBalancer Service type exposes the application externally through a Loadbalancer object provided either by a cloud service or a locally optimized one like MetalLB.
- A LoadBalancer is like a superset of the NodePort service, it uses the loadbalancer implementation platform resources to assign an IP address to the Service.
- In setting up a LoadBalancer service, the implemetation platform usually requires that the user add some metadata in the form of annotations to the YAML manifest file.
- A simplified version of the YAML manifest file of a Loadbalancer service is:
```YAML
apiVersion: v1
kind: Service
metadata:
  name: example-loadbalancer
  annotations: ## optional
spec:
  type: LoadBalancer
  clusterIP: ## optional
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: nginx
    environment: production
```

#### ExternalName Service
- This Service maps the Service to a DNS name.
- The request that comes through the Service is responded to with a CNAME record that points to the a DNS name that was set.
- In this Service, the request is not proxied or forwarded, it is simply redirected. The redirection happens at the DNS level.
- The configuration of this Service is as follows:
```YAML
apiVersion: v1
kind: Service
metadata:
  name: example-externalName
spec:
  type: ExternalName
  externalName: my.example.site.com
```
- A use-case for ExternalName service is when you are migration production workload to a new Kubernetes cluster. You move the stateless applications first, then to make sure these stateless applications still access the other production services like databases, and other APIs, you provide an ExternalName service so that the Pods in the new cluster can still access resources from the old cluster.
- ExternalName serivce is good for referencing services that exist off the platform, on other clusters, or locally.

- A resource that is associated with Services is **Ingress**.
- The Ingress is a resource that acts as a middleman between the internet and the Services running on a cluster.
- By definition, the Ingress is an object that defines the rules that are used to manage external access to the Services in a Kubernetes cluster.