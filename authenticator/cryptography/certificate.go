package cryptography

import (
	".authenticator/internal/utils"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"os"
)

type Certificate struct {
	KeyPair    *KeyPair
	MacAddress []string
	Text       string
}

func CreateCertificate(txt string) {
	CreateSymmetricKey()
	key := ReadKeyFile()

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	outfile, err := os.OpenFile("ciphertext.bin", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes
	//	buf := make([]byte, 1024)
	bMsg := ToBytes(Certificate{
		Text: txt,
	})
	n := len(bMsg)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(bMsg, bMsg[:n])
	// Write into file
	outfile.Write(bMsg[:n])

	// Append the IV
	outfile.Write(iv)
}

func ReadCertificate() string {
	infile, err := os.Open("ciphertext.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	// The key should be 16 bytes (AES-128), 24 bytes (AES-192) or
	// 32 bytes (AES-256)

	key := ReadKeyFile()
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	fi, err := infile.Stat()
	if err != nil {
		log.Fatal(err)
	}

	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = infile.ReadAt(iv, msgLen)
	if err != nil {
		log.Fatal(err)
	}

	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			// The last bytes are the IV, don't belong the original message
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
	return FromBytes(buf[0 : fi.Size()-int64(len(iv))]).Text
}

func CreateSymmetricKey() []byte {
	token := make([]byte, 16)
	rand.Read(token)

	keyfile, err := os.OpenFile("key.txt", os.O_RDWR|os.O_CREATE, 0777)
	utils.Handle(err)
	keyfile.Write(token)
	return token
}

func ReadKeyFile() []byte {
	bytes := make([]byte, 16)
	infile, _ := os.Open("key.txt")
	_, err := infile.Read(bytes)
	utils.Handle(err)
	return bytes
}