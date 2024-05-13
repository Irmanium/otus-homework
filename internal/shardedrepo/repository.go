package shardedrepo

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"otus-homework/internal/domain"
	"sort"
	"strings"
)

type Repo struct {
	*pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{
		Pool: pool,
	}
}

func getDialogID(firstID, secondID string) string {
	ids := make([]string, 0, 2)
	ids = append(ids, firstID, secondID)
	sort.Strings(ids)

	return strings.Join(ids, "")
}

func (r *Repo) SendMessage(ctx context.Context, id string, message domain.Message) error {
	_, err := r.Pool.Exec(ctx, sendMessageQuery, id, getDialogID(message.From, message.To), message.From, message.To, message.Text)
	return err
}

func (r *Repo) GetDialog(ctx context.Context, userID, interlocutorID string) ([]domain.Message, error) {
	rows, err := r.Pool.Query(ctx, getDialogQuery, getDialogID(userID, interlocutorID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var message domain.Message
		err = rows.Scan(&message.From, &message.To, &message.Text)
		messages = append(messages, message)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
