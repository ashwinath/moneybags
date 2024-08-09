package db

type ClearAndBulkAdder interface {
	Clear() error
	BulkAdd(objs interface{}) error
}
