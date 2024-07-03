package deco

import (
	"fmt"
	"math"
)

// LEM implements one of the many decompression models described by
// NEDU TR 18-05 December 2018 title "THALMANN ALGORITHM PARAMETER SETS 
// FOR SUPPORT OF CONSTANT 1.3 ATM PO2 HE-O2 DIVING TO 300 FSW".
// Authors: David J. Doolette, f. Gregory Murphy and Waine A. Gerth 
// NEDU satisfies the interface Model

type LEM struct {
	// Configuration
	drate float64
	arate float64

	name string
	comp []tissueComp
}

const (
	highPPO2 = 1.3
	lowPPO2  = 0.7
	maxDepth = 320
	numMPT   = maxDepth / 10
	PVCO2    = 2.3
	PVO2     = 2
	PH20     = 0
	PBOVP    = 0
	Pfvg     = PVO2 + PVCO2 + PH20
)

type tissueComp struct {
	p          float64
	half       float64
	sdr        float64
	halfChange bool
	linear     bool
	MPT        [numMPT]float64
	bo         float64
	mo         float64
	b1         float64
}

func (tc *tissueComp) String() string {
	return fmt.Sprintf("C[%v]=%f sdr=%v MPT=%v",
		tc.half, tc.p, tc.sdr, tc.MPT)
}

func (tc *tissueComp) dive(depth float64, time float64, pgas float64) {
	k := math.Log(2) / tc.half
	Pcompk := tc.p + (pgas-tc.p)*(1-math.Exp(-k*time))
	Pcomp := tc.p + (pgas-tc.p)*(1-math.Pow(2, -time/tc.half))
	if math.Abs(Pcompk-Pcomp) > 0.00001 {
		fmt.Println(Pcompk, Pcomp, math.Abs(Pcompk-Pcomp))
		panic("dive")
	}

	if tc.linear {
		Pamb := depth + 33
		if Pamb < tc.p+Pfvg-PBOVP {
			// Look at NEDU report EQ 12 and 13
			Pa := Pamb - (42.9 + 1.5 + 0)
			PcompL := tc.p + (Pa-Pamb+Pfvg-PBOVP)*time*k
			tc.p = PcompL
			if false{
				debug("[%v][linear]Pamb/Pa %v/%4.3f Ptis %4.3f Ptis+Pfvg-PBOVP %4.3f \n",
					int(tc.half), Pamb, Pa, tc.p, tc.p+Pfvg-PBOVP)
			}
			return
		} else {
			if false {
				debug("[%v][exp] Pamb/pgas %v/%4.3f Ptis %4.3f Ptis+Pfvg-PBOVP %4.3f \n",
					int(tc.half), Pamb, pgas, tc.p, tc.p+Pfvg-PBOVP)
			}
		}
	}
	tc.p = Pcomp
}

//TODO: what happens if ascending and the tissue go supersaturate on this segment
// we need to change the funcion here also???
func (tc *tissueComp) descend(time float64, pgas float64, rate float64) {
	k := math.Log(2) / tc.half
	Pcomp := pgas + rate*(time-1/k) - (pgas-tc.p-rate/k)*math.Exp(-k*time)
	tc.p = Pcomp
	debug("descend[%4.3f] ptis %4.3f\n", tc.half, tc.p)
}

func (tc *tissueComp) ceiling() float64 {
	if tc.sdr != 1.0 {
		if tc.halfChange == false {
			debug("ceiling: halfChange %v\n", tc.half/tc.sdr)
			tc.half /= tc.sdr
			tc.halfChange = true
		}
	}
	// TODO: Check how high is the supersaturation
	if tc.linear == false {
		debug("ceiling: changing to linear half %v\n", tc.half)
		tc.linear = true
	}
	for i, mpt := range tc.MPT {
		if false {
			debug("P=%v, depth =%v, MPT=%v\n", tc.p, (i+1)*10, mpt)
		}
		if mpt > tc.p {
			debug("ceiling for %v is %v\n", tc.half,i*10)
			return float64(i * 10)
		}
	}
	return 0.0
}


func (m *LEM) String() string {
	return fmt.Sprintf("[%v maxdepth=%v comp=%v MPT=%v]",
		m.name, maxDepth, len(m.comp), len(m.comp[0].MPT))
}

func (m *LEM) Print(total bool, label string) {
        if label == "" {
                label = "model"
        }
	for i := range m.comp {
		fmt.Fprintf(mwriter,"%v: C[%6.4v]=%.4v\n",label,m.comp[i].half,m.comp[i].p)
	}
}
func (m *LEM) UseMeters() {
	panic("Meters not supporte by LEM")
}

