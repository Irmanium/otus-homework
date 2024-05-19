package tarantoolrepo

import (
	"context"
	"sort"
	"strings"

	"github.com/tarantool/go-tarantool/v2"
	"otus-homework/internal/domain"
)

type Repo struct {
	*tarantool.Connection
}

func NewAndMigrate(conn *tarantool.Connection) (*Repo, error) {
	c := &Repo{Connection: conn}

	_, err := conn.Do(tarantool.NewEvalRequest(migrateMessagesQuery)).Get()
	if err != nil {
		return nil, err
	}
	_, err = conn.Do(tarantool.NewEvalRequest(migrateMessagesPrimaryIndexQuery)).Get()
	if err != nil {
		return nil, err
	}
	_, err = conn.Do(tarantool.NewEvalRequest(migrateMessagesIndexQuery)).Get()
	if err != nil {
		return nil, err
	}
	_, err = conn.Do(tarantool.NewEvalRequest(migrateGetFuncQuery)).Get()
	if err != nil {
		return nil, err
	}
	_, err = conn.Do(tarantool.NewEvalRequest(migrateInsertFuncQuery)).Get()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getDialogID(firstID, secondID string) string {
	ids := make([]string, 0, 2)
	ids = append(ids, firstID, secondID)
	sort.Strings(ids)

	return strings.Join(ids, "")
}

func (r *Repo) SendMessage(_ context.Context, id string, message domain.Message) error {
	_, err := r.Do(tarantool.NewCallRequest(insertDialogFuncName).Args([]interface{}{id, getDialogID(message.From, message.To), message.From, message.To, message.Text})).Get()
	return err
}

func (r *Repo) GetDialog(_ context.Context, userID, interlocutorID string) ([]domain.Message, error) {
	data, err := r.Do(tarantool.NewCallRequest(getDialogFuncName).
		Args([]interface{}{getDialogID(userID, interlocutorID)}),
	).Get()
	if err != nil {
		return nil, err
	}

	rows := data[0].([]interface{})
	messages := make([]domain.Message, 0, len(rows))
	for _, row := range rows {
		values := row.([]interface{})

		messages = append(messages, domain.Message{
			From: values[2].(string),
			To:   values[3].(string),
			Text: values[4].(string),
		})
	}

	return messages, nil
}
