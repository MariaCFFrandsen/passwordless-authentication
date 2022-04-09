package cryptography

import (
	"bytes"
	"encoding/gob"
	"net"
)

func getMacAddr() ([]string, error) { //change to []byte
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
	//enc := gob.NewEncoder(fp)
	//enc.Encode(data)
}

func getMacAddr2() ([]byte, error) {
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}

	buffer := &bytes.Buffer{}

	gob.NewEncoder(buffer).Encode(as)
	byteSlice := buffer.Bytes()
	return byteSlice, nil
}

func ToBytes(certificate *Certificate) []byte {
	return nil
}

func FromBytes(b []byte) *Certificate {
	return nil
}

func SToBytes(s string) []byte {
	return nil
}

func SFromBytes(b []byte) string {
	return "nil"
}
