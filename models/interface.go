package models

import "context"

type RunInTransactionFunc func(c context.Context) error

type TransactionManager interface {
	WithTransaction(c context.Context, fn RunInTransactionFunc) error
}
