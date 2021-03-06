/*
   Copyright The containerd Authors.

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

package typeurl

import (
	"reflect"
	"testing"
)

type test struct {
	Name string
	Age  int
}

func clear() {
	registry = make(map[reflect.Type]string)
}

func TestRegisterPointerGetPointer(t *testing.T) {
	clear()
	expected := "test"
	Register(&test{}, "test")

	url, err := TypeURL(&test{})
	if err != nil {
		t.Fatal(err)
	}
	if url != expected {
		t.Fatalf("expected %q but received %q", expected, url)
	}
}

func TestMarshal(t *testing.T) {
	clear()
	expected := "test"
	Register(&test{}, "test")

	v := &test{
		Name: "koye",
		Age:  6,
	}
	any, err := MarshalAny(v)
	if err != nil {
		t.Fatal(err)
	}
	if any.TypeUrl != expected {
		t.Fatalf("expected %q but received %q", expected, any.TypeUrl)
	}

	// marshal it again and make sure we get the same thing back.
	newany, err := MarshalAny(any)
	if err != nil {
		t.Fatal(err)
	}

	if newany != any { // you that right: we want the same *pointer*!
		t.Fatalf("expected to get back same object: %v != %v", newany, any)
	}

}

func TestMarshalUnmarshal(t *testing.T) {
	clear()
	Register(&test{}, "test")

	v := &test{
		Name: "koye",
		Age:  6,
	}
	any, err := MarshalAny(v)
	if err != nil {
		t.Fatal(err)
	}
	nv, err := UnmarshalAny(any)
	if err != nil {
		t.Fatal(err)
	}
	td, ok := nv.(*test)
	if !ok {
		t.Fatal("expected value to cast to *test")
	}
	if td.Name != "koye" {
		t.Fatal("invalid name")
	}
	if td.Age != 6 {
		t.Fatal("invalid age")
	}
}

func TestIs(t *testing.T) {
	clear()
	Register(&test{}, "test")

	v := &test{
		Name: "koye",
		Age:  6,
	}
	any, err := MarshalAny(v)
	if err != nil {
		t.Fatal(err)
	}
	if !Is(any, &test{}) {
		t.Fatal("Is(any, test{}) should be true")
	}
}

func TestRegisterDiffUrls(t *testing.T) {
	clear()
	defer func() {
		if err := recover(); err == nil {
			t.Error("registering the same type with different urls should panic")
		}
	}()
	Register(&test{}, "test")
	Register(&test{}, "test", "two")
}
