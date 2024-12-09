package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/ADEMOLA200/ECommerce/cmd/models"
	"github.com/ADEMOLA200/ECommerce/cmd/router"
)

type ServerInterface interface {
	CheckAndKillProcess(port int)
}

type Server struct {
	App    *models.Application
	router router.Router
}

func (s *Server) Server() error {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.App.Config.Port),
		Handler:           s.router.Handler(),
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	s.App.InfoLog.Printf("Starting http server in %s mode on port %d", s.App.Config.Env, s.App.Config.Port)

	return server.ListenAndServe()
}

// checkAndKillProcess ensures the given port is not already in use by another process.
func (s *Server) CheckAndKillProcess(port int) {
	address := fmt.Sprintf(":%d", port)
	conn, err := net.Listen("tcp", address)
	if err == nil {
		_ = conn.Close()
		return
	}

	fmt.Printf("Port %d is in use. Attempting to free it...\n", port)

	cmd := exec.Command("lsof", "-t", fmt.Sprintf("-i:%d", port))
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Failed to check port %d: %v", port, err)
	}

	pids := string(output)
	if pids == "" {
		log.Fatalf("No processes found using port %d, but it appears to still be blocked. Check system activity.", port)
	}

	for _, pid := range strings.Fields(pids) {
		killCmd := exec.Command("kill", "-9", pid)
		err = killCmd.Run()
		if err != nil {
			log.Printf("Failed to kill process %s using port %d: %v", pid, port, err)
		} else {
			fmt.Printf("Successfully killed process %s using port %d.\n", pid, port)
		}
	}

	conn, err = net.Listen("tcp", address)
	if err == nil {
		_ = conn.Close()
		fmt.Printf("Port %d is now free and ready to use.\n", port)
	} else {
		log.Fatalf("Port %d is still in use after cleanup attempts. Please investigate manually.", port)
	}
}
