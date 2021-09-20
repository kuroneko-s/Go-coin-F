package blockchain

import (
	"errors"
	"time"

	"github.com/goLangCoin/utils"
	"github.com/goLangCoin/wallet"
)

const (
	minerReward int = 50
)

//mempool -> 거래가 성립되지 않는 주소가 대기하는 값
type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

// 처음에는 코인베이스가 마이너한테 지급해야한다.
type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	TxID       string `json:"txID"`       // 해당 트랜잭션의 InOuts 중에서
	Index      int    `json:"index"`      // 몇번째에 있는지를 알려줄게
	Sigunature string `json:"sigunature"` // Owner -> Sigunature
}

type TxOut struct {
	Address string `json:"address"` // Owner -> Adress
	Amount  int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"txID"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

func isOnMempool(uTxOut *UTxOut) bool {
	exists := false
Outer:
	for _, tx := range Mempool.Txs {
		for _, input := range tx.TxIns {
			if input.TxID == uTxOut.TxID && input.Index == uTxOut.Index {
				exists = true
				break Outer
			}
		}
	}

	return exists
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	TxOut := []*TxOut{
		{address, minerReward},
	}

	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    TxOut,
	}
	tx.getId()
	return &tx
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

func (t *Tx) sign() {
	for _, txIn := range t.TxIns {
		txIn.Sigunature = wallet.Sign(t.Id, wallet.Wallet())
		// t.Id -> Tx Input을 만들때 사용했던 tx에 대한 Id 정보
		// publicKey -> Tx Input이 가지고 있는 txOut에 대한 정보를 기반으로 Address를 가지고 올 수 있고
		// privateKey -> 내가 가지고 있고
		// 결과적으로 암호화를 하고 싶은 데이터는 txId
	}
}

func validate(tx *Tx) bool {
	valid := true
	for _, txIn := range tx.TxIns {
		prevTx := FindTx(Blockchain(), txIn.TxID)
		if prevTx == nil {
			valid = false
			break
		}
		address := prevTx.TxOuts[txIn.Index].Address
		valid = wallet.Verify(txIn.Sigunature, tx.Id, address)
		if !valid {
			break
		}
	}

	return valid
}

var ErrorNoMoney = errors.New("not enoguh money")
var ErrorNoValid = errors.New("Tx No Valid")

func makeTx(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, Blockchain()) < amount {
		return nil, ErrorNoMoney
	}
	var txOuts []*TxOut
	var txIns []*TxIn
	total := 0
	uTxOuts := UTxOutsByAddress(from, Blockchain())
	for _, uTxOut := range uTxOuts {
		if total >= amount {
			break
		}

		txIn := &TxIn{uTxOut.TxID, uTxOut.Index, from}
		txIns = append(txIns, txIn)
		total += uTxOut.Amount
	}

	if change := total - amount; change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}

	txOut := &TxOut{to, amount}
	txOuts = append(txOuts, txOut)
	tx := &Tx{"", int(time.Now().Unix()), txIns, txOuts}
	tx.getId()
	tx.sign()
	valid := validate(tx)
	if !valid {
		return nil, ErrorNoValid
	}
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx(wallet.Wallet().Address, to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

//get tx to confirm
func (m *mempool) TxToConfirm() []*Tx {
	// tx를 다 받아서 처리하고 mempool을 비워준다.
	coinbase := makeCoinbaseTx(wallet.Wallet().Address)
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
