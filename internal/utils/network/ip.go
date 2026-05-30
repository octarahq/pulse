package network

import (
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	psnet "github.com/shirou/gopsutil/v4/net"
)

func GetNetIpLocal(interName string) string {
	if interName == "" {
		return ""
	}

	ifaces, err := psnet.Interfaces()
	if err != nil {
		return ""
	}

	ok := false
	for _, iface := range ifaces {
		if iface.Name == interName {
			ok = true
			break
		}
	}
	if !ok {
		return ""
	}

	iface, err := net.InterfaceByName(interName)
	if err != nil {
		return ""
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		switch v := addr.(type) {
		case *net.IPNet:
			if ip := v.IP.To4(); ip != nil && !ip.IsLoopback() {
				return ip.String()
			}
		case *net.IPAddr:
			if ip := v.IP.To4(); ip != nil && !ip.IsLoopback() {
				return ip.String()
			}
		}
	}

	return ""
}

func GetNetIpPublic() string {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(body))
}

func GetNetIpGateway() string {
	data, err := os.ReadFile("/proc/net/route")
	if err != nil {
		return ""
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[1] != "00000000" {
			continue
		}

		gw, err := strconv.ParseUint(fields[2], 16, 32)
		if err != nil {
			continue
		}

		return net.IPv4(byte(gw), byte(gw>>8), byte(gw>>16), byte(gw>>24)).String()
	}

	return ""
}
