package ctxLogger

import (
	"context"
	"github.com/sirupsen/logrus"
	"reflect"
	"unsafe"
)

type CtxHook struct {
}

func (h *CtxHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *CtxHook) Fire(e *logrus.Entry) error {
	for _, key := range GetContextKeys(e.Context) {
		e.Data["X-"+key] = e.Context.Value(key)
	}
	return nil
}

const (
	_contextField = "Context"
	_keyField     = "key"
)

func GetContextKeys(ctx context.Context) (l []string) {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == _contextField {
				c, ok := reflectValue.Interface().(context.Context)
				if ok == true {
					l = append(l, GetContextKeys(c)...)
				}
			} else if reflectField.Name == _keyField {
				s, ok := reflectValue.Interface().(string)
				if ok == true {
					l = append(l, s)
				}
			}
		}
	}
	return
}
