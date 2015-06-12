# Copyright 2010 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

PEGNAME = dts
EXT = .exe

$(PEGNAME)$(EXT): $(PEGNAME).peg.go main.go
	go build -o $@

$(PEGNAME).peg.go: $(PEGNAME).peg
	peg -switch -inline $(PEGNAME).peg

clean:
	rm -f $(PEGNAME)$(EXT) *.peg.go
