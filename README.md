# Traefik Startup Probe sidecar

Traefik 1.x sidecar container for [Kubernetes Startup Probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/#define-startup-probes)

## Background

Traefik provides a `/ping` endpoint that can be used for Kubernetes Liveness and Readiness Probes.  This endpoint starts returning 200 before all the ingress resources have been loaded from k8s, which could yield errors when a traefik deployment scales up.

This is particularly noticeable when using a [Global Default Backend](https://doc.traefik.io/traefik/v1.7/configuration/backends/kubernetes/#global-default-backend-ingresses) for configuring wildcard TLS certificates: in certain conditions the default certificate loaded onto the Traefik pods is returned before all ingresses have been loaded.

This has been tested on Traefik 1.7.24.

## Usage

1. Add the sidecar container to the Traefik deployment:

        - env:
          - name: LOG_LEVEL
            value: debug
          image: registry.example.com/traefik-startup-probe:latest
          name: startup-probe

2. Add a `startupProbe` to the traefik container:

        startupProbe:
          failureThreshold: 30
          httpGet:
            path: /healthz
            port: 8083
            scheme: HTTP
          periodSeconds: 3
          successThreshold: 1
          timeoutSeconds: 1
