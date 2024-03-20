package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {
	proxyUser := flag.String("proxy-user", "", "The remote (proxy) SSH username")
	proxyHost := flag.String("proxy-host", "", "The remote (proxy) SSH host")
	proxySshPort := flag.String("proxy-ssh-port", "22", "The remote (proxy) SSH port")
	proxyPort := flag.String("proxy-port", "11111", "The remote port to forward")
	localPort := flag.String("local-port", "21112", "The local (machine) port to be forwarded to from proxyPort")
	flag.Parse()

	keyPath, err := generateSSHKey()
	if err != nil {
		log.Fatalf("Failed to generate or load SSH key: %v", err)
	}

	if *proxyHost == "" {
		log.Fatal("proxy-host is required")
	}

	client, err := connectToSSH(*proxyUser, *proxyHost, *proxySshPort, keyPath)
	if err != nil {
		log.Fatalf("Failed to connect over SSH: %v", err)
	}
	defer client.Close()

	err = forwardPort(client, *proxyHost, *proxyPort, *localPort)
	if err != nil {
		log.Fatalf("Failed to forward port: %v", err)
	}

	select {} // Prevents the application from exiting immediately
}

func generateSSHKey() (string, error) {
	keyPath := "id_ed25519"
	pubKeyPath := keyPath + ".pub"

	// Check if the key already exists
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		// Generate a new ED25519 private/public key pair
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return "", err
		}

		// Marshal the private key into a PEM block
		privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return "", err
		}
		privPEM := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privKeyBytes,
		}

		// Write the PEM to file
		file, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return "", err
		}
		defer file.Close()
		if err := pem.Encode(file, privPEM); err != nil {
			return "", err
		}

		// Convert the ed25519 public key into the format for the authorized_keys file
		pubKeySSH, err := ssh.NewPublicKey(publicKey)
		if err != nil {
			return "", err
		}
		pubKeyBytes := ssh.MarshalAuthorizedKey(pubKeySSH)

		// Write the public key to a file
		if err := ioutil.WriteFile(pubKeyPath, pubKeyBytes, 0644); err != nil {
			return "", err
		}
	}

	return keyPath, nil
}

func connectToSSH(user, host, port, keyPath string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Note: In production, replace with a more secure option
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func forwardPort(client *ssh.Client, proxyHost, proxyPort, localPort string) error {
	// Listen on remote server port
	listener, err := client.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", proxyPort))
	if err != nil {
		return err
	}
	defer listener.Close()

	// Handle incoming connections on remote port
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection: %s", err)
			continue
		}

		// Handle the connection in a new goroutine
		go func() {
			defer conn.Close()

			// Connect to the target (local Windows machine)
			localConn, err := net.Dial("tcp", fmt.Sprintf("0.0.0.0:%s", localPort))
			if err != nil {
				log.Printf("Failed to connect to local port: %s", err)
				return
			}
			defer localConn.Close()
			// Copy data from the remote connection to the local connection
			go func() {
				_, err := io.Copy(localConn, conn)
				if err != nil {
					log.Printf("Error copying data from remote to local: %v", err)
				}
			}()

			// Copy data from the local connection to the remote connection
			_, err = io.Copy(conn, localConn)
			if err != nil {
				log.Printf("Error copying data from local to remote: %v", err)
			}
		}()
	}
	return nil
}
