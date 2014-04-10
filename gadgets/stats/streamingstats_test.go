//Basic Coverage of StreamingStats
package stats

import (
	"fmt"
	"testing"
)


var sampIn = []int{1,1,3,3,5,5,7,7,9,9}
var sampOut = []float64{0,0,1.1547005383792517,1.1547005383792515,1.6733200530681511,1.7888543819998317,2.2253945610567474,2.390457218668787,2.788866755113585,2.9814239699997196}
var sampOutP = []float64{0,0,0.9428090415820634,1,1.4966629547095764,1.632993161855452,2.0603150145508513,2.23606797749979,2.629368792488718,2.8284271247461903}

//input matches expected
func TestStDev(t *testing.T) {

	ss := &StreamStats{}

	//Check outputs after every cumulative push
	for i,s := range sampIn {
		ss.Push( float64(s) )

		v,err := ss.StandardDeviation(STDev)
		if err != nil {
			t.Error("fail")
		}
		if v != sampOut[i] {
			t.Errorf("!=", i, v, sampOut[i])
		}

		vp,err := ss.StandardDeviation(STDevP)
		if err != nil {
			t.Error("fail")
		}
		if vp != sampOutP[i] {
			t.Errorf("!= %d %f17 %f17", i, vp, sampOutP[i])
		}
	}

	ss.Reset() //ready for a new set of inputs
}

//sample output
func ExampleStDev() {
	ss := &StreamStats{}

	//we can get to stdev's after any push.
	for _,s := range sampIn {
		ss.Push( float64(s) )
	}

	fmt.Println("Mean:", ss.Mean())
	if v,err := ss.Variance(STDev); err == nil {
		fmt.Println("Variance:", v)
	}
	if v,err := ss.Variance(STDevP); err == nil {
		fmt.Println("VarianceP:", v)
	}
	if d,err := ss.StandardDeviation(STDev); err == nil {
		fmt.Println("Stdev:", d)
	}
	if d,err := ss.StandardDeviation(STDevP); err == nil {
		fmt.Println("StdevP:", d)
	}
	ss.Reset() //ready for new set of inputs
	if d,err := ss.StandardDeviation(STDev); err == nil {
		fmt.Println("Stdev:", d)
	}
	//Output:
	// Mean: 5
	// Variance: 8.88888888888889
	// VarianceP: 8
	// Stdev: 2.9814239699997196
	// StdevP: 2.8284271247461903
	// Stdev: 0

}
