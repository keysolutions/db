// Copyright (c) 2012-present The upper.io/db authors. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package sqlbuilder

import (
	"database/sql"
	"fmt"
	"reflect"

	"upper.io/db.v3"
)

type scanner struct {
	v db.Unmarshaler
}

func (u scanner) Scan(v interface{}) error {
	return u.v.UnmarshalDB(v)
}

type nullableScanner struct {
	scanner func(v interface{}) error
}

func newNullableScanner(v interface{}) (nullableScanner, error) {
	if s, ok := v.(sql.Scanner); ok {
		return nullableScanner{
			scanner: s.Scan,
		}, nil
	}

	v1 := reflect.Indirect(reflect.ValueOf(v))
	if !v1.CanSet() {
		return nullableScanner{}, fmt.Errorf("%v must be assignable", v1)
	}
	return nullableScanner{
		scanner: func(v interface{}) error {
			v2 := reflect.Indirect(reflect.ValueOf(v))
			if v1.Type() != v2.Type() {
				return fmt.Errorf("unable to convert from %v to %v", v1, v2)
			}
			v1.Set(v2)
			return nil
		},
	}, nil
}

func (n nullableScanner) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	return n.scanner(v)
}

var _ sql.Scanner = scanner{}
var _ sql.Scanner = nullableScanner{}
