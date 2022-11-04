package table

// 키 & 값을 읽고 쓸수 있음
import (
	"io"
	"sync/atomic"

	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/syndtr/goleveldb/leveldb/cache"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/table"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type iStorage struct {
	storage.Storage
	read  uint64
	write uint64
}

/*
 테이블은 하나 이상의 블록으로 구성
 account & transaction 으로 구성되어있음
 매개변수를 유지하는 특수 블록은
 account - address,type,info(기타내용),노드 수
 transaction - tx hash,block number,from address,to address,value,timestamp

*/

// 테이블의 기본정보
type tabl_info struct {
	fd       storage.FileDesc // 저장 파일
	seekLeft int32
	size     int64
	// imin, imax internalKey >>> 키에 관한건 추후 따로 추가할것
}

type Writer struct {
	writer io.Writer
	err    error
	// Options
	cmp         comparer.Comparer
	filter      filter.Filter
	compression opt.Compression
	blockSize   int

	dataBlock   blockWriter
	indexBlock  blockWriter
	filterBlock filterWriter
	pendingBH   blockHandle
	offset      uint64
	nEntries    int
	// Scratch allocated enough for 5 uvarint. Block writer should not use
	// first 20-bytes since it will be used to encode block handle, which
	// then passed to the block writer itself.
	scratch            [50]byte
	comparerScratch    []byte
	compressionScratch []byte
}

// 한 검색을 수행하고 현재 검색을 왼쪽으로 반환합니다.
func (t *tabl_info) consumeSeek() int32 {
	return atomic.AddInt32(&t.seekLeft, -1)
}

func newTableFile(fd storage.FileDesc, size int64) *tabl_info {
	f := &tabl_info{
		fd:   fd,
		size: size,
	}

	f.seekLeft = int32(size / 16384)
	if f.seekLeft < 100 {
		f.seekLeft = 100
	}

	return f
}

type tWriter struct {
	t *tOps

	fd storage.FileDesc
	w  storage.Writer
	tw *table.Writer

	first, last []byte
}

type tOps struct {
	s            *session
	noSync       bool
	evictRemoved bool
	cache        *cache.Cache
	bcache       *cache.Cache
	bpool        *util.BufferPool
}

func (s *session) allocFileNum() int64 {
	return atomic.AddInt64(&s.stNextFileNum, 1) - 1
}

// NewWriter는 파일에 대해 초기화된 새 테이블 작성기를 만듭니다.

func NewWriter(f io.Writer, o *opt.Options) *Writer {
	w := &Writer{
		writer:          f,
		cmp:             o.GetComparer(),
		filter:          o.GetFilter(),
		compression:     o.GetCompression(),
		blockSize:       o.GetBlockSize(),
		comparerScratch: make([]byte, 0),
	}
	// data block
	w.dataBlock.restartInterval = o.GetBlockRestartInterval()
	// The first 20-bytes are used for encoding block handle.
	w.dataBlock.scratch = w.scratch[20:]
	// index block
	w.indexBlock.restartInterval = 1
	w.indexBlock.scratch = w.scratch[20:]
	// filter block
	if w.filter != nil {
		w.filterBlock.generator = w.filter.NewGenerator()
		w.filterBlock.flush(0)
	}
	return w
}

func (t *tOps) create() (*tWriter, error) {
	fd := storage.FileDesc{Type: storage.TypeTable, Num: t.s.allocFileNum()}
	fw, err := t.s.stor.Create(fd)
	if err != nil {
		return nil, err
	}
	return &tWriter{
		t:  t,
		fd: fd,
		w:  fw,
		tw: table.NewWriter(fw, t.s.o.Options),
	}, nil
}
func main() {

}
