// 指示: miu200521358
package logging

// LogLevel はログレベルを表す。
type LogLevel int

const (
	// LOG_LEVEL_VERBOSE は冗長ログを表す。
	LOG_LEVEL_VERBOSE LogLevel = 0
	// LOG_LEVEL_DEBUG はデバッグログを表す。
	LOG_LEVEL_DEBUG LogLevel = 10
	// LOG_LEVEL_INFO は情報ログを表す。
	LOG_LEVEL_INFO LogLevel = 20
	// LOG_LEVEL_WARN は警告ログを表す。
	LOG_LEVEL_WARN LogLevel = 30
	// LOG_LEVEL_ERROR はエラーログを表す。
	LOG_LEVEL_ERROR LogLevel = 40
	// LOG_LEVEL_FATAL は致命ログを表す。
	LOG_LEVEL_FATAL LogLevel = 50
)

// VerboseIndex は冗長ログの対象種別。
type VerboseIndex int

const (
	// VERBOSE_INDEX_MOTION はモーション冗長ログ。
	VERBOSE_INDEX_MOTION VerboseIndex = 230
	// VERBOSE_INDEX_IK はIK冗長ログ。
	VERBOSE_INDEX_IK VerboseIndex = 320
	// VERBOSE_INDEX_PHYSICS は物理冗長ログ。
	VERBOSE_INDEX_PHYSICS VerboseIndex = 330
	// VERBOSE_INDEX_VIEWER はビューワー冗長ログ。
	VERBOSE_INDEX_VIEWER VerboseIndex = 540
)

// IMessageBuffer はメッセージ欄バッファI/F。
type IMessageBuffer interface {
	Text() string
	Lines() []string
	Clear()
}

// IVerboseSink は冗長ログ出力のI/F。
type IVerboseSink interface {
	WriteLine(text string) error
	Close() error
}

// ILogger はログ出力I/F。
type ILogger interface {
	Level() LogLevel
	SetLevel(level LogLevel)
	EnableVerbose(idx VerboseIndex)
	DisableVerbose(idx VerboseIndex)
	IsVerboseEnabled(idx VerboseIndex) bool
	MessageBuffer() IMessageBuffer
	AttachVerboseSink(idx VerboseIndex, sink IVerboseSink)
	Verbose(idx VerboseIndex, msg string, params ...any)
	Debug(msg string, params ...any)
	Info(msg string, params ...any)
	Warn(msg string, params ...any)
	Error(msg string, params ...any)
	Fatal(msg string, params ...any)
}
