package messengers

import (
	"io"
)

type Messenger interface {
	SendMessage(string, string) error
	// SendFile(string, string, io.Reader) error
	// SendAudio(string, string, io.Reader) error
	SendPhoto(string, string, io.Reader) error
	// SendVideo(string, string, io.Reader) error
}
