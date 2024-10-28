package repl

type Handler interface {
	Prompt() string
	Eval(buffer string) string
	Tab(buffer string) string
}
