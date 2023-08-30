package nuts

import (
	"github.com/xujiajun/nutsdb"
	"go.uber.org/zap"
	"pear-admin-go/app/core/log"
)

var nuts *INuts

func Instance() *INuts {
	if nuts == nil {
		log.Instance().Error("No nuts DB Clint")
		return nil
	}
	return nuts
}

func InitNuts() {
	nuts = NewINuts("nutsdb0", "runtime/ndb").Open()
}

type INuts struct {
	nuts   *nutsdb.DB
	bucket string
	dir    string
}

func NewINuts(bucket string, dir string) *INuts {
	return &INuts{bucket: bucket, dir: dir}
}

func (this *INuts) Open() *INuts {
	opt := nutsdb.DefaultOptions
	opt.Dir = this.dir
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Instance().Fatal("Nuts.Open", zap.Error(err))
	}
	this.nuts = db
	return this
}

func (this *INuts) Get(key string) (string, error) {
	data := ""
	if err := this.nuts.View(func(t *nutsdb.Tx) error {
		if e, err := t.Get(this.bucket, []byte(key)); err != nil {
			if err == nutsdb.ErrNotFoundKey {
				return nil
			} else {
				return err
			}
		} else {
			data = string(e.Value)
		}
		return nil
	}); err != nil {
		return data, err
	}
	return data, nil
}

func (this *INuts) Set(key, val string, ttl ...uint32) error {
	var tl uint32
	if len(ttl) > 0 {
		tl = ttl[0]
	}
	if err := this.nuts.Update(func(t *nutsdb.Tx) error {
		if err := t.Put(this.bucket, []byte(key), []byte(val), tl); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (this *INuts) Add(key, val string) *INuts {
	_ = this.nuts.Update(func(t *nutsdb.Tx) error {
		if err := t.Put(this.bucket, []byte(key), []byte(val), 0); err != nil {
			log.Instance().Error("NutsDB.Add", zap.Error(err))
		}
		return nil
	})
	return this
}

func (this *INuts) Delete(key string) error {
	if err := this.nuts.View(func(t *nutsdb.Tx) error {
		if err := t.Delete(this.bucket, []byte(key)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (this *INuts) GetDatas(prefix string, offset, limit int) map[string]string {
	data := make(map[string]string)
	if err := this.nuts.View(
		func(t *nutsdb.Tx) error {
			if entries, _, err := t.PrefixSearchScan(this.bucket, []byte(prefix), "", offset, limit); err != nil {
				return err
			} else {
				for _, entry := range entries {
					data[string(entry.Key)] = string(entry.Value)
				}
			}
			return nil
		}); err != nil {
		log.Instance().Error("Nuts.FindByPage", zap.Error(err))
		return data
	}
	return data
}

func (this *INuts) FindByPage(start, end string) map[string]string {
	data := make(map[string]string)
	if err := this.nuts.View(
		func(t *nutsdb.Tx) error {
			if entries, err := t.RangeScan(this.bucket, []byte(start), []byte(end)); err != nil {
				return err
			} else {
				for _, entry := range entries {
					data[string(entry.Key)] = string(entry.Value)
				}
			}
			return nil
		}); err != nil {
		log.Instance().Error("Nuts.FindByPage", zap.Error(err))
		return data
	}
	return data
}

func (this *INuts) Close() {
	_ = this.nuts.Close()
}
