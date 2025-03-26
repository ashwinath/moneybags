package db

type ClearAndBulkAdder interface {
	Clear() error
	BulkAdd(objs any) error
}

type Counter interface {
	Count() (int64, error)
}
