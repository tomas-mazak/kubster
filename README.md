Kubster
=======

Kubster is a simple http server that can be used to demonstrate k8s liveness and readiness probes'
concept.

API
---

  * `GET /` - always returns `200 OK` and print nice k8s logo
  * `GET /live` - returns `200 OK` (default) or `503 Service Unavailable` (if set)
  * `GET /ready` - returns `503 Service Unavailable` after start, changes to `200 OK` after 
    `KUBSTER_READYDELAY` number of seconds (default), or can be overriden
  * `GET /set?live=<true|false>&ready=<true|false>` - override `/live` and `/ready` endpoint 
    responses (true for `200 OK`, false for `503 Service Unavailable`)

Configuration
-------------

Environment variables:

  - `KUBSTER_BIND` - host:port to bind to and listen on
  - `KUBSTER_READYDELAY` - number of seconds before the `/ready` endpoint starts succeeding


Kubernetes
----------

Clone this repo and make it your current working directory first.

```
# Create the pod first, the pod name will be kubster-pod
kubectl apply -f manifests/pod.yaml

# Run multiple times to observe how pods gets ready after 30 seconds from launching
kubectl describe pod kubster-pod

# You can use port forwarding to access the http server and view/modify the responses
kubectl port-forward kubster-pod 8080:8080 &
curl -s localhost:8080/
curl -sv localhost:8080/live
curl -sv localhost:8080/ready

# Make kubernetes restart the pod
curl -s localhost:8080/set?live=false
```