func dilPressureFSW(depth float64) float64 {
	return dilPressureFSW2(depth, highPPO2)

}
func dilPressureFSW2(depth, ppo2 float64) float64 {
	atm := depth/33 + 1
	atm -= 0.063
	fo2 := ppo2 / atm
	fdil := 1 - fo2
	dildepth := fdil * (depth + 33)

	Pamb := depth + 33
	Pa := Pamb - (42.9 + 1.5 + 0)

	if false {
		debug("dilPressureFSW: depth %4.3f Pis %4.3f Pa %4.3f\n", depth, dildepth, Pa)
	}

	return dildepth
}

func (m *LEM) LevelOff(time float64, depth float64, mix *Mix) {

	debug("LevelOff: %v min %v ft\n", time, depth)
	Pis := dilPressureFSW(depth)
	debug("leveloff: Pis %v time %v\n", Pis, time)
	for i := range m.comp {
		m.comp[i].dive(depth, time, Pis)
	}
	if pmodel {
		s := fmt.Sprintf("leveloff")
                m.Print(true, s)
	}
}

func (m *LEM) Descend(to float64, mix *Mix) float64 {
	ppo2 := lowPPO2
	debug("Descend CCR(%v/%v) to %v\n", lowPPO2,highPPO2, to) 	
	for n := 1.0; n <= to; n++ {
		time := 1 / m.drate
		Pis := dilPressureFSW2(n, ppo2)
		if false {
			debug("Descend CCR(%v) n %v time %v Pis %v\n", ppo2, n, time, Pis)
		}
		for i := range m.comp {
			m.comp[i].dive(n, time, Pis)
		}
		if ppo2 == lowPPO2 && n > 20 {
			ppo2 = highPPO2
		}
	}
        if pmodel {
                s := fmt.Sprintf("descend")
                m.Print(true, s)
        }
	return to / m.drate
}

func (m *LEM) Ascend(from, to float64, mix *Mix) float64 {
	ppo2 := highPPO2
	debug("Ascent CCR (%v) from %v to %v\n", highPPO2, from, to)                        
	for n := from - 1; n >= to; n-- {
		time := 1 / m.arate
		Pis := dilPressureFSW2(n, ppo2)
		if false {
			debug("Ascend(%v) n %v time %v Pis %v\n", ppo2, n, time, Pis)
		}
		for i := range m.comp {
			m.comp[i].dive(n, time, Pis)
		}
	}
        if pmodel {
                s := fmt.Sprintf("ascend")
                m.Print(true, s)
        }
	return ((from - 1) - (to - 1)) / m.arate
}

func (m *LEM) Ceiling() float64 {
	maxP := 0.0
	half := 0.0
	for i := range m.comp {
		ceil := m.comp[i].ceiling()
		if ceil > maxP {
			maxP = ceil
			half = m.comp[i].half
		}
	}
	if half == 0.0 {
		for i := range m.comp {
			half = m.comp[i].half
		}
	}
	debug("Ceiling is %v controlling %v min\n", maxP, half)
	return maxP
}

func (m *LEM) Decompress(from float64, g *Mix, halfdepth bool) ([]Segment, float64) {
	return nil, 0
}
func (m *LEM) DecompressOCTech(from float64, g *Mix, dg []*Mix,
	laststopdeep bool, suggested_inc float64) ([]Segment, float64) {
	return nil, 0
}

func (m *LEM) DecompressCCR(from float64, g *Mix, ccrPPO2 float64,
	laststopdeep bool, suggested_inc float64) ([]Segment, float64) {

	total := 0.0
	stoplen := 0.0
	inc := 1.0

        if pmodel {
                m.Print(true, "deco "+"start ")
        }

	ceil := m.Ceiling()
	stop := ceil
	if ceil == 0 {
		return nil, 0
	}

	total += m.Ascend(from, ceil, nil)

	if nceil := m.Ceiling(); nceil != ceil {
		stop = nceil
		debug("ceil: resetting from %v to %v\n", ceil, nceil)
	}

	deco := make([]Segment, 0, 100)
	for i := 0; i < 10000; i++ {
		m.LevelOff(inc, stop, nil)
		ceil = m.Ceiling()
		if ceil == stop || stop == 20 && ceil != 0 {
			stoplen += inc
		} else {
			stoplen += inc
                        if pmodel {
                                s := fmt.Sprintf("deco %v %v", stop, stoplen)
                                m.Print(true, s)
                        }
			deco = append(deco, Segment{stop, stoplen})
			total += stoplen

			if false {
				fmt.Printf(">%v %v\n", stop, stoplen)
			}

			if ceil == 0 {
				return deco, total
			}

			laststop := stop
			if stop-ceil > 10 {
				stop -= 10
			} else {
				stop = ceil
			}
			if true {
				total += m.Ascend(laststop, stop, nil)
			}
			stoplen = 0
		}
	}
	panic(0)

	return nil, 0
}

