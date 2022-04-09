package cryptography

import (
	".authenticator/internal/utils"
	"bytes"
	"encoding/gob"
	"net"
)

func GetMacAddr() []string { //change to []byte(?)
	ifas, err := net.Interfaces()
	utils.Handle(err)
	if err != nil {
		return nil
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as
}

func ToBytes(certificate Certificate) []byte {
	buffer := &bytes.Buffer{}
	gob.NewEncoder(buffer).Encode(certificate)
	byteSlice := buffer.Bytes()
	return byteSlice
}

func FromBytes(byteSlice []byte) Certificate {
	var unmarshal Certificate
	bf := bytes.NewBuffer(byteSlice)
	gob.NewDecoder(bf).Decode(&unmarshal)

	/*
		privateKey, _ := crypto.ParsePKCS1PrivateKey(unmarshal.PrivateKey)
		publickey, _ := crypto.ParsePKCS1PublicKey(unmarshal.PublicKey)

		unmarshal.KeyPair = &KeyPair{
			PrivateKey: &PrivateKey{
				PrivateKey: privateKey,
			},
			PublicKey: &PublicKey{
				PublicKey: publickey,
			},
		}

	*/

	return unmarshal
}
