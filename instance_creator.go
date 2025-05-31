package main

import (
	"fmt"
	"log"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	dbInstances := make(map[string]*leveldb.DB)

	for i := 1; i <= 5; i++ {
		path := fmt.Sprintf("./data/db%d", i)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create directory for %s: %v", path, err)
		}

		db, err := leveldb.OpenFile(path, nil)
		if err != nil {
			log.Fatalf("Failed to open LevelDB at %s: %v", path, err)
		}

		nodeID := fmt.Sprintf("node%d", i)
		dbInstances[nodeID] = db

		log.Printf("LevelDB instance created for %s at %s\n", nodeID, path)
	}

	// (Optional) Clean up at the end if this is a test
	// for _, db := range dbInstances {
	//     db.Close()
	// }
}
