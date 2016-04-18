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

package remotecommand

import (
	"encoding/json"
	"net/http"
	"sync"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/util/httpstream"
	"k8s.io/kubernetes/pkg/util/runtime"
)

// streamProtocolV3 implements version 3 of the streaming protocol for attach
// and exec. This version adds support for resizing the container's terminal.
type streamProtocolV3 struct {
	*streamProtocolV2

	resizeStream httpstream.Stream
}

var _ streamProtocolHandler = &streamProtocolV3{}

func newStreamProtocolV3(options StreamOptions) streamProtocolHandler {
	return &streamProtocolV3{
		streamProtocolV2: newStreamProtocolV2(options).(*streamProtocolV2),
	}
}

func (e *streamProtocolV3) createStreams(conn streamCreator) error {
	// set up the streams from v2
	if err := e.streamProtocolV2.createStreams(conn); err != nil {
		return err
	}

	// set up resize stream
	if e.Tty {
		headers := http.Header{}
		headers.Set(api.StreamType, api.StreamTypeResize)
		var err error
		e.resizeStream, err = conn.CreateStream(headers)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *streamProtocolV3) handleResizes() {
	if e.resizeStream == nil {
		return
	}

	go func() {
		defer runtime.HandleCrash()

		encoder := json.NewEncoder(e.resizeStream)
		for {
			size := e.TerminalSizeQueue.Next()
			if size == nil {
				return
			}
			if err := encoder.Encode(&size); err != nil {
				runtime.HandleError(err)
			}
		}
	}()
}

func (e *streamProtocolV3) stream(conn streamCreator) error {
	// set up all the streams first
	if err := e.createStreams(conn); err != nil {
		return err
	}

	// now that all the streams have been created, proceed with reading & copying

	// always read from errorStream
	errorChan := e.setupErrorStreamReading()

	e.handleResizes()

	e.copyStdin()

	var wg sync.WaitGroup
	e.copyStdout(&wg)
	e.copyStderr(&wg)

	// we're waiting for stdout/stderr to finish copying
	wg.Wait()

	// waits for errorStream to finish reading with an error or nil
	return <-errorChan
}
