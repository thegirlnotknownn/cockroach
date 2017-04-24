// Copyright 2017 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package settings_test

import (
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/testutils"
)

const mb = int64(1024 * 1024)

var boolTA = settings.RegisterBoolSetting("bool.t", "", true)
var boolFA = settings.RegisterBoolSetting("bool.f", "", false)
var strFooA = settings.RegisterStringSetting("str.foo", "", "")
var strBarA = settings.RegisterStringSetting("str.bar", "", "bar")
var i1A = settings.RegisterIntSetting("i.1", "", 0)
var i2A = settings.RegisterIntSetting("i.2", "", 5)
var fA = settings.RegisterFloatSetting("f", "", 5.4)
var dA = settings.RegisterDurationSetting("d", "", time.Second)
var eA = settings.RegisterEnumSetting("e", "", "foo", map[int64]string{1: "foo", 2: "bar", 3: "baz"})
var byteSize = settings.RegisterByteSizeSetting("zzz", "", mb)
var _ = settings.RegisterBoolSetting("sekretz", "", false)

func init() {
	settings.Hide("sekretz")
}

func TestCache(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		if expected, actual := false, boolFA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, boolTA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "", strFooA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "bar", strBarA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := int64(0), i1A.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := int64(5), i2A.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := 5.4, fA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := time.Second, dA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := mb, byteSize.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := int64(1), eA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("lookup", func(t *testing.T) {
		if actual, _, ok := settings.Lookup("i.1"); !ok || i1A != actual {
			t.Fatalf("expected %v, got %v (exists: %v)", i1A, actual, ok)
		}
		if actual, _, ok := settings.Lookup("f"); !ok || fA != actual {
			t.Fatalf("expected %v, got %v (exists: %v)", fA, actual, ok)
		}
		if actual, _, ok := settings.Lookup("d"); !ok || dA != actual {
			t.Fatalf("expected %v, got %v (exists: %v)", dA, actual, ok)
		}
		if actual, _, ok := settings.Lookup("e"); !ok || eA != actual {
			t.Fatalf("expected %v, got %v (exists: %v)", eA, actual, ok)
		}
		if actual, _, ok := settings.Lookup("dne"); ok {
			t.Fatalf("expected nothing, got %v", actual)
		}
	})

	t.Run("read and write each type", func(t *testing.T) {
		u := settings.MakeUpdater()
		if err := u.Set("bool.t", settings.EncodeBool(false), "b"); err != nil {
			t.Fatal(err)
		}
		if err := u.Set("bool.f", settings.EncodeBool(true), "b"); err != nil {
			t.Fatal(err)
		}
		if err := u.Set("str.foo", "baz", "s"); err != nil {
			t.Fatal(err)
		}
		if err := u.Set("i.2", settings.EncodeInt(3), "i"); err != nil {
			t.Fatal(err)
		}
		if err := u.Set("f", settings.EncodeFloat(3.1), "f"); err != nil {
			t.Fatal(err)
		}
		if err := u.Set("d", settings.EncodeDuration(2*time.Hour), "d"); err != nil {
			t.Fatal(err)
		}
		if err := u.Set("zzz", settings.EncodeInt(mb*5), "z"); err != nil {
			t.Fatal(err)
		}
		if err := u.Set("e", settings.EncodeInt(2), "e"); err != nil {
			t.Fatal(err)
		}
		if expected, err := "strconv.Atoi: parsing \"notAValidValue\": invalid syntax",
			u.Set("e", "notAValidValue", "e"); !testutils.IsError(err, expected) {
			t.Fatalf("expected '%s' != actual error '%s'", expected, err)
		}
		u.Done()

		if expected, actual := false, boolTA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := true, boolFA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := "baz", strFooA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := int64(3), i2A.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := 3.1, fA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := 2*time.Hour, dA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := int64(2), eA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		if expected, actual := mb*5, byteSize.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}

		// We didn't change this one, so should still see the default.
		if expected, actual := "bar", strBarA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("any setting not included in an Updater reverts to default", func(t *testing.T) {
		{
			u := settings.MakeUpdater()
			if err := u.Set("bool.f", settings.EncodeBool(true), "b"); err != nil {
				t.Fatal(err)
			}
			if err := u.Set("i.1", settings.EncodeInt(1), "i"); err != nil {
				t.Fatal(err)
			}
			if err := u.Set("i.2", settings.EncodeInt(7), "i"); err != nil {
				t.Fatal(err)
			}
			u.Done()
		}

		if expected, actual := true, boolFA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
		// If the updater doesn't have a key, e.g. if the setting has been deleted,
		// Doneing it from the cache.
		settings.MakeUpdater().Done()

		if expected, actual := false, boolFA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}

		if expected, actual := false, boolFA.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("an invalid update to a given setting preserves its previously set value", func(t *testing.T) {
		{
			u := settings.MakeUpdater()
			if err := u.Set("i.2", settings.EncodeInt(9), "i"); err != nil {
				t.Fatal(err)
			}
			u.Done()
		}
		before := i2A.Get()

		// Doneing after attempting to set with wrong type preserves current value.
		{
			u := settings.MakeUpdater()
			// We don't use testutils.IsError, to avoid the import.
			if err := u.Set("i.2", settings.EncodeBool(false), "b"); !testutils.IsError(err,
				"setting 'i.2' defined as type i, not b",
			) {
				t.Fatal(err)
			}
			u.Done()
		}

		if expected, actual := before, i2A.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}

		// Doneing after attempting to set with invalid format does too.
		{
			u := settings.MakeUpdater()
			if err := u.Set("i.2", settings.EncodeBool(false), "i"); !testutils.IsError(err,
				"strconv.Atoi: parsing \"false\": invalid syntax",
			) {
				t.Fatal(err)
			}
			u.Done()
		}

		if expected, actual := before, i2A.Get(); expected != actual {
			t.Fatalf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("mocks", func(t *testing.T) {
		{
			f := settings.TestingSetBool(&boolFA, true)
			if expected, actual := true, boolFA.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
			f()
			if expected, actual := false, boolFA.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
		}

		{
			f := settings.TestingSetString(&strBarA, "override")
			if expected, actual := "override", strBarA.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
			f()
			if expected, actual := "bar", strBarA.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
		}

		{
			f := settings.TestingSetInt(&i1A, 64)
			if expected, actual := int64(64), i1A.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
			f()
			if expected, actual := int64(0), i1A.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
		}

		{
			f := settings.TestingSetFloat(&fA, 6.7)
			if expected, actual := 6.7, fA.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
			f()
			if expected, actual := 5.4, fA.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
		}

		{
			f := settings.TestingSetDuration(&dA, 10*time.Hour)
			if expected, actual := 10*time.Hour, dA.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
			f()
			if expected, actual := time.Second, dA.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
		}

		{
			f := settings.TestingSetEnum(&eA, 3)
			if expected, actual := int64(3), eA.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
			f()
			if expected, actual := int64(1), eA.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
		}

		{
			f := settings.TestingSetByteSize(&byteSize, mb*7)
			if expected, actual := mb*7, byteSize.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
			f()
			if expected, actual := mb, byteSize.Get(); expected != actual {
				t.Fatalf("expected %v, got %v", expected, actual)
			}
		}
	})
}

func TestHide(t *testing.T) {
	keys := make(map[string]struct{})
	for _, k := range settings.Keys() {
		keys[k] = struct{}{}
	}
	if _, ok := keys["bool.t"]; !ok {
		t.Errorf("expected 'bool.t' to be unhidden")
	}
	if _, ok := keys["sekretz"]; ok {
		t.Errorf("expected 'sekretz' to be hidden")
	}
}
