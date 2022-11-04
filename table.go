package table

// 키 & 값을 읽고 쓸수 있음
import (
	"encoding/binary"
)

/*
 테이블은 하나 이상의 블록으로 구성
 account & transaction 으로 구성되어있음
 매개변수를 유지하는 특수 블록은 
 account - address,type,info(기타내용),노드 수
 transaction - tx hash,block number,from address,to address,value,timestamp

*/