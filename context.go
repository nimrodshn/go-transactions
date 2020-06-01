package main

import (
	"context"
	"database/sql"

	"github.com/golang/glog"

	errors "github.com/zgalor/weberr"
)

type contextKey int

const (
	transactionKey contextKey = iota
)

// NewContext returns a new context with transaction stored in it
func NewContext(ctx context.Context) (context.Context, error) {
	tx, err := new()
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, transactionKey, tx), nil
}

// FromContext Retrieves the transaction from the context.
func FromContext(ctx context.Context) (*sql.Tx, error) {
	transaction, err := getTransaction(ctx)
	if err != nil {
		return nil, err
	}
	return transaction.tx, err
}

func getTransaction(ctx context.Context) (*transaction, error) {
	transaction, ok := ctx.Value(transactionKey).(*transaction)
	if !ok {
		return nil, errors.Errorf("Could not retrieve transaction from context")
	}
	return transaction, nil
}

// Resolve resolves the current transaction according to the rollback flag.
func Resolve(ctx context.Context) error {
	transaction, err := getTransaction(ctx)
	if err != nil {
		glog.Errorf("%v", err)
		return err
	}

	if transaction.rollback {
		if err := rollback(ctx); err != nil {
			glog.Errorf("Could not rollback transaction: %v", err)
			return err
		}
		glog.Info("Rolled back transaction")
	} else {
		if err := commit(ctx); err != nil {
			glog.Errorf("Could not commit transaction: %v", err)
			return err
		}
	}
	return nil
}

// MarkForRollback flags the transaction stored in the context for rollback.
func MarkForRollback(ctx context.Context) {
	transaction, err := getTransaction(ctx)
	if err != nil {
		glog.Errorf("failed to mark transaction for rollback: %v", err)
		return
	}
	glog.Info("Marked transaction for rollback.")
	transaction.rollback = true
}

// AddPostCommitCallback adds a callback function that will be executed after the transaction is
// successfully committed. If multiple callbacks functions are added then the order of their
// execution isn't guaranteed, and they may run in parallel in different goroutines.
func AddPostCommitCallback(ctx context.Context, callback func()) error {
	transaction, err := getTransaction(ctx)
	if err != nil {
		return err
	}
	transaction.postCommitCallbacks = append(transaction.postCommitCallbacks, callback)
	return nil
}

// AddPostRollbackCallback adds a callback function that will be executed after the transaction is
// rolled back. If multiple callbacks functions are added then the order of their
// execution isn't guaranteed, and they may run in parallel in different goroutines.
func AddPostRollbackCallback(ctx context.Context, callback func()) error {
	transaction, err := getTransaction(ctx)
	if err != nil {
		return err
	}
	transaction.postRollbackCallbacks = append(transaction.postRollbackCallbacks, callback)
	return nil
}

// commit commits the transaction stored in context or returns an err if one occurred.
func commit(ctx context.Context) error {
	transaction, err := getTransaction(ctx)
	if err != nil {
		return err
	}

	// Commit the transaction:
	err = transaction.tx.Commit()
	if err != nil {
		return err
	}

	// Run the post commit callbacks:
	for _, callback := range transaction.postCommitCallbacks {
		if callback != nil {
			go callback()
		}
	}

	return nil
}

// rollback rollbacks the transaction stored in context or returns an err if one occurred..
func rollback(ctx context.Context) error {
	transaction, err := getTransaction(ctx)
	if err != nil {
		return err
	}

	// Rollback the transaction:
	err = transaction.tx.Rollback()
	if err != nil {
		return err
	}

	// Run the post rollback callbacks:
	for _, callback := range transaction.postRollbackCallbacks {
		if callback != nil {
			go callback()
		}
	}

	return nil
}
