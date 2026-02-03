# CMPE 273 – Week 1 Lab 1: Your First Distributed System (Starter)

This starter provides two implementation tracks:
- `python-http/` (Flask + requests)
- `go-http/` (net/http)

Pick **one** track for Week 1.

## Lab Goal
Build **two services** that communicate over the network:
- **Service A** (port 8080): `/health`, `/echo?msg=...`
- **Service B** (port 8081): `/health`, `/call-echo?msg=...` calls Service A

Minimum requirements:
- Two independent processes
- HTTP (or gRPC if you choose stretch)
- Basic logging per request (service name, endpoint, status, latency)
- Timeout handling in Service B
- Demonstrate independent failure (stop A; B returns 503 and logs error)

## How to Run Locally (Go Track)

**Terminal 1 – Start Service A:**
```bash
cd go-http/service-a && go run main.go
```

**Terminal 2 – Start Service B:**
```bash
cd go-http/service-b && go run main.go
```

## Test Output

**Success case (both services running):**
```
$ curl "http://127.0.0.1:8081/call-echo?msg=hello"
{"from_a":{"echo":"hello"},"service":"B"}
```

**Failure case (Service A stopped):**
```
$ curl -w "\nHTTP Status: %{http_code}\n" "http://127.0.0.1:8081/call-echo?msg=hello"
{"error":"ServiceA unavailable"}
HTTP Status: 503
```

Service B logs the error:
```
[ServiceB] ERROR calling ServiceA: Get "http://localhost:8080/echo?msg=hello": dial tcp [::1]:8080: connect: connection refused
```

## What Makes This Distributed?

This system is distributed because it consists of two independent processes (Service A and Service B) that communicate over the network via HTTP rather than sharing memory or running in the same process. Each service can be started, stopped, and scaled independently. Service B does not crash when Service A becomes unavailable—instead, it gracefully handles the failure by returning a 503 status code and logging the error. This demonstrates a key property of distributed systems: partial failure, where one component can fail while others continue operating.

## Failure Handling

### What happens on timeout?

Service B uses an HTTP client with a 2-second timeout (`http.Client{Timeout: 2 * time.Second}`). If Service A takes longer than 2 seconds to respond, the request is cancelled and Service B returns HTTP 503 with `{"error":"ServiceA unavailable"}`. The timeout error is logged:
```
[ServiceB] ERROR calling ServiceA: Get "http://localhost:8080/echo?msg=hello": context deadline exceeded
```

### What happens if Service A is down?

If Service A is not running, Service B receives a "connection refused" error. It handles this gracefully by:
1. Logging the error with full details
2. Returning HTTP 503 (Service Unavailable) to the client
3. Continuing to serve other requests normally

```
[ServiceB] ERROR calling ServiceA: Get "http://localhost:8080/echo?msg=hello": dial tcp [::1]:8080: connect: connection refused
```

## Logging and Debugging

### What do the logs show?

Every request is logged with: service name, HTTP method, path, status code, and latency.

**Service A logs:**
```
[ServiceA] listening on :8080
[ServiceA] GET /health 200 45.208µs
[ServiceA] GET /echo 200 28.125µs
```

**Service B logs:**
```
[ServiceB] listening on :8081
[ServiceB] GET /health 200 31.042µs
[ServiceB] GET /call-echo 200 1.234ms
[ServiceB] ERROR calling ServiceA: ...
[ServiceB] GET /call-echo 503 52.167µs
```

### How to debug?

1. **Check if services are running:** `curl http://127.0.0.1:8080/health` and `curl http://127.0.0.1:8081/health`
2. **Check ports:** `lsof -i:8080` and `lsof -i:8081`
3. **Watch logs in real-time:** Run each service in a separate terminal to see logs as requests come in
4. **Test Service A directly:** `curl "http://127.0.0.1:8080/echo?msg=test"` to isolate issues
5. **Check for timeout vs connection refused:** Timeout errors indicate Service A is slow; connection refused means it's not running
