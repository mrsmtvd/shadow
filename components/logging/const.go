package logging

const (
	ComponentName    = "logging"
	ComponentVersion = "3.0.0"

	ModeDevelopment = "dev"
	ModeProduction  = "prod"

	EncoderTypeJSON    = "json"
	EncoderTypeConsole = "console"

	EncoderTimeISO8601 = "iso8601"
	EncoderTimeMillis  = "millis"
	EncoderTimeNanos   = "nanos"
	EncoderTimeSeconds = "seconds"

	EncoderDurationSeconds = "seconds"
	EncoderDurationNanos   = "nanos"
	EncoderDurationString  = "string"

	EncoderCallerFull  = "full"
	EncoderCallerShort = "short"
)
