package main

import (
	"flag"
	"log"
	mrand "math/rand"
	_ "net/http/pprof"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var (
	dbPath                 = path.Join(os.TempDir(), "goleveldb-testdb")
	openFilesCacheCapacity = 500
	keyLen                 = 63
	valueLen               = 256
	numKeys                = arrayInt{100000, 1332, 531, 1234, 9553, 1024, 35743}
	httpProf               = "http://163.152.26.69/"
	transactionProb        = 0.5
	enableBlockCache       = false
	enableCompression      = false
	enableBufferPool       = false
	maxManifestFileSize    = opt.DefaultMaxManifestFileSize

	wg         = new(sync.WaitGroup)
	done, fail uint32

	bpool *util.BufferPool
)

type arrayInt []int

func init() {
	flag.StringVar(&dbPath, "db", dbPath, "testdb path")
	flag.IntVar(&openFilesCacheCapacity, "openfilescachecap", openFilesCacheCapacity, "open files cache capacity")
	flag.IntVar(&keyLen, "keylen", keyLen, "key length")
	flag.IntVar(&valueLen, "valuelen", valueLen, "value length")
	flag.Var(&numKeys, "numkeys", "num keys")
	flag.StringVar(&httpProf, "httpprof", httpProf, "http pprof listen addr")
	flag.Float64Var(&transactionProb, "transactionprob", transactionProb, "probablity of writes using transaction")
	flag.BoolVar(&enableBufferPool, "enablebufferpool", enableBufferPool, "enable buffer pool")
	flag.BoolVar(&enableBlockCache, "enableblockcache", enableBlockCache, "enable block cache")
	flag.BoolVar(&enableCompression, "enablecompression", enableCompression, "enable block compression")
	flag.Int64Var(&maxManifestFileSize, "maxManifestFileSize", maxManifestFileSize, "max manifest file size")
}

func (a arrayInt) String() string {
	var str string
	for i, n := range a {
		if i > 0 {
			str += ","
		}
		str += strconv.Itoa(n)
	}
	return str
}

//error 처리
func (a *arrayInt) Set(str string) error {
	var na arrayInt
	for _, s := range strings.Split(str, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			n, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			na = append(na, n)
		}
	}
	*a = na
	return nil
}

// 초기화

func init() {
	flag.StringVar(&dbPath, "db", dbPath, "testdb path")
	flag.IntVar(&openFilesCacheCapacity, "openfilescachecap", openFilesCacheCapacity, "open files cache capacity")
	flag.IntVar(&keyLen, "keylen", keyLen, "key length")
	flag.IntVar(&valueLen, "valuelen", valueLen, "value length")
	flag.Var(&numKeys, "numkeys", "num keys")
	flag.StringVar(&httpProf, "httpprof", httpProf, "http pprof listen addr")
	flag.Float64Var(&transactionProb, "transactionprob", transactionProb, "probablity of writes using transaction")
	flag.BoolVar(&enableBufferPool, "enablebufferpool", enableBufferPool, "enable buffer pool")
	flag.BoolVar(&enableBlockCache, "enableblockcache", enableBlockCache, "enable block cache")
	flag.BoolVar(&enableCompression, "enablecompression", enableCompression, "enable block compression")
	flag.Int64Var(&maxManifestFileSize, "maxManifestFileSize", maxManifestFileSize, "max manifest file size")
}

func main() {

	// 반환된 DB는 동시호출 가능
	// DB생성 & 열기

	ahmiadb, err := leveldb.OpenFile("/home/covert/Desktop/blockchain/DB/ahmia", nil)
	BitcoinAbusedb, err := leveldb.OpenFile("/home/covert/Desktop/blockchain/DB/BitcoinAbuse", nil)
	etherscandb, err := leveldb.OpenFile("/home/covert/Desktop/blockchain/DB/etherscan", nil)

	defer ahmiadb.Close()
	defer BitcoinAbusedb.Close()
	defer etherscandb.Close()

	// DB콘텐츠 읽기 & 수정
	data, err := ahmiadb.Get([]byte("key"), nil)

	err = ahmiadb.Put([]byte("key"), []byte("value"), nil)

	err = ahmiadb.Delete([]byte("key"), nil)
	iter := ahmiadb.NewIterator(nil, nil)

	iter.Release()
	err = iter.Error()
	defer ahmiadb.Close()

	go func() {
		for b := range writeReq {

			var err error
			if mrand.Float64() < transactionProb {
				log.Print("> Write using transaction")
				gTrasactionStat.start()
				var tr *leveldb.Transaction
				if tr, err = db.OpenTransaction(); err == nil {
					if err = tr.Write(b, nil); err == nil {
						if err = tr.Commit(); err == nil {
							gTrasactionStat.record(b.Len())
						}
					} else {
						tr.Discard()
					}
				}
			} else {
				gWriteStat.start()
				if err = db.Write(b, nil); err == nil {
					gWriteStat.record(b.Len())
				}
			}
			writeAck <- err
			<-writeAckAck
		}
	}()

}
