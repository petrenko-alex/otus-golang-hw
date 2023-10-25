package logger_test

import (
	"bytes"
	"testing"

	"github.com/petrenko-alex/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	const msg = "log msg"

	t.Run("log with exact level", func(t *testing.T) {
		out := &bytes.Buffer{}

		logg := logger.New(logger.Warning, out)
		logg.Warning(msg)

		require.Contains(t, out.String(), msg)
	})

	t.Run("log with higher level", func(t *testing.T) {
		out := &bytes.Buffer{}

		logg := logger.New(logger.Info, out)
		logg.Warning(msg)

		require.Contains(t, out.String(), msg)
	})

	t.Run("log with lower level", func(t *testing.T) {
		out := &bytes.Buffer{}

		logg := logger.New(logger.Error, out)
		logg.Warning(msg)

		require.Empty(t, out.String())
	})
}
