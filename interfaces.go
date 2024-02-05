package SignalHandler

type Logger interface {
	LogEmerge(structure, function, msg string, id int, vars ...any)
	LogNotice(structure, function, msg string, id int, vars ...any)
	LogInfo(structure, function, msg string, id int, vars ...any)
	LogDebug(structure, function, msg string, id int, vars ...any)
	LogTrace(structure, function, msg string, id int, vars ...any)
}
