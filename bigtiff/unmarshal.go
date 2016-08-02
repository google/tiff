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

package bigtiff // import "jonathanpittman.com/tiff/bigtiff"

import (
	"fmt"
	"log"
	"reflect"

	"jonathanpittman.com/tiff"
)

func UnmarshalSubIFDs(ifd tiff.IFD, br tiff.BReader, tsp tiff.TagSpace, out interface{}) error {
	if br == nil {
		return fmt.Errorf("bigtiff: UnmarshalSubIFDs: no BReader available")
	}
	if len(ifd.Fields()) == 0 {
		return fmt.Errorf("bigtiff: UnmarshalSubIFDs: ifd has no fields")
	}
	v := reflect.ValueOf(out).Elem()
	structType := v.Type()

	for i := 0; i < v.NumField(); i++ {
		stField := structType.Field(i)
		sTag := tiff.ParseTiffStructTag(stField.Tag.Get("tiff"))
		if sTag == nil || sTag.Type != "subifd" {
			continue
		}
		vf := v.Field(i)
		vft := vf.Type()
		vftk := vft.Kind()

		if vftk != reflect.Ptr && vftk != reflect.Struct {
			log.Printf("bigtiff: UnmarshalSubIFDs: using a tiff SubIFD struct tag is only supported for structs (not a %v)\n", vftk)
			continue
		}

		siTag := tiff.ParseTiffSubIFDStructTag(sTag.Data)
		if siTag == nil {
			log.Printf("bigtiff: UnmarshalSubIFDs: skipping struct field %q due to malformed tiff subifd struct tag (%q).\n", stField.Name, sTag.Data)
			continue
		}
		if siTag.Tag == nil {
			log.Printf("bigtiff: UnmarshalSubIFDs: skipping struct field %q due to missing \"tag\" key in tiff subifd struct tag.\n", stField.Name)
			continue
		}
		if !ifd.HasField(*siTag.Tag) {
			// Default values do not work here.  We have to simply skip it.
			continue
		}
		ifdField := ifd.GetField(*siTag.Tag)
		ifdFT := ifdField.Type()
		fvBytes := ifdField.Value().Bytes()
		fvBo := ifdField.Value().Order()
		if ifdField.Count() == 0 {
			continue
		}

		var offsets []uint64
		switch ifdFT.ReflectType().Kind() {
		case reflect.Uint32:
			offsets = make([]uint64, ifdField.Count())
			for i := range offsets {
				offsets[i] = uint64(fvBo.Uint32(fvBytes))
				fvBytes = fvBytes[4:]
			}
		case reflect.Uint64:
			offsets = make([]uint64, ifdField.Count())
			for i := range offsets {
				offsets[i] = fvBo.Uint64(fvBytes)
				fvBytes = fvBytes[8:]
			}
		default:
			continue
		}

		off := offsets[0]
		if siTag.Index != nil {
			if *siTag.Index >= len(offsets) || *siTag.Index < 0 {
				// log a warning
				continue
			}
			off = offsets[*siTag.Index]
		}
		if tsp == nil {
			if siTag.TagSpace != nil {
				newSpace := tiff.GetTagSpace(*siTag.TagSpace)
				if newSpace != nil {
					tsp = newSpace
				}
			}
			if tsp == nil {
				tsp = tiff.DefaultTagSpace
			}
		}

		subIFD, err := ParseIFD(br, uint64(off), tsp, nil)
		if err != nil {
			return err
		}
		switch vftk {
		case reflect.Ptr:
			// We do not support recursive unmarshaling when
			// the field points back to the enclosing struct.
			if vf.Elem() != v && vft.Elem().Kind() == reflect.Struct {
				newStruct := reflect.New(vft.Elem())
				if err := tiff.UnmarshalIFD(subIFD, newStruct.Interface()); err != nil {
					return err
				}
				vf.Set(newStruct)
			}
		case reflect.Struct:
			embStructPtr := v.Field(i).Addr().Interface()
			if err := tiff.UnmarshalIFD(subIFD, embStructPtr); err != nil {
				return err
			}
		}
	}
	return nil
}
