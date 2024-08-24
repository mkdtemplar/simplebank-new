package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/mkdtemplar/simplebank-new/db/sqlc"
)

type TaskProccessor interface {
	Start() error
	ProccessTaskVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProccessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProccessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProccessor {
	server := asynq.NewServer(redisOpt, asynq.Config{})

	return &RedisTaskProccessor{
		server: server,
		store:  store,
	}
}

// Start implements TaskProccessor.
func (proccessor *RedisTaskProccessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, proccessor.ProccessTaskVerifyEmail)

	return proccessor.server.Start(mux)
}
