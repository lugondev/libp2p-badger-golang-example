package fsm

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v2"
	"github.com/hashicorp/raft"
	"io"
	"libp2p-badger/env"
	"os"
	"strings"
)

type BadgerFSM struct {
	db   *badger.DB
	conf *env.Config
}

// Get fetch data from badgerDB
func (b BadgerFSM) Get(key string) (interface{}, error) {
	var keyByte = []byte(key)
	var data interface{}

	txn := b.db.NewTransaction(false)
	defer func() {
		_ = txn.Commit()
	}()

	item, err := txn.Get(keyByte)
	if err != nil {
		data = map[string]interface{}{}
		return data, err
	}

	var value = make([]byte, 0)
	err = item.Value(func(val []byte) error {
		value = append(value, val...)
		return nil
	})

	if err != nil {
		data = map[string]interface{}{}
		return data, err
	}

	if value != nil && len(value) > 0 {
		err = json.Unmarshal(value, &data)
	}

	if err != nil {
		data = map[string]interface{}{}
	}

	return data, err
}

// Set store data to badgerDB
func (b BadgerFSM) Set(key string, value interface{}) error {
	var data = make([]byte, 0)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if data == nil || len(data) <= 0 {
		return nil
	}

	txn := b.db.NewTransaction(true)
	err = txn.Set([]byte(key), data)
	if err != nil {
		txn.Discard()
		return err
	}

	return txn.Commit()
}

// SetArr store [data] to badgerDB
func (b BadgerFSM) SetArr(key string, value interface{}) error {
	var data = make([]byte, 0)
	data, err := json.Marshal([]interface{}{value})
	if err != nil {
		return err
	}

	if data == nil || len(data) <= 0 {
		return nil
	}
	existValue, err := b.Get(key)
	if err != nil {
		fmt.Println("get exists value:", err)
	}
	if existValue != nil && err == nil {
		//fmt.Println("existValue:", existValue)
		//fmt.Println("value:", value)
		newData := append(existValue.([]interface{}), value)
		data, err = json.Marshal(newData)
		if err != nil {
			return err
		}

	}
	txn := b.db.NewTransaction(true)

	err = txn.Set([]byte(key), data)
	if err != nil {
		txn.Discard()
		return err
	}

	return txn.Commit()
}

// Delete remove data from badgerDB
func (b BadgerFSM) Delete(key string) error {
	var keyByte = []byte(key)

	txn := b.db.NewTransaction(true)
	err := txn.Delete(keyByte)
	if err != nil {
		return err
	}

	return txn.Commit()
}

// Apply log is invoked once a log entry is committed.
// It returns a value which will be made available in the
// ApplyFuture returned by Raft.Apply method if that
// method was called on the same Raft node as the FSM.
func (b BadgerFSM) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		var payload = CommandPayload{}
		if err := json.Unmarshal(log.Data, &payload); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error marshalling store payload %s\n", err.Error())
			return nil
		}

		op := strings.ToUpper(strings.TrimSpace(payload.Operation))
		switch op {
		case "SET":
			return &ApplyResponse{
				Error: b.Set(payload.Key, payload.Value),
				Data:  payload.Value,
			}

		case "SET_ARR":
			return &ApplyResponse{
				Error: b.SetArr(payload.Key, payload.Value),
				Data:  payload.Value,
			}
		case "GET":
			data, err := b.Get(payload.Key)
			return &ApplyResponse{
				Error: err,
				Data:  data,
			}

		case "DELETE":
			return &ApplyResponse{
				Error: b.Delete(payload.Key),
				Data:  nil,
			}
		}
	}

	_, _ = fmt.Fprintf(os.Stderr, "not raft log command type\n")
	return nil
}

// Snapshot will be called during make snapshot.
// Snapshot is used to support log compaction.
// No need to call snapshot since it already persisted in disk (using BadgerDB) when raft calling Apply function.
func (b BadgerFSM) Snapshot() (raft.FSMSnapshot, error) {
	return newSnapshotNoop()
}

// Restore is used to restore an FSM from a Snapshot. It is not called
// concurrently with any other command. The FSM must discard all previous
// state.
// Restore will update all data in BadgerDB
func (b BadgerFSM) Restore(rClose io.ReadCloser) error {
	defer func() {
		if err := rClose.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[FINALLY RESTORE] close error %s\n", err.Error())
		}
	}()

	_, _ = fmt.Fprintf(os.Stdout, "[START RESTORE] read all message from snapshot\n")
	var totalRestored int

	decoder := json.NewDecoder(rClose)
	for decoder.More() {
		var data = &CommandPayload{}
		err := decoder.Decode(data)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error decode data %s\n", err.Error())
			return err
		}

		if err := b.Set(data.Key, data.Value); err != nil {
			_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error persist data %s\n", err.Error())
			return err
		}

		totalRestored++
	}

	// read closing bracket
	_, err := decoder.Token()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] error %s\n", err.Error())
		return err
	}

	_, _ = fmt.Fprintf(os.Stdout, "[END RESTORE] success restore %d messages in snapshot\n", totalRestored)
	return nil
}

// NewBadger implementation using badgerDB
func NewBadger(badgerDB *badger.DB, conf *env.Config) *BadgerFSM {
	return &BadgerFSM{
		db:   badgerDB,
		conf: conf,
	}
}
