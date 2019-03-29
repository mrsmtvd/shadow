package session

import (
	"sync"
)

type AutoExpireFlashBag struct {
	lock     sync.Mutex
	messages map[string][]string
}

func NewAutoExpireFlashBag() *AutoExpireFlashBag {
	return &AutoExpireFlashBag{
		messages: make(map[string][]string),
	}
}

func (b *AutoExpireFlashBag) Notice(message string) {
	b.Add(LevelNotice, message)
}

func (b *AutoExpireFlashBag) Info(message string) {
	b.Add(LevelInfo, message)
}

func (b *AutoExpireFlashBag) Success(message string) {
	b.Add(LevelSuccess, message)
}

func (b *AutoExpireFlashBag) Error(message string) {
	b.Add(LevelError, message)
}

func (b *AutoExpireFlashBag) Add(level, message string) {
	b.lock.Lock()
	defer b.lock.Unlock()

	if _, ok := b.messages[level]; !ok {
		b.messages[level] = make([]string, 0)
	}

	b.messages[level] = append(b.messages[level], message)
}

func (b *AutoExpireFlashBag) Get(level string) []string {
	b.lock.Lock()
	defer b.lock.Unlock()

	messages, ok := b.messages[level]
	if ok {
		b.messages[level] = make([]string, 0)
		return messages
	}

	return nil
}

func (b *AutoExpireFlashBag) All() map[string][]string {
	b.lock.Lock()
	defer b.lock.Unlock()

	messages := b.messages
	b.messages = make(map[string][]string)

	return messages
}
