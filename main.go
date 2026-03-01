package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type VMInfo struct {
	IPAddress  string `json:"ip-address"`
	MacAddress string `json:"mac-address"`
	Hostname   string `json:"hostname"`
	ClientID   string `json:"client-id"`
	ExpiryTime int64  `json:"expiry-time"`
}

func getVMIP(vmName, bridgeName string) (string, error) {
	statusFile := fmt.Sprintf("/var/lib/libvirt/dnsmasq/%s.status", bridgeName)
	data, err := os.ReadFile(statusFile)
	if err != nil {
		return "", fmt.Errorf("failed to read status file: %v", err)
	}

	var vmInfos []VMInfo
	if err := json.Unmarshal(data, &vmInfos); err != nil {
		return "", fmt.Errorf("failed to parse status file: %v", err)
	}

	for _, info := range vmInfos {
		if info.Hostname == vmName {
			return info.IPAddress, nil
		}
	}

	return "", fmt.Errorf("VM not found: %s", vmName)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: kvm-ssh [options] <vm_name>\n\n")
	fmt.Fprintf(os.Stderr, "Options:\n")
	flag.PrintDefaults()
}

func main() {
	defaultUser := os.Getenv("USER")

	bridge := flag.String("bridge", "virbr0", "Bridge name")
	user := flag.String("user", defaultUser, "SSH user")
	local := flag.String("local", "", "Local forward ports (comma-separated)")
	remote := flag.String("remote", "", "Remote forward ports (comma-separated)")
	sshOpts := flag.String("ssh-opts", "", "Additional SSH options")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	vmName := flag.Arg(0)

	ip, err := getVMIP(vmName, *bridge)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var sshArgs []string

	if *local != "" {
		for _, port := range strings.Split(*local, ",") {
			port = strings.TrimSpace(port)
			sshArgs = append(sshArgs, "-L", fmt.Sprintf("localhost:%s:localhost:%s", port, port))
		}
	}

	if *remote != "" {
		for _, port := range strings.Split(*remote, ",") {
			port = strings.TrimSpace(port)
			sshArgs = append(sshArgs, "-R", fmt.Sprintf("localhost:%s:localhost:%s", port, port))
		}
	}

	if *sshOpts != "" {
		sshArgs = append(sshArgs, strings.Fields(*sshOpts)...)
	}

	sshArgs = append(sshArgs, fmt.Sprintf("%s@%s", *user, ip))

	fmt.Printf("Executing: ssh %s\n", strings.Join(sshArgs, " "))

	sshCmd := exec.Command("ssh", sshArgs...)
	sshCmd.Stdin = os.Stdin
	sshCmd.Stdout = os.Stdout
	sshCmd.Stderr = os.Stderr

	if err := sshCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
