package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
)

func main() {

	// style := lipgloss.NewStyle().
    //     Foreground(lipgloss.Color("205")).
    //     Background(lipgloss.Color("15")).
    //     Padding(1).
    //     MarginRight(2).
    //     Border(lipgloss.RoundedBorder()).
    //     BorderForeground(lipgloss.Color("205"))
	// Prompt user for server details
	reader := bufio.NewReader(os.Stdin)


	fmt.Print("Please input the IP of your VPS: ")
	serverIP, _ := reader.ReadString('\n')
	serverIP = strings.TrimSpace(serverIP)
	domain := serverIP

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password or path to SSH key: ")
	passwordOrKeyPath, _ := reader.ReadString('\n')
	passwordOrKeyPath = strings.TrimSpace(passwordOrKeyPath)

	fmt.Print("Have you pointed a domain to the server (must if you want SSL) (y/n): ")
	domainExist, _ := reader.ReadString('\n')
	domainExist = strings.TrimSpace(domainExist)

	if domainExist == "y" {
		fmt.Print("Enter the domain name: ")
		domainName, _ := reader.ReadString('\n')
		domainName = strings.TrimSpace(domainName)
		domain = domainName
	}

	
	

	var authMethod ssh.AuthMethod
	if strings.HasPrefix(passwordOrKeyPath, "/") {
		// Use SSH key
		key, err := os.ReadFile(passwordOrKeyPath)
		if err != nil {
			fmt.Println("Failed to read SSH key:", err)
			return
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			fmt.Println("Failed to parse SSH key:", err)
			return
		}
		authMethod = ssh.PublicKeys(signer)
	} else {
		// Use password
		password := strings.TrimSpace(passwordOrKeyPath)
		authMethod = ssh.Password(password)
	}

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, 22), config)
	if err != nil {
		fmt.Println("Failed to dial:", err)
		return
	}
	defer client.Close()
	bar := progressbar.Default(7)

	// Execute commands sequentially
	commands := []string{
		"sudo apt-get install unzip -y",
		"mkdir /var/www && mkdir /var/www/pocketbase && cd /var/www/pocketbase && wget https://github.com/pocketbase/pocketbase/releases/download/v0.22.4/pocketbase_0.22.4_linux_amd64.zip && unzip pocketbase_0.22.4_linux_amd64.zip && rm pocketbase_0.22.4_linux_amd64.zip",
		"cd /lib/systemd/system && touch pocketbase.service",
		`echo '[Unit]
		Description = pocketbase
	
		[Service]
		Type           = simple
		User           = root
		Group          = root
		LimitNOFILE    = 4096
		Restart        = always
		RestartSec     = 5s
		StandardOutput = append:/var/www/pocketbase/errors.log
		StandardError  = append:/var/www/pocketbase/errors.log
		ExecStart      = /var/www/pocketbase/pocketbase serve --http ` + domain + `:80
		
		[Install]
		WantedBy = multi-user.target' > /lib/systemd/system/pocketbase.service`,
		"sudo systemctl unmask pocketbase.service",

		"sudo systemctl enable pocketbase",
		"sudo systemctl start pocketbase",
		
	}

	for _, cmd := range commands {
		// Start session for each command
		session, err := client.NewSession()
		if err != nil {
			fmt.Println("Failed to create session:", err)
			return
		}
		defer session.Close()

		// Set up pipes for session input/output
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
		session.Stdin = os.Stdin

		// Execute command
		if err := session.Run(cmd); err != nil {
			fmt.Printf("Failed to run command \"%s\": %v\n", cmd, err)
			return
		}
		fmt.Printf("Command \"%s\" executed successfully.\n", cmd)
		time.Sleep(1 * time.Second)
		bar.Add(1)

	}

	fmt.Println("All commands executed successfully.")

	// print made a directory for pocketbase
	fmt.Println("You are all set, go to http://" + domain + "/_. It's safe to close the terminal.")
	
}
