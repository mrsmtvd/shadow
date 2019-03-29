package session

type FlashBag interface {
	Add(level, message string)
	All() map[string][]string
	Get(level string) []string
}
