Goal: Build a tiny, locally distributed system with two services that communicate over the network, include basic logging, and demonstrate independent failure.

What You Will Build
Service A (Echo API) on localhost:8080
GET /health → {"status":"ok"}
GET /echo?msg=hello → {"echo":"hello"}
Service B (Client) on localhost:8081
GET /health → {"status":"ok"}
GET /call-echo?msg=hello → calls Service A /echo and returns a combined response
Requirements
Two independent processes (run in separate terminals)
HTTP (or gRPC for stretch)
Basic request logging: service name, endpoint, status, latency
Service B must use a timeout when calling Service A
Demonstrate failure: stop Service A; Service B returns 503 and logs an error
Install / Setup
Git
One runtime: Python 3.10+ or Go 1.21+ (or Java 17+ if your team prefers)
Docker Desktop (recommended for future labs)
Starter Repo
https://github.com/ranjanr/cmpe273-week1-lab1-starter

Use the above-provided starter and pick one track:

python-http/ (Flask + requests)
go-http/ (net/http)
How to Test
Success:

curl "http://127.0.0.1:8081/call-echo?msg=hello"
Failure: Stop Service A (Ctrl+C), then rerun the same curl command. Expect HTTP 503 and a clear error log.


