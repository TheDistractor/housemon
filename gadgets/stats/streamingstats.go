//Simple Streaming Stats.
//By: TheDistractor
//License: MIT  see:https://thedistractor.github.io/LICENSE.md
//
//STDev (sample set)
//STDevP (population set)
package stats

import (
	"errors"
	"math"
)

//compute some primitive running statistics
type StreamStats struct {
	idx float64
	pM  float64
	pS  float64
	cM  float64
	cS  float64
}

const (
	STDev = iota
	STDevP
)

func (s *StreamStats) Reset() {
	s.idx = 0
	s.pM = 0
	s.pS = 0
	s.cM = 0
	s.cS = 0
}

func (s *StreamStats) Push(val float64) {

	s.idx += 1

	if s.idx == 1 {
		s.pM, s.cM = val, val
		s.pS = 0.0
	} else {
		s.cM = s.pM + (val-s.pM)/s.idx
		s.cS = s.pS + (val-s.pM)*(val-s.cM)

		s.pM = s.cM
		s.pS = s.cS
	}

}

func (s *StreamStats) Mean() float64 {
	if s.idx > 0 {
		return s.cM
	}
	return 0.0
}

func (s *StreamStats) Variance(typ int) (float64, error) {
	if s.idx > 1 {
		switch typ {
		case STDev:
			return s.cS / (s.idx - 1), nil
		case STDevP:
			return s.cS / (s.idx ), nil
		default:
			return 0, errors.New("Invalid Deviation Type")
		}
	}
	return 0.0, nil
}

func (s *StreamStats) StandardDeviation(typ int) (float64, error) {
	sd, err := s.Variance(typ)
	if err != nil {
		return 0.0, err
	}
	return math.Sqrt(sd), nil
}
