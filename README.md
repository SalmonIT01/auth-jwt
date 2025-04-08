# Auth JWT API

API สำหรับระบบยืนยันตัวตนใช้ JWT Token พัฒนาด้วย Go และ MongoDB

## ความต้องการระบบ

- [Docker](https://www.docker.com/products/docker-desktop/)
- [Docker Compose](https://docs.docker.com/compose/install/) (มักจะมาพร้อมกับ Docker Desktop)

## วิธีการรัน

1. Clone โปรเจคนี้:
```bash
git clone https://github.com/SalmonIT01/auth-jwt.git
cd auth-jwt
```

2. สร้างไฟล์ `.env` (ตัวอย่าง): ******ไม่ต้องทำก็ได้********
```bash
cp .env.example .env
```

3. แก้ไขไฟล์ `.env` ตามต้องการ (เช่น เปลี่ยน JWT_SECRET)

4. รันแอปพลิเคชันด้วย Docker:
```bash
docker-compose up -d
```

5. ตรวจสอบว่าแอปทำงานได้:
```bash
docker ps
```

6. เข้าถึง API ได้ที่ `http://localhost:8080`

## เส้นทาง API

### สาธารณะ (Public)
- `POST /api/register` - ลงทะเบียนผู้ใช้ใหม่
- `POST /api/login` - เข้าสู่ระบบ

### ต้องยืนยันตัวตน (Protected)
- `GET /api/profile` - ดูข้อมูลโปรไฟล์
- `PUT /api/profile` - อัปเดตข้อมูลโปรไฟล์
- `DELETE /api/profile` - ลบบัญชีผู้ใช้

## ตัวอย่างการใช้งาน

### ลงทะเบียนผู้ใช้ใหม่
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123","fullname":"Test User","tel":"0812345678"}'
```

### เข้าสู่ระบบ
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### ดูโปรไฟล์ (ต้องมี token)
```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## การหยุดการทำงาน

```bash
docker-compose down
```

ถ้าต้องการลบ volumes ด้วย:
```bash
docker-compose down -v
```
