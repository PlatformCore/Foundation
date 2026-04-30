# Enterprise Transport + Adapter Completion Report

เพิ่มโดยไม่ลบโค้ดเดิม:

## 1. transport/ กลาง
- core
- http
- http3
- grpc
- nats
- kafka
- mqtt
- redis

## 2. adapter/ แยกชัด
- nethttp
- http3
- grpc
- nats
- kafka
- mqtt
- redis

## 3. middleware เพิ่มแบบแยก transport
- middleware/logging/http.go, grpc.go, nats.go, kafka.go, mqtt.go, redis.go
- middleware/recovery
- middleware/ratelimit
- middleware/tracing
- middleware/timeout
- middleware/idempotency
- middleware/metrics

หลักการจัดวาง:
- transport/ = interface/contract กลาง
- adapter/ = framework/protocol bridge
- middleware/ = flow control / cross-cutting behavior
- clients/ = outbound client wrapper ของเดิม
