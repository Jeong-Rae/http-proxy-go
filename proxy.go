package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var fileLogger *log.Logger

func setupLogger() {
	logsPath := "./logs"
	fileName := "info.log"
	fullPath := logsPath + "/" + fileName

	// logs 디렉토리 생성
	if err := os.MkdirAll(logsPath, 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		return
	}

	// 로그 파일 열기
	logFile, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}

	fileLogger = log.New(logFile, "", 0) // 로그 파일을 위한 로거 생성, 타임스탬프 없음
}

func logRequest(r *http.Request, bodyBytes []byte) {
    requestLine := fmt.Sprintf("%s %s %s\n", r.Method, r.URL.Path, r.Proto)

    // 헤더를 순회하며 출력 형식을 조정
    var headerStr string
    for name, values := range r.Header {
        for _, value := range values {
            headerStr += fmt.Sprintf("%s: %s\n", name, value)
        }
    }

    bodyStr := string(bodyBytes)
    if len(bodyStr) > 0 {
        bodyStr += "\n\n" 
    }

    completeLog := fmt.Sprintf("%s%s\n%s", requestLine, headerStr, bodyStr)

    // 로그를 비동기적으로 기록
    go func(logData string) {
        fmt.Print(logData) // 표준 출력에 로그 출력
        fileLogger.Print(logData) // 파일에 로그 출력
    }(completeLog)
}


func getTargetURL() string {
    targetHost := os.Getenv("LOGGING_TARGET_HOST")
    if targetHost == "" {
        targetHost = "localhost:5050"
    }
    return "http://" + targetHost
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	targetURL, _ := url.Parse(getTargetURL())
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	logRequest(r, bodyBytes)

	proxy.ServeHTTP(w, r)
}



func main() {
	setupLogger()
	http.HandleFunc("/", handleProxy)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
