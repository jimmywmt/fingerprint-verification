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
	socketPath = "/run/fingerprint.sock"
)

func main() {
	fingerprint, err := tools.GenerateFingerprint()
	if err != nil {
		panic(err)
	}
	fmt.Println("📌 Host Fingerprint:", fingerprint)

	os.Remove(socketPath)
	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	if err := os.Chmod(socketPath, 0600); err != nil {
		panic("failed to chmod socket: " + err.Error())
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

			nonce, err := hex.DecodeString(strings.TrimSpace(nonceHex))
			if err != nil {
				return
			}

			key := sha256.Sum256([]byte(GetSharedSecret()))
			block, err := aes.NewCipher(key[:])
			if err != nil {
				return
			}
			gcm, err := cipher.NewGCM(block)
			if err != nil {
				return
			}
			ciphertext := gcm.Seal(nil, nonce, []byte(fingerprint+"@@"+GetSharedSecret()), nil)
			encoded := hex.EncodeToString(ciphertext)
			c.Write([]byte("RESPONSE|" + encoded + "\n"))
		}(conn)
	}
}
