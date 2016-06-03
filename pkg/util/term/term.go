/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package term

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/pkg/term"
	"k8s.io/kubernetes/pkg/util/interrupt"
	"k8s.io/kubernetes/pkg/util/runtime"
)

// SafeFunc is a function to be invoked by TTY.
type SafeFunc func() error

// TTY helps invoke a function and preserve the state of the terminal, even if the
// process is terminated during execution.
type TTY struct {
	// In is a reader used to set the terminal to raw mode.
	In io.Reader
	// Out is a writer to check for a terminal.
	Out io.Writer
	// Raw is true if the terminal should be set raw.
	Raw bool
	// TryDev indicates the TTY should try to open /dev/tty if the provided input
	// is not a file descriptor.
	TryDev bool
	// Parent is an optional interrupt handler provided to this function - if provided
	// it will be invoked after the terminal state is restored. If it is not provided,
	// a signal received during the TTY will result in os.Exit(0) being invoked.
	Parent *interrupt.Handler
	// sizeQueue is set after a call to MonitorSize() and is used to monitor SIGWINCH signals when the
	// user's terminal resizes.
	sizeQueue *sizeQueue
}

// fd returns a file descriptor for a given object.
type fd interface {
	Fd() uintptr
}

// IsTerminal returns true if the provided input is a terminal. Does not check /dev/tty
// even if TryDev is set.
func (t TTY) IsTerminal() bool {
	return IsTerminal(t.In)
}

// Safe invokes the provided function and will attempt to ensure that when the
// function returns (or a termination signal is sent) that the terminal state
// is reset to the condition it was in prior to the function being invoked. If
// t.Raw is true the terminal will be put into raw mode prior to calling the function.
// If the input file descriptor is not a TTY and TryDev is true, the /dev/tty file
// will be opened (if available).
func (t TTY) Safe(fn SafeFunc) error {
	in := t.In

	var hasFd bool
	var inFd uintptr
	if desc, ok := in.(fd); ok && in != nil {
		inFd = desc.Fd()
		hasFd = true
	}
	if t.TryDev && (!hasFd || !term.IsTerminal(inFd)) {
		if f, err := os.Open("/dev/tty"); err == nil {
			defer f.Close()
			inFd = f.Fd()
			hasFd = true
		}
	}
	if !hasFd || !term.IsTerminal(inFd) {
		return fn()
	}

	var state *term.State
	var err error
	if t.Raw {
		state, err = term.MakeRaw(inFd)
	} else {
		state, err = term.SaveState(inFd)
	}
	if err != nil {
		return err
	}
	return interrupt.Chain(t.Parent, func() {
		if t.sizeQueue != nil {
			t.sizeQueue.stop()
		}

		term.RestoreTerminal(inFd, state)
	}).Run(fn)
}

// IsTerminal returns whether the passed io.Reader is a terminal or not
func IsTerminal(r io.Reader) bool {
	file, ok := r.(fd)
	return ok && term.IsTerminal(file.Fd())
}

// Size represents the width and height of a terminal.
type Size struct {
	Width  uint16
	Height uint16
}

// GetSize returns the current size of the user's terminal. If it isn't a terminal,
// nil is returned.
func (t TTY) GetSize() *Size {
	if t.Out == nil {
		return nil
	}

	out := t.Out

	desc, ok := out.(fd)
	if !ok {
		return nil
	}

	return GetSize(desc.Fd())
}

// GetSize returns the current size of the terminal associated with fd.
func GetSize(fd uintptr) *Size {
	winsize, err := term.GetWinsize(fd)
	if err != nil {
		runtime.HandleError(fmt.Errorf("unable to get terminal size: %v", err))
		return nil
	}

	return &Size{Width: winsize.Width, Height: winsize.Height}
}

// SetSize sets the terminal associated with fd.
func SetSize(fd uintptr, size Size) error {
	return term.SetWinsize(fd, &term.Winsize{Height: size.Height, Width: size.Width})
}

// MonitorSize monitors the terminal's size. It returns a TerminalSizeQueue primed with
// initialSizes, or nil if there's no TTY present.
func (t *TTY) MonitorSize(initialSizes ...*Size) TerminalSizeQueue {
	if !t.IsTerminal() {
		return nil
	}

	t.sizeQueue = &sizeQueue{
		t: *t,
		// make it buffered so we can send the initial terminal sizes without blocking, prior to starting
		// the streaming below
		resizeChan:   make(chan Size, len(initialSizes)),
		stopResizing: make(chan struct{}),
	}

	t.sizeQueue.monitorSize(initialSizes...)

	return t.sizeQueue
}

// TerminalSizeQueue is capable of returning terminal resize events as they occur.
type TerminalSizeQueue interface {
	// Next returns the new terminal size after the terminal has been resized. It returns nil when
	// monitoring has been stopped.
	Next() *Size
}

// sizeQueue implements TerminalSizeQueue
type sizeQueue struct {
	t TTY
	// resizeChan receives a Size each time the user's terminal is resized.
	resizeChan   chan Size
	stopResizing chan struct{}
}

// make sure sizeQueue implements the TerminalSizeQueue interface
var _ TerminalSizeQueue = &sizeQueue{}

// monitorSize primes resizeChan with initialSizes and then monitors for resize events. With each
// new event, it sends the current terminal size to resizeChan.
func (s *sizeQueue) monitorSize(initialSizes ...*Size) {
	// send the initial sizes
	for i := range initialSizes {
		if initialSizes[i] != nil {
			s.resizeChan <- *initialSizes[i]
		}
	}

	resizeEvents := make(chan Size, 1)

	// TTY.MonitorSize ensures we have a terminal, so it's safe to assume that s.t.Out implements fd
	monitorResizeEvents(s.t.Out.(fd).Fd(), resizeEvents, s.stopResizing)

	// listen for resize events in the background
	go func() {
		defer runtime.HandleCrash()

		for {
			select {
			case size := <-resizeEvents:
				select {
				// try to send the size to resizeChan, but don't block
				case s.resizeChan <- size:
					// send successful
				default:
					// unable to send / no-op
				}
			case <-s.stopResizing:
				return
			}
		}
	}()
}

// Next returns the new terminal size after the terminal has been resized. It returns nil when
// monitoring has been stopped.
func (s *sizeQueue) Next() *Size {
	size, ok := <-s.resizeChan
	if !ok {
		return nil
	}
	return &size
}

// stop stops the background goroutine that is monitoring for terminal resizes.
func (s *sizeQueue) stop() {
	close(s.stopResizing)
}
