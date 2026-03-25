In go, if we pass a address, but the argument of a func is a value, Go will dereference and create a copy.

Example

in main():
RunApp((&JsonLogger{fileName: "log.json"}))

The correct version (in C) is:
func (c *JsonLogger) Log(level LogLevel, message string)

but this is also true:
func (c JsonLogger) Log(level LogLevel, message string)
