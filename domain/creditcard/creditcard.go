package creditcard

import (
	"awesomeProject/log"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"math/rand"
)

type CreditCard struct {
	Id int
	HolderId int
	CardId int
	LimitId int
}

func CreateTable() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "credit_card",
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
			"card_id": &memdb.IndexSchema{
				Name:    "card_id",
				Unique:  true,
				Indexer: &memdb.IntFieldIndex{Field: "CardId"},
			},
			"limit_id": &memdb.IndexSchema{
				Name:    "limit_id",
				Unique:  true,
				Indexer: &memdb.IntFieldIndex{Field: "LimitId"},
			},
		},
	}
}

func (c CreditCard) Insert(db *memdb.MemDB, requestId uuid.UUID) int {
	txn := db.Txn(true)

	if c.Id == 0 {
		c.Id = rand.Int()
	}

	if err := txn.Insert("credit_card", c); err != nil {
		log.Error(fmt.Sprintf("erro ao criar credit card: %v", err), requestId)
	}

	txn.Commit()
	return c.Id
}

func FindByHolderId(db *memdb.MemDB, id int, requestId uuid.UUID) *CreditCard {
	txn := db.Txn(false)
	defer txn.Abort()

	raw, err := txn.First("credit_card", "holder_id", id)
	if err != nil {
		log.Error(err.Error(), requestId)
		panic(err)
	}

	if raw == nil {
		return nil
	}

	card := raw.(CreditCard)
	return &card
}