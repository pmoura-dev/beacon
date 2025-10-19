package subscribers

import (
	"reflect"
	"testing"

	"github.com/pmoura-dev/beacon"
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

func Test_extractParamsFromMQTTTopic(t *testing.T) {
	type testCase struct {
		rawTopic  string
		mqttTopic string
		expected  *beacon.TopicMatch
	}

	tests := map[string]testCase{
		"Simple - one segment": {
			rawTopic:  "foo",
			mqttTopic: "foo",
			expected:  beacon.NewTopicMatch("foo", map[string]string{}),
		},
		"Simple - multiple segments": {
			rawTopic:  "foo/bar/baz",
			mqttTopic: "foo/bar/baz",
			expected:  beacon.NewTopicMatch("foo/bar/baz", map[string]string{}),
		},
		"Single level wildcard - single": {
			rawTopic:  "{foo_id}",
			mqttTopic: "12345",
			expected: beacon.NewTopicMatch("12345", map[string]string{
				"foo_id": "12345",
			}),
		},
		"Single level wildcard - beginning": {
			rawTopic:  "{foo_id}/bar",
			mqttTopic: "12345/bar",
			expected: beacon.NewTopicMatch("12345/bar", map[string]string{
				"foo_id": "12345",
			}),
		},
		"Single level wildcard - middle": {
			rawTopic:  "foo/{foo_id}/bar",
			mqttTopic: "foo/12345/bar",
			expected: beacon.NewTopicMatch("foo/12345/bar", map[string]string{
				"foo_id": "12345",
			}),
		},
		"Single level wildcard - end": {
			rawTopic:  "foo/{foo_id}",
			mqttTopic: "foo/12345",
			expected: beacon.NewTopicMatch("foo/12345", map[string]string{
				"foo_id": "12345",
			}),
		},
		"Single level wildcard - multiple": {
			rawTopic:  "foo/{foo_id}/bar/{bar_id}",
			mqttTopic: "foo/12345/bar/abcde",
			expected: beacon.NewTopicMatch("foo/12345/bar/abcde", map[string]string{
				"foo_id": "12345",
				"bar_id": "abcde",
			}),
		},
		"Multi level wildcard - root": {
			rawTopic:  "*",
			mqttTopic: "random/segment",
			expected:  beacon.NewTopicMatch("random/segment", map[string]string{}),
		},
		"Multi level wildcard - end": {
			rawTopic:  "foo/*",
			mqttTopic: "foo/random/segment",
			expected:  beacon.NewTopicMatch("foo/random/segment", map[string]string{}),
		},
		"Multi level wildcard - with single level wildcard before": {
			rawTopic:  "foo/{foo_id}/*",
			mqttTopic: "foo/12345/random/segment",
			expected: beacon.NewTopicMatch("foo/12345/random/segment", map[string]string{
				"foo_id": "12345",
			}),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			topic, _ := beacon.NewTopic(test.rawTopic)

			got := extractParamsFromMQTTTopic(topic, test.mqttTopic)

			if !reflect.DeepEqual(got, test.expected) {
				t.Fatalf("Test failed! Expected: %s, got: %s", test.expected, got)
			}
		})
	}
}
