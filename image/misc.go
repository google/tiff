// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image // import "jonathanpittman.com/tiff/image"

// uint16Slice implements the sort.Sorter interface for a []uint16.
type uint16Slice []uint16

func (p uint16Slice) Len() int           { return len(p) }
func (p uint16Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p uint16Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
