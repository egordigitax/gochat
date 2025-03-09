package base

import "github.com/google/uuid"

type AggregateRoot struct {
	uid string
}

func (a *AggregateRoot) GetUid() string {
	return a.uid
}

func (a *AggregateRoot) SetRoot() {
	a.uid = uuid.New().String()
}
