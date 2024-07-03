package deco

import (
	"fmt"
	"math"
)

// ZHL16 implements one of the many decompression models created by
// Dr. Albert A. Bühlmann. See  Bühlmann, Albert A (1984). Decompression-Decompression Sickness.
// Berlin New York: Springer-Verlag. ISBN 0-387-13308-9.
// ZHL16 satisfies the interface Model
type ZHL16 struct {
	checkmix  bool
	ismix     bool
	useMeters bool

	// Gradient Factor
	low      float64
	high     float64
	gradient float64
	ginc     float64

	// Configuration
	drate float64
	arate float64

	name       string
	gasCompSet map[string][]compartment
}

type compartment struct {
	p    float64
	half float64
	a    float64
	b    float64
}

func (c *compartment) String() string {
	return fmt.Sprintf("C[%v]=%f a=%v b=%v",
		c.half, c.p, c.a, c.b)
}

func (c *compartment) dive(time float64, pgas float64) {

	k := math.Log(2) / c.half
	Pcompk := c.p + (pgas-c.p)*(1-math.Exp(-k*time))
	Pcomp := c.p + (pgas-c.p)*(1-math.Pow(2, -time/c.half))
	if math.Abs(Pcompk-Pcomp) > 0.00001 {
		fmt.Println(Pcompk, Pcomp, math.Abs(Pcompk-Pcomp))
		panic("dive")
	}
	c.p = Pcomp
}

func (c *compartment) descend(time float64, pgas float64, rate float64) {
	k := math.Log(2) / c.half
	Pcomp := pgas + rate*(time-1/k) - (pgas-c.p-rate/k)*math.Exp(-k*time)
	c.p = Pcomp
}

func (c *compartment) ceiling(gf float64) float64 {
	if gf == 1.0 {
		Pambtol := (c.p - c.a) * c.b
		return Pambtol
	}
	A := c.a * gf
	B := gf/c.b + 1.0 - gf
	return (c.p - A) / B
}

func (m *ZHL16) String() string {
	return fmt.Sprintf("[%v GF=[%v,%v,%v] ismix=%v]",
		m.name, m.gradient*100, m.low*100, m.high*100, m.ismix)
}

func (m *ZHL16) Print(total bool, label string) {

	if label == "" {
		label = "model"
	}

	if total == true {
		max := 0
		for _, gasComp := range m.gasCompSet {
			n := 0
			for range gasComp {
				n++
			}
			if n > max {
				max = n
			}
		}
		for i := 0; i < max; i++ {
			p := 0.0

			for _, gasComp := range m.gasCompSet {
				p += gasComp[i].p
			}
			if false {
				fmt.Fprintf(mwriter,"%v: C[%v]=%v\n", label, i, p)
			} else {
				fmt.Fprintf(mwriter,"%v: C[%v]=%f ", label, i, p)
				//
				// Doing it this way as the order can change
				// as walking the map with range
				//
				fmt.Fprintf(mwriter,"%v=%f ", "N2", m.gasCompSet["N2"][i].p)
				fmt.Fprintf(mwriter,"%v=%f ", "He", m.gasCompSet["He"][i].p)
				fmt.Fprintf(mwriter,"\n")
			}
		}
		return
	}

	for g, gasComp := range m.gasCompSet {
		for i, _ := range gasComp {
			fmt.Printf("%v: %v %v\n", label, g, gasComp[i])
		}
	}
}

func (m *ZHL16) UseMeters() {
	m.useMeters = true
}

func (m *ZHL16) setDeltaGradient(n float64) {
	if m.ginc != 0 {
		return
	}
	if m.low < m.high && n != 0 {
		inc := (m.high - m.low) / n
		m.ginc = inc
		m.gradient = m.low + inc
	}
}

func (m *ZHL16) increaseGradient() {
	if m.low < m.high && m.ginc != 0 {
		m.gradient += m.ginc
	}
}

