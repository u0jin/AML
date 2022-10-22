package main

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {

	// 반환된 DB는 동시호출 가능
	// DB생성 & 열기

	ahmiadb, err := leveldb.OpenFile("/home/covert/Desktop/blockchain/DB/ahmia", nil)
	defer ahmiadb.Close() // 에러처리후 클린업 작업함
	if err != nil {
		fmt.Println("Error DB")
	}

	BitcoinAbusedb, err := leveldb.OpenFile("/home/covert/Desktop/blockchain/DB/BitcoinAbuse", nil)
	defer BitcoinAbusedb.Close()
	if err != nil {
		fmt.Println("Error DB")
	}

	etherscandb, err := leveldb.OpenFile("/home/covert/Desktop/blockchain/DB/etherscan", nil)
	defer etherscandb.Close()
	if err != nil {
		fmt.Println("Error DB")
	}

	// ahmiaDB데이터 저장 & 리턴
	err = ahmiadb.Put([]byte("key"), []byte("value"), nil)
	// 주어진 키에 지정된 값을 저장할때 사용함 - 키 없으면 새로 생서으 키 있으면 수정하여 저장
	// 키는 string ,byte[], int 타입가능
	ahmiadata, err := ahmiadb.Get([]byte("key"), nil)
	//키에 매핑되어 있는 값을 리턴하는 기능
	fmt.Println("ahmia Value: ", string(ahmiadata))

	// BitcoinAbuseDB데이터 저장 & 리턴
	err = BitcoinAbusedb.Put([]byte("key"), []byte("value"), nil)
	// 주어진 키에 지정된 값을 저장할때 사용함 - 키 없으면 새로 생서으 키 있으면 수정하여 저장
	// 키는 string ,byte[], int 타입가능
	BitcoinAbusedata, err := BitcoinAbusedb.Get([]byte("key"), nil)
	//키에 매핑되어 있는 값을 리턴하는 기능
	fmt.Println("BitcoinAbuse Value: ", string(BitcoinAbusedata))

	// etherscanDB데이터 저장 & 리턴
	err = etherscandb.Put([]byte("key"), []byte("value"), nil)
	// 주어진 키에 지정된 값을 저장할때 사용함 - 키 없으면 새로 생서으 키 있으면 수정하여 저장
	// 키는 string ,byte[], int 타입가능
	etherscandata, err := etherscandb.Get([]byte("key"), nil)
	//키에 매핑되어 있는 값을 리턴하는 기능
	fmt.Println("etherscan Value: ", string(etherscandata))

	err = ahmiadb.Delete([]byte("key"), nil)

	// DB 데이터 반복 - 모두 가져오기
	iter := ahmiadb.NewIterator(nil, nil)

	iter.Release()
	err = iter.Error()

}
