package dashboard

type PanicError struct {
	Error interface{}
	Stack string
	File  string
	Line  int
}
