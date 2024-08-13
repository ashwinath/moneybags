package financials

type Loader interface {
	Load() error
	Name() string
}
