package db

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type ReplayDB struct {
	db *bolt.DB
}

func NewReplayDB() (*ReplayDB, error) {
	dbopen, err := bolt.Open("./ReplayDB", 0600, nil)
	if err != nil {
		return nil, err
	}

	return &ReplayDB{db: dbopen}, nil
}

func (replay *ReplayDB) StoreRequest(path, method, reqBody, respBody string) (err error) {
	replay.db.Update(func(tx *bolt.Tx) error {
		pathBucket, err := tx.CreateBucketIfNotExists([]byte(path))
		if err != nil {
			err = fmt.Errorf("create bucket: %s", err)
			return err
		}
		methodBucket, err := pathBucket.CreateBucketIfNotExists([]byte(method))
		if len(reqBody) == 0 {
			reqBody = "request"
		}
		err = methodBucket.Put([]byte(reqBody), []byte(respBody))
		if err != nil {
			err = fmt.Errorf("store bucket error: %s", err)
			return err
		}
		return nil
	})
	return
}

func (replay *ReplayDB) GetResponse(path, method, reqBody string) (resp []byte, err error) {
	replay.db.View(func(tx *bolt.Tx) error {
		pathBucket := tx.Bucket([]byte(path))
		if pathBucket == nil {
			err = fmt.Errorf("No response found for path:%s", path)
			return err
		}
		methodBucket := pathBucket.Bucket([]byte(method))
		if methodBucket == nil {
			err = fmt.Errorf("No response for path and method:%s", path+method)
			return err
		}
		resp = methodBucket.Get([]byte(reqBody))
		if resp == nil {
			err = fmt.Errorf("No response for path and method:%s", path+method)
		}
		return nil
	})
	return
}
