// +build windows

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

func monitorResizeEvents(in uintptr, resizeEvents chan<- Size, stop chan struct{}) {
	// FIXME(ncdc) I've been told this doesn't work for Windows, and I don't have access to a Windows
	// system, so I'm commenting this out for now. Hopefully someone who does have Windows can
	// implement this properly in the future.
	/*
		go func() {
			defer runtime.HandleCrash()

			var lastSize Size

			for {
				// see if we need to stop running
				select {
				case <-stop:
					return
				default:
				}

				size := GetSize(in)
				if size == nil {
					return
				}

				if size.Height != lastSize.Height || size.Width != lastSize.Width {
					lastSize.Height = size.Height
					lastSize.Width = size.Width
					resizeEvents <- *size
				}

				time.Sleep(250 * time.Millisecond)
			}
		}()
	*/
}
