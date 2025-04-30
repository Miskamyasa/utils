//TestExecAsyncSuccess: Verifies the core functionality. It runs a function that returns a string after a short delay, then calls Await() and checks if the correct result is returned. It also checks if Await() actually blocked for roughly the expected duration.
//TestExecAsyncDifferentTypes: Uses a table-driven approach to ensure the mechanism works correctly with various common return types like int, a custom struct, nil, and even an error returned as an interface{}.
//TestExecAsyncMultipleAwaits: Checks an important property: calling Await() multiple times on the same Future should return the cached result without re-running the original function. It uses a counter (executionCount) to verify the function runs only once and checks that subsequent calls to Await() are much faster.
//TestExecAsyncNoDelay: Ensures that even if the background goroutine finishes almost instantly (before Await might even be called), Await still correctly retrieves the result.

package async

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestExecAsyncSuccess tests the basic happy path: executing a function
// asynchronously and awaiting its successful result.
func TestExecAsyncSuccess(t *testing.T) {
	expectedResult := "success value"
	delay := 20 * time.Millisecond // Simulate some work

	future := ExecAsync(func() interface{} {
		time.Sleep(delay)
		return expectedResult
	})

	// Measure time to ensure Await blocks appropriately
	startTime := time.Now()
	result := future.Await()
	duration := time.Since(startTime)

	assert.Equal(t, expectedResult, result, "Await should return the correct result")
	assert.GreaterOrEqual(t, duration, delay, "Await should block for at least the duration of the async function")
	// Add a reasonable upper bound to catch potential hangs or extreme delays
	assert.Less(t, duration, delay*5, "Await took unexpectedly long")
}

// TestExecAsyncDifferentTypes tests that ExecAsync and Await work correctly
// with different return types (int, struct, nil).
func TestExecAsyncDifferentTypes(t *testing.T) {
	type customStruct struct {
		Name string
		Age  int
	}

	testCases := []struct {
		name           string
		inputFunc      func() interface{}
		expectedResult interface{}
	}{
		{
			name: "Integer Result",
			inputFunc: func() interface{} {
				time.Sleep(5 * time.Millisecond)
				return 12345
			},
			expectedResult: 12345,
		},
		{
			name: "Struct Result",
			inputFunc: func() interface{} {
				time.Sleep(5 * time.Millisecond)
				return customStruct{Name: "async test", Age: 99}
			},
			expectedResult: customStruct{Name: "async test", Age: 99},
		},
		{
			name: "Nil Result",
			inputFunc: func() interface{} {
				time.Sleep(5 * time.Millisecond)
				return nil
			},
			expectedResult: nil,
		},
		{
			name: "Error Result (as interface{})",
			inputFunc: func() interface{} {
				time.Sleep(5 * time.Millisecond)
				return errors.New("simulated error")
			},
			expectedResult: errors.New("simulated error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			future := ExecAsync(tc.inputFunc)
			result := future.Await()
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

// TestExecAsyncMultipleAwaits tests that calling Await multiple times on the
// same Future returns the same result without re-executing the function.
func TestExecAsyncMultipleAwaits(t *testing.T) {
	expectedResult := "consistent result"
	executionCount := 0
	delay := 15 * time.Millisecond

	future := ExecAsync(func() interface{} {
		executionCount++ // Increment counter on execution
		time.Sleep(delay)
		return expectedResult
	})

	// First Await
	result1 := future.Await()
	assert.Equal(t, expectedResult, result1, "First Await should return the correct result")
	assert.Equal(t, 1, executionCount, "Function should be executed exactly once after first Await")

	// Second Await - should be much faster and not increment counter
	startTime := time.Now()
	result2 := future.Await()
	duration := time.Since(startTime)

	assert.Equal(t, expectedResult, result2, "Second Await should return the same result")
	assert.Equal(t, 1, executionCount, "Function should not be executed again on second Await")
	assert.Less(t, duration, 5*time.Millisecond, "Second Await should be very fast")

	// Third Await
	result3 := future.Await()
	assert.Equal(t, expectedResult, result3, "Third Await should return the same result")
	assert.Equal(t, 1, executionCount, "Function should not be executed again on third Await")
}

// TestExecAsyncNoDelay tests the case where the async function completes very quickly.
func TestExecAsyncNoDelay(t *testing.T) {
	expectedResult := "immediate"

	future := ExecAsync(func() interface{} {
		// No delay
		return expectedResult
	})

	// Await should still work correctly, even if the goroutine finished before Await was called.
	result := future.Await()
	assert.Equal(t, expectedResult, result)
}

// Note: Testing context cancellation is not directly possible with the public
// Await() method as it always uses context.Background(). If context propagation
// was a requirement, the Await method would need to accept a context.Context.
