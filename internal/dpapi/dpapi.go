package dpapi

import (
	"encoding/base64"
	"log"
	"syscall"
	"unsafe"
)

var (
	crypt32            = syscall.NewLazyDLL("crypt32.dll")
	procCryptProtect   = crypt32.NewProc("CryptProtectData")
	procCryptUnprotect = crypt32.NewProc("CryptUnprotectData")
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	procLocalFree      = kernel32.NewProc("LocalFree")
)

type dataBlob struct {
	cbData uint32
	pbData *byte
}

const (
	// CRYPTPROTECT_LOCAL_MACHINE encrypts for the local machine scope,
	// allowing any process on this machine (including LocalSystem services) to decrypt.
	CRYPTPROTECT_LOCAL_MACHINE = 0x00000001
)

// Encrypt encrypts plaintext using Windows DPAPI (local machine scope)
// and returns base64-encoded ciphertext with "DPAPI:" prefix.
// Using CRYPTPROTECT_LOCAL_MACHINE so that services running as LocalSystem
// or other accounts can decrypt the data.
func Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	log.Printf("[dpapi] Encrypt: encrypting data (length=%d, scope=LOCAL_MACHINE)", len(plaintext))

	input := []byte(plaintext)
	inBlob := dataBlob{
		cbData: uint32(len(input)),
		pbData: &input[0],
	}

	var outBlob dataBlob

	ret, _, err := procCryptProtect.Call(
		uintptr(unsafe.Pointer(&inBlob)),
		0,
		0,
		0,
		0,
		uintptr(CRYPTPROTECT_LOCAL_MACHINE),
		uintptr(unsafe.Pointer(&outBlob)),
	)
	if ret == 0 {
		log.Printf("[dpapi] Encrypt: CryptProtectData failed: %v", err)
		return "", err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outBlob.pbData)))

	encrypted := unsafe.Slice(outBlob.pbData, outBlob.cbData)
	result := "DPAPI:" + base64.StdEncoding.EncodeToString(encrypted)
	log.Printf("[dpapi] Encrypt: success (encrypted length=%d)", len(result))
	return result, nil
}

// Decrypt decrypts a DPAPI-encrypted ciphertext (with "DPAPI:" prefix).
// Returns the original plaintext string.
func Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	prefix := "DPAPI:"
	if len(ciphertext) <= len(prefix) {
		return ciphertext, nil
	}
	if ciphertext[:len(prefix)] != prefix {
		return ciphertext, nil
	}

	log.Printf("[dpapi] Decrypt: decrypting data (scope=LOCAL_MACHINE)")

	data, err := base64.StdEncoding.DecodeString(ciphertext[len(prefix):])
	if err != nil {
		log.Printf("[dpapi] Decrypt: base64 decode failed: %v", err)
		return "", err
	}

	inBlob := dataBlob{
		cbData: uint32(len(data)),
		pbData: &data[0],
	}

	var outBlob dataBlob

	ret, _, err := procCryptUnprotect.Call(
		uintptr(unsafe.Pointer(&inBlob)),
		0,
		0,
		0,
		0,
		uintptr(CRYPTPROTECT_LOCAL_MACHINE),
		uintptr(unsafe.Pointer(&outBlob)),
	)
	if ret == 0 {
		log.Printf("[dpapi] Decrypt: CryptUnprotectData failed: %v", err)
		return "", err
	}
	defer procLocalFree.Call(uintptr(unsafe.Pointer(outBlob.pbData)))

	decrypted := unsafe.Slice(outBlob.pbData, outBlob.cbData)
	log.Printf("[dpapi] Decrypt: success (decrypted length=%d)", len(decrypted))
	return string(decrypted), nil
}

// IsEncrypted checks if a string was encrypted by our DPAPI Encrypt function.
func IsEncrypted(s string) bool {
	return len(s) > 6 && s[:6] == "DPAPI:"
}
