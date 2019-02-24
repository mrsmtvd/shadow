package shadow

// enumer:json
type componentStatus int64

const (
	ComponentStatusUnknown componentStatus = iota
	// готов к работе, но Run еще не завершился
	ComponentStatusReady
	// Run завершился с ошибкой
	ComponentStatusRunFailed
	// Run завершился успешно
	// для долгоиграющих компонентов, типа dashboard такой статус не будет никогда установлен,
	// так как Run блокируется
	ComponentStatusFinished
	// Остановлен через функцию Shutdown
	ComponentStatusShutdown
)

func (i componentStatus) Int64() int64 {
	if i < 0 || i >= componentStatus(len(_componentStatusIndex)-1) {
		return -1
	}

	return int64(i)
}
