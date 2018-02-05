package messengers

type Messenger interface {
	SendMessage(string, string) error
}
