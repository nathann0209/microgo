package handlers

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Post struct {
	l *log.Logger
}

type LogRecord struct {
	IP           string
	Endpoint     string
	Method       string
	RequestSize  int64
	ResponseSize int
	StatusCode   int
	Timestamp    time.Time
	duration     time.Duration
}

func NewLogRecord(IpAddr string, Endp string, Meth string, ReqSize int64,
	RespSize int, StatCode int, ts time.Time, duration time.Duration) *LogRecord {
	return &LogRecord{IpAddr, Endp, Meth, ReqSize, RespSize, StatCode, ts, duration}
}

func NewPost(l *log.Logger) *Post {
	return &Post{l}
}

func (p *Post) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	p.l.Println("Hello from your POST handler!")

	// Extract the information about request
	start := time.Now()

	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.RemoteAddr
	}
	endpoint := r.URL.Path
	method := r.Method
	requestSize := r.ContentLength
	if requestSize == -1 {
		requestSize = 0
	}
	timestamp := time.Now()

	// Write response
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "Oops, something went wrong reading the request body.", http.StatusBadRequest)
		return
	}
	response := fmt.Sprintf("Posted for request data: %s!", string(data))
	rw.Write([]byte(response))

	duration := time.Since(start)

	// Prepare Reqeuest/Response data for logging
	logRecord := NewLogRecord(
		ip,
		endpoint,
		method,
		requestSize,
		len(response),
		http.StatusOK,
		timestamp,
		duration,
	)

	record := []string{
		logRecord.IP,
		logRecord.Endpoint,
		logRecord.Method,
		strconv.FormatInt(logRecord.RequestSize, 10),
		strconv.Itoa(logRecord.ResponseSize),
		strconv.Itoa(logRecord.StatusCode),
		logRecord.Timestamp.Format(time.RFC3339),
		logRecord.duration.String(),
	}

	// Check if the CSV file exists
	fileExists := true
	if _, err := os.Stat("logs/request_metrics.csv"); os.IsNotExist(err) {
		fileExists = false
	}

	// Open or create the CSV file in append mode
	file, err := os.OpenFile("logs/request_metrics.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to write records to CSV File: %s", err)
	}

	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers if the file did not exist
	if !fileExists {
		headers := []string{"IP", "Endpoint", "Method", "RequestSize", "ResponseSize", "StatusCode", "Timestamp", "Duration"}
		err = writer.Write(headers)
		if err != nil {
			log.Fatalf("Failed to write headers to CSV file: %s", err)
		}
	}

	// Write the data to the CSV file
	writeErr := writer.Write(record)
	if writeErr != nil {
		log.Fatalf("Failed to write record to CSV file: %s", writeErr)
	}

}
