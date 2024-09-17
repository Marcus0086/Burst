package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
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

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Received:", line)
		_, err := conn.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// Read timed out due to read deadline, context was canceled
			fmt.Println("Connection closed due to server shutdown")
			return
		}
		if err == io.EOF {
			// Client closed the connection
			return
		}
		fmt.Println("Error reading from connection:", err)
	}
}
