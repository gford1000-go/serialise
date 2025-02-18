package serialise

import (
	"bytes"
	"testing"
)

func TestWithAESGCMEncryption(t *testing.T) {

	createOpt := func(key []byte) *SerialisationOptions {
		applyEncryption := WithAESGCMEncryption(key)

		opt := SerialisationOptions{}
		applyEncryption(&opt)

		return &opt
	}

	var key = []byte("01234567890123456789012345678912")

	opt := createOpt(key)

	tests := []string{
		"",
		"A",
		"This is a test string",
		"的一是在不了有和人这",
		"ᨕᨗᨒᨁᨒᨗᨁᨚ.",
	}

	for _, text := range tests {
		data := []byte(text)

		ciphertext, err := opt.Encryptor(data)
		if err != nil {
			t.Fatalf("Unexpected error during encryption of '%s': %v", text, err)
		}

		plaintext, err := opt.Decryptor(ciphertext)
		if err != nil {
			t.Fatalf("Unexpected error during decryption of '%s': %v", text, err)
		}

		if !bytes.Equal(data, plaintext) {
			t.Fatalf("Mismatch between original and decrypted data for '%s'", text)
		}
	}
}

func TestWithAESGCMEncryption_1(t *testing.T) {

	createOpt := func(key []byte) *SerialisationOptions {
		applyEncryption := WithAESGCMEncryption(key)

		opt := SerialisationOptions{}
		applyEncryption(&opt)

		return &opt
	}

	var encKey = []byte("01234567890123456789012345678912")
	var badKey = []byte("01234567890123456789012345678911")

	optEnc := createOpt(encKey)
	optDec := createOpt(badKey)

	var data = []byte("This is a test string")

	ciphertext, err := optEnc.Encryptor(data)
	if err != nil {
		t.Fatalf("Unexpected error during encryption: %v", err)
	}

	plaintext, err := optDec.Decryptor(ciphertext)
	if plaintext != nil {
		t.Fatal("Unexpected return of plaintext when expected nil")
	}

	if err == nil {
		t.Fatal("Unexpected success when expected error")
	}
}

func TestWithAESGCMEncryption_2(t *testing.T) {

	createOpt := func(key []byte) *SerialisationOptions {
		applyEncryption := WithAESGCMEncryption(key)

		opt := SerialisationOptions{}
		applyEncryption(&opt)

		return &opt
	}

	var key = []byte("01234567890123456789012345678912")

	opt := createOpt(key)

	var data = []byte("This is a test string")

	ciphertext, err := opt.Encryptor(data)
	if err != nil {
		t.Fatalf("Unexpected error during encryption: %v", err)
	}

	ciphertext[len(ciphertext)-1] = ciphertext[0]

	plaintext, err := opt.Decryptor(ciphertext)
	if plaintext != nil {
		t.Fatal("Unexpected return of plaintext when expected nil")
	}

	if err == nil {
		t.Fatal("Unexpected success when expected error")
	}
}
