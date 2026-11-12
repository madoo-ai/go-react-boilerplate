package security

import (
	"bytes"
	"testing"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func TestHasherRoundTrip(t *testing.T) {
	hasher := NewHasher()
	hash, err := hasher.Hash("correct horse battery staple")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	if !hasher.Verify(hash, "correct horse battery staple") {
		t.Fatal("expected the password to verify")
	}
	if hasher.Verify(hash, "wrong password") {
		t.Fatal("expected the wrong password to fail")
	}
	if hasher.Verify("not-a-phc-string", "anything") {
		t.Fatal("expected a malformed hash to fail")
	}
}

func TestTokens(t *testing.T) {
	tokens := NewTokens()
	first, err := tokens.NewToken()
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	second, err := tokens.NewToken()
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	if first == second {
		t.Fatal("tokens must be unique")
	}
	firstHash := tokens.HashToken(first)
	secondHash := tokens.HashToken(first)
	if firstHash != secondHash {
		t.Fatal("token hashing must be deterministic")
	}
	if firstHash == first {
		t.Fatal("the stored hash must not be the raw token")
	}
}

func TestBackupCodes(t *testing.T) {
	tokens := NewTokens()
	code, err := tokens.NewBackupCode()
	if err != nil {
		t.Fatalf("backup code: %v", err)
	}
	if len(code) != 10 {
		t.Fatalf("backup code length = %d, want 10", len(code))
	}
	want := tokens.HashBackupCode(code)
	if got := tokens.HashBackupCode(" " + code[:5] + "-" + code[5:] + " "); got != want {
		t.Fatal("backup code normalization should ignore spaces, case and dashes")
	}
}

func TestCipherRoundTrip(t *testing.T) {
	cipher, err := NewCipher([]byte("0123456789abcdef0123456789abcdef"))
	if err != nil {
		t.Fatalf("cipher: %v", err)
	}
	encrypted, err := cipher.Encrypt("JBSWY3DPEHPK3PXP")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if encrypted == "JBSWY3DPEHPK3PXP" {
		t.Fatal("the secret must not be stored in plain text")
	}
	plain, err := cipher.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if plain != "JBSWY3DPEHPK3PXP" {
		t.Fatalf("round trip = %q", plain)
	}
	other, _ := NewCipher([]byte("abcdef0123456789abcdef0123456789"))
	if _, err := other.Decrypt(encrypted); err == nil {
		t.Fatal("a different key must not decrypt the payload")
	}
}

func TestTOTP(t *testing.T) {
	service := NewTOTP()
	secret, uri, err := service.Generate("Test App", "person@example.com")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if secret == "" || uri == "" {
		t.Fatal("generate returned empty values")
	}

	fixed := time.Date(2026, 11, 12, 9, 0, 0, 0, time.UTC)
	previous := timeNow
	timeNow = func() time.Time { return fixed }
	defer func() { timeNow = previous }()

	code, err := totp.GenerateCodeCustom(secret, fixed, totp.ValidateOpts{
		Period: 30, Skew: 1, Digits: otp.DigitsSix, Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		t.Fatalf("code: %v", err)
	}
	if !service.Validate(secret, code) {
		t.Fatal("expected the generated code to validate")
	}
	if service.Validate(secret, "000000") {
		t.Fatal("expected a wrong code to fail")
	}

	png, err := service.QRCode(uri)
	if err != nil {
		t.Fatalf("qr: %v", err)
	}
	if !bytes.HasPrefix(png, []byte{0x89, 'P', 'N', 'G'}) {
		t.Fatal("expected a PNG payload")
	}
}
