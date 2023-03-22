# echo-server

This is a simple HTTP server that will answer a 200 OK message to every request.
For every request that arrives the server will print the sender's address, the path used, and the body's content.
You can change the listening address by setting the LISTEN_ADDRESS and LISTEN_PORT environment variables.

Available endpoints:

- `/metrics/`  Prometheus metrics
- `/health/`   returns always 200. Could change in the future
- `/liveness/` returns always 200. Could change in the future
- `/<path>/`   returns 200. The request details get printed to stdout.

This server can be used for example to debug if Alertmanager has been properly configured and it's sending alerts to the webhooks.

Just point Alertmanager to the `echo-server` and you should see the payload in `stdout` every time Alertmanager sends an alert.

## Building

### Binary

```console
go build
```

### Container image

We use `ko` to build a multi-arch container image with SBOM included.

```console
KO_DOCKER_REPO=ralgozino ko build --platform=all -B
```

## Deploying into a Kubernetes cluster

You can use the following manifest:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: notifications-receiver
  name: notifications-receiver
  namespace: monitoring
spec:
  replicas: 1
  selector:
    matchLabels:
      app: notifications-receiver
  strategy: {}
  template:
    metadata:
      labels:
        app: notifications-receiver
    spec:
      containers:
        - image: ralgozino/echo
          name: echo
          resources:
            requests:
              cpu: 100m
              memory: 256Mi
            limits:
              cpu: 1000m
              memory: 1Gi
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          readinessProbe:
            httpGet:
              path: /health
              port: http
          livenessProbe:
            httpGet:
              path: /liveness
              port: http
---
apiVersion: v1
kind: Service
metadata:
  name: notifications-receiver
  namespace: monitoring
spec:
  selector:
    app: notifications-receiver
  ports:
    - port: 80
      targetPort: http
```

Now simply configure Alertmanager to send alerts to `http://notificationes-receiver.monitoring.svc`. You should see the payloads in the `notifications-receiver` pod.
