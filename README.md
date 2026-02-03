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
