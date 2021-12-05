package event

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	db "github.com/tuhalang/quiz-server/app/db/sqlc"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// QuizEvent handles event of contract
type QuizEvent struct {
	store *db.Store
}

// NewQuizEvent returns a new QuizEvent
func NewQuizEvent(store *db.Store) (*QuizEvent, error) {
	return &QuizEvent{store: store}, nil
}

func initQuery(wss, address string, blockNumber int64) (*ethclient.Client, *ethereum.FilterQuery, error) {
	client, err := ethclient.Dial(wss)
	if err != nil {
		return nil, nil, err
	}
	contractAddress := common.HexToAddress(address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		FromBlock: big.NewInt(blockNumber),
	}
	return client, &query, nil
}

// Filter gets all event of contract from blockNumber
func (event *QuizEvent) Filter(wss, address string, blockNumber, step int64) ([]types.Log, error) {
	client, query, err := initQuery(wss, address, blockNumber)
	if err != nil {
		return nil, err
	}
	if blockNumber != 0 && step > 0 {
		query.ToBlock = big.NewInt(blockNumber + step)
	}
	logs, err := client.FilterLogs(context.Background(), *query)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// Subscriber listens all event of contract
func (event *QuizEvent) Subscriber(wss, address string, blockNumber int64, fn func(types.Log) error) {

RETRY:
	client, query, err := initQuery(wss, address, blockNumber)

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), *query, logs)
	if err != nil {
		log.Println(err)
		goto RETRY
	}

	for {
		select {
		case err := <-sub.Err():
			log.Println(err)
			goto RETRY
		case vLog := <-logs:
			err := fn(vLog)
			if err != nil {
				log.Println("error when handle event: ", err)
			}
		}
	}
}
