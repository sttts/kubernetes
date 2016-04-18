/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package remotecommand

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/util/httpstream"
)

// streamProtocolV1 implements the first version of the streaming exec & attach
// protocol. This version has some bugs, such as not being able to detect when
// non-interactive stdin data has ended. See http://issues.k8s.io/13394 and
// http://issues.k8s.io/13395 for more details.
type streamProtocolV1 struct {
	StreamOptions

	errorStream  httpstream.Stream
	remoteStdin  httpstream.Stream
	remoteStdout httpstream.Stream
	remoteStderr httpstream.Stream
}

var _ streamProtocolHandler = &streamProtocolV1{}

func newStreamProtocolV1(options StreamOptions) streamProtocolHandler {
	return &streamProtocolV1{
		StreamOptions: options,
	}
}

func (e *streamProtocolV1) stream(conn streamCreator) error {
	doneChan := make(chan struct{}, 2)
	errorChan := make(chan error)

	cp := func(s string, dst io.Writer, src io.Reader) {
		glog.V(6).Infof("Copying %s", s)
		defer glog.V(6).Infof("Done copying %s", s)
		if _, err := io.Copy(dst, src); err != nil && err != io.EOF {
			glog.Errorf("Error copying %s: %v", s, err)
		}
		if s == api.StreamTypeStdout || s == api.StreamTypeStderr {
			doneChan <- struct{}{}
		}
	}

	// set up all the streams first
	var err error
	headers := http.Header{}
	headers.Set(api.StreamType, api.StreamTypeError)
	e.errorStream, err = conn.CreateStream(headers)
	if err != nil {
		return err
	}
	defer e.errorStream.Reset()

	// Create all the streams first, then start the copy goroutines. The server doesn't start its copy
	// goroutines until it's received all of the streams. If the client creates the stdin stream and
	// immediately begins copying stdin data to the server, it's possible to overwhelm and wedge the
	// spdy frame handler in the server so that it is full of unprocessed frames. The frames aren't
	// getting processed because the server hasn't started its copying, and it won't do that until it
	// gets all the streams. By creating all the streams first, we ensure that the server is ready to
	// process data before the client starts sending any. See https://issues.k8s.io/16373 for more info.
	if e.Stdin != nil {
		headers.Set(api.StreamType, api.StreamTypeStdin)
		e.remoteStdin, err = conn.CreateStream(headers)
		if err != nil {
			return err
		}
		defer e.remoteStdin.Reset()
	}

	if e.Stdout != nil {
		headers.Set(api.StreamType, api.StreamTypeStdout)
		e.remoteStdout, err = conn.CreateStream(headers)
		if err != nil {
			return err
		}
		defer e.remoteStdout.Reset()
	}

	if e.Stderr != nil && !e.Tty {
		headers.Set(api.StreamType, api.StreamTypeStderr)
		e.remoteStderr, err = conn.CreateStream(headers)
		if err != nil {
			return err
		}
		defer e.remoteStderr.Reset()
	}

	// now that all the streams have been created, proceed with reading & copying

	// always read from errorStream
	go func() {
		message, err := ioutil.ReadAll(e.errorStream)
		if err != nil && err != io.EOF {
			errorChan <- fmt.Errorf("Error reading from error stream: %s", err)
			return
		}
		if len(message) > 0 {
			errorChan <- fmt.Errorf("Error executing remote command: %s", message)
			return
		}
	}()

	if e.Stdin != nil {
		// TODO this goroutine will never exit cleanly (the io.Copy never unblocks)
		// because stdin is not closed until the process exits. If we try to call
		// stdin.Close(), it returns no error but doesn't unblock the copy. It will
		// exit when the process exits, instead.
		go cp(api.StreamTypeStdin, e.remoteStdin, e.Stdin)
	}

	waitCount := 0
	completedStreams := 0

	if e.Stdout != nil {
		waitCount++
		go cp(api.StreamTypeStdout, e.Stdout, e.remoteStdout)
	}

	if e.Stderr != nil && !e.Tty {
		waitCount++
		go cp(api.StreamTypeStderr, e.Stderr, e.remoteStderr)
	}

Loop:
	for {
		select {
		case <-doneChan:
			completedStreams++
			if completedStreams == waitCount {
				break Loop
			}
		case err := <-errorChan:
			return err
		}
	}

	return nil
}
