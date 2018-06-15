//+build unit

package kafka

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestMurmurHasher(t *testing.T) {
	var expected = new(murmurHash)
	var actual = MurmurHasher()
	assert.Equal(t, actual, expected)
}

type writeTest struct {
	value        string
	expectedLen  int
	expectedHash int32
	expectedErr  error
}

var writeTestCases = []writeTest{
	{
		value:        "foo",
		expectedLen:  3,
		expectedHash: 597841616,
		expectedErr:  nil,
	},
	{
		value:        "foobar",
		expectedLen:  6,
		expectedHash: -790332482,
		expectedErr:  nil,
	},
}

func TestWrite(t *testing.T) {
	for _, v := range writeTestCases {
		//Arrange
		var m murmurHash

		//Act
		actualLen, actualErr := m.Write([]byte(v.value))

		//Assert
		assert.Equal(t, actualLen, v.expectedLen)
		assert.Equal(t, m.v, v.expectedHash)
		assert.Equal(t, actualErr, v.expectedErr)
	}
}

func TestReset(t *testing.T) {
	//Arrange
	var m murmurHash
	m.v = 42
	var expected int32
	expected = 0

	//Act
	m.Reset()

	//Assert
	assert.Equal(t, m.v, expected)
}

func TestSize(t *testing.T) {
	//Arrange
	var m murmurHash
	expected := 32

	//Act
	actual := m.Size()

	//Assert
	assert.Equal(t, actual, expected)
}

func TestBlockSize(t *testing.T) {
	//Arrange
	var m murmurHash
	expected := 4

	//Act
	actual := m.BlockSize()

	//Assert
	assert.Equal(t, actual, expected)
}

func TestSum(t *testing.T) {
	//Arrange
	var m murmurHash
	b := []byte("foo")
	expected := b

	//Act
	actual := m.Sum(b)

	//Assert
	assert.Equal(t, expected, actual)
}

type sum32Test struct {
	value    int32
	expected uint32
}

var sum32TestCases = []sum32Test{
	{
		value:    42,
		expected: 42,
	},
	{
		value:    -42,
		expected: 2147483606,
	},
}

func TestSum32(t *testing.T) {
	for _, v := range sum32TestCases {
		//Arrange
		var m murmurHash
		m.v = v.value

		//Act
		actual := m.Sum32()

		//Assert
		assert.Equal(t, actual, v.expected)
	}
}

type murmur2Test struct {
	value    string
	expected int32
}

var murmur2TestCases = []murmur2Test{
	{
		value:    "foo",
		expected: 597841616,
	},
	{
		value:    "foobar",
		expected: -790332482,
	},
	{
		value:    "42",
		expected: 417700972,
	},
}

func TestMurmur2(t *testing.T) {
	for _, v := range murmur2TestCases {
		//Arrange
		b := []byte(v.value)

		//Act
		actual := murmur2(b)

		//Assert
		assert.Equal(t, actual, v.expected)
	}
}

type toPositiveTest struct {
	value    int32
	expected int32
}

var toPositiveTestCases = []toPositiveTest{
	{
		value:    42,
		expected: 42,
	},
	{
		value:    -42,
		expected: 2147483606,
	},
}

func TestToPositive(t *testing.T) {
	for _, v := range toPositiveTestCases {
		//Act
		actual := toPositive(v.value)

		//Assert
		assert.Equal(t, actual, v.expected)
	}
}
