// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image // import "jonathanpittman.com/tiff/image"

import (
	"fmt"
	"sync"
)

type CompressionError struct {
	Method  string
	Message string
}

func (ce CompressionError) Error() string {
	return fmt.Sprintf("tiff/image: compression: %s: %s", ce.Method, ce.Message)
}

type CompressionNotSupported struct {
	Method uint16
}

func (cni CompressionNotSupported) Error() string {
	return fmt.Sprintf("tiff/image: compression: unsupported type %d", cni.Method)
}

type Compression interface {
	ID() uint16
	Name() string
	Compress([]byte) ([]byte, error)
	Decompress([]byte) ([]byte, error)
}

func NewCompression(id uint16, name string, comp, decomp func([]byte) ([]byte, error)) Compression {
	return &compression{
		id:         id,
		name:       name,
		compress:   comp,
		decompress: decomp,
	}
}

type compression struct {
	id         uint16
	name       string
	compress   func([]byte) ([]byte, error)
	decompress func([]byte) ([]byte, error)
}

func (c *compression) ID() uint16 {
	return c.id
}

func (c *compression) Name() string {
	return c.name
}

func (c *compression) Compress(in []byte) ([]byte, error) {
	return c.compress(in)
}

func (c *compression) Decompress(in []byte) ([]byte, error) {
	return c.decompress(in)
}

/* Uncompressed */

// compUncompressed is the function representing both halves of the Uncompressed
// compression method.
func compUncompressed(in []byte) ([]byte, error) {
	return in, nil
}

/* PackBits Compression and Decompression */
// decompPackBits decompresses a
func decompPackBits(in []byte) ([]byte, error) {
	buf := in[:]
	out := make([]byte, 0, 1024)
	for len(buf) > 0 {
		n := int(int8(buf[0]))
		buf = buf[1:]
		switch {
		case n >= 0:
			if len(buf) < n+1 {
				return nil, CompressionError{"PackBits", "not enough bytes to complete decompression"}
			}
			out = append(out, buf[:n+1]...)
			buf = buf[n+1:]
		case n == -128:
			// NOOP: Skip this byte.
			continue
		default:
			if len(buf) == 0 {
				return nil, CompressionError{"PackBits", "not enough bytes to complete decompression"}
			} else {
				tmp := make([]byte, -n+1)
				b := buf[0]
				for i := range tmp {
					tmp[i] = b
				}
				out = append(out, tmp...)
			}
			buf = buf[1:]
		}
	}
	return out, nil
}

func compPackBits(in []byte) ([]byte, error) {
	return nil, CompressionError{"PackBits", "compressing not implemented."}
}

var (
	uncompressedCompression = &compression{
		id:         1,
		name:       "Uncompressed",
		compress:   compUncompressed,
		decompress: compUncompressed,
	}

	packbitsCompression = &compression{
		id:         32773,
		name:       "PackBits",
		compress:   compPackBits,
		decompress: decompPackBits,
	}
)

/* Registration pieces for Compression methods */
var allCompressions = struct {
	mu   sync.RWMutex
	list map[uint16]Compression
}{
	list: make(map[uint16]Compression, 1),
}

func RegisterCompression(c Compression) {
	allCompressions.mu.Lock()
	allCompressions.list[c.ID()] = c
	allCompressions.mu.Unlock()
}

func GetCompression(id uint16) Compression {
	allCompressions.mu.RLock()
	defer allCompressions.mu.RUnlock()
	return allCompressions.list[id]
}

func init() {
	RegisterCompression(uncompressedCompression)
	RegisterCompression(packbitsCompression)
}
