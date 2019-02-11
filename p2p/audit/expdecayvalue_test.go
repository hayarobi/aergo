/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"
)

func Test_makeDecaySlice(t *testing.T) {
	type args struct {
		meanTime int
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{"T15Sec", args{15}, getDecaySlice(15)},
		{"T1Min", args{60}, getDecaySlice(60)},
		{"T15Min", args{900}, getDecaySlice(900)},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeDecaySlice(tt.args.meanTime)
			//fmt.Println("Slice ", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeDecaySlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExponentDecayValue_AddValue(t *testing.T) {
	type args struct {
		timeSec int64
		n       int
	}
	tests := []struct {
		name   string
		meanTime int
		args   args
	}{
		{"T900Sec", 900, args{}},
		{"T15Sec", 15, args{}},
		{"T5Sec", 5, args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewExponentDecayValue(tt.meanTime)
			b := NewExponentDecayValue(tt.meanTime)
			initialTime := time.Now().Unix()
			initialValue := float64(1000000)
			a.AddValue(initialTime, initialValue)
			b.AddValue(initialTime, initialValue)

			for i:=0;i<tt.meanTime; i++ {
				initialTime ++
				a.AddValue(initialTime, 0)
			}

			fmt.Println("After meantime", a.RawValue())
			valB := b.Value(initialTime)

			mathError := math.Abs(valB - a.RawValue())
			fmt.Println("Math error ",mathError)
			if mathError/a.RawValue() > 1.0e-10 {
				t.Errorf("unexpected value %f , want %f",a.RawValue(), valB)
			}
		})
	}
}
