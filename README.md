Kubster
=======

Kubster is a simple web server that can be used to demonstrate k8s liveness and readiness probes'
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
