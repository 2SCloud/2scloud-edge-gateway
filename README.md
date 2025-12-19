# Roadmap v1.0.0 (Envoy + K8s)

- [ ] Kubernetes running (local or server)
- [ ] The Edge-Gateway is the only entrypoint
- [ ] Test backend exposed via TLS
- [ ] TLS operationnal
- [ ] `app.example.com` -> service
- [ ] Observability OK
- [ ] Minimal Control Plane in Golang
- [ ] Rust structure ready for future modules

## Deployment
```bash
make deploy
```

## Project struct explained:

- `cmd/edge-gateway/main.go` -> Entrypoint

- `internal/server` -> HTTP/TLS listeners
- `internal/proxy/` -> Reverse proxy logic
- `internal/config/` -> Config loading (file/env/api)
- `internal/runtime` -> Module runtime (WASM)
- `internal/middleware` -> Glue between `Go` & `Rust`
- `internal/observability` -> All metrics about the Edge-Gateway

- `modules/` -> Rust security modules
- `modules/waf` -> Web Application Firewall
- `modules/firewall` -> L3/L4/L7 filtering
- `modules/ratelimit` -> Rate limiting / quotas
- `modules/bot-detection` -> Bot / scraper detection
- `modules/ddos` -> DDoS heuristics
- `modules/authz` -> AuthZ (JWT, OAuth scopes)
- `modules/authn` -> AuthN (mTLS, tokens)
- `modules/ip-reputation` -> IP reputation / geo / ASN
- `modules/header-sanitizer` -> IP reputation / geo / ASN
- `modules/request-validator` -> Schema / contract validation
- `modules/response-filter` -> Response masking / DLP
- `modules/anomaly-detection` -> Behavioral anomalies
- `modules/circuit-breaker` -> Backend protection
- `modules/shadow-mode` -> Observe-only security
- `modules/sandbox` -> Untrusted plugins

- `api/admin/openapi.yaml` -> Admin API

- `deploy/docker/Dockerfile` -> Deploy on docker
- `modules/kubernetes/edge-gateway.yaml` -> Deploy on kubernetes

- `scripts/build-wasm.sh`
- `scripts/run-dev.sh`


# Roadmap v2.0.0

## Next steps

- Add Rust modules for WAF, rate-limiting, bot detection
- Implement dynamic rule push from Go control-plane to Envoy / Rust modules
- Advanced observability (Prometheus, Grafana, alerting)
- Authentication and policy management