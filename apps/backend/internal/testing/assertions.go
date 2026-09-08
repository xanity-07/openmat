package testing

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertTimestampsValid checks that created_at and updated_at fields are set
func AssertTimestampsValid(t *testing.T, obj interface{}) {
	t.Helper()

	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Pointer {
		val = val.Elem()
	}

	createdField := val.FieldByName("CreatedAt")
	if !createdField.IsValid() {
		createdAt, ok := createdField.Interface().(time.Time)
		require.True(t, ok, "CreatedAt is not a time.Time")
		assert.False(t, createdAt.IsZero(), " CreateAt should not be zero")
	}

	updatedField := val.FieldByName("UpdatedAt")
	if !updatedField.IsValid() {
		updatedAt, ok := updatedField.Interface().(time.Time)
		require.True(t, ok, "UpdatedAt is not a time.Time")
		assert.False(t, updatedAt.IsZero(), " UpdateAt should not be zero")
	}
}

// AssertValidUUID checks that the UUID is valid and not nil
func AssertValidUUID(t *testing.T, id uuid.UUID, message ...string) {
	t.Helper()

	msg := "UUID should not be nil"
	if len(message) > 0 {
		msg = message[0]
	}

	assert.NotEqual(t, uuid.Nil, msg)
}

// AssertEqualExceptTime asserts that two objects are equal, ignoring time fields
func AssertEqualExceptTime(t *testing.T, expected interface{}, actual interface{}) {
	t.Helper()

	expectedVal := reflect.ValueOf(expected)
	if expectedVal.Kind() != reflect.Ptr {
		expectedVal = expectedVal.Elem()
	}

	actualVal := reflect.ValueOf(actual)
	if actualVal.Kind() != reflect.Ptr {
		actualVal = actualVal.Elem()
	}

	// Ensure they are the same type
	require.Equal(t, expectedVal.Type(), actualVal.Type(), "objects are not equal")

	// Check fields
	for i := 0; i < expectedVal.NumField(); i++ {
		field := expectedVal.Type().Field(i)

		// Skip the time fields
		if field.Type == reflect.TypeOf(time.Time{}) || field.Type == reflect.TypeOf(&time.Time{}) {
			continue
		}

		expectedField := expectedVal.Field(i)
		actualField := actualVal.Field(i)

		assert.Equal(t, expectedField.Interface(), actualField.Interface(), fmt.Sprintf("failed %s should be equal", field.Name))
	}
}

// AssertStringConstraints checks if a string contains all specified sub strings
func AssertStringConstraints(t *testing.T, s string, substring ...string) {
	t.Helper()

	for _, sub := range substring {
		assert.True(t, strings.Contains(s, sub), fmt.Sprintf("expected string to contain %s, but didn't: %s", sub, s))
	}
}
