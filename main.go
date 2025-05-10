package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/jimmywmt/fingerprint-verification/tools"
)

const (
	socketPath   = "/run/fingerprint.sock"
	sharedSecret = "{{SHARED_SECRET}}"
)

func main() {
	fingerprint, err := tools.GenerateFingerprint()
	if err != nil {
		panic(err)
	}
	fmt.Println("📌 Host Fingerprint:", fingerprint)
	fmt.Println("🔑 Shared Secret:", sharedSecret)

	os.Remove(socketPath)
	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go func(c net.Conn) {
			defer c.Close()
			reader := bufio.NewReader(c)
			nonceHex, err := reader.ReadString('\n')
			if err != nil {
				return
			}

			fmt.Println("Received nonceHex from container:", nonceHex)

			nonce, err := hex.DecodeString(strings.TrimSpace(nonceHex))
			if err != nil {
				return
			}

			key := sha256.Sum256([]byte(sharedSecret))
			block, err := aes.NewCipher(key[:])
			if err != nil {
				return
			}
			gcm, err := cipher.NewGCM(block)
			if err != nil {
				return
			}
			ciphertext := gcm.Seal(nil, nonce, []byte(fingerprint), nil)
			encoded := hex.EncodeToString(ciphertext)
			c.Write([]byte("RESPONSE|" + encoded + "\n"))
		}(conn)
	}
}
