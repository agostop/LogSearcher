package db

import (
    "fmt"
    "os"
    "sync"
    "log"

    "github.com/pkg/errors"
    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/opt"
    "github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDbHandle struct {
    db        *leveldb.DB
    writeLock sync.Mutex
}

var instance *LevelDbHandle

func GetHandleInstance() *LevelDbHandle {
    return instance
}

func InitLevelDb(dbFolder string) *LevelDbHandle {

    err := createDbPathIfNotExist(dbFolder)
    if err != nil {
        panic(fmt.Sprintf("Error create dir %s by leveldbHandle: %s", dbFolder, err))
    }

    db, err := leveldb.OpenFile(dbFolder, nil)
    if err != nil {
        panic(fmt.Sprintf("the db file open failed: %s", err))
    }

    instance = &LevelDbHandle{
        db: db,
    }

    return instance
}

func (h *LevelDbHandle) Put(key []byte, value []byte) error {
    if value == nil {
        log.Print("the value is nil.")
        return errors.New("the value is nil.")
    }
    err := h.db.Put(key, value, &opt.WriteOptions{Sync: false})
    if err != nil {
        log.Printf("writing failed. key [%#v]", key)
        return errors.Wrapf(err, "error writing leveldb. key [%#v]", key)
    }

    return err
}

func (h *LevelDbHandle) Delete(key []byte) error {
    err := h.db.Delete(key, &opt.WriteOptions{Sync: false})
    if err != nil {
        log.Printf("deleting key failed, key: [%#v]", key)
        return errors.Wrapf(err, "error deleting leveldb, key: [%#v]", key)
    }
    return err
}

func (h *LevelDbHandle) Get(key []byte) ([]byte, error) {
    value, err := h.db.Get(key, nil)
    if err == leveldb.ErrNotFound {
        value = nil
        err = nil
    }
    if err != nil {
        log.Printf("getting leveldbprovider key [%#v], err:%s", key, err.Error())
        return nil, errors.Wrapf(err, "error getting leveldbhandle key [%#v]", key)
    }
    return value, nil
}

func (h *LevelDbHandle) BatchPut(batch *leveldb.Batch) error {

    h.writeLock.Lock()
    defer h.writeLock.Unlock()

    if err := h.db.Write(batch, nil); err != nil {
        log.Printf("write batch to leveldb failed.")
        return errors.Wrap(err, "error write batch to leveldb.")
    }

    return nil
}

func (h *LevelDbHandle) IteratorWithPrefix(prefix []byte) ([]string, error) {
    if len(prefix) == 0 {
        return nil, errors.Errorf("iterator prefix should not be empty key.")
    }

    r := util.BytesPrefix(prefix)
    return h.IteratorWithRange(r)
}

func (h *LevelDbHandle) IteratorWithRange(r *util.Range) ([]string, error) {
    if r == nil {
        return nil, errors.Errorf("iterator prefix should not be empty key.")
    }

    result := []string{}
    keyRange := &util.Range{Start: r.Start, Limit: r.Limit}
    it := h.db.NewIterator(keyRange, nil)
    defer it.Release()
    b := it.Last()
    if !b {
        return nil, errors.Errorf("doesn't have any key.")
    }
    // 使用倒序输出
    for it.Prev() {
        result = append(result, string(it.Value()))
    }
    return result, nil
}

func createDbPathIfNotExist(path string) error {
    _, err := os.Stat(path)
    if err == nil {
        return nil
    }

    if os.IsNotExist(err) {
        err := os.MkdirAll(path, os.ModePerm)
        if err != nil {
            return err
        }
    }

    return nil
}

func (h *LevelDbHandle) CleanAll() error {
    it := h.db.NewIterator(nil, nil)
    batch := &leveldb.Batch{}
    for it.Next() {
        batch.Delete(it.Key())
    }
    return h.BatchPut(batch)
}
