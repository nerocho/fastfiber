package snowflake

import (
	"fmt"
	"testing"
)

func Test_gen(t *testing.T) {
	worker, err := NewIdWorker(1, 2, twepoch)
	if err != nil {
		t.Error(err)
	}

	id, err := worker.NextId()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("生成Id: %d\n", id)
	fmt.Printf("GetGenTimestamp: %d,GetGenTime: %s\n", worker.GetGenTimestamp(id), worker.GetFormatedGenTime(id))

	workerId, datacenterId := worker.GetDeviceID(id)
	fmt.Printf("%d,%d\n", workerId, datacenterId)

}
