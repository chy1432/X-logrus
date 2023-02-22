package main_test

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zhumeme/X-logrus/ctxLogger"
	"path"
	"runtime"
	"testing"
	"time"
)

func Test_CtxHook(t *testing.T) {
	log := logrus.StandardLogger()
	log.AddHook(&ctxLogger.CtxHook{})
	log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.DateTime,
		FullTimestamp:   true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			return fmt.Sprintf("[%s]", f.Function), fmt.Sprintf("[%s:%d]", path.Base(f.File), f.Line)
		},
		ForceQuote:       true,
		QuoteEmptyFields: true,
		PadLevelText:     true,
	})
	log.SetReportCaller(true)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "foo", "foo")

	log.WithContext(ctx).WithFields(logrus.Fields{
		"error": fmt.Errorf("err"),
	}).Info("test")

	ctx = context.WithValue(ctx, "bar", "bar")
	log.WithContext(ctx).WithFields(logrus.Fields{
		"error": fmt.Errorf("err"),
	}).Info("test")
}
