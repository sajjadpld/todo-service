package registry

//go:generate mockgen -source=./contract.go -destination=./mocks/registry_mock.go -package=registry_mock
type IRegistry interface {
	Init()
	Parse(interface{})
}
