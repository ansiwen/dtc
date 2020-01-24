package utils

import "C"
import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/asn1"
	"fmt"
	"math/big"
)

// from github.com/Thalesignite/crypto11
var CurveNameToASN1 = map[string]asn1.ObjectIdentifier{
	"P-224": {1, 3, 132, 0, 33},
	"P-256": {1, 2, 840, 10045, 3, 1, 7},
	"P-384": {1, 3, 132, 0, 34},
	"P-521": {1, 3, 132, 0, 35},
}

func CurveNameToASN1Bytes(curve string) ([]byte, error) {
	obj, ok := CurveNameToASN1[curve]
	if !ok {
		return nil, fmt.Errorf("curve unsupported")
	}
	return asn1.Marshal(obj)
}

func ASN1ToCurveName(b []byte) (string, error) {
	var v asn1.ObjectIdentifier
	extra, err := asn1.Unmarshal(b, &v)
	if len(extra) > 0 {
		return "", fmt.Errorf("extra data in params")
	}
	if err != nil {
		return "", fmt.Errorf("error decrypting params: %s", err)
	}
	for name, item := range CurveNameToASN1 {
		if v.Equal(item) {
			return name, nil
		}
	}
	return "", fmt.Errorf("curve unsupported")
}

func PubKeyToASN1Bytes(pk *ecdsa.PublicKey) ([]byte, error) {
	numSize := (pk.Curve.Params().BitSize + 7) / 8
	ecPointBytes := make([]byte, 2*numSize)
	copy(ecPointBytes[numSize-len(pk.X.Bytes()):numSize], pk.X.Bytes())
	copy(ecPointBytes[(2*numSize)-len(pk.Y.Bytes()):(2*numSize)], pk.Y.Bytes())
	ecPointASN1, err := asn1.Marshal(ecPointBytes)
	if err != nil {
		return nil, err
	}
	return ecPointASN1, nil
}

func ASN1BytesToPubKey(c elliptic.Curve, b []byte) (*ecdsa.PublicKey, error) {
	numSize := (c.Params().BitSize + 7) / 8
	var pointBytes []byte
	rest, err := asn1.Unmarshal(b, &pointBytes)
	if err != nil {
		return nil, fmt.Errorf("error decoding ec pubkey: %s", err.Error())
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("error decoding ec pubkey: rest length is greater than zero")
	}
	if len(pointBytes) != 2*numSize {
		return nil, fmt.Errorf("error decoding ec pubkey: length of point bytes is not 2 * curveSize")
	}
	x, y := pointBytes[:numSize], pointBytes[numSize:]
	return &ecdsa.PublicKey{
		Curve: c,
		X:     new(big.Int).SetBytes(x),
		Y:     new(big.Int).SetBytes(y),
	}, nil
}
