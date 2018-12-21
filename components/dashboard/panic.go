package dashboard

type PanicError struct {
	Error interface{}
	Stack []byte
	File  string
	Line  int
}
