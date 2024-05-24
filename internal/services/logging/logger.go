package logging

type Logger interface {
	LogInfo(msg string)
	LogWarn(msg string)
	LogError(msg string)
}
