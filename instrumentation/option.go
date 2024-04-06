package instrumentation

import "github.com/glocurrency/commons/logger"

type NoticeOption interface {
	Apply(*logger.Entry) *logger.Entry
}

type withField struct {
	key   string
	value interface{}
}

func (w withField) Apply(entry *logger.Entry) *logger.Entry {
	return entry.EWithFields(map[string]interface{}{w.key: w.value})
}

func WithField(key string, value interface{}) NoticeOption {
	return withField{key: key, value: value}
}

type withFields struct {
	fields map[string]interface{}
}

func (w withFields) Apply(entry *logger.Entry) *logger.Entry {
	return entry.EWithFields(w.fields)
}

func WithFields(fields map[string]interface{}) NoticeOption {
	return withFields{fields: fields}
}
