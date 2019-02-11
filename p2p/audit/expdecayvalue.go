/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"math"
	"sync"
)

const (
	decaySliceLength = 32
)

var (
	// key is mean lifetime in seconds, and value is slice of pre-calculated time factor
	decayMatrix map[int][]float64
	matrixLock sync.Mutex
)

func init() {
	decayMatrix = make(map[int][]float64)
	// create matrix with frequently used decay factor
	decayMatrix[15] = makeDecaySlice(15)
	decayMatrix[60] = makeDecaySlice(60)
	decayMatrix[900] = makeDecaySlice(900)
}

func makeDecaySlice(meanTime int) []float64 {
	decayPerSec := math.Exp(-1/float64(meanTime))
	decaySlice := make([]float64,decaySliceLength)
	decaySlice[0] = 1.0 * decayPerSec
	for i:= 1; i < decaySliceLength ; i++ {
		decaySlice[i] = math.Pow(decayPerSec,float64(i+1))
	}
	return decaySlice
}

func getDecaySlice(meanTime int) []float64 {
	matrixLock.Lock()
	defer matrixLock.Unlock()
	if slice, found := decayMatrix[meanTime]; found {
		return slice
	} else {
		decayMatrix[meanTime] = makeDecaySlice(meanTime)
		return decayMatrix[meanTime]
	}
}

// ExponentDecayValue store exponentially decayed value by second precision.
// This is not thread-safe.
type ExponentDecayValue struct {
	decaySlice  []float64
	truncateCount int

	value       float64
	lastTimeSec int64
}


// NewExponentDecayValue create exponentMetric with given mean lifetime seconds
func NewExponentDecayValue(meanTime int) *ExponentDecayValue {
	return newExponentDecayValue(getDecaySlice(meanTime), meanTime<<5)
}

func newExponentDecayValue(decaySlice []float64, truncateCount int) *ExponentDecayValue {
	return &ExponentDecayValue{decaySlice: decaySlice, truncateCount:truncateCount}
}

// Update adds value n, timeSec is unix timestamp
func (a *ExponentDecayValue) AddValue(timeSec int64, n float64) {
	passed := timeSec - a.lastTimeSec
	for passed > 0 {
		// Decay current value
		// if time is too past, value wnet to effectively zero.
		if passed > (decaySliceLength<<4) {
			a.value = 0.0
			break
		} else if passed > decaySliceLength {
			a.value *= a.decaySlice[decaySliceLength-1]
			passed -= decaySliceLength
			continue
		} else {
			a.value *= a.decaySlice[passed-1]
			break
		}
	}
	a.lastTimeSec = timeSec
	a.value += float64(n)
}

// RawValue returns current value without time-correction
func (a *ExponentDecayValue) RawValue() float64 {
	return a.value
}

// Value returns current value corrected to given time
func (a *ExponentDecayValue) Value(timeSec int64) float64 {
	passed := timeSec - a.lastTimeSec
	for passed > 0 {
		// Decay current value
		if passed > decaySliceLength {
			a.value *= a.decaySlice[decaySliceLength-1]
			passed -= decaySliceLength
			continue
		} else {
			a.value *= a.decaySlice[passed-1]
			break
		}
	}
	a.lastTimeSec = timeSec
	return a.value
}
