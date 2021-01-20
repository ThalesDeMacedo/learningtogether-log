package person

import (
	"awesomeProject/log"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"math/rand"
)

type Person struct {
	Id int
	Name *string
	MotherName *string
	Age int
	Score int
}

func CreateTable() *memdb.TableSchema {
	return &memdb.TableSchema{
		Name: "person",
		Indexes: map[string]*memdb.IndexSchema{
			"id": &memdb.IndexSchema{
				Name:    "id",
				Unique:  true,
				Indexer: &memdb.IntFieldIndex{Field: "Id"},
			},
			"name": &memdb.IndexSchema{
				Name:    "name",
				Indexer: &memdb.StringFieldIndex{Field: "Name"},
				AllowMissing: false,
			},
			"mother_name": &memdb.IndexSchema{
				Name:    "mother_name",
				Indexer: &memdb.StringFieldIndex{Field: "MotherName"},
				AllowMissing: true,
			},
			"age": &memdb.IndexSchema{
				Name:    "age",
				Indexer: &memdb.IntFieldIndex{Field: "Age"},
			},
		},
	}
}

func (p Person) Insert(db *memdb.MemDB, requestId uuid.UUID) int {
	log.Info("abrir transaction para salvar pessoa", requestId)
	txn := db.Txn(true)

	if p.Id == 0 {
		log.Info("gerando novo Id para pessoa", requestId)
		p.Id = rand.Int()
	}

	log.Info("criando pessoa", requestId)
	if err := txn.Insert("person", p); err != nil {
		log.Error(fmt.Sprintf("erro ao criar person: %v", err), requestId)
	}
	log.Info("pessoa criada com sucesso", requestId)

	txn.Commit()
	return p.Id
}

func FindByName(db *memdb.MemDB, name string, requestId uuid.UUID) *Person {
	log.Info("abrir transaction de leitura", requestId)
	txn := db.Txn(false)
	defer txn.Abort()

	log.Info("Buscando pessoa by name", requestId)
	raw, err := txn.First("person", "name", name)
	if err != nil {
		log.Error("Problema ao buscar person", requestId)
	}

	if raw == nil {
		log.Info("pessoa não encontrada", requestId)
		return nil
	}

	log.Info("pessoa encontrada", requestId)
	person := raw.(Person)
	return &person
}

func FindById(db *memdb.MemDB, id int, requestId uuid.UUID) *Person {
	log.Info("abrir transaction de leitura", requestId)
	txn := db.Txn(false)
	defer txn.Abort()

	log.Info("Buscando pessoa by Id", requestId)
	raw, err := txn.First("person", "id", id)
	if err != nil {
		log.Error("Problema ao buscar person", requestId)
	}

	if raw == nil {
		log.Info("pessoa não encontrada", requestId)
		return nil
	}

	log.Info("pessoa encontrada", requestId)
	person := raw.(Person)
	return &person
}