package network

import (
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var pingRegex = regexp.MustCompile(`time=([\d\.]+)\s*ms`)

func PingTarget(target string) (string, error) {
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		start := time.Now()
		client := http.Client{Timeout: 3 * time.Second}
		resp, err := client.Get(target)
		if err != nil || resp.StatusCode >= 500 {
			return "", fmt.Errorf("DOWN")
		}
		duration := time.Since(start)
		return fmt.Sprintf("%d ms", duration.Milliseconds()), nil
	}

	if strings.Contains(target, ":") {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", target, 3*time.Second)
		if err != nil {
			return "", fmt.Errorf("DOWN")
		}
		conn.Close()
		duration := time.Since(start)
		return fmt.Sprintf("%d ms", duration.Milliseconds()), nil
	}

	out, err := exec.Command("ping", "-c", "1", "-W", "1", target).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("DOWN")
	}

	matches := pingRegex.FindStringSubmatch(string(out))
	if len(matches) >= 2 {
		valFloat, _ := strconv.ParseFloat(matches[1], 64)
		return fmt.Sprintf("%d ms", int(valFloat)), nil
	}
	
	return "", fmt.Errorf("DOWN")
}
