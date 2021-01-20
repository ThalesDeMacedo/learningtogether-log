package limit

import (
	"awesomeProject/domain/person"
	"awesomeProject/log"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"math"
	"math/rand"
)

type Limit struct {
	Id int
	Value int32
}

func CreateTable() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "limit",
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:    "id",
				Unique:  true,
				Indexer: &memdb.IntFieldIndex{Field: "Id"},
			},
			"value": &memdb.IndexSchema{
				Name:    "value",
				Unique:  true,
				Indexer: &memdb.IntFieldIndex{Field: "Value"},
			},
		},
	}
}

func New(person person.Person, requestId uuid.UUID) Limit {
	v := int32(math.Abs(float64(div(sum(person.Age, person.Score, requestId), 100, requestId))))
	if v > 500 {
		log.Error(fmt.Sprintf("limite inválido: %v", v), requestId)
	}
	return Limit{
		Id:    rand.Int(),
		Value: v,
	}
}

func sum(a int, b int, requestId uuid.UUID) int {
	log.Info(fmt.Sprintf("sum(%v, %v)", a, b), requestId)
	r := a + b
	log.Info(fmt.Sprintf("sum result: %v", r), requestId)
	return r
}

func div(a int, b int, requestId uuid.UUID) int {
	log.Info(fmt.Sprintf("div(%v, %v)", a, b), requestId)
	r := a + b
	log.Info(fmt.Sprintf("div result: %v", r), requestId)
	return r
}

func (l Limit) Insert(db *memdb.MemDB, requestId uuid.UUID) int {
	txn := db.Txn(true)

	if l.Id == 0 {
		l.Id = rand.Int()
	}

	if err := txn.Insert("limit", l); err != nil {
		log.Error(fmt.Sprintf("erro ao criar limit: %v", err), requestId)
	}

	txn.Commit()
	return l.Id
}