// LevelOff exposes the model to a constant depth, depth, for a number
// of seconds, time, while using a mix of gases, mix
func (m *ZHL16) LevelOff(time float64, depth float64, mix *Mix) {

	atm := Feet2ATM(depth)
	if m.useMeters {
		atm = Meter2ATM(depth)
	}

	gases := 0
	debug("leveloff: g %v time %v\n", mix, time)
	for g, gasComp := range m.gasCompSet {
		pgas := (*mix).gas[g] * (atm - 0.063)
		debug("leveloff: g %v time %v pgas %v\n", g, time, pgas)
		if pgas != 0 {
			gases++
		}
		for i, _ := range gasComp {
			gasComp[i].dive(time, pgas)
		}
	}
	if m.checkmix {
		m.checkmix = false
		m.ismix = false
		if mix.ismix {
			m.ismix = true
		}
	}

	if pmodel {
		s := fmt.Sprintf("leveloff")
		m.Print(true, s)
	}
}

// Descend exposes the model to a descent to pressure in depth, to,while using a mix of gases, mix.
// It return the number of second that took to descent
func (m *ZHL16) Descend(depth float64, mix *Mix) float64 {
	to := Feet2ATM(depth)
	if m.useMeters {
		to = Meter2ATM(depth)
	}

	gases := 0
	rate := m.drate
	time := (to - 1) / rate

	debug("descend: g %v time %v rate %v\n", mix, time, rate)
	for g, gasComp := range m.gasCompSet {
		pgas := (*mix).gas[g] * (to - 0.063)
		rgas := (*mix).gas[g] * rate
		debug("descend: g %v time %v pgas %v rgas %v \n", g, time, pgas, rgas)
		if pgas != 0 {
			gases++
		}
		for i, _ := range gasComp {
			gasComp[i].descend(time, pgas, rgas)
		}
	}
	if m.checkmix {
		m.checkmix = false
		m.ismix = false
		if mix.ismix {
			m.ismix = true
		}
	}

	if pmodel {
		s := fmt.Sprintf("descend")
		m.Print(true, s)
	}

	return time
}

// Ascend exposes the model to an ascent from two depths, depth1 and depth2,
// ,while using a mix of gases, mix. It return the number of second that took to move
func (m *ZHL16) Ascend(depth1, depth2 float64, mix *Mix) float64 {

	from := Feet2ATM(depth1)
	to := Feet2ATM(depth2)
	if m.useMeters {
		from = Meter2ATM(depth1)
		to = Meter2ATM(depth2)
	}

	gases := 0
	time := ((from - 1) - (to - 1)) / m.arate
	rate := -m.arate

	debug("ascend: g %v from %v to %v time %v rate %v\n", mix, from, to, time, rate)

	for g, gasComp := range m.gasCompSet {
		pgas := (*mix).gas[g] * (to - 0.063)
		rgas := (*mix).gas[g] * rate
		debug("ascend: g %v time %v pgas %v rgas %v \n", g, time, pgas, rgas)
		if pgas != 0 {
			gases++
		}
		for i, _ := range gasComp {
			gasComp[i].descend(time, pgas, rgas)
		}
	}
	if m.checkmix {
		m.checkmix = false
		m.ismix = false
		if mix.ismix {
			m.ismix = true
		}
	}

	if pmodel {
		s := fmt.Sprintf("ascend")
		m.Print(true, s)
	}

	return time
}

func roundNext10ftStop(atm float64) float64 {

	// TODO: allow to ascent a bit before 1.0 atm?
	if false && atm < 1.03 {
		return 0.0
	}

	ft := ATM2Feet(atm)
	if ft <= 0.0 {
		return 0.0
	}

	var depth float64
	for depth = 10; depth < ft; depth += 10 {
	}

	return depth
}

func roundNext3mStop(atm float64) float64 {

	// TODO: allow to ascent a bit before 1.0 atm?
	if false && atm < 1.03 {
		return 0.0
	}

	mt := ATM2Meter(atm)
	if mt <= 0.0 {
		return 0.0
	}

	var depth float64
	for depth = 3; depth < mt; depth += 3 {
	}

	return depth
}

