package toolkit

import (
	"crypto/md5"
	crypt_rand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

//就是拼接字符串，
func GetBuildStr(str ...string) string {
	var build strings.Builder
	for _, v := range str {
		build.WriteString(v)
	}
	return build.String()
}



//生成一定范围内随机数,
// isLoop 是否在循环之内,如果在，请设为ture 因为程序太快了，会生成一样的数字
func RandInt64(min, max int64, isLoop bool) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	if isLoop {
		time.Sleep(10 * time.Nanosecond)
	}

	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

func GUID() string {
	objID := NewObjectId()
	return objID.Hex()
}


var objectIdCounter uint32 = 0

// machineId stores machine id generated once and used in subsequent calls
// to NewObjectId function.
var machineId = readMachineId()

// ObjectId is a unique ID identifying a BSON value. It must be exactly 12 bytes
// long. MongoDB objects by default have such a property set in their "_id"
// property.
//
// http://www.mongodb.org/display/DOCS/Object+IDs
type ObjectId string



// readMachineId generates machine id and puts it into the machineId global
// variable. If this function fails to get the hostname, it will cause
// a runtime error.
func readMachineId() []byte {
	var sum [3]byte
	id := sum[:]
	hostname, err1 := os.Hostname()
	if err1 != nil {
		_, err2 := io.ReadFull(crypt_rand.Reader, id)

		if err2 != nil {
			panic(fmt.Errorf("cannot get hostname: %v; %v", err1, err2))
		}
		return id
	}
	hw := md5.New()
	hw.Write([]byte(hostname))
	copy(id, hw.Sum(nil))
	//fmt.Println("readMachineId:" + string(id))
	return id
}

// NewObjectId returns a new unique ObjectId.
// 4byte 时间，
// 3byte 机器ID
// 2byte pid
// 3byte 自增ID
func NewObjectId() ObjectId {
	var b [12]byte
	// Timestamp, 4 bytes, big endian
	binary.BigEndian.PutUint32(b[:], uint32(time.Now().Unix()))
	// Machine, first 3 bytes of md5(hostname)
	b[4] = machineId[0]
	b[5] = machineId[1]
	b[6] = machineId[2]
	// Pid, 2 bytes, specs don't specify endianness, but we use big endian.
	pid := os.Getpid()
	b[7] = byte(pid >> 8)
	b[8] = byte(pid)
	// Increment, 3 bytes, big endian
	i := atomic.AddUint32(&objectIdCounter, 1)
	b[9] = byte(i >> 16)
	b[10] = byte(i >> 8)
	b[11] = byte(i)
	return ObjectId(b[:])
}

// Hex returns a hex representation of the ObjectId.
// 返回16进制对应的字符串
func (id ObjectId) Hex() string {
	return hex.EncodeToString([]byte(id))
}
