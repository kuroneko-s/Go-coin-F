package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/goLangCoin/utils"
)

/*
	1) hashing the msg.
	"i love you" -> hash(x) -> "hashed_message"
	2) generate key pair
	KeyPair (privateK, publicK)
	private Key - sign, public Key - verify
	save private key to a file ( private key를 파일로써 저장한다 )-> 이게 wallet 만들떄 주는 파일
	3) sign the hash (해쉬를 서명한다.)
	hashed_message + privateK -> "signature"
	4) verify
	hashed_message + signature + publicK -> true / false
*/

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

const (
	filename string = "nomadcoin.wallet"
)

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	//key의 parsing과 marshalling을 책임지는 x509 library
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(filename, bytes, 0644)
	utils.HandleErr(err)
}

// return named -> very shot function에서 사용하는 것을 권장
func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsBytes, err := os.ReadFile(filename)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return
}

// 16byte + 16byte -> 32byte type string
func encodeBigInts(a, b []byte) string {
	result := append(a, b...)
	return fmt.Sprintf("%x", result)
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

// 서명해주는 얘
// payload -> Sing할떄 필요한 message에 해당함
func Sign(payload string, w *wallet) string {
	// 그냥 []byte(payload)해도 동작은 되지만 string의 값에 문제가 있는지 우선적으로 확인하는 차원에서 진행
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

// String to big.Int
func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

// sign의 결과물을 검증해주는 함수
func Verify(signature, payload, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	// 그냥 별도로 만들어줌
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

// 사람들이 보내는 코인의 주소는 public key여야 한다.
func Wallet() *wallet {
	if w == nil {
		w = &wallet{} // instance
		// has a wallet already?
		// yes - restore from file
		// no - create private key, save to file
		if hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivKey()
			// false is don't have the key file
			persistKey(key)    // make the file by key
			w.privateKey = key // binde the key in wallet
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}
