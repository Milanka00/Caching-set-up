FROM envoyproxy/envoy-alpine:v1.20-latest
COPY envoy.yaml /etc/envoy/envoy.yaml
CMD ["envoy", "-c", "/etc/envoy/envoy.yaml", "--log-level", "info"]
