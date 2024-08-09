package db

type ClearAndBulkAdder interface {
	Clear() error
	BulkAdd(objs interface{}) error
}

type Counter interface {
	Count() (int64, error)
}
