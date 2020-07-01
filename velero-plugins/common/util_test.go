// tests for util.go

package common

import (
	"testing"
	"errors"
	"reflect"
)


// ReplaceImageRefPrefix replaces an image reference prefix with newPrefix.
// If the input image reference does not start with oldPrefix, an error is returned
func TestReplaceImageRefPrefix(t *testing.T) {
	tests := map[string]struct {
		s, oldPrefix, newPrefix string
		namespaceMapping map[string]string
                exp string
		expErr error
        }{
		"1": {s: "foo/baz/cat", oldPrefix: "foo", newPrefix: "bar", namespaceMapping: map[string]string{"baz": "qux"}, exp: "bar/qux/cat", expErr: nil},
		"2": {s: "foo/baz/cat", oldPrefix: "foo", newPrefix: "bar", namespaceMapping: map[string]string{}, exp: "bar/baz/cat", expErr: nil},
		"3": {s: "foo/baz", oldPrefix: "foo", newPrefix: "bar", namespaceMapping: map[string]string{}, exp: "bar/baz", expErr: nil},
		"4": {s: "foo/baz", oldPrefix: "foob", newPrefix: "bar", namespaceMapping: map[string]string{}, exp: "", expErr: errors.New("")},
		"5": {s: "foo", oldPrefix: "fo", newPrefix: "bar", namespaceMapping: map[string]string{}, exp: "", expErr: errors.New("")},
		"6": {s: "foo/openshift/cat@swan", oldPrefix: "foo", newPrefix: "bar", namespaceMapping: map[string]string{}, exp: "bar/openshift/cat", expErr: nil},
        }

        for name, tc := range tests {
                t.Run(name, func(t *testing.T) {
                        got, err := ReplaceImageRefPrefix(tc.s, tc.oldPrefix, tc.newPrefix, tc.namespaceMapping)
			if tc.expErr == nil && got != tc.exp {
                                t.Fatalf("expected: %v, got: %v", tc.exp, got)
                        }
			if tc.expErr != nil && err == nil {
				t.Fatalf("expected error, got no error")
			}
                })
        }
}


// HasImageRefPrefix returns true if the input image reference begins with
// the input prefix followed by "/"
func TestHasImageRefPrefix(t *testing.T) {
	tests := map[string]struct {
		s, prefix string
		want	  bool
	}{
		"1": {s: "cat/", prefix: "cat", want: true},
		"2": {s: "catt/", prefix: "cat", want: false},
		"3": {s: "cat/dog/spider", prefix: "cat", want: true},
		"4": {s: "//ss", prefix: "", want: true},
		"5": {s: "abc", prefix: "", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := HasImageRefPrefix(tc.s, tc.prefix)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}


// ParseLocalImageReference
func TestParseLocalImageReference(t *testing.T) {
        tests := map[string]struct {
                s, prefix string
                exp       *LocalImageReference
		expErr	  error
        }{
		"1": {s: "reg/ns/name@dig", prefix: "reg", exp: &LocalImageReference{Registry: "reg", Namespace: "ns", Name: "name", Digest: "dig"}, expErr: nil},
		"2": {s: "reg/ns/name@dig:est", prefix: "reg", exp: &LocalImageReference{Registry: "reg", Namespace: "ns", Name: "name", Digest: "dig:est"}, expErr: nil},
		"3": {s: "reg/ns/name:tg", prefix: "reg", exp: &LocalImageReference{Registry: "reg", Namespace: "ns", Name: "name", Tag: "tg"}, expErr: nil},
		"4": {s: "reg/ns/name@dig", prefix: "cat", exp: nil, expErr: errors.New("")},
		"5": {s: "reg/cat", prefix: "reg", exp: nil, expErr: errors.New("")},
		"6": {s: "reg/ns/name/dig", prefix: "reg", exp: nil, expErr: errors.New("")},
		"7": {s: "reg/ns/name@dig@est", prefix: "reg", exp: nil, expErr: errors.New("")},
		"8": {s: "reg/ns/name:ta:g", prefix: "reg", exp: nil, expErr: errors.New("")},
		"9": {s: "reg/ns/name", prefix: "reg", exp: &LocalImageReference{Registry: "reg", Namespace: "ns", Name: "name"}, expErr: nil},
        }

        for name, tc := range tests {
                t.Run(name, func(t *testing.T) {
                        got, err := ParseLocalImageReference(tc.s, tc.prefix)
			if tc.expErr == nil && !reflect.DeepEqual(got, tc.exp) {
                                t.Fatalf("expected: %v, got: %v", tc.exp, got)
                        }
			if tc.expErr != nil && err == nil {
				t.Fatalf("expected error, got no error")
			}
                })
        }
}






