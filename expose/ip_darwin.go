package expose

import (
	"fmt"
	"os/exec"
	"strings"
)

// we find the interface by parsing the "route" cmd output
func defaultIf() (string, error) {
	cmd := exec.Command("route", "-n", "get", "default")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(out), "\n")[2:] {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		if fields[0] == "interface:" {
			return fields[1], nil
		}
	}

	return "", fmt.Errorf("default interface not found!")
}

// small wrapper around ifconfig that returns an interface's ip
func ifconfig(iface string) (string, error) {
	cmd := exec.Command("ifconfig", iface)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(out), "\n")[1:] {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if fields[0] == "inet" {
			return fields[1], nil
		}
	}

	return "", fmt.Errorf("no valid ipv4 found for iface %s", iface)
}

// returns the IP of the default interface
func localIP() (string, error) {
	iface, err := defaultIf()
	if err != nil {
		return "", err
	}

	return ifconfig(iface)
}
