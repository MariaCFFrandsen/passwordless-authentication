package api

import (
	".authenticator/cryptography"
	".authenticator/internal/blockchain/block"
	blockchain2 ".authenticator/internal/blockchain/chain"
	".authenticator/internal/utils"
	crypto "crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	addBlock, err := blockchain2.InitBlockChain("filepath/").AddBlock("new block from api", &cryptography.PublicKey{PublicKey: publicKey})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if addBlock != nil && err == nil {
		json.NewEncoder(w).Encode("Created")
	} else {
		json.NewEncoder(w).Encode("Created") //change
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func AuthenticateUser(w http.ResponseWriter, r *http.Request) { //post for now,
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	utils.Handle(err)
	var user User
	err = json.Unmarshal(body, &user)
	//search the blockchain
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted) //define which status codes should mean what
	json.NewEncoder(w).Encode("true")
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("pong")
}
