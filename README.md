#노마드 코인     
https://nomadcoders.co/nomadcoin/lobby      
노마드 코더님의 인터넷 강의

        
                
                
                
                
                
<hr/>
Tx
TxIn[]
이 거래를 실행하기 이전의 금액
TxOut[]
거래가 끝났을때 각각의 사람들이 갖고 있는 금액

    just chage owner
    send 5$
    TxIn[$5(Me)]
    TxOut[$0(Me), $5(You)]

    send 5$ but i have 10$
    TxIn[$10(Me)]
    TxOut[$5(Me), $5(You)]

=> UTXO 방식

코인베이스 거래 ( 채굴한 보상으로 금액을 지불 )
TxIn[$10(blockchain)]
TxOut[$10(miner), $0 blockchain]
