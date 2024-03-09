package txcore

type Transaction interface {
	Commit() error
	Rollback() error
}
