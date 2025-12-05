package logafa

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

type LogLevel slog.Level

var (
	LogFile *os.File
)

type LogafaHandler struct {
	handler slog.Handler // 內部用官方 handler 輸出結構化
}

func NewLogafaHandler(opts *slog.HandlerOptions) *LogafaHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	}
	// 讓官方 handler 幫我們加 source（它會正確抓 caller）
	opts.AddSource = true

	return &LogafaHandler{
		handler: slog.NewTextHandler(os.Stdout, opts), // 或 NewJSONHandler
	}
}

func (h *LogafaHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *LogafaHandler) Handle(_ context.Context, r slog.Record) error {
	// timestamp
	location := "unknown:0"
	if r.PC != 0 {
		frame, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()
		location = fmt.Sprintf("%s:%d", frame.File, frame.Line)
	}

	// --- 1. 彩色 Console ---
	var colorize *color.Color
	switch r.Level {
	case slog.LevelDebug:
		colorize = color.New(color.FgCyan)
	case slog.LevelInfo:
		colorize = color.New(color.FgGreen)
	case slog.LevelWarn:
		colorize = color.New(color.FgYellow)
	case slog.LevelError:
		colorize = color.New(color.FgRed)
	}

	// Attributes
	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		attrs += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
		return true
	})

	// Console(彩色)
	consoleLine := fmt.Sprintf("[%s] [%s] [%s] %s%s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		levelString(r.Level),
		location,
		colorize.Sprint(r.Message),
		attrs,
	)
	_, _ = os.Stdout.WriteString(consoleLine)

	// --- 3. 寫入檔案（乾淨格式） ---
	if LogFile != nil {
		fileLine := fmt.Sprintf("time=%s level=%s msg=%q %s file=%q \n",
			r.Time.Format(time.RFC3339),
			levelString(r.Level),
			r.Message,
			attrs,
			location,
		)
		LogFile.WriteString(fileLine)
	}

	return nil
}

func (h *LogafaHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LogafaHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h *LogafaHandler) WithGroup(name string) slog.Handler {
	return &LogafaHandler{handler: h.handler.WithGroup(name)}
}

func Debug(msg string, args ...any) { logf(slog.LevelDebug, msg, args...) }
func Info(msg string, args ...any)  { logf(slog.LevelInfo, msg, args...) }
func Warn(msg string, args ...any)  { logf(slog.LevelWarn, msg, args...) }
func Error(msg string, args ...any) { logf(slog.LevelError, msg, args...) }

func logf(level slog.Level, msg string, args ...any) {
	if !slog.Default().Enabled(context.Background(), level) {
		return
	}

	var pc uintptr
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:]) // 跳過：logf → Info/Debug → 真正呼叫處
	pc = pcs[0]

	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)

	_ = slog.Default().Handler().Handle(context.Background(), r)
}

func levelString(l slog.Level) string {
	switch l {
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARN"
	case slog.LevelError:
		return "ERROR"
	default:
		// 處理自訂 level，例如 DEBUG+4 → DEBUG+4
		s := l.String()
		if len(s) > 5 && s[:5] == "Level" {
			// 舊版格式：Level(12) → 12
			return fmt.Sprintf("%d", l)
		}
		// 新版格式：直接回傳 DEBUG+4 之類的
		return strings.ToUpper(s)
	}
}
