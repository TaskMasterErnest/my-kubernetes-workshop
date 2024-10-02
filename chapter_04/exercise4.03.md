# Enabling and Disabling API Groups and Versions on KinD Cluster
- As a the time of writing this, my Kubernetes version is 1.31 (the latest) hence almost everything is stable.
- The following example will walk you through how to enable and disable API groups and versions if you (or I ) ever come across them.

1. Start the KinD cluster with the following configuration:
```bash
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
runtimeConfig:
  "batch/v2alpha1" : "true"
nodes:
- role: control-plane
- role: worker
- role: worker
```

2. Check the details of the kube-apiserver pod. Use the describe command to filter the results using the runtime keyword: `kubectl describe pod kube-apiserver-k8s-workshop -n kube-system | grep runtime`.
  - the response will be this: `--runtime-config=batch/v2alpha1`.
  - another way to do this will be to get the api versions that have been enabled: `kubectl api-versions | grep batch/v2alpha1`.

3. You can create a CronJob manifest and apply it. It should run on the Kubernetes cluster.

4. Modify the API spec of the cluster in the API server pod. Use the command: `kubectl exec -it pod/kube-apiserver-k8s-workshop -n kube-system -- bash` to get into the API server pod.
  - run the command to view the manifest file for the API server; `sudo vi /etc/kubernetes/manifests/kube-apiserver.yaml`.
  - for the line that contains the runtime config, disable the batch/v2alpha1 modifying it to be like this `--runtime-config=batch/v2alpha1=false`. Save the modified file.

5. For the changes to take effect, we need to restart the API server and the controller-manager pods. 
  - they are stateless deployments hence a new version will be redeployed when we delete them.
  - delete the API server resources with this command: `kubectl delete pods -n kube-system -l component=kube-apiserver`. This will delete the pods that have the label "component=kube-apiserver".
  - delete the controller manager resources with this command: `kubectl delete pods -n kube-system -l component=kube-controller-manager`.

6. The list of api versions will not longer contain the `batch/v2alpha1` api version. Double check with the command: `kubectl api-versions | grep batch/v2alpha1`