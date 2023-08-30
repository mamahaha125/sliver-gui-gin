package nuts

import (
	"fmt"
	"github.com/xujiajun/nutsdb"
	"strconv"
	"testing"
)

func init() {
	InitNuts()
}

func TestInitNuts(t *testing.T) {
	err := Instance().nuts.Update(func(t *nutsdb.Tx) error {
		key := []byte("name")
		val := []byte("val")
		bucket := "table1"
		if err := t.Put(bucket, key, val, 0); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestNutsDB(t *testing.T) {
	for i := 0; i < 1000; i++ {
		Instance().Add(fmt.Sprintf("bucket:0000000000%d", i), strconv.Itoa(i))
	}
	datas := Instance().GetDatas("bucket", 0, 5)
	for k, v := range datas {
		fmt.Println(k, v)
	}
}
