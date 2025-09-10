package server

import (
	"io"
	"log"
	"net"
	"redisgo/internal/config"
	"redisgo/internal/core"
	"redisgo/internal/core/io_multiplexing"
	"syscall"
	"time"
)

func readCommand(fd int) (*core.Command, error) {
	buf := make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}
	return core.ParseCmd(buf)
}

func RunIOMultiplexingServer() {
	log.Println("Starting an I/O Multiplexing TCP server on", config.Port)
	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Get the file descriptor from the listener
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		log.Fatal("Listener is not a TCPListener")
	}
	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal(err)
	}
	defer listenerFile.Close()

	serverFD := int(listenerFile.Fd())

	// Create an I/O Multiplexing instance (epoll in Linux, kqueue in MacOS)
	ioMultiplexer, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		log.Fatal(err)
	}
	defer ioMultiplexer.Close()

	// Monitor "read" events on the Server FD
	if err = ioMultiplexer.Monitor(io_multiplexing.Event{
		Fd: serverFD,
		Op: io_multiplexing.OpRead,
	}); err != nil {
		log.Fatal(err)
	}

	lastActiveExpireExecTime := time.Now()
	for {
		if time.Now().After(lastActiveExpireExecTime.Add(core.ACTIVE_EXPIRE_FREQUENCY)) {
			core.ActiveDeleteExpiredKeys()
			lastActiveExpireExecTime = time.Now()
		}

		// Wait for file descriptors in the monitoring list to be ready for I/O
		events, err := ioMultiplexer.Wait() // Blocking call
		if err != nil {
			continue
		}
		for _, event := range events {
			if event.Fd == serverFD {
				log.Print("New client is trying to connect")
				connFD, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("Err: ", err)
					continue
				}

				log.Println("Setting up a new connection", connFD)
				if err = ioMultiplexer.Monitor(io_multiplexing.Event{
					Fd: connFD,
					Op: io_multiplexing.OpRead,
				}); err != nil {
					log.Fatal("Err: ", err)
				}
			} else {
				cmd, err := readCommand(event.Fd)
				if err != nil {
					if err == io.EOF || err == syscall.ECONNRESET {
						log.Println("Client disconnected", event.Fd)
						_ = syscall.Close(event.Fd)
						continue
					}
					log.Println("Read command error: ", err)
					continue
				}
				if err = core.ExecuteAndResponse(cmd, event.Fd); err != nil {
					log.Println("Error execute and response:", err)
				}
			}
		}
	}
}
