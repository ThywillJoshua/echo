// server.go
package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

func StartHTTPServer(message *string) {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}
	execDir := filepath.Dir(execPath)
	distPath := filepath.Join(execDir, "dist", "ui", "browser")

	fs := http.FileServer(http.Dir(distPath))
	http.Handle("/", fs)

	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		messageHandler(w, r, message)
	})
	http.HandleFunc("/commit", commitHandler)
	http.HandleFunc("/fcommit", fakeStreamHandler)
	http.HandleFunc("/close", closeHandler)
	fmt.Println("Starting server on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func messageHandler(w http.ResponseWriter, r *http.Request, message *string) {
	if err := json.NewEncoder(w).Encode(*message); err != nil {
		http.Error(w, "Error generating diff message", http.StatusInternalServerError)
	}
}

func closeHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func fakeStreamHandler(w http.ResponseWriter, r *http.Request) {
	// Set the response header to indicate chunked transfer encoding.
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Get the http.Flusher interface to enable streaming.
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported by the server", http.StatusInternalServerError)
		return
	}

	// Simulate streaming data in chunks with a delay between each chunk.
	for i := 1; i <= 5; i++ {
		// Write a chunk of data to the response.
		fmt.Fprintf(w, "Chunk %d: The quick brown fox jumps over the lazy dog.\n", i)

		// Flush the buffer to send the chunk to the client immediately.
		flusher.Flush()

		// Simulate a delay between chunks to mimic streaming.
		time.Sleep(1 * time.Second)
	}

	// Indicate the end of the response stream.
	fmt.Fprintln(w, "End of streaming data.")
}

func commitHandler(w http.ResponseWriter, r *http.Request) {
	type CommitRequest struct {
		Message string `json:"message"`
	}

	var commitReq CommitRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &commitReq); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	message := commitReq.Message

	// Set headers for streaming
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")

	// Create a flusher to handle streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Start the command
	cmd := exec.Command("git", "commit", "-m", message)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obtaining stdout pipe: %v", err), http.StatusInternalServerError)
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error obtaining stderr pipe: %v", err), http.StatusInternalServerError)
		return
	}

	if err := cmd.Start(); err != nil {
		http.Error(w, fmt.Sprintf("Error starting command: %v", err), http.StatusInternalServerError)
		return
	}

	// Use a WaitGroup to wait for both goroutines to complete
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine for reading stdout
	go func() {
		defer wg.Done()
		defer stdoutPipe.Close()
		stdoutScanner := bufio.NewScanner(stdoutPipe)
		for stdoutScanner.Scan() {
			fmt.Fprintf(w, "STDOUT: %s\n", stdoutScanner.Text())
			flusher.Flush()
		}
		if err := stdoutScanner.Err(); err != nil {
			fmt.Fprintf(w, "Error reading stdout: %v\n", err)
			flusher.Flush()
		}
	}()

	// Goroutine for reading stderr
	go func() {
		defer wg.Done()
		defer stderrPipe.Close()
		stderrScanner := bufio.NewScanner(stderrPipe)
		for stderrScanner.Scan() {
			fmt.Fprintf(w, "STDERR: %s\n", stderrScanner.Text())
			flusher.Flush()
		}
		if err := stderrScanner.Err(); err != nil {
			fmt.Fprintf(w, "Error reading stderr: %v\n", err)
			flusher.Flush()
		}
	}()

	// Wait for both stdout and stderr to be fully read
	wg.Wait()

	// Wait for the command to exit
	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(w, "Command execution error: %v\n", err)
		flusher.Flush()
		return
	}

	// Command succeeded
	fmt.Fprintln(w, "Command executed successfully")
	flusher.Flush()
}