// Ceiling return the smaller pressure in atm the model can tolerate based on
// current gas loading
func (m *ZHL16) Ceiling() float64 {

	roundNextStop := roundNext10ftStop
	if m.useMeters {
		roundNextStop = roundNext3mStop
	}
	stop := 1.0
	top := -1

	if m.ismix {
		var A, B, P float64
		i := 0
		size := 0

		for {
			A = 0.0
			B = 0.0
			P = 0.0

			for _, gasComp := range m.gasCompSet {
				if size == 0 {
					size = len(gasComp)
				}
				c := gasComp[i]

				P += c.p
				A += c.a * c.p
				B += c.b * c.p
			}

			A /= P
			B /= P
			Pambtol := (P - A) * B

			if m.gradient != 1.0 {
				A = A * m.gradient
				B = m.gradient/B + 1.0 - m.gradient
				Pambtol = (P - A) / B
			}

			if Pambtol > stop {
				debug("gradient %v\n", m.gradient)
				for g, gasComp := range m.gasCompSet {
					c := gasComp[i]
					debug("ceil[%v]: %v %v %v\n", i, g, c.String(), Pambtol)
				}
				stop = Pambtol
				top = i
			}
			if i++; i >= size {
				break
			}
		}
	} else {
		for g, gasComp := range m.gasCompSet {
			for i, _ := range gasComp {
				c := gasComp[i]
				if ceil := c.ceiling(m.gradient); ceil > stop {
					debug("ceil[%v]: %v %v %v\n", i, g, c.String(), ceil)
					stop = ceil
					top = i
				}
			}
		}
	}

	if top > 0 {
		for g, gasComp := range m.gasCompSet {
			c := gasComp[top]
			debug("top[%v]: %v %v %v\n", top, g, c.String(), stop)
		}
	}
	return roundNextStop(stop)
}

// Decompress creates a very simple decompression ascent, storing it on Segment.
// It starts at pressure from, and uses mix g during the whole ascent.
// If halfdepth is true, it forces the deco to start at least at half of the depth
// It returns the total ascent time as second argument
func (m *ZHL16) Decompress(from float64, g *Mix, halfdepth bool) ([]Segment, float64) {
	if m.useMeters {
		return m.decompress(from, g, halfdepth, 3.0, Meter2ATM, ATM2Meter, roundNext3mStop)
	}
	return m.decompress(from, g, halfdepth, 10.0, Feet2ATM, ATM2Feet, roundNext10ftStop)
}

func (m *ZHL16) decompress(from float64, g *Mix, halfdepth bool,
	move float64, Depth2ATM, ATM2Depth, roundNextStop func(float64) float64) ([]Segment, float64) {

	total := 0.0
	inc := 1.0
	stoplen := 0.0

	if pmodel {
		m.Print(true, "")
	}
	ceil := m.Ceiling()
	stop := ceil
	debug("ceil: %v depth %v\n", ceil, stop)
	if stop == 0 {
		return nil, 0
	}

	if halfdepth && stop < from/2 {
		stop = roundNextStop(Depth2ATM(from / 2))
	}

	total += m.Ascend(from, stop, g)

	m.setDeltaGradient(stop / move)
	deco := make([]Segment, 0, 100)
	for i := 0; i < 10000; i++ {
		m.LevelOff(inc, float64(stop), g)

		ceil = m.Ceiling()
		nextstop := ceil
		debug("ceil: %v depth %v\n", ceil, nextstop)

		if stop == nextstop {
			stoplen += inc
		} else {
			if pmodel {
				m.Print(true, "")
			}
			stoplen += inc

			deco = append(deco, Segment{stop, stoplen})
			total += stoplen

			if nextstop == 0 {
				return deco, total
			}

			laststop := stop
			if stop-move == nextstop {
				stop = nextstop
			} else {
				stop = stop - move
			}
			if true {
				total += m.Ascend(laststop, stop, g)
			}
			m.increaseGradient()

			stoplen = 0
		}
	}
	panic("decompress")
}

