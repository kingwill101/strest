// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2010 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

/*
Package proto converts data structures to and from the wire format of
protocol buffers.  It works in concert with the Go source code generated
for .proto files by the protocol compiler.

A summary of the properties of the protocol buffer interface
for a protocol buffer variable v:

  - Names are turned from camel_case to CamelCase for export.
  - There are no methods on v to set fields; just treat
	them as structure fields.
  - There are getters that return a field's value if set,
	and return the field's default value if unset.
	The getters work even if the receiver is a nil message.
  - The zero value for a struct is its correct initialization state.
	All desired fields must be set before marshaling.
  - A Reset() method will restore a protobuf struct to its zero state.
  - Non-repeated fields are pointers to the values; nil means unset.
	That is, optional or required field int32 f becomes F *int32.
  - Repeated fields are slices.
  - Helper functions are available to aid the setting of fields.
	msg.Foo = proto.String("hello") // set field
  - Constants are defined to hold the default values of all fields that
	have them.  They have the form Default_StructName_FieldName.
	Because the getter methods handle defaulted values,
	direct use of these constants should be rare.
  - Enums are given type names and maps from names to values.
	Enum values are prefixed by the enclosing message's name, or by the
	enum's type name if it is a top-level enum. Enum types have a String
	method, and a Enum method to assist in message construction.
  - Nested messages, groups and enums have type names prefixed with the name of
	the surrounding message type.
  - Extensions are given descriptor names that start with E_,
	followed by an underscore-delimited list of the nested messages
	that contain it (if any) followed by the CamelCased name of the
	extension field itself.  HasExtension, ClearExtension, GetExtension
	and SetExtension are functions for manipulating extensions.
  - Oneof field sets are given a single field in their message,
	with distinguished wrapper types for each possible field value.
  - Marshal and Unmarshal are functions to encode and decode the wire format.

When the .proto file specifies `syntax="proto3"`, there are some differences:

  - Non-repeated fields of non-message type are values instead of pointers.
  - Enum types do not get an Enum method.

The simplest way to describe this is to see an example.
Given file test.proto, containing

	package example;

	enum FOO { X = 17; }

	message Test {
	  required string label = 1;
	  optional int32 type = 2 [default=77];
	  repeated int64 reps = 3;
	  optional group OptionalGroup = 4 {
	    required string RequiredField = 5;
	  }
	  oneof union {
	    int32 number = 6;
	    string name = 7;
	  }
	}

The resulting file, test.pb.go, is:

	package example

	import proto "github.com/golang/protobuf/proto"
	import math "math"

	type FOO int32
	const (
		FOO_X FOO = 17
	)
	var FOO_name = map[int32]string{
		17: "X",
	}
	var FOO_value = map[string]int32{
		"X": 17,
	}

	func (x FOO) Enum() *FOO {
		p := new(FOO)
		*p = x
		return p
	}
	func (x FOO) String() string {
		return proto.EnumName(FOO_name, int32(x))
	}
	func (x *FOO) UnmarshalJSON(data []byte) error {
		value, err := proto.UnmarshalJSONEnum(FOO_value, data)
		if err != nil {
			return err
		}
		*x = FOO(value)
		return nil
	}

	type Test struct {
		Label         *string             `protobuf:"bytes,1,req,name=label" json:"label,omitempty"`
		Type          *int32              `protobuf:"varint,2,opt,name=type,def=77" json:"type,omitempty"`
		Reps          []int64             `protobuf:"varint,3,rep,name=reps" json:"reps,omitempty"`
		Optionalgroup *Test_OptionalGroup `protobuf:"group,4,opt,name=OptionalGroup" json:"optionalgroup,omitempty"`
		// Types that are valid to be assigned to Union:
		//	*Test_Number
		//	*Test_Name
		Union            isTest_Union `protobuf_oneof:"union"`
		XXX_unrecognized []byte       `json:"-"`
	}
	func (m *Test) Reset()         { *m = Test{} }
	func (m *Test) String() string { return proto.CompactTextString(m) }
	func (*Test) ProtoMessage() {}

	type isTest_Union interface {
		isTest_Union()
	}

	type Test_Number struct {
		Number int32 `protobuf:"varint,6,opt,name=number"`
	}
	type Test_Name struct {
		Name string `protobuf:"bytes,7,opt,name=name"`
	}

	func (*Test_Number) isTest_Union() {}
	func (*Test_Name) isTest_Union()   {}

	func (m *Test) GetUnion() isTest_Union {
		if m != nil {
			return m.Union
		}
		return nil
	}
	const Default_Test_Type int32 = 77

	func (m *Test) GetLabel() string {
		if m != nil && m.Label != nil {
			return *m.Label
		}
		return ""
	}

	func (m *Test) GetType() int32 {
		if m != nil && m.Type != nil {
			return *m.Type
		}
		return Default_Test_Type
	}

	func (m *Test) GetOptionalgroup() *Test_OptionalGroup {
		if m != nil {
			return m.Optionalgroup
		}
		return nil
	}

	type Test_OptionalGroup struct {
		RequiredField *string `protobuf:"bytes,5,req" json:"RequiredField,omitempty"`
	}
	func (m *Test_OptionalGroup) Reset()         { *m = Test_OptionalGroup{} }
	func (m *Test_OptionalGroup) String() string { return proto.CompactTextString(m) }

	func (m *Test_OptionalGroup) GetRequiredField() string {
		if m != nil && m.RequiredField != nil {
			return *m.RequiredField
		}
		return ""
	}

	func (m *Test) GetNumber() int32 {
		if x, ok := m.GetUnion().(*Test_Number); ok {
			return x.Number
		}
		return 0
	}

	func (m *Test) GetName() string {
		if x, ok := m.GetUnion().(*Test_Name); ok {
			return x.Name
		}
		return ""
	}

	func init() {
		proto.RegisterEnum("example.FOO", FOO_name, FOO_value)
	}

To create and play with a Test object:

	package main

	import (
		"log"

		"github.com/golang/protobuf/proto"
		pb "./example.pb"
	)

	func main() {
		test := &pb.Test{
			Label: proto.String("hello"),
			Type:  proto.Int32(17),
			Reps:  []int64{1, 2, 3},
			Optionalgroup: &pb.Test_OptionalGroup{
				RequiredField: proto.String("good bye"),
			},
			Union: &pb.Test_Name{"fred"},
		}
		data, err := proto.Marshal(test)
		if err != nil {
			log.Fatal("marshaling error: ", err)
		}
		newTest := &pb.Test{}
		err = proto.Unmarshal(data, newTest)
		if err != nil {
			log.Fatal("unmarshaling error: ", err)
		}
		// Now test and newTest contain the same data.
		if test.GetLabel() != newTest.GetLabel() {
			log.Fatalf("data mismatch %q != %q", test.GetLabel(), newTest.GetLabel())
		}
		// Use a type switch to determine which oneof was set.
		switch u := test.Union.(type) {
		case *pb.Test_Number: // u.Number contains the number.
		case *pb.Test_Name: // u.Name contains the string.
		}
		// etc.
	}
*/
package proto

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"sync"
)

