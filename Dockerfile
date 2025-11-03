# ---- Giai đoạn 1: Build ----
# Sử dụng một image Golang chính thức làm builder. Chúng ta chỉ định phiên bản để đảm bảo tính tái lập.
FROM golang:1.25-alpine AS builder

# Thiết lập thư mục làm việc bên trong container.
WORKDIR /app

# Sao chép các tệp module Go và tải về các dependency.
# Việc này được thực hiện trong một layer riêng để tận dụng cache của Docker.
# Các dependency chỉ được tải lại nếu go.mod hoặc go.sum thay đổi.
COPY go.mod ./
RUN go mod download

# Sao chép phần còn lại của mã nguồn.
COPY . .

# Biên dịch ứng dụng, tạo ra một file binary được liên kết tĩnh.
# CGO_ENABLED=0 vô hiệu hóa Cgo, cần thiết cho một bản build tĩnh.
# -ldflags="-w -s" loại bỏ thông tin gỡ lỗi, giúp giảm kích thước file binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /quote-api

# ---- Giai đoạn 2: Runtime ----
# Sử dụng một base image tối giản, không phải root. 'distroless' là một lựa chọn tuyệt vời
# vì nó chỉ chứa ứng dụng của chúng ta và các dependency runtime của nó.
# Nó không chứa trình quản lý gói, shell, hoặc các chương trình khác, giúp giảm đáng kể bề mặt tấn công.
FROM gcr.io/distroless/static-debian11

# Sao chép file binary đã biên dịch từ giai đoạn 'builder'.
COPY --from=builder /quote-api /quote-api

# Mở cổng mà ứng dụng sẽ chạy trên đó.
EXPOSE 8084

# Thiết lập người dùng là một người dùng không phải root để tăng cường bảo mật.
USER nonroot:nonroot

# Lệnh sẽ được chạy khi container khởi động.
ENTRYPOINT ["/quote-api"]
