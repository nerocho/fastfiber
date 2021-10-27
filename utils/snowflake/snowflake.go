package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	twepoch            = int64(1635299694073)
	workerIdBits       = uint(5)
	datacenterIdBits   = uint(5)
	maxWorkerId        = -1 ^ (-1 << workerIdBits)
	maxDatacenterId    = -1 ^ (-1 << datacenterIdBits)
	sequenceBits       = uint(12)
	workerIdShift      = sequenceBits
	datacenterIdShift  = sequenceBits + workerIdBits
	timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits
	sequenceMask       = -1 ^ (-1 << sequenceBits)
	maxNextIdsNum      = 100
)

type IdWorker struct {
	sequence      int64
	lastTimestamp int64
	workerId      int64
	twepoch       int64
	datacenterId  int64
	mutex         sync.Mutex
}

func NewIdWorker(workerId, datacenterId int64, twepoch int64) (*IdWorker, error) {
	idWorker := &IdWorker{}
	if workerId > maxWorkerId || workerId < 0 {
		return nil, fmt.Errorf("worker Id 不能大于 %d 或 小于 0", maxWorkerId)
	}
	if datacenterId > maxDatacenterId || datacenterId < 0 {
		return nil, fmt.Errorf("datacenter Id 不能大于 %d 或 小于 0", maxDatacenterId)
	}
	if twepoch > 0 {
		idWorker.twepoch = twepoch
	}

	idWorker.workerId = workerId
	idWorker.datacenterId = datacenterId
	idWorker.lastTimestamp = -1
	idWorker.sequence = 0
	idWorker.mutex = sync.Mutex{}
	return idWorker, nil
}

func timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

func (id *IdWorker) NextId() (int64, error) {
	id.mutex.Lock()
	defer id.mutex.Unlock()
	timestamp := timeGen()
	if timestamp < id.lastTimestamp {
		return 0, fmt.Errorf("时钟回拨, %d milliseconds 后再试", id.lastTimestamp-timestamp)
	}
	if id.lastTimestamp == timestamp {
		id.sequence = (id.sequence + 1) & sequenceMask
		if id.sequence == 0 {
			timestamp = tilNextMillis(id.lastTimestamp)
		}
	} else {
		id.sequence = 0
	}
	id.lastTimestamp = timestamp
	return ((timestamp - id.twepoch) << timestampLeftShift) | (id.datacenterId << datacenterIdShift) | (id.workerId << workerIdShift) | id.sequence, nil
}

func (id *IdWorker) NextIds(num int) ([]int64, error) {
	if num > maxNextIdsNum || num < 0 {
		return nil, fmt.Errorf("NextIds num 不能大于 %d 或 小于 0", maxNextIdsNum)
	}
	ids := make([]int64, num)
	id.mutex.Lock()
	defer id.mutex.Unlock()
	for i := 0; i < num; i++ {
		timestamp := timeGen()
		if timestamp < id.lastTimestamp {
			return nil, fmt.Errorf("时钟回拨, %d milliseconds 后再试", id.lastTimestamp-timestamp)
		}
		if id.lastTimestamp == timestamp {
			id.sequence = (id.sequence + 1) & sequenceMask
			if id.sequence == 0 {
				timestamp = tilNextMillis(id.lastTimestamp)
			}
		} else {
			id.sequence = 0
		}
		id.lastTimestamp = timestamp
		ids[i] = ((timestamp - id.twepoch) << timestampLeftShift) | (id.datacenterId << datacenterIdShift) | (id.workerId << workerIdShift) | id.sequence
	}
	return ids, nil
}
