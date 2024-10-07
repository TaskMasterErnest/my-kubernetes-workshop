- Imagine you are working with a team of developers who have built an awesome application that they want you to deploy in a pod. The application has a process that starts up and takes approximately 20 seconds to load all the required assets.
- Once the application starts up, it's ready to start receiving requests. If, for some reason, the application crashes, you would want the pod to restart itself as well. They have given you the task of creating the pod using a configuration that will satisfy these needs for the application developers in the best way possible.

- How I chose to tackle this:
1. Create a new namespace for this. activity-example
2. Create a Pod configuration.
```YAML
apiVersion: v1
kind: Pod
metadata:
  name: custom-application-pod
  namespace: activity-example
spec:
  restartPolicy: Always
  containers:
  - name: custom-application-container
    image: packtworkshops/the-kubernetes-workshop:custom-application-for-pods-chapter
    readinessProbe:
      exec:
        - cat
        - /tmp/health
      initialDelaySeconds: 20
      periodSeconds: 10
```