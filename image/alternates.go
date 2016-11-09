// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"bytes"
	"sort"
	"sync"

	"github.com/google/tiff"
)

var (
	// Handlers based on the value in the Make tag (tag id 271)
	// If a manufacturer specific package exists and is imported, one would
	// hope it would know best how to handle decoding the tiff.
	altTIFFHandlersFromMakeTagValue = struct {
		mu            sync.RWMutex
		makeToHandler map[string]TIFFHandler
	}{
		makeToHandler: make(map[string]TIFFHandler, 1),
	}

	// Handlers based on the presence of a tag id in IFD0.
	altTIFFHandlersFromTagPresence = struct {
		mu           sync.RWMutex
		tagToHandler map[uint16]TIFFHandler
	}{
		tagToHandler: make(map[uint16]TIFFHandler, 1),
	}
)

func findAlternateIFDHandler(ifd tiff.IFD) IFDHandler {
	return nil
}

// For now, findAlternates will only check against two concepts.  One, is the
// "Make" tag and the other checks the presence of tags.
func findAlternateTIFFHandler(t tiff.TIFF) TIFFHandler {
	ifd0 := t.IFDs()[0]
	// Do tag presence check first.  This is useful for identifying
	// tiff files that conform to a certain specification that uses
	// a tag to identify that specification.  For example, tiff/ep
	// uses the "TIFF/EPStandardID" tag (id 37398) to indicate the
	// version number of tiff/ep in use for this tiff file.  DNG has
	// the "DNGVersion" tag (id 50706) for the same.  And Leaf .MOS
	// files do not indicate make/model, but they do have a custom
	// private tag that they use for certain data (tag id 34310).
	// The idea here is that a tiff/ep package or dng package or
	// leaf mos package would be able to better handle processing
	// the tiff than the generic package.
	for _, tagID := range ListRegisteredTagPresenceIDs() {
		if ifd0.HasField(tagID) {
			hndlr := GetHandlerByTagPresence(tagID)
			if hndlr != nil && hndlr.CanHandle(t) {
				return hndlr
			}
		}
	}

	// Check for a specific make next
	// Tag ID 271 is a baseline tag for the "Make" or maker or
	// manufacturer of a device that created this tiff.  The idea
	// here is that a separate manufacturer package may be more
	// directly aware of how to process a specific tiff as opposed
	// to playing guessing games in a generic package.
	if ifd0.HasField(271) {
		f := ifd0.GetField(271)
		maker := string(bytes.TrimRight(f.Value().Bytes(), " \x00"))
		hndlr := GetHandlerByMake(maker)
		if hndlr != nil && hndlr.CanHandle(t) {
			return hndlr
		}
	}
	return nil
}

func RegisterHandlerByMake(m string, h TIFFHandler) {
	altTIFFHandlersFromMakeTagValue.mu.Lock()
	defer altTIFFHandlersFromMakeTagValue.mu.Unlock()
	altTIFFHandlersFromMakeTagValue.makeToHandler[m] = h
}

func GetHandlerByMake(m string) TIFFHandler {
	altTIFFHandlersFromMakeTagValue.mu.RLock()
	defer altTIFFHandlersFromMakeTagValue.mu.RUnlock()
	return altTIFFHandlersFromMakeTagValue.makeToHandler[m]
}

func RegisterHandlerByTagPresence(t uint16, h TIFFHandler) {
	altTIFFHandlersFromTagPresence.mu.Lock()
	defer altTIFFHandlersFromTagPresence.mu.Unlock()
	altTIFFHandlersFromTagPresence.tagToHandler[t] = h
}

func GetHandlerByTagPresence(t uint16) TIFFHandler {
	altTIFFHandlersFromTagPresence.mu.RLock()
	defer altTIFFHandlersFromTagPresence.mu.RUnlock()
	return altTIFFHandlersFromTagPresence.tagToHandler[t]
}

func ListRegisteredTagPresenceIDs() []uint16 {
	altTIFFHandlersFromTagPresence.mu.RLock()
	defer altTIFFHandlersFromTagPresence.mu.RUnlock()
	ids := make([]uint16, 0, len(altTIFFHandlersFromTagPresence.tagToHandler))
	for k := range altTIFFHandlersFromTagPresence.tagToHandler {
		ids = append(ids, k)
	}
	sort.Sort(uint16Slice(ids))
	return ids
}
