package ui

type ApplicationPrinter interface {
	Name() string
	Organization() string
	Space() string
}
