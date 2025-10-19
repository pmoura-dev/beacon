package publishers

import (
	"testing"
)

func Test_toMQTTTopic(t *testing.T) {
	type testCase struct {
		topic    string
		expected string
	}

	tests := map[string]testCase{
		"Simple - one segment": {
			topic:    "foo",
			expected: "foo",
		},
		"Simple - multiple segments": {
			topic:    "foo/bar/baz",
			expected: "foo/bar/baz",
		},
		"Single level wildcard - single": {
			topic:    "{foo_id}",
			expected: "+",
		},
		"Single level wildcard - beginning": {
			topic:    "{foo_id}/bar",
			expected: "+/bar",
		},
		"Single level wildcard - middle": {
			topic:    "foo/{foo_id}/bar",
			expected: "foo/+/bar",
		},
		"Single level wildcard - end": {
			topic:    "foo/{foo_id}",
			expected: "foo/+",
		},
		"Single level wildcard - multiple": {
			topic:    "foo/{foo_id}/bar/{bar_id}",
			expected: "foo/+/bar/+",
		},
		"Multi level wildcard - root": {
			topic:    "*",
			expected: "#",
		},
		"Multi level wildcard - end": {
			topic:    "foo/*",
			expected: "foo/#",
		},
		"Multi level wildcard - with single level wildcard before": {
			topic:    "foo/{foo_id}/*",
			expected: "foo/+/#",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := toMQTTTopic(test.topic)

			if test.expected != got {
				t.Fatalf("Test failed! Expected: %s, got: %s", test.expected, got)
			}
		})
	}
}