// RequiredNotSetError is an error type returned by either Marshal or Unmarshal.
// Marshal reports this when a required field is not initialized.
// Unmarshal reports this when a required field is missing from the wire data.
type RequiredNotSetError struct{ field string }

func (e *RequiredNotSetError) Error() string {
	if e.field == "" {
		return fmt.Sprintf("proto: required field not set")
	}
	return fmt.Sprintf("proto: required field %q not set", e.field)
}
func (e *RequiredNotSetError) RequiredNotSet() bool {
	return true
}

type invalidUTF8Error struct{ field string }

func (e *invalidUTF8Error) Error() string {
	if e.field == "" {
		return "proto: invalid UTF-8 detected"
	}
	return fmt.Sprintf("proto: field %q contains invalid UTF-8", e.field)
}
func (e *invalidUTF8Error) InvalidUTF8() bool {
	return true
}

// errInvalidUTF8 is a sentinel error to identify fields with invalid UTF-8.
// This error should not be exposed to the external API as such errors should
// be recreated with the field information.
var errInvalidUTF8 = &invalidUTF8Error{}

// isNonFatal reports whether the error is either a RequiredNotSet error
// or a InvalidUTF8 error.
func isNonFatal(err error) bool {
	if re, ok := err.(interface{ RequiredNotSet() bool }); ok && re.RequiredNotSet() {
		return true
	}
	if re, ok := err.(interface{ InvalidUTF8() bool }); ok && re.InvalidUTF8() {
		return true
	}
	return false
}

type nonFatal struct{ E error }

// Merge merges err into nf and reports whether it was successful.
// Otherwise it returns false for any fatal non-nil errors.
func (nf *nonFatal) Merge(err error) (ok bool) {
	if err == nil {
		return true // not an error
	}
	if !isNonFatal(err) {
		return false // fatal error
	}
	if nf.E == nil {
		nf.E = err // store first instance of non-fatal error
	}
	return true
}

// Message is implemented by generated protocol buffer messages.
type Message interface {
	Reset()
	String() string
	ProtoMessage()
}

// Stats records allocation details ab