package main

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func testhash() {
	h_md5 := md5.New()
	h_sha := sha256.New()

	io.WriteString(h_md5, "Welcome to Go Language Secure Coding Practices")
	io.WriteString(h_sha, "Welcome to Go Language Secure Coding Practices")
	fmt.Printf("MD5 : %x\n", h_md5.Sum(nil))
	fmt.Printf("SHA256: %x\n", h_sha.Sum(nil))
}

func testEncrypt()  {
	key := []byte("Encryption key should be 32 bit ")
	data := []byte("Welcome to Go language secure coding practices")

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	encrypted_data := aesgcm.Seal(nil, nonce, data, nil)
	fmt.Printf("Encrypted: %x\n", encrypted_data)
	decrypted_data, err := aesgcm.Open(nil, nonce, encrypted_data, nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Decrypted: %s\n", decrypted_data)

}

func main() {

	testhash()
	testEncrypt()


}
