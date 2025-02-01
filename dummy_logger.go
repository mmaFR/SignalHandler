package SignalHandler

type dummyLogger struct{}

func (dummyLogger) LogEmerge(_, _, _ string, _ int, _ ...any) {}
func (dummyLogger) LogNotice(_, _, _ string, _ int, _ ...any) {}
func (dummyLogger) LogInfo(_, _, _ string, _ int, _ ...any)   {}
func (dummyLogger) LogDebug(_, _, _ string, _ int, _ ...any)  {}
func (dummyLogger) LogTrace(_, _, _ string, _ int, _ ...any)  {}
