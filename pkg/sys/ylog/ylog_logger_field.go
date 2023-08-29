package ylog

type Field struct {
	Key string
	Val any
}

func Any(key string, val any) Field {
	return Field{Key: key, Val: val}
}

// With 创建一个子logger, 添加到该子logger的字段不会影响父级, 反之亦然
func (l *Logger) With(fields ...Field) ILogger {
	_l := l.copy()
	_l.fields = append(_l.fields, fields...)
	return _l
}
