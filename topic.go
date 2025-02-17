package beacon

import (
	"errors"
	"slices"
	"strings"
)

var (
	ErrEmptySingleLevelWildcard          = errors.New("single level wildcard is empty")
	ErrDuplicatedSingleLevelWildcard     = errors.New("single level wildcard is duplicated")
	ErrInvalidMultiLevelWildcardPosition = errors.New("multi-level wildcard '*' must be the last level")
)

type Topic struct {
	raw      string
	segments []string
	params   []string
}

func NewTopic(raw string) (*Topic, error) {
	segments := strings.Split(raw, "/")

	var params []string
	for i, s := range segments {

		if strings.Trim(s, " ") == "*" && i != len(segments)-1 {
			return nil, ErrInvalidMultiLevelWildcardPosition
		}

		if isWildcard(s) {
			param := strings.Trim(s[1:len(s)-1], " ")
			if param == "" {
				return nil, ErrEmptySingleLevelWildcard
			}

			if slices.Contains(params, param) {
				return nil, ErrDuplicatedSingleLevelWildcard
			}

			params = append(params, param)
		}
	}

	return &Topic{
		raw:      raw,
		segments: segments,
		params:   params,
	}, nil
}

func (t *Topic) Raw() string {
	return t.raw
}

func (t *Topic) Segments() []string {
	return t.segments
}

func (t *Topic) Params() []string {
	return t.params
}

type TopicMatch struct {
	fullName string
	params   map[string]string
}

func NewTopicMatch(fullName string, params map[string]string) *TopicMatch {
	return &TopicMatch{
		fullName: fullName,
		params:   params,
	}
}

func (m *TopicMatch) FullName() string {
	return m.fullName
}

func (m *TopicMatch) Params() map[string]string {
	return m.params
}

func isWildcard(s string) bool {
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}
