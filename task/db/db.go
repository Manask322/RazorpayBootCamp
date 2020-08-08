package db

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB
var err error

// ConnetToDatabase is used to Connect to Database
func ConnetToDatabase(databaseName string) error {
	db, err = bolt.Open("my.db", 0600, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: false})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("tasks"))
		return err
	})
}

// CreateRecord is create new Task
func CreateRecord(task string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		id, seqErr := b.NextSequence()
		if seqErr != nil {
			return seqErr
		}
		taskJSON := "{\"task\": " + "\"" + task + "\"" + ",\"done\": 0}"
		putErr := b.Put(itob(int(id)), []byte(taskJSON))
		return putErr
	})
	if err != nil {
		return err
	}
	return nil
}

//DeleteRecord is delete a Record
func DeleteRecord(ID int) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		err := b.Delete(itob(ID))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//Task is the struct used to store task in batabase
type Task struct {
	Task string `json:"task"`
	Done int    `json:"done"`
}

//ListAllTasks is lists all the incomplete tasks
func ListAllTasks() error {
	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("tasks"))
		var t Task
		c := b.Cursor()
		i := 1
		for k, v := c.First(); k != nil; k, v = c.Next() {
			_ = json.Unmarshal(v, &t)
			if t.Done == 0 {
				fmt.Println(i, ". ", t.Task, "( id = ", int(binary.BigEndian.Uint64(k)), ")")
				i++
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//MarkAsDone is marks task as done
func MarkAsDone(ID int) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		id, seqErr := b.NextSequence()
		if seqErr != nil {
			return seqErr
		}
		var t Task
		c := b.Cursor()
		task := ""
		for k, v := c.First(); k != nil; k, v = c.Next() {
			_ = json.Unmarshal(v, &t)
			if int(binary.BigEndian.Uint64(k)) == ID && t.Done != 1 {
				task = t.Task
			}
		}
		if task == "" {
			return errors.New("Task Not Found")
		}
		err = tx.Bucket([]byte("tasks")).Delete(itob(ID))
		if err != nil {
			return err
		}
		taskJSON := "{\"task\": " + "\"" + task + "\"" + ",\"done\": 1}"
		err = b.Put(itob(int(id)), []byte(taskJSON))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//CompletedTasks is used to list all the completed tasks.
func CompletedTasks() error {
	err := db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("tasks"))
		var t Task
		c := b.Cursor()
		i := 1
		for k, v := c.First(); k != nil; k, v = c.Next() {
			_ = json.Unmarshal(v, &t)
			if t.Done == 1 {
				fmt.Println(i, ". ", t.Task, "( id = ", int(binary.BigEndian.Uint64(k)), ")")
				i++
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//Helper function to convert int to bytes
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
