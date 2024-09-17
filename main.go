package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const maxWorkers = 10

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan net.Conn)
	var wg sync.WaitGroup

	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())

	// Start worker goroutines
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(ctx, jobs, &wg)
	}

	// Set up channel to listen for OS interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Handle graceful shutdown in a goroutine
	go func() {
		<-stop
		fmt.Println("\nShutting down server...")
		listener.Close() // This will cause listener.Accept() to return an error
		close(jobs)      // Close the jobs channel to signal workers to exit
		cancel()         // Cancel the context to signal connections to close
	}()

	// Accept connections in the main goroutine
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Check if the error is due to the listener being closed
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				// Listener has been closed, exit the accept loop
				break
			}
			log.Println("Error accepting connection:", err)
			continue
		}
		jobs <- conn
	}

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("Server gracefully stopped.")
}

func worker(ctx context.Context, jobs chan net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	for conn := range jobs {
		handleConnection(ctx, conn)
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	// Goroutine to set read deadline when context is canceled
	go func() {
		<-ctx.Done()
		conn.SetReadDeadline(time.Now()) // Unblock reads
	}()

	reader:= bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Error reading request line:", err)
		return
	}

	fmt.Println("Received:", requestLine)
	
	method, uri, version, ok := parseRequestLine(requestLine)
	if !ok {
		log.Println("Invalid request line:", requestLine)
		return
	}

	headers, err := parseHeaders(reader)
	if err != nil {
		log.Println("Error parsing headers:", err)
		return
	}

	fmt.Printf("Method: %s\nURI: %s\nVersion: %s\nHeaders: %v\n", method, uri, version, headers)

	if uri == "/" {
		sendHttpResponse(conn, 200, "OK", "Hello, World!", headers)
	} else {
		sendHttpResponse(conn, 404, "Not Found", "404 Not Found", headers)
	}
}

func parseRequestLine(requestLine string) (method, uri, version string, ok bool) {
	parts:= strings.Split(strings.TrimSpace(requestLine), " ")
	if len(parts) != 3 {
		return "", "", "", false
	}
	return parts[0], parts[1], parts[2], true
}

func parseHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if line == "\r\n" {
			break
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}
		headers[parts[0]] = parts[1]
	}
	return headers, nil
}

func sendHttpResponse(conn net.Conn, statusCode int, statusText string,body string, headers map[string]string) {
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText)
	response += "Content-Type: text/html\r\n"
	response += fmt.Sprintf("Content-Length: %d\r\n", len(body))
	response += "\r\n"
	response += body
	conn.Write([]byte(response))
}