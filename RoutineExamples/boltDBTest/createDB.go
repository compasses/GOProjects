package main

import (
	"time"
	"log"
	"fmt"
	"os"

    "encoding/json" 
	"github.com/boltdb/bolt"
)

func showStats(db *bolt.DB) {
	// Grab the initial stats.
	prev := db.Stats()


	for {
		// Wait for 10s.
		time.Sleep(10 * time.Second)

		// Grab the current stats and diff them.
		stats := db.Stats()
		diff := stats.Sub(&prev)

		// Encode stats to JSON and print to STDERR.
		json.NewEncoder(os.Stderr).Encode(diff)

		// Save stats for the next loop.
		prev = stats
	}
}

func main() {
	log.SetFlags(0);
	db, err := bolt.Open("./newDB", 0666, &bolt.Options{NoGrowSync: true, ReadOnly: false});
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	fmt.Println("DB Info: ", db.Info())
	// Start a writable transaction.
	tx, err := db.Begin(true)

	if err != nil {
		fmt.Println("Got error ", err)
	}

	defer tx.Rollback()
	
	// Use the transaction...
	b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
	if err != nil {
		fmt.Println("Got error ", err)
	}
	
	data := make([]byte, 4096)
	timeNow := time.Now()
	str := fmt.Sprintf("time%d", timeNow)
	
	b.Put([]byte(str), data)
	
	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		fmt.Println("Got error ", err)
	}
	showStats(db)
//	db.Update(func(tx *bolt.Tx) error {
//        // Retrieve the users bucket.
//        // This should be created when the DB is first opened.
//        b := tx.Bucket([]byte("users"))

//        // Generate ID for the user.
//        // This returns an error only if the Tx is closed or not writeable.
//        // That can't happen in an Update() call so I ignore the error check.
//        id, _ = b.NextSequence()
//        u.ID = int(id)

//        // Marshal user data into bytes.
//        buf, err := json.Marshal(u)
//        if err != nil {
//            return err
//        }

//        // Persist bytes to users bucket.
//        return b.Put(itob(u.ID), buf)
//    })
	
}