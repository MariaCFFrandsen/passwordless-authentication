package api

import (
	crypto "crypto/x509"
	"encoding/json"
	"github.com/passwordless-authentication/authenticator/cryptography"
	"github.com/passwordless-authentication/authenticator/internal/blockchain/block"
	blockchain "github.com/passwordless-authentication/authenticator/internal/blockchain/chain"
	"github.com/passwordless-authentication/authenticator/internal/utils"
	"io/ioutil"
	"net/http"
)

const (
	dbPath = "C:\\Users\\Bruger\\goprojects\\passwordless-authentication\\authenticator\\tmp\\blocks"
	//dbPath = "..\\tmp\\blocks"
)

type Service interface { //this is the API the server should have
	CreateUser(data string, key cryptography.PublicKey) (*block.Block, error)
	AuthenticateUser(hash []byte) bool
}

type User struct {
	Username string `json:"username"`
	PK       []byte `json:"publickey"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) { //post
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	utils.Handle(err)
	var user User
	err = json.Unmarshal(body, &user)
	publicKey, err := crypto.ParsePKCS1PublicKey(user.PK)
	bc := blockchain.InitBlockChain(dbPath)
	addBlock, err := bc.AddBlock("block-from-api", &cryptography.PublicKey{PublicKey: publicKey})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if addBlock != nil && err == nil {
		json.NewEncoder(w).Encode("Created")
	} else {
		json.NewEncoder(w).Encode("Not created") //change
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	bc.Iterator().Database.Close()
}

func AuthenticateUser(w http.ResponseWriter, r *http.Request) { //post for now,
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	utils.Handle(err)
	var user User
	err = json.Unmarshal(body, &user)
	key, err := crypto.ParsePKCS1PublicKey(user.PK)
	iterator := blockchain.InitBlockChain(dbPath).Iterator()
	_, b := iterator.SearchBlockchainByPublicKey(&cryptography.PublicKey{PublicKey: key})
	//defer db close
	if b {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted) //define which status codes should mean what
		json.NewEncoder(w).Encode("true")
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable) //define which status codes should mean what
		json.NewEncoder(w).Encode("true")
	}
	iterator.Database.Close()
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("pong")
}
