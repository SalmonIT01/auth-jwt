# ใช้ golang:latest เพื่อได้เวอร์ชันล่าสุด
FROM golang:latest AS builder

# ติดตั้ง dependencies ที่จำเป็น
RUN apk update && apk add --no-cache git || apt-get update && apt-get install -y git

# ตั้งค่า working directory
WORKDIR /app

# คัดลอก go.mod และ go.sum ก่อน
COPY go.mod go.sum ./ 

# ดาวน์โหลด dependencies
RUN go mod download

# คัดลอกโค้ดทั้งหมด
COPY . .

# Build แอปพลิเคชัน
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-api ./cmd/main.go

# ขั้นตอนการสร้างภาพสุดท้าย
FROM alpine:latest

# ติดตั้งแพคเกจที่จำเป็น
RUN apk --no-cache add ca-certificates tzdata

# ตั้งค่า timezone
ENV TZ=Asia/Bangkok

# ตั้งค่า working directory
WORKDIR /root/

# คัดลอกไฟล์ที่ build มาจากขั้นตอนก่อนหน้า
COPY --from=builder /app/auth-api .

# เปิด port
EXPOSE 8080

# คำสั่งที่จะทำงานเมื่อรัน container
CMD ["./auth-api"]
