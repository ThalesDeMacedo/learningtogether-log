package service

import (
	"awesomeProject/domain/card"
	"awesomeProject/domain/creditcard"
	"awesomeProject/domain/limit"
	"awesomeProject/domain/person"
	"awesomeProject/log"
	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	"math/rand"
	"time"
)

func Process(p person.Person) {

	rand.Seed(time.Now().UnixNano())
	requestId := uuid.New()

	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"card": card.CreateTable(),
			"person": person.CreateTable(),
			"credit_card": creditcard.CreateTable(),
			"limit":limit.CreateTable(),
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		log.Error(err.Error(), requestId)
		panic(err)
	}

	log.Info("Buscando pessoa", requestId)
	personSaved := person.FindByName(db, *p.Name, requestId)
	if personSaved == nil {
		log.Info("pessoa não encontrada", requestId)
		log.Info("gravar pessoa", requestId)
		id := p.Insert(db, requestId)
		log.Info("pessoa salva", requestId)
		log.Info("buscar pessoa by id", requestId)
		personSaved = person.FindById(db, id, requestId)
	}
	if personSaved != nil {
		log.Info("buscar cartão de crédito", requestId)
		if cc := creditcard.FindByHolderId(db, personSaved.Id, requestId); cc == nil {
			log.Info("pessoa não possui cartão de crédito, criando um novo", requestId)
			l := limit.New(*personSaved, requestId)
			var cardId int
			if c := card.FindByHolderId(db, personSaved.Id); c == nil {
				cardId = card.New(personSaved.Id).Insert(db, requestId)
			}
			creditcard.CreditCard{
				Id:       rand.Int(),
				HolderId: personSaved.Id,
				CardId:   cardId,
				LimitId:  l.Id,
			}.Insert(db, requestId)
		}
		log.Info("cartão criado com sucesso", requestId)
	} else {
		log.Error("pessoa não foi criada", requestId)
	}

}