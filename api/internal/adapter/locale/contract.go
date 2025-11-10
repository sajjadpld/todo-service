package locale

//go:generate mockgen -source=./contract.go -destination=./mocks/locale_mock.go -package=locale_mock
type ILocale interface {
	Init()
	Get(key string) string
	Plural(key string, params map[string]string) string
}
