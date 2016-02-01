package expose

import (
	"fmt"
	"os/exec"
	"strings"
)

// we parse the output of "ip route" to retrieve the default interface
func ip_route(res map[string]string) error {
	// parse the output of the ip route command to get the default gateway
	cmd := exec.Command("ip", "route")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		// ignore empty lines
		if len(fields) < 1 {
			continue
		}

		// parse output
		name := fields[0]
		route_data := make(map[string]string)
		var key string
		for _, field := range fields[1:] {
			if key == "" {
				key = field
				continue
			}
			route_data[key] = field
			key = ""
		}
		res[name] = route_data["dev"]
	}

	return nil
}

// parse the output of "route" to find the default interface
func route(res map[string]string) error {
	// parse the output of the route command to get the default gateway
	cmd := exec.Command("route")
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(out), "\n")[2:] {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		res[fields[0]] = fields[len(fields)-1]
	}

	return nil
}

// get the default interface from either "ip route" or "route"
func get_routes() (map[string]string, error) {
	res := make(map[string]string)

	_, err := exec.LookPath("ip")
	if err != nil {
		_, err = exec.LookPath("route")
		if err != nil {
			return res, fmt.Errorf("nor ip nor route are available")
		}

		// run the route command
		err = route(res)
		return res, err
	}

	// otherwise run ip route
	err = ip_route(res)
	return res, err
}

func default_interface() (string, error) {
	routes, err := get_routes()
	if err != nil {
		return "", nil
	}

	if iface, ok := routes["default"]; ok {
		return iface, nil
	}

	return "", fmt.Errorf("default interface not found!")
}

func ip_addr(iface string) (string, error) {
	cmd := exec.Command("ip", "addr", "show", iface)
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
			return strings.Split(fields[1], "/")[0], nil
		}
	}

	return "", fmt.Errorf("no valid ipv4 found for iface %s", iface)
}

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

func local_ip() (string, error) {
	iface, err := default_interface()
	if err != nil {
		return "", err
	}

	_, err = exec.LookPath("ip")
	if err != nil {
		_, err = exec.LookPath("ifconfig")
		if err != nil {
			return "", fmt.Errorf("nor ip nor ifconfig are available")
		}

		return ifconfig(iface)
	}

	return ip_addr(iface)
}
