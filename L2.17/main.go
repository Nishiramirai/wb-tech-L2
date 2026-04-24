package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Реализовать простой telnet-клиент с возможностью соединяться к TCP-серверу и взаимодействовать с ним

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: telnet-go [--timeout=10s] <host> <port>\n")
		os.Exit(1)
	}

	address := net.JoinHostPort(args[0], args[1])

	conn, err := net.DialTimeout("tcp", address, *timeout)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to connect to %s: %v\n", address, err)
		os.Exit(1)
	}

	fmt.Printf("Connected to %s\n", address)

	// closeOnce гарантирует, что conn.Close() будет вызван ровно один раз,
	// независимо от того, кто инициировал завершение — сигнал или горутина.
	var closeOnce sync.Once
	closeConn := func() {
		closeOnce.Do(func() { _ = conn.Close() })
	}
	defer closeConn()

	var wg sync.WaitGroup

	// stdin → socket
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer closeConn() // EOF в stdin (Ctrl+D) → закрываем conn → будит вторую горутину
		_, _ = io.Copy(conn, os.Stdin)
	}()

	// socket → stdout
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer closeConn() // EOF от сервера → закрываем conn → будит первую горутину
		_, _ = io.Copy(os.Stdout, conn)
	}()

	// Ждём сигнал — он тоже вызывает closeConn
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigCh
		fmt.Printf("\nReceived signal: %v. Shutting down...\n", sig)
		closeConn()
	}()

	wg.Wait()
}
