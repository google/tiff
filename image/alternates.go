package image

import (
	"bytes"
	"image"
	"sort"
	"sync"

	"github.com/jonathanpittman/tiff"
)

type Handler interface {
	Process(*tiff.TIFF) (image.Image, error)
	GetConfig(*tiff.TIFF) (image.Config, error)
	CanHandle(*tiff.TIFF) bool
}

var (
	// Handlers based on the value in the Make tag (tag id 271)
	altHandlersFromMakeTagValue = struct {
		mu            sync.RWMutex
		makeToHandler map[string]Handler
	}{
		makeToHandler: make(map[string]Handler, 1),
	}

	// Handlers based on the presence of a tag id in IFD0.
	altHandlersFromTagPresence = struct {
		mu           sync.RWMutex
		tagToHandler map[uint16]Handler
	}{
		tagToHandler: make(map[uint16]Handler, 1),
	}
)

// For now, findAlternates will only check against two concepts.  One, is the
// "Make" tag and the other checks the presence of tags.
func findAlternates(t *tiff.TIFF) Handler {
	ifd0 := t.IFDs[0]
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

func RegisterHandlerByMake(m string, h Handler) {
	altHandlersFromMakeTagValue.mu.Lock()
	defer altHandlersFromMakeTagValue.mu.Unlock()
	altHandlersFromMakeTagValue.makeToHandler[m] = h
}

func GetHandlerByMake(m string) Handler {
	altHandlersFromMakeTagValue.mu.RLock()
	defer altHandlersFromMakeTagValue.mu.RUnlock()
	return altHandlersFromMakeTagValue.makeToHandler[m]
}

func RegisterHandlerByTagPresence(t uint16, h Handler) {
	altHandlersFromTagPresence.mu.Lock()
	defer altHandlersFromTagPresence.mu.Unlock()
	altHandlersFromTagPresence.tagToHandler[t] = h
}

func GetHandlerByTagPresence(t uint16) Handler {
	altHandlersFromTagPresence.mu.RLock()
	defer altHandlersFromTagPresence.mu.RUnlock()
	return altHandlersFromTagPresence.tagToHandler[t]
}

func ListRegisteredTagPresenceIDs() []uint16 {
	altHandlersFromTagPresence.mu.RLock()
	defer altHandlersFromTagPresence.mu.RUnlock()
	ids := make([]uint16, 0, len(altHandlersFromTagPresence.tagToHandler))
	for k := range altHandlersFromTagPresence.tagToHandler {
		ids = append(ids, k)
	}
	sort.Sort(uint16Slice(ids))
	return ids
}
