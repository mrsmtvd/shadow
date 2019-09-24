package logging

const (
	ConfigMode            = ComponentName + ".mode"
	ConfigLevel           = ComponentName + ".level"
	ConfigFields          = ComponentName + ".fields"
	ConfigStacktraceLevel = ComponentName + ".stacktrace-level"
	ConfigEncoderType     = ComponentName + ".encoder.type"
	ConfigEncoderTime     = ComponentName + ".encoder.time"
	ConfigEncoderDuration = ComponentName + ".encoder.duration"
	ConfigEncoderCaller   = ComponentName + ".encoder.caller"
	ConfigSentryEnabled   = ComponentName + ".sentry.enabled"
	ConfigSentryDSN       = ComponentName + ".sentry.dsn"
)
