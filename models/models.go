package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"

	"github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"

	"github.com/containerops/vessel/modules/setting"
)

var (
	LedisDB *ledis.DB

	fakeValue = []byte("0")
)

var (
	ErrObjectNotExist = errors.New("Object does not exist")
)

type SetType string

const (
	SET_TYPE_FLOW     SetType = "GLOBAL_FLOW"
	SET_TYPE_PIPELINE SetType = "GLOBAL_PIPELINE"
	SET_TYPE_STAGE    SetType = "GLOBAL_STAGE"
	SET_TYPE_JOB      SetType = "GLOBAL_JOB"
)

func InitDb() error {
	opt := &config.Config{
		DataDir: path.Join(setting.DataDir, "vessel.db"),
	}

	l, err := ledis.Open(opt)
	if err != nil {
		return fmt.Errorf("open Ledis DB: %v", err)
	}

	db := 0
	LedisDB, err = l.Select(db)
	if err != nil {
		return fmt.Errorf("select Ledis DB '%d': %v", db, err)
	}

	return nil
}

func getSetName(obj interface{}) (SetType, error) {
	var setName SetType
	switch tp := obj.(type) {
	case *Flow:
		setName = SET_TYPE_FLOW
	case *Pipeline:
		setName = SET_TYPE_PIPELINE
	default:
		return "", fmt.Errorf("unknown type: %v", tp)
	}
	return setName, nil
}

// Save saves an object with given UUID.
func Save(uuid string, obj interface{}) error {
	setName, err := getSetName(obj)
	if err != nil {
		return err
	}

	value, err := json.Marshal(obj)
	if err != nil {
		return fmt.Errorf("encoding JSON: %v", err)
	}

	if _, err = LedisDB.HSet([]byte(setName), []byte(uuid), fakeValue); err != nil {
		return fmt.Errorf("HSet: %v", err)
	} else if err = LedisDB.Set([]byte(uuid), value); err != nil {
		return fmt.Errorf("Set: %v", err)
	}

	return nil
}

// Retrieve reads and returns an object with given UUID.
func Retrieve(uuid string, obj interface{}) error {
	value, err := LedisDB.Get([]byte(uuid))
	if err != nil {
		return err
	} else if len(value) == 0 {
		return ErrObjectNotExist
	}

	return json.Unmarshal(value, obj)
}

// Delete deletes an object with given UUID.
func Delete(uuid string, setName SetType) (err error) {
	if _, err = LedisDB.HDel([]byte(setName), []byte(uuid)); err != nil {
		return err
	}
	_, err = LedisDB.Del([]byte(uuid))
	return err
}