// DecompressOCTech creates a technical decompression ascent, storing it on Segment.
// It starts at depth from, and uses deco mixes available in dg
// It adjust gradients factor during the stop of the decompression
//
// laststopdeep controls if the last stop, 10ft/3m, is skipped during decompression
// suggested_inc adjusts the stop len multiplier to make stops a bit longer and even
//
// It returns the total ascent time as second argument
func (m *ZHL16) DecompressOCTech(from float64, g *Mix, dg []*Mix,
	laststopdeep bool, suggested_inc float64) ([]Segment, float64) {

	if m.useMeters {
		return m.decompressOCTech(from, g, dg,
			laststopdeep, suggested_inc,
			3.0, 21.0, Meter2ATM, ATM2Meter, roundNext3mStop)
	}
	return m.decompressOCTech(from, g, dg,
		laststopdeep, suggested_inc,
		10.0, 70.0, Feet2ATM, ATM2Feet, roundNext10ftStop)
}

func (m *ZHL16) decompressOCTech(from float64, g *Mix, dg []*Mix,
	laststopdeep bool, suggested_inc float64,
	move float64, minstop float64,
	Depth2ATM, ATM2Depth, roundNextStop func(float64) float64) ([]Segment, float64) {

	// TODO: We might want to have a slide of starting depth for each mix, as well
	// as a slide of suggestion for the length of stop for each mix. In the meantime
	// we estimated it
	mod := make([]float64, len(dg))
	if len(dg) != 0 {
		max := 1.56
		if move == 3 {
			max = 1.5
		}
		for i, m := range dg {
			mod[i] = roundNextStop(max / m.o2)
			if m.o2 == 0.32 {
				mod[i] = roundNextStop(max / 0.5)
			}
			if mod[i] > minstop {
				minstop = mod[i]
			}
		}
	}

	total := 0.0
	stoplen := 0.0
	inc := 1.0
	if pmodel {
		m.Print(true, "deco "+"start ")
	}

	ceil := m.Ceiling()
	stop := ceil
	debug("ceil: %v depth %v maxdepth@%f:%v\n", ceil, stop,
		m.gradient, stop)
	if stop == 0 {
		return nil, 0
	}

	if stop < minstop {
		stop = minstop
	}

	total += m.Ascend(from, stop, g)

	if true {
		//Ceiling might have moved during the ascent
		if nceil := m.Ceiling(); nceil != ceil {
			ceil = nceil
			stop = ceil
			debug("ceil: resetting %v depth %v maxdepth@%f:%v\n", nceil, stop,
				m.gradient, stop)
		}
		if stop == 0 {
			return nil, 0
		}
		if stop < minstop {
			stop = minstop
		}
	}

	selectDG := func() {
		if len(mod) > 0 && len(dg) > 0 {
			for i, _ := range dg {
				if stop <= mod[i] {
					g = dg[i]
					inc = suggested_inc
				}
			}
		}
	}
	selectDG()

	if laststopdeep {
		m.setDeltaGradient(stop/move - 1)
	} else {
		m.setDeltaGradient(stop / move)
	}

	deco := make([]Segment, 0, 100)
	for i := 0; i < 10000; i++ {
		m.LevelOff(inc, float64(stop), g)

		ceil = m.Ceiling()
		nextstop := ceil
		debug("ceil: %v depth %v %v@%f:%v\n", ceil, nextstop, stop,
			m.gradient, nextstop)
		if stop == nextstop || (laststopdeep == true && stop == (2*move) && nextstop == move) {
			stoplen += inc
		} else {
			stoplen += inc

			if pmodel {
				s := fmt.Sprintf("deco %v %v", stop, stoplen)
				m.Print(true, s)
			}

			deco = append(deco, Segment{stop, stoplen})

			total += stoplen

			if nextstop == 0 {
				return deco, total
			}

			laststop := stop
			if stop-nextstop > move {
				stop -= move
			} else {
				stop = nextstop
			}
			if true {
				total += m.Ascend(laststop, stop, g)
			}

			m.increaseGradient()

			stoplen = 0

			selectDG()
		}
	}
	panic("decompress")

}

