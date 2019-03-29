package session

const (
	LevelNotice  = "notice"
	LevelInfo    = "info"
	LevelSuccess = "success"
	LevelError   = "error"
)

type FlashBag interface {
	Notice(message string)
	Info(message string)
	Success(message string)
	Error(message string)

	Add(level, message string)
	All() map[string][]string
	Get(level string) []string
}
