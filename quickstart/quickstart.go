package quickstart

import (
	"context"
	"os/exec"
)

type QuickStarter interface {
	Start(path string, onComplete func()) error
	Stop()
}

type quickStarter struct {
	cancel context.CancelFunc
}

func NewQuickStarter() QuickStarter {
	return quickStarter {}
}

func (qs quickStarter) Start(path string, onComplete func()) error {
	if qs.cancel != nil {
		qs.Stop()
	}

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, path)
	err := cmd.Start()
	if err == nil {
		qs.cancel = cancel

		go awaitCompletion(&qs, ctx, cmd, onComplete)
	} else {
		cancel()
	}

	return err
}


// Note - currently, calling Stop to interrupt / kill the main process does not function
// Not sure why, but we don't need this behavior yet.   Leaving this in place but removing the button to call it.
func (qs quickStarter) Stop() {
	if qs.cancel != nil {
		qs.cancel()

		qs.cancel = nil
	}
}

func awaitCompletion(qs *quickStarter, ctx context.Context, cmd *exec.Cmd, onComplete func()) {
	cmd.Wait()

	// If it was cancelled, dont Stop and dont onComplete - assume external cancellation was intentional
	if ctx.Err() == nil {
		qs.Stop()

		onComplete()
	}
}