// DecompressCCR creates a CCR decompression ascent, storing it on Segment.
// It starts at depth from, and uses ccrPPO2 to change the breathing mix on ascent
// It adjusts gradients factor during the stop of the decompression
//
// laststopdeep controls if the last stop, 10ft/3m, is skipped during decompression
// suggested_inc adjusts the stop len multiplier to make stops a bit longer and even
//
// It returns the total ascent time as second argument
func (m *ZHL16) DecompressCCR(from float64, g *Mix, ccrPPO2 float64,
	laststopdeep bool, suggested_inc float64) ([]Segment, float64) {

	if m.useMeters {
		return m.decompressCCR(from, g, ccrPPO2, laststopdeep, suggested_inc,
			3.0, 21.0, Meter2ATM, roundNext3mStop)
	}
	return m.decompressCCR(from, g, ccrPPO2, laststopdeep, suggested_inc,
		10.0, 70.0, Feet2ATM, roundNext10ftStop)
}

func (m *ZHL16) decompressCCR(from float64, g *Mix, ccrPPO2 float64,
	laststopdeep bool, suggested_inc float64,
	move float64, minstop float64,
	Depth2ATM, roundNextStop func(float64) float64) ([]Segment, float64) {

	total := 0.0
	stoplen := 0.0
	inc := 1.0

	if pmodel {
		m.Print(true, "deco "+"start ")
	}
	ceil := m.Ceiling()
	stop := ceil
	debug("ceil: %v feet %v\n", ceil, stop)
	if stop == 0 {
		return nil, 0
	}

	if stop < minstop {
		stop = minstop
	}

	// TODO: should we the bottom ppo2 to the first stop
	// and then switch at the first deco stop
	nmix := CurrentCCRMix(g, Depth2ATM(stop), ccrPPO2)

	total += m.Ascend(from, stop, nmix)

	if true {
		//Ceiling might have moved during the ascent
		if nceil := m.Ceiling(); nceil != ceil {
			ceil = nceil
			stop = ceil
			debug("ceil: resetting %v depth %v maxdepth@%f:%v\n", nceil, stop,
				m.gradient, stop)
		}
		if stop == 0 {
			return nil, 0
		}
		if stop < minstop {
			stop = minstop
		}
	}

	if stop <= (2 * move) {
		inc = suggested_inc
	} else if stop <= minstop {
		inc = suggested_inc
	}

	if laststopdeep {
		m.setDeltaGradient(stop/10 - 1)
	} else {
		m.setDeltaGradient(stop / 10)
	}

	deco := make([]Segment, 0, 100)
	for i := 0; i < 10000; i++ {
		nmix := CurrentCCRMix(g, Depth2ATM(stop), ccrPPO2)
		m.LevelOff(inc, float64(stop), nmix)

		ceil = m.Ceiling()
		nextstop := ceil
		debug("ceil: %v feet %v\n", ceil, stop)

		if stop == nextstop || (laststopdeep && stop == (2*move) && nextstop == move) {
			stoplen += inc
		} else {
			stoplen += inc

			if pmodel {
				s := fmt.Sprintf("deco %v %v", stop, stoplen)
				m.Print(true, s)
			}

			deco = append(deco, Segment{stop, stoplen})

			total += stoplen

			if nextstop == 0 {
				return deco, total
			}

			laststop := stop
			if stop-nextstop > move {
				stop -= move
			} else {
				stop = nextstop
			}
			if true {
				total += m.Ascend(laststop, stop, g)
			}

			m.increaseGradient()

			stoplen = 0

			if stop <= (2 * move) {
				inc = suggested_inc
			} else if stop <= minstop {
				inc = suggested_inc
			}
		}
	}
	panic("decompress")
}
