package webhook

import "testing"

func TestSignAndVerify(t *testing.T) {
	payload := []byte(`{"event":"email.delivered","email_id":"123"}`)
	secret := "whsec_test_secret_key"

	signature := Sign(payload, secret)
	if signature == "" {
		t.Fatal("Sign() returned empty signature")
	}

	if !Verify(payload, secret, signature) {
		t.Error("Verify() returned false for valid signature")
	}
}

func TestVerify_WrongSignature(t *testing.T) {
	payload := []byte(`{"event":"email.delivered"}`)
	secret := "whsec_test_secret_key"

	if Verify(payload, secret, "0000000000000000000000000000000000000000000000000000000000000000") {
		t.Error("Verify() should return false for wrong signature")
	}
}

func TestVerify_WrongSecret(t *testing.T) {
	payload := []byte(`{"event":"email.delivered"}`)
	secret1 := "whsec_correct_secret"
	secret2 := "whsec_wrong_secret"

	signature := Sign(payload, secret1)

	if Verify(payload, secret2, signature) {
		t.Error("Verify() should return false when verified with wrong secret")
	}
}

func TestSign_DifferentPayloadsProduceDifferentSignatures(t *testing.T) {
	secret := "whsec_test_secret"
	sig1 := Sign([]byte(`{"a":1}`), secret)
	sig2 := Sign([]byte(`{"a":2}`), secret)

	if sig1 == sig2 {
		t.Error("Sign() should produce different signatures for different payloads")
	}
}

func TestSign_DifferentSecretsProduceDifferentSignatures(t *testing.T) {
	payload := []byte(`{"event":"test"}`)
	sig1 := Sign(payload, "secret-one")
	sig2 := Sign(payload, "secret-two")

	if sig1 == sig2 {
		t.Error("Sign() should produce different signatures for different secrets")
	}
}
