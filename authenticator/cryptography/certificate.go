package cryptography // Package cryptography should probably be inside a pkg folder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	crypto "crypto/x509"
	"fmt"
	"io"
	"log"
	"os"
)

type Certificate struct {
	PRK        *rsa.PrivateKey
	PUK        *rsa.PublicKey
	PrivateKey []byte
	PublicKey  []byte
	MacAddress []string
	Text       string //this should not be filename
}

type CertificateService interface {
	CreateCertificate(pair KeyPair) Certificate
	SaveCertificate(certificate Certificate)
	RetrieveCertificate() Certificate
	CreateSymmetricKey() []byte
	RetrieveSymmetricKey() []byte
}

func CreateCertificate(pair KeyPair) Certificate {
	pk := crypto.MarshalPKCS1PrivateKey(pair.PrivateKey.PrivateKey)
	puk := crypto.MarshalPKCS1PublicKey(pair.PublicKey.PublicKey)

	return Certificate{
		PrivateKey: pk,
		PublicKey:  puk,
		MacAddress: GetMacAddr(),
		Text:       "we try",
	}
}

func SaveCertificate(certificate Certificate, path ...string) {
	var exists bool
	//check if key and certificate exists
	if path != nil {
		exists = FileExists(fmt.Sprintf("%s-certificate.bin"))
		CreateSymmetricKey(path[0])
	} else {
		exists = FileExists("certificate.bin")
		CreateSymmetricKey()
	}
	if exists {
		return //this should probably return an error
	}

	key := RetrieveSymmetricKey(path[0])

	block, err := aes.NewCipher(key)
	Handle(err)

	iv := make([]byte, block.BlockSize())
	_, err = io.ReadFull(rand.Reader, iv)
	Handle(err)
	var outfile *os.File
	if path != nil {
		outfile, err = os.OpenFile(fmt.Sprintf("%s-certificate.bin", path[0]), os.O_RDWR|os.O_CREATE, 0777) //should be less permissions
	} else {
		outfile, err = os.OpenFile("certificate.bin", os.O_RDWR|os.O_CREATE, 0777) //should be less permissions
	}
	Handle(err)
	defer outfile.Close()

	bMsg := ToBytes(certificate)
	n := len(bMsg)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(bMsg, bMsg[:n])
	outfile.Write(bMsg[:n])
	outfile.Write(iv)
}

func RetrieveCertificate(path ...string) Certificate {
	var infile *os.File
	var key []byte
	if path != nil {
		infile, _ = os.Open(fmt.Sprintf("%s-certificate.bin", path[0]))
		key = RetrieveSymmetricKey(path[0])
	} else {
		infile, _ = os.Open("certificate.bin")
		key = RetrieveSymmetricKey()
	}
	defer infile.Close()

	block, err := aes.NewCipher(key)
	Handle(err)

	fi, err := infile.Stat()
	Handle(err)

	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = infile.ReadAt(iv, msgLen)
	Handle(err)

	buf := make([]byte, 4096)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			if n > int(msgLen) {
				n = int(msgLen)
			}
			msgLen -= int64(n)
			stream.XORKeyStream(buf, buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	return FromBytes(buf[0 : fi.Size()-int64(len(iv))])
}

func CreateSymmetricKey(path ...string) []byte { //all these should have a check that either zero or no arguments is given
	token := make([]byte, 16)
	rand.Read(token) //this could also be seeded

	var keyfile *os.File
	if path != nil {
		keyfile, _ = os.OpenFile(fmt.Sprintf("%s-key.txt", path[0]), os.O_RDWR|os.O_CREATE, 0777)

	} else {
		keyfile, _ = os.OpenFile("key.txt", os.O_RDWR|os.O_CREATE, 0777) //permission 0700 ? s?? der er ikke read exe
	}
	keyfile.Write(token)
	return token
}

func RetrieveSymmetricKey(path ...string) []byte {
	var infile *os.File
	bytes := make([]byte, 16)

	if path != nil {
		infile, _ = os.OpenFile(fmt.Sprintf("%s-key.txt", path[0]), os.O_RDWR|os.O_CREATE, 0777)

	} else {
		infile, _ = os.OpenFile("key.txt", os.O_RDWR|os.O_CREATE, 0777) //permission 0700 ? s?? der er ikke read exe
	}
	_, err := infile.Read(bytes)
	Handle(err)
	return bytes
}
