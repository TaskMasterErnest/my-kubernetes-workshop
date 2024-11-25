# ConfigMaps And Secrets
- A truth in having containerized applications is to decouple the applications. With this decoupling comes the way to provide configuration data to each container so they can work.
- Kubernetes provides a way to associate environment-specific data with containerized applications without making changes to the container image.
- The ways to do this are: (1) provide command-line arguments to the Pods, (2) provide environment variables to the Pod and (3) mount configuration files in the containers.
- For configuration data, Kubernetes provides **ConfigMaps** and for sensitive data, Kubernetes provides **Secrets** - that is the main difference between them, otherwise they work the same way by injecting data into the Pod containers.


## ConfigMaps
- A configmap is a Kubernetes object that is used to define application configuration-related data.
- It allows the client/admin to inject configuration data to multiple containers running in different environments.
- ConfigMaps are mainly used to store non-sensitive configuration data such as config files or environment variables.
- The data in ConfigMaps is loaded into Pod containers as read-only data.
- ConfigMaps can be used to hold configuration data for system applications such as Kubernetes operators and Kubernetes controllers.
- ConfigMaps can be created. The command to create configmap objects in Kubernetes is `kubectl create configmap <MAP_NAME> <DATA_SOURCE>`.
  - the map_name is the name of the configmap and the data source is the value, file or directory to draw the data from.
  - the data source takes in a key-value pair format.
- ConfigMaps can be created for a single value, a list of values, or values from an entire file or directory.
  - To see how to specify these commands, run the `kubectl create configmap --help` command.

- There are various ways to use ConfigMaps, One of them is to define a configmap from a file and then mount that file inside a Pod container.
  - this method is useful when we have external configuration data that differs between environments. We then load the correct configuration into the right environment using a ConfigMap.
  - an example will be that a packaged application container has to connect to a specific database URL when it moves from dev to test environments. We can configure the application to read the specific database URL from a file. We then package this file as a ConfigMap that can be mounted unto the container.
  - Also we can use ConfigMaps to set out specific config parameters that can be mounted when needed.

- In the `configmap-as-volume.yaml`, we are mounting a ConfigMap as a volume. We created the configmap from a file `application.properties`.


## Secrets
- A Secret is a Kubernetes resource used to store sensitive data.
- The differences between a ConfigMap and a Secret are:
  - (1) a secret is intended to store a small amount of sensitive data (<=1MB). It is alse base64-encoded hence cannot be treated as secure. It can store binary data - such as public and/or private keys.
  - (2) Kubernetes makes sure that Secrets are only passed on to the nodes that run the Pods that need these Secrets.
- Secrets are not secure, the way to make sure they are is to encrypt them. There are Kubernetes solutions that you can use to encrypt secrets.
- Secrets can be mounted as environment variables or as files to the Pods that need them.
- There are three(3) types of Secrets, they can be found when you run the `kubectl create secret --help` command.
  - (1) **generic** - this holds any custom key-value pair
  - (2) **tls** - this secret is for holding a public-private key pair for communication with the TLS protocol.
  - (3) **docker-registry** - this stores the username, password, and email address to access a Docker registry.