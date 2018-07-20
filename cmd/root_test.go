//+build unit

package cmd

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

type parametersTest struct {
	params   parameters
	expected error
}

var parametersTestCases = []parametersTest{
	{
		params: parameters{
			fromBrokers: "foo",
			fromTopic:   "bar",
			toTopic:     "foobar",
			hasher:      "murmur2",
		},
		expected: nil,
	},
	{
		params: parameters{
			fromBrokers: "foo",
			toTopic:     "foobar",
			hasher:      "murmur2",
		},
		expected: errMissingSourceTopic,
	},
	{
		params: parameters{
			fromBrokers: "foo",
			fromTopic:   "bar",
			hasher:      "murmur2",
		},
		expected: errMissingTargetTopic,
	},
	{
		params: parameters{
			fromTopic: "bar",
			toTopic:   "foobar",
			hasher:    "murmur2",
		},
		expected: errMissingSourceBrokers,
	},
	{
		params: parameters{
			fromBrokers: "foo",
			fromTopic:   "bar",
			toTopic:     "foobar",
			hasher:      "murmur2",
			loop:        true,
		},
		expected: errLoopCloningWithTarget,
	},
	{
		params: parameters{
			fromBrokers: "foo",
			fromTopic:   "bar",
			toTopic:     "bar",
			hasher:      "murmur2",
		},
		expected: errLoopRequired,
	},
	{
		params: parameters{
			fromBrokers: "foo",
			fromTopic:   "bar",
			toBrokers:   "foo",
			toTopic:     "foobar",
			hasher:      "murmur2",
		},
		expected: errSourceBrokersIsTarget,
	},
	{
		params: parameters{
			fromBrokers: "foo",
			fromTopic:   "bar",
			toTopic:     "foobar",
			hasher:      "murmur",
		},
		expected: errUnknownHasher,
	},
}

func TestValidateParameters(t *testing.T) {
	for _, v := range parametersTestCases {
		//Act
		actual := v.params.validate()

		//Assert
		assert.Equal(t, actual, v.expected)
	}
}

func TestGetBrokers(t *testing.T) {
	//Arrange
	params.fromBrokers = "localhost1:9092;localhost2:9092;localhost3:9092"
	params.toBrokers = "distanthost1:9092;distanthost2:9092;distanthost3:9092"
	expectedFromBrokers := []string{"localhost1:9092", "localhost2:9092", "localhost3:9092"}
	expectedToBrokers := []string{"distanthost1:9092", "distanthost2:9092", "distanthost3:9092"}

	//Act
	actualFromBrokers, actualToBrokers := getBrokers()

	//Assert
	assert.Equal(t, actualFromBrokers, expectedFromBrokers)
	assert.Equal(t, actualToBrokers, expectedToBrokers)
}

type containsTest struct {
	s        []string
	e        string
	expected bool
}

var containsTestCases = []containsTest{
	{
		s:        []string{"foo", "bar", "foobar"},
		e:        "foo",
		expected: true,
	},
	{
		s:        []string{"foo", "bar", "foobar"},
		e:        "bar",
		expected: true,
	},
	{
		s:        []string{"foo", "bar", "foobar"},
		e:        "foobar",
		expected: true,
	},
	{
		s:        []string{"foo", "bar", "foobar"},
		e:        "barfoo",
		expected: false,
	},
}

func TestContains(t *testing.T) {
	for _, v := range containsTestCases {
		//Act
		actual := contains(v.s, v.e)

		//Assert
		assert.Equal(t, actual, v.expected)
	}
}
