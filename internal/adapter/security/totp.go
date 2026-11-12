package security

import (
	"bytes"
	"image/png"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// TOTP implements the time-based one-time password port.
type TOTP struct{}

func NewTOTP() *TOTP { return &TOTP{} }

func (t *TOTP) Generate(issuer, account string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer: issuer, AccountName: account, Period: 30, SecretSize: 20,
		Digits: otp.DigitsSix, Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", "", err
	}
	return key.Secret(), key.URL(), nil
}

func (t *TOTP) Validate(secret, code string) bool {
	valid, err := totp.ValidateCustom(code, secret, timeNow(), totp.ValidateOpts{
		Period: 30, Skew: 1, Digits: otp.DigitsSix, Algorithm: otp.AlgorithmSHA1,
	})
	return err == nil && valid
}

func (t *TOTP) QRCode(uri string) ([]byte, error) {
	code, err := qrcode.New(uri, qrcode.Medium)
	if err != nil {
		return nil, err
	}
	code.DisableBorder = false
	var buffer bytes.Buffer
	if err := png.Encode(&buffer, code.Image(256)); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
