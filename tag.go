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

package tiff // import "jonathanpittman.com/tiff"

import (
	"fmt"
)

type Tag interface {
	ID() uint16
	Name() string
	Interpreter() FieldInterpreter
	//ValidFieldTypes() []FieldType
}

func NewTag(id uint16, name string, fi FieldInterpreter) Tag {
	return &tag{id: id, name: name, fi: fi}
}

type tag struct {
	id   uint16
	name string
	fi   FieldInterpreter
}

func (t *tag) ID() uint16 {
	return t.id
}

func (t *tag) Name() string {
	if len(t.name) > 0 {
		return t.name
	}
	return fmt.Sprintf("UNNAMED_TAG_%d", t.id)
}

func (t *tag) Interpreter() FieldInterpreter {
	if t.fi == nil {
		return defaultFieldInterpreter
	}
	return t.fi
}

type FieldInterpreter func(Field) string

func defaultFieldInterpreter(f Field) string {
	return ""
}

/*
type FieldInterpreter interface {
	Interpret(Field) string
	Register(reg interface{}) bool
	// Reg, in this case, could be anything, but usually a struct or map with a key/value pair.
	// The idea is that the FieldInterpreter should know what to do with the thing or return false.
	//Register(key, val reflect.Value) bool
}
*/
