# my-kubernetes-workshop
A repo containing my deep dive into how Kubernetes functions in detail.

## Why Am I Doing This?
- I have found my Kubernetes knowledge lacking in certain areas and I need to bridge those gaps.
- Kubernetes is now a constant tool in application deployment so I have to get up to speed with it.
- It will be fun to learn (all over again but in more detail) how Kubernetes works (almost not as the first point lol!)
- Maybe I can get a job working with Kubernetes if someone comes across this repo. (I will never know)

## NOTE:
- I will be using Podman to do the Docker-related work in this repo because I don't want to install Docker, not out of spite but I run a Fedora Workstation and it comes with Podman pre-installed. So I just decided to go with it.
- Podman and Docker have almost similar commands, easily transferable/translatable syntax.
- For my Kubernetes, I installed Docker finally and used that to power my local KinD cluster. Podman does rootless, giving no access to privileged ports, making it a headache to configure rootful mode for my experiment.
- For the parts where I have to expose my application (without using NGINX Ingress), I use MetalLB to configure a LoadBalancer to do so. [Here](https://metallb.org/)
