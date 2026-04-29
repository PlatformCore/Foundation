# libpackage Enterprise Reorganization

จัดกลุ่มใหม่แบบ enterprise โดยคัดลอกจาก repo เดิมทั้งหมด:
- ไม่ rewrite logic ในไฟล์ .go
- `compat/gometrics` ถูกยกเป็นตัวหลักที่ `observability/gometrics`
- `observability/metrics` เดิมถูกเก็บเป็น `observability/metrics/basic`
- `compat/go...` อื่น ๆ ถูกย้ายเข้า domain ที่ตรงกับหน้าที่ เช่น security/resilience/ratelimit/observability/core
- ไฟล์ audit/docs เดิมเก็บใน `docs/original`
- ดูรายการครบทุกไฟล์ที่ `ENTERPRISE_REORG_FULL_FILE_TREE.txt`
- ดู mapping ไฟล์เก่า -> ไฟล์ใหม่ที่ `ENTERPRISE_REORG_MAPPING.csv`

หมายเหตุ:
- ยังไม่ได้ rewrite import path ภายใน .go files
- หลังย้ายจริงควรทำ pass ที่ 2: ปรับ `go.mod`, `go.work`, และ import path ให้ตรง module ใหม่
