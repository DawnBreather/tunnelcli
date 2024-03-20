# TunnelCLI
## What
TunnelCLI is a command-line tool designed to facilitate secure SSH tunneling. It allows users to easily set up port forwarding from a local machine through a remote SSH server to another local port. This tool is particularly useful for securely exposing local services to the internet or accessing services on a private network.

## Why
1. Don't have a public IP address on the server where you are hosting your application, yet want to make it accessible from the public internet?
2. If your transport protocol is TCP or above (i.e. HTTP), and you have access to SSH server with a public IP address with SSH tunnelling enabled - use this tool.

## Features
* Easy SSH Key Generation: Automatically generates an SSH key pair for authentication if one doesn't already exist.
* Simple Port Forwarding Setup: Streamlines the process of setting up SSH port forwarding.
* Flexible Configuration: Supports custom SSH usernames, hosts, and ports for versatile network configurations.
* Cross-platform: Yeah, it's Go
* Self-sufficient: No dependencies, just this compiled binary is all you need to make things work

## Requirements
* Go 1.15 or higher.
* Access to a remote SSH server where you have permission to authenticate and create SSH tunnels.
* The remote SSH server must have GatewayPorts enabled in the sshd_config file to allow binding to non-localhost addresses.

## Installation

To get started with TunnelCli, clone this repository and build the application:
```sh
git clone https://github.com/DawnBreather/tunnelcli.git
cd tunnelcli
go build
```

## Usage

To use TunnelCLI, you'll need to provide the remote SSH username, host, SSH port, the remote port to forward, and the local port to be forwarded.
```sh
./tunnelcli --proxy-user <SSH Username> --proxy-host <SSH Host> --proxy-ssh-port <SSH Port> --proxy-port <Remote Port to Forward> --local-port <Local Port>
```

### Examples

Forwarding local port 8080 to remote port 2112 through an SSH server:
```sh
./tunnelcli --proxy-user ec2-user --proxy-host 177.17.68.8 --proxy-ssh-port 22 --proxy-port 2112 --local-port 8080
```

## Key Generation
Upon the first run, if an SSH key does not exist, TunnelCLI will generate a new ED25519 SSH key pair and store it locally. The public key should be added to the authorized_keys file on the SSH server to allow for password-less authentication.
Contributing

Contributions to TunnelCLI are welcome! Feel free to fork the repository, make your changes, and submit a pull request.

## License
TunnelCli is dedicated to the public domain. The author has waived all rights to the work worldwide under copyright law, including all related and neighboring rights, to the extent allowed by law.
You can copy, modify, distribute, and perform the work, even for commercial purposes, all without asking permission. For more information, see [Creative Commons Zero v1.0 Universal](https://creativecommons.org/publicdomain/zero/1.0/deed.en).