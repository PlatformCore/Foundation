# Integrated single-source middleware layout

รอบนี้ปรับแนวคิดเป็น **รวมระบบเดียว** ไม่ใช่เพิ่มชุดใหม่ซ้อนชุดเก่า

## หลักการ

- logic หนักของเดิมยังเป็นตัวหลัก:
  - `observability/logging`, `observability/logging/gologger`
  - `observability/gometrics`
  - `resilience/retry/goretry`
  - `resilience/timeout/gotimeout`
  - `ratelimit`
  - `messaging/idempotency`
- unified middleware ทำหน้าที่เป็น adapter เข้า `transport/core.Middleware`
- ของเดิมยังอยู่เพื่อ compat แต่ path ใหม่เรียก engine เดิม ไม่ได้สร้าง retry/logger/metrics/ratelimit ชุดที่สอง

## Flow

```txt
adapter/{http,grpc,nats,kafka,mqtt,redis}
  -> transport/core.Context
  -> middleware/chain
  -> middleware wrappers using original heavy engines
  -> handler/service
```

## Main preset

ใช้ `middleware/preset.UnifiedEnterprise(...)` เป็น chain หลักได้เลย
