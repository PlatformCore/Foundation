# Foundation Release Notes (TH) - Dev + PM

วันที่ออก: 2026-04-29  
รีโป: https://github.com/PlatformCore/Foundation

## ภาพรวมรอบนี้
- ปรับ `module path` ของ Foundation ทั้งหมดให้เป็น `github.com/PlatformCore/Foundation/...`
- จัดโครงสร้างโมดูลแบบติดตั้งแยกได้ (multi-module)
- เพิ่มระบบจัดการเวอร์ชันผ่าน `tools/module_versions.json`
- เพิ่มสคริปต์ออกเวอร์ชันแบบคำสั่งเดียว (bump + commit + push + tag)

## เวอร์ชันโมดูล (รอบนี้)
| โมดูล | Path | Version | สถานะ | เหมาะสำหรับ |
|---|---|---|---|---|
| clients | `clients` | v0.1.0 | stable | HTTP/gRPC/NATS client layer |
| config | `config` | v0.1.0 | stable | config loading/env/file/watcher |
| core | `core` | v0.1.0 | stable | contracts/types/errors/logger พื้นฐาน |
| messaging | `messaging` | v0.1.0 | stable | inbox/outbox/replay/dlq |
| middleware | `middleware` | v0.1.0 | stable | HTTP/GRPC middleware กลาง |
| observability | `observability` | v0.1.0 | stable | logging/metrics/tracing/profiler |
| orchestration | `orchestration` | v0.1.0 | stable | DI/workflow/orchestration |
| persistence | `persistence` | v0.1.0 | stable | uow/tx/distributed lock |
| platform | `platform` | v0.1.0 | stable | eventbus/featureflag/servicemesh |
| plugins | `plugins` | v0.1.0 | stable | plugin engine/registry/hooks |
| ratelimit | `ratelimit` | v0.1.0 | stable | rate limiting policy/store |
| resilience | `resilience` | v0.1.0 | stable | retry/circuit/timeout/bulkhead |
| runtime | `runtime` | v0.1.0 | stable | lifecycle/shutdown/signals |
| security | `security` | v0.1.0 | stable | auth/jwt/oauth2/encryption |
| tools | `tools` | v0.1.0 | beta | เครื่องมือช่วย generate/release |
| validation | `validation` | v0.1.0 | stable | binding/schema validation |

## วิธีอัปเวอร์ชันแบบคำสั่งเดียว
ใช้สคริปต์:
`./scripts/release_one_command.ps1 -ModulePath platform -Version v0.2.0`

สิ่งที่สคริปต์ทำ:
1. อัปเดตเวอร์ชันใน `tools/module_versions.json`
2. commit เข้า `main`
3. push `main` ไป `origin`
4. สร้างและ push tag รูปแบบ `modulePath/version` เช่น `platform/v0.2.0`

## หมายเหตุสำหรับ PM
- ทุกโมดูลออกเวอร์ชันแยกได้ จึงวางแผน rollout เป็นรายโดเมนได้
- แนะนำคง `v0.x` ช่วง hardening และยก `v1.0.0` เมื่อ API/behavior คงที่
