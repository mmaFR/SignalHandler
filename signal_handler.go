package SignalHandler

import (
	"context"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

const structure string = "SignalHandler"

type SignalHandler struct {
	signalChan    chan os.Signal
	callbackFuncs []func(os.Signal, Logger)
	logger        Logger
	ctx           context.Context
	cancel        context.CancelFunc
	running       bool
	wg            *sync.WaitGroup
}

func (sh *SignalHandler) watch() {
	const function string = "watch"
	sh.logger.LogNotice(structure, function, "started", -1)
	var sig os.Signal
	sh.logger.LogDebug(structure, function, "waiting for a signal or a context cancel", -1)
	select {
	case sig = <-sh.signalChan:
		sh.logger.LogNotice(structure, function, "got the signal %s", -1, sig.String())
		var fName string
		for _, f := range sh.callbackFuncs {
			fName = runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
			sh.logger.LogDebug(structure, function, "execute the callback function named %s", -1, fName)
			f(sig, sh.logger)
			sh.logger.LogTrace(structure, function, "callback function named %s executed", -1, fName)
		}
		sh.cancel()
	case <-sh.ctx.Done():
		sh.logger.LogNotice(structure, function, "local context cancelled", -1)
		signal.Reset()
	}
	sh.logger.LogNotice(structure, function, "stopped", -1)
	sh.wg.Done()
}
func (sh *SignalHandler) RegisterCallback(f func(os.Signal, Logger)) error {
	if sh.running {
		return ErrAlreadyStarted
	} else {
		sh.callbackFuncs = append(sh.callbackFuncs, f)
		sh.logger.LogInfo(structure, "RegisterCallback", "callback named %s added", -1, runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())
		return nil
	}
}
func (sh *SignalHandler) StartOn(signals []os.Signal) error {
	if sh.running {
		return ErrAlreadyStarted
	} else {
		sh.wg.Add(1)
		sh.running = true
		var sigsName string
		for _, s := range signals {
			sigsName += s.String() + ", "
		}
		sigsName = strings.Trim(sigsName, ", ")
		sh.logger.LogNotice(structure, "StartOn", "starting to watch over the signals %s", -1, sigsName)
		sh.ctx, sh.cancel = context.WithCancel(context.Background())
		signal.Notify(sh.signalChan, signals...)
		go sh.watch()
		return nil
	}
}
func (sh *SignalHandler) Stop() error {
	if !sh.running {
		return ErrNotRunning
	} else {
		sh.logger.LogDebug(structure, "Stop", "cancelling the local context", -1)
		sh.cancel()
		sh.running = false
		return nil
	}
}
func (sh *SignalHandler) Wait() {
	sh.wg.Wait()
}

func NewSignalHandler(logger Logger) *SignalHandler {
	if logger == nil {
		logger = new(dummyLogger)
	}
	return &SignalHandler{
		signalChan:    make(chan os.Signal, 1),
		callbackFuncs: make([]func(os.Signal, Logger), 0),
		logger:        logger,
		ctx:           nil,
		cancel:        nil,
		running:       false,
		wg:            &sync.WaitGroup{},
	}
}
