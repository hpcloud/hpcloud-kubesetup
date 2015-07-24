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

package misc_test

import (
	"encoding/json"
	"testing"
	"time"

	"git.openstack.org/stackforge/golang-client.git/misc"
	"git.openstack.org/stackforge/golang-client.git/testUtil"
)

var milliSecondsTestValue = `{"created_at":"2015-01-26T22:47:27.851022"}`
var milliSecondsTestTime, _ = time.Parse(`"2006-01-02T15:04:05.999999"`, `"2015-01-26T22:47:27.851022"`)
var milliSecondsTimeTestValue = timeTest{CreatedAt: misc.NewDateTime(milliSecondsTestTime)}

var testValue = `{"created_at":"2014-09-29T14:44:31"}`
var testTime, _ = time.Parse(`"2006-01-02T15:04:05"`, `"2014-09-29T14:44:31"`)
var timeTestValue = timeTest{CreatedAt: misc.NewDateTime(testTime)}
var duration250ms, _ = time.ParseDuration("250ms")

func TestMarshalTimeTestMilliseconds(t *testing.T) {
	bytes, _ := json.Marshal(milliSecondsTimeTestValue)

	testUtil.Equals(t, milliSecondsTestValue, string(bytes))
}

func TestUnmarshalValidTimeTestMilliseconds(t *testing.T) {
	val := timeTest{}
	err := json.Unmarshal([]byte(milliSecondsTestValue), &val)
	testUtil.IsNil(t, err)
	testUtil.Equals(t, milliSecondsTimeTestValue.CreatedAt, val.CreatedAt)
}

func TestMarshalTimeTest(t *testing.T) {
	bytes, _ := json.Marshal(timeTestValue)

	testUtil.Equals(t, testValue, string(bytes))
}

func TestUnmarshalValidTimeTest(t *testing.T) {
	val := timeTest{}
	err := json.Unmarshal([]byte(testValue), &val)
	testUtil.IsNil(t, err)
	testUtil.Equals(t, timeTestValue.CreatedAt, val.CreatedAt)
}

func TestUnmarshalInvalidDataFormatTimeTest(t *testing.T) {
	val := timeTest{}
	err := json.Unmarshal([]byte("something other than date time"), &val)
	testUtil.Assert(t, err != nil, "expected an error")
}

// Added this test to ensure that its understood that
// only one specific format is supported at this time.
func TestUnmarshalInvalidDateTimeFormatTimeTest(t *testing.T) {
	val := timeTest{}
	err := json.Unmarshal([]byte("2014-09-29T14:44"), &val)
	testUtil.Assert(t, err != nil, "expected an error")
}

func TestBeforeThanIsTrue(t *testing.T) {
	newTime := testTime.Add(duration250ms)
	isBefore := misc.NewDateTime(testTime).Before(misc.NewDateTime(newTime))
	testUtil.Assert(t, isBefore, "expected value to be before")
}

func TestAfterThanIsTrue(t *testing.T) {
	timeAfterOriginal := testTime.Add(duration250ms)
	isAfter := misc.NewDateTime(timeAfterOriginal).After(misc.NewDateTime(testTime))
	testUtil.Assert(t, isAfter, "expected value to be after")
}

func TestEqualIsTrue(t *testing.T) {
	isEqual := misc.NewDateTime(testTime).Equal(misc.NewDateTime(testTime))
	testUtil.Assert(t, isEqual, "expected values to be equal")
}

type timeTest struct {
	CreatedAt misc.RFC8601DateTime `json:"created_at"`
}
