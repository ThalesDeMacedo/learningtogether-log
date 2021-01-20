package card

import (
	"awesomeProject/log"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"math/rand"
)

type Card struct {
	Id int
	HolderId int
	Number int
}

func CreateTable() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "card",
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:    "id",
				Unique:  true,
				Indexer: &memdb.IntFieldIndex{Field: "Id"},
			},
			"holder_id": &memdb.IndexSchema{
				Name:    "holder_id",
				Unique:  true,
				Indexer: &memdb.IntFieldIndex{Field: "HolderId"},
			},
			"number": &memdb.IndexSchema{
				Name:    "number",
				Unique:  true,
				Indexer: &memdb.IntFieldIndex{Field: "Number"},
			},
		},
	}
}

func New(holderId int)Card {
	return Card{
		Id:       rand.Int(),
		HolderId: holderId,
		Number:   rand.Int(),
	}
}


func (c Card) Insert(db *memdb.MemDB, requestId uuid.UUID) int {
	txn := db.Txn(true)

	if c.Id == 0 {
		c.Id = rand.Int()
	}

	if err := txn.Insert("card", c); err != nil {
		log.Error(fmt.Sprintf("erro ao criar card: %v", err), requestId)
	}

	txn.Commit()
	return c.Id
}

func FindByHolderId(db *memdb.MemDB, id int) *Card {
	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("card", "holder_id", id)
	if err != nil {
		panic(err)
	}
	if raw == nil {
		return nil
	}
	c := raw.(Card)
	return &c
}