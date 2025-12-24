package tcphttp

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

func Client() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.Write([]byte("hello"))
}

func HttpClient() {
	resp, err := http.Get("http://127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	res, err := http.Post("http://127.0.0.1:8080/hello", "application/json", nil)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("body: %s \n", body)

	largeTextFile, err := os.OpenFile("large_text.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer largeTextFile.Close()

	fileInfo, err := largeTextFile.Stat()
	if err != nil {
		panic(err)
	}
	fileSize := fileInfo.Size()
	fmt.Printf("fileSize: %d \n", fileSize)

	client := &http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodPost, "127.0.0.1", largeTextFile)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", fileSize))
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	readerLimit := io.LimitedReader{R: resp.Request.Body, N: 1024 * 1024}
	responseBody, err := io.ReadAll(readerLimit)
	if err != nil {
		panic(err)
	}
	fmt.Printf("responseBody: %s \n", responseBody)

	_, err = io.Copy(io.Discard, resp.Body)

	if err != nil {
		panic(err)
	}

}

func sendLargeMessage(url string, filePath string) error {

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	req, err := http.NewRequest(http.MethodPost, url, file)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	req.ContentLength = fileInfo.Size()
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return nil
	//req.Header.Set("Content-Length", fmt.Sprintf("%d", fileSize))

}

func handleLargeMessage(w http.ResponseWriter, r *http.Request) {

	tmpFile, err := os.CreateTemp("", "large_text_*.txt")
	if err != nil {
		panic(err)
	}
	defer tmpFile.Close()

	defer os.Remove(tmpFile.Name())

	defer tmpFile.Close()

	written, err := io.Copy(tmpFile, r.Body)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("written: %d", written)))
}

func handleLargeMessageCustom(w http.ResponseWriter, r *http.Request) {

	tmpFile, err := os.CreateTemp("", "large_text_*.txt")
	if err != nil {
		panic(err)
	}
	defer tmpFile.Close()

	defer os.Remove(tmpFile.Name())

	defer tmpFile.Close()

	chunkSize := 1024 * 64
	buf := make([]byte, chunkSize)
	written, err := io.CopyBuffer(tmpFile, r.Body, buf)
	if err != nil && err != io.EOF {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("written: %d", written)))

}

func handleLargeMessageCustomTransport(w http.ResponseWriter, r *http.Request) {
	customTransport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     time.Second * 30,
		TLSHandshakeTimeout: time.Second * 10,
	}

	client := &http.Client{
		Transport: customTransport,
		Timeout:   time.Second * 5,
	}
	client.Transport.RoundTrip(r)

	resp, err := client.Get("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	proxy := httputil.NewSingleHostReverseProxy(r.URL)
}
