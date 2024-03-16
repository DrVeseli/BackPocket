package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {
	// Prompt user for server details
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter server IP address: ")
	serverIP, _ := reader.ReadString('\n')
	serverIP = strings.TrimSpace(serverIP)

	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// SSH client configuration
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
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

	// Execute commands sequentially
	commands := []string{
		//"sudo apt-get install nginx -y",
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
		ExecStart      = /var/www/pocketbase/pocketbase serve --http ` + serverIP + `:80
		
		[Install]
		WantedBy = multi-user.target' > /lib/systemd/system/pocketbase.service`,
		"sudo systemctl unmask pocketbase.service",

		"sudo systemctl enable pocketbase",
		"sudo systemctl start pocketbase",
		//"cd /etc/nginx/sites-available && touch pocketbase",
		// `echo 'server {
		// 	      listen 80;
		// 	      server_name ` + serverIP + `;
		// 	      client_max_body_size 10M;

		// 	      location / {
		// 	              # check http://nginx.org/en/docs/http/ngx_http_upstream_module.html#keepalive
		// 	              proxy_set_header Connection '';
		// 	              proxy_http_version 1.1;
		// 	              proxy_read_timeout 360s;

		// 	              proxy_set_header Host \$host;
		// 	              proxy_set_header X-Real-IP \$remote_addr;
		// 	             proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
		// 	              proxy_set_header X-Forwarded-Proto \$scheme;

		// 	              # enable if you are serving under a subpath location
		// 	              # rewrite /yourSubpath/(.*) /\$1  break;

		// 	              proxy_pass http://127.0.0.1:8090;
		// 	      }
		// }' >> pocketbase.service`,
		// "sudo ln -s /etc/nginx/sites-available/pocketbase /etc/nginx/sites-enabled",
		// "sudo systemctl restart nginx",
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
		time.Sleep(5 * time.Second)
	}

	fmt.Println("All commands executed successfully.")

	// print made a directory for pocketbase
	fmt.Println("A directory for pocketbase has been successfully created on the remote server.")

	// // download pocketbase from github releases and unzip
	// _ = session.Run("wget https://github.com/pocketbase/pocketbase/releases/download/v0.22.4/pocketbase_0.22.4_linux_amd64.zip")
	// _ = session.Run("unzip pocketbase_0.22.4_linux_amd64.zip")
	// //prin
	// fmt.Println("Pocketbase has been successfully downloaded and unzipped on the remote server.")
	// //remove the zip file
	// _ = session.Run("rm pocketbase_0.22.4_linux_amd64.zip")

	// // make pocketbase executable
	// _ = session.Run("chmod +x pocketbase")

	// // cd into /lib/systemd/system and create a pocketbase.service file
	// _ = session.Run("cd /lib/systemd/system")
	// _ = session.Run("touch pocketbase.service")

	// // write the pocketbase.service file
	// _ = session.Run(`echo '[Unit]
	// 	Description = pocketbase

	// [Service]
	// 	Type           = simple
	// 	User           = root
	// 	Group          = root
	// 	LimitNOFILE    = 4096
	// 	Restart        = always
	// 	RestartSec     = 5s
	// 	StandardOutput = append:/root/pocketbase/errors.log
	// 	StandardError  = append:/root/pocketbase/errors.log
	// 	ExecStart      = /root/pocketbase/pocketbase serve ` + serverIP + `

	// [Install]
	// 	WantedBy = multi-user.target' >> pocketbase.service"`)

	// // enable the pocketbase service
	// _ = session.Run("systemctl enable pocketbase")

	// // start the pocketbase service
	// _ = session.Run("systemctl start pocketbase")

	// // check the status of the pocketbase service
	// _ = session.Run("systemctl status pocketbase")

	// // Print success message
	// fmt.Println("Pocketbase has been successfully installed and started on the remote server.")

	// // cd into nginx sites-available and create a pocketbase file
	// _ = session.Run("cd /etc/nginx/sites-available")
	// _ = session.Run("touch pocketbase")

	// // write the pocketbase file
	// _ = session.Run(`echo 'server {
	// 	listen 80;
	// 	server_name example.com;
	// 	client_max_body_size 10M;

	// 	location / {
	// 		# check http://nginx.org/en/docs/http/ngx_http_upstream_module.html#keepalive
	// 		proxy_set_header Connection '';
	// 		proxy_http_version 1.1;
	// 		proxy_read_timeout 360s;

	// 		proxy_set_header Host $host;
	// 		proxy_set_header X-Real-IP $remote_addr;
	// 		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	// 		proxy_set_header X-Forwarded-Proto $scheme;

	// 		# enable if you are serving under a subpath location
	// 		# rewrite /yourSubpath/(.*) /$1  break;

	// 		proxy_pass http://127.0.0.1:8090;
	// 	}
	// }' >> pocketbase.service"`)

	// // create a symbolic link to the pocketbase file
	// _ = session.Run("ln -s /etc/nginx/sites-available/pocketbase /etc/nginx/sites-enabled")

	// // restart nginx
	// _ = session.Run("systemctl restart nginx")

	// // Print success message
	// fmt.Println("Pocketbase has been successfully configured with Nginx on the remote server.")

	// // Print success message
	// fmt.Println("Pocketbase has been successfully installed and started on the remote server.")

	// // Prompt user to press enter to exit
	// fmt.Println("Press enter to exit.")
	// _, _ = reader.ReadString('\n')

}
