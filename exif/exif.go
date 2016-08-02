/*
Copyright 2016 Google Inc.

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

package exif // import "jonathanpittman.com/tiff/exif"

import (
	"fmt"
	"io"
	"log"

	"jonathanpittman.com/tiff"
)

func Parse(r io.Reader) (eIFD, gIFD, ioIFD tiff.IFD, err error) {
	rars := tiff.NewReadAtReadSeeker(r)
	var two [2]byte
	if _, err = rars.Read(two[:]); err != nil {
		return
	}

	if _, err = rars.Seek(0, 0); err != nil {
		return
	}

	switch string(two[:]) {
	case "MM", "II": // likely a tiff
		var t tiff.TIFF
		if t, err = tiff.Parse(rars, nil, nil); err != nil {
			return
		}
		for _, tIFD := range t.IFDs() {
			if tIFD.HasField(ExifIFDTagID) {
				eFld := tIFD.GetField(ExifIFDTagID)
				offset := eFld.Type().Valuer()(eFld.Value().Bytes(), eFld.Value().Order()).Uint()
				if eIFD, err = tiff.ParseIFD(t.R(), offset, ExifTagSpace, nil); err != nil {
					return
				}
				if tIFD.HasField(GPSIFDTagID) {
					gFld := tIFD.GetField(GPSIFDTagID)
					offset = gFld.Type().Valuer()(gFld.Value().Bytes(), gFld.Value().Order()).Uint()
					if gIFD, err = tiff.ParseIFD(t.R(), offset, GPSTagSpace, nil); err != nil {
						log.Printf("exif: GPS IFD found, but had trouble retrieving it from offset %d: %v\n", offset, err)
					}
				}
				if tIFD.HasField(InteroperabilityIFDTagID) {
					ioFld := tIFD.GetField(InteroperabilityIFDTagID)
					offset = ioFld.Type().Valuer()(ioFld.Value().Bytes(), ioFld.Value().Order()).Uint()
					if ioIFD, err = tiff.ParseIFD(t.R(), offset, IOPTagSpace, nil); err != nil {
						log.Printf("exif: IOP IFD found, but had trouble retrieving it from offset %d: %v\n", offset, err)
					}
				}
				return
			}
		}
		err = fmt.Errorf("exif: no exif ifd found in tiff")
		return
	case "\xff\xd8": // likely a jpeg
		err = fmt.Errorf("exif: still working on jpeg support")
		return
	}
	// Anything else is currently unsupported.
	err = fmt.Errorf("exif: unsupported header: %q", two[:])
	return
}
