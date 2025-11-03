package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// Quote định nghĩa cấu trúc dữ liệu cho một câu trích dẫn.
type Quote struct {
	Text   string `json:"quote"`
	Author string `json:"author"`
}

// quotes là một slice chứa các câu trích dẫn mẫu.
// Trong các mô-đun sau, chúng ta sẽ thay thế nó bằng một cơ sở dữ liệu thực sự.
var quotes = []Quote{
	{"The greatest glory in living lies not in never falling, but in rising every time we fall.", "Nelson Mandela"},
	{"The way to get started is to quit talking and begin doing.", "Walt Disney"},
	{"Your time is limited, so don't waste it living someone else's life.", "Steve Jobs"},
	{"If life were predictable it would cease to be life, and be without flavor.", "Eleanor Roosevelt"},
}

// quoteHandler xử lý các yêu cầu đến và trả về một câu trích dẫn ngẫu nhiên.
func quoteHandler(w http.ResponseWriter, r *http.Request) {
	// Chọn một câu trích dẫn ngẫu nhiên từ slice.
	rand.Seed(time.Now().UnixNano())
	quote := quotes[rand.Intn(len(quotes))]

	// Thiết lập header Content-Type là application/json.
	w.Header().Set("Content-Type", "application/json")
	// Mã hóa câu trích dẫn thành JSON và gửi về cho client.
	json.NewEncoder(w).Encode(quote)
}

func main() {
	// Lấy cổng từ biến môi trường PORT, nếu không có thì mặc định là 8080.
	// Cloud Run sẽ tự động cung cấp biến môi trường này.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	// Đăng ký handler cho đường dẫn gốc "/".
	http.HandleFunc("/", quoteHandler)

	log.Printf("Server starting on port %s\n", port)
	// Khởi động máy chủ web.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
