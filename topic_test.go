package beacon

import (
	"reflect"
	"testing"
)

func Test_NewTopic(t *testing.T) {
	type testCase struct {
		raw         string
		expected    *Topic
		wantErr     bool
		expectedErr error
	}

	tests := map[string]testCase{
		"Simple - one segment": {
			raw: "foo",
			expected: &Topic{
				raw:      "foo",
				segments: []string{"foo"},
			},
		},
		"Simple - multiple segments": {
			raw: "foo/bar/baz",
			expected: &Topic{
				raw:      "foo/bar/baz",
				segments: []string{"foo", "bar", "baz"},
			},
		},
		"Single level wildcard - single": {
			raw: "{foo_id}",
			expected: &Topic{
				raw:      "{foo_id}",
				segments: []string{"{foo_id}"},
				params:   []string{"foo_id"},
			},
		},
		"Single level wildcard - beginning": {
			raw: "{foo_id}/bar",
			expected: &Topic{
				raw:      "{foo_id}/bar",
				segments: []string{"{foo_id}", "bar"},
				params:   []string{"foo_id"},
			},
		},
		"Single level wildcard - middle": {
			raw: "foo/{foo_id}/bar",
			expected: &Topic{
				raw:      "foo/{foo_id}/bar",
				segments: []string{"foo", "{foo_id}", "bar"},
				params:   []string{"foo_id"},
			},
		},
		"Single level wildcard - end": {
			raw: "foo/{foo_id}",
			expected: &Topic{
				raw:      "foo/{foo_id}",
				segments: []string{"foo", "{foo_id}"},
				params:   []string{"foo_id"},
			},
		},
		"Single level wildcard - multiple": {
			raw: "foo/{foo_id}/bar/{bar_id}",
			expected: &Topic{
				raw:      "foo/{foo_id}/bar/{bar_id}",
				segments: []string{"foo", "{foo_id}", "bar", "{bar_id}"},
				params:   []string{"foo_id", "bar_id"},
			},
		},
		"Multi level wildcard - root": {
			raw: "*",
			expected: &Topic{
				raw:      "*",
				segments: []string{"*"},
			},
		},
		"Multi level wildcard - end": {
			raw: "foo/*",
			expected: &Topic{
				raw:      "foo/*",
				segments: []string{"foo", "*"},
			},
		},
		"Multi level wildcard - with single level wildcard before": {
			raw: "foo/{foo_id}/*",
			expected: &Topic{
				raw:      "foo/{foo_id}/*",
				segments: []string{"foo", "{foo_id}", "*"},
				params:   []string{"foo_id"},
			},
		},
		"Error - Empty single level wildcard": {
			raw:         "foo/{}/bar",
			wantErr:     true,
			expectedErr: ErrEmptySingleLevelWildcard,
		},
		"Error - Duplicated single level wildcard": {
			raw:         "foo/{foo_id}/{foo_id}",
			wantErr:     true,
			expectedErr: ErrDuplicatedSingleLevelWildcard,
		},
		"Error - Invalid multi-level wildcard position": {
			raw:         "foo/*/bar",
			wantErr:     true,
			expectedErr: ErrInvalidMultiLevelWildcardPosition,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := NewTopic(test.raw)

			if test.wantErr {
				if err != test.expectedErr {
					if err != test.expectedErr {
						t.Fatalf("Test failed! Expected error: %v, got: %v", test.expectedErr, err)
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("Test failed! Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(*test.expected, *got) {
				t.Fatalf("Test failed! Expected: %v, got: %v", *test.expected, *got)
			}
		})
	}
}
