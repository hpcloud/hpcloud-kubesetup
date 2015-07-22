// Copyright (c) 2014 Hewlett-Packard Development Company, L.P.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package misc

import "time"

// RFC8601DateTime is a type for decoding and encoding json
// date times that follow RFC 8601 format. The type currently
// decodes and encodes with exactly precision to seconds. If more
// formats of RFC8601 need to be supported additional work
// will be needed.
type RFC8601DateTime struct {
	time time.Time
}

// NewDateTime creates a new RFC8601DateTime from a time.Time
func NewDateTime(input time.Time) RFC8601DateTime {
	return RFC8601DateTime{time: input}
}

// NewDateTimeFromString creates a new RFC8601DateTime taking a string as input.
// It must follow the "2006-01-02T15:04:05" pattern.
func NewDateTimeFromString(input string) (val RFC8601DateTime, err error) {
	val = RFC8601DateTime{}
	err = val.parseValue(input)
	return val, err
}

// UnmarshalJSON converts the bytes give to a RFC8601DateTime
// Errors will occur if the bytes when converted to a string
// don't match the format "2006-01-02T15:04:05".
func (r *RFC8601DateTime) UnmarshalJSON(data []byte) error {
	return r.parseValue(string(data))
}

// MarshalJSON converts a RFC8601DateTime to a []byte.
func (r RFC8601DateTime) MarshalJSON() ([]byte, error) {
	var val string
	if r.time.Nanosecond() > 0 {
		val = r.time.Format(format2)
	} else {
		val = r.time.Format(format)
	}
	return []byte(val), nil
}

// Time will return the embedded Time instance
func (r *RFC8601DateTime) Time() time.Time {
	return r.time
}

// Before compares the DateTime to the other value and
// returns true if its before than the other value or false if not.
func (r RFC8601DateTime) Before(other RFC8601DateTime) bool {
	return r.time.Before(other.time)
}

// After compares the DateTime to the other value and
// returns true if its after than the other value or false if not.
func (r RFC8601DateTime) After(other RFC8601DateTime) bool {
	return r.time.After(other.time)
}

// Equal compares the DateTime to the other value and
// returns true if its equal to the other value or false if not.
func (r RFC8601DateTime) Equal(other RFC8601DateTime) bool {
	return r.time.Equal(other.time)
}

func (r *RFC8601DateTime) parseValue(input string) (err error) {
	numTimeStr := len(input)
	var timeVal time.Time

	switch numTimeStr {
	case len(format4):
		timeVal, err = time.Parse(format4, input)
	case len(format3):
		timeVal, err = time.Parse(format3, input)
	case len(format2):
		timeVal, err = time.Parse(format2, input)
	default:
		timeVal, err = time.Parse(format, input)
	}

	r.time = timeVal
	return err
}

const format = `"2006-01-02T15:04:05"`
const format2 = `"2006-01-02T15:04:05.999999"`
const format3 = `2006-01-02T15:04:05`
const format4 = `2006-01-02T15:04:05.999999`
