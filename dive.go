package deco

import (
	"errors"
	"fmt"
	"io"
	"math"
)

// Mix struct is used to represent a diving mix, Nitrox, Trimix, Heliox. They are usually create
// by calling the NewMix function for the type of gas, like for example, NewNitrox
type Mix struct {
	o2    float64
	gas   map[string]float64
	ismix bool
}

func (m *Mix) String() string {
	return fmt.Sprintf("%v/%v/%v",
		m.o2*100, m.gas["N2"]*100, m.gas["He"]*100)
}

// NewNitrox creates a Nitrox mix giving its o2 content
func NewNitrox(O2 float64) *Mix {

	m := &Mix{}
	m.o2 = O2 / 100
	m.gas = make(map[string]float64)
	m.gas["He"] = 0
	m.gas["N2"] = (100 - O2) / 100
	m.ismix = false

	return m
}

// NewHeliox creates a Heliox mix giving its o2 content
func NewHeliox(O2 float64) *Mix {

	m := &Mix{}
	m.o2 = O2 / 100
	m.gas = make(map[string]float64)
	m.gas["He"] = (100 - O2) / 100
	m.gas["N2"] = 0
	m.ismix = true

	return m
}

// NewTrimix creates a Trimix mix giving its o2 content and He content
func NewTrimix(O2, He float64) *Mix {

	m := &Mix{}
	m.o2 = O2 / 100
	m.gas = make(map[string]float64)
	m.gas["He"] = He / 100
	m.gas["N2"] = (100 - O2 - He) / 100
	m.ismix = true

	return m
}

// Segment struct is used to hold a section of a dive. They are used as []Segment by
// function acting on dive profiles, as for example, Decompress families return []Segment.
type Segment struct {
	Depth float64
	Time  float64
}

func (s *Segment) String() string {
	return fmt.Sprintf("%v %v", s.Depth, s.Time)
}

// Profile struct is used to hold a dive profile, which is described by a series of
// segments, and some basic parameters
type Profile struct {
	Segment  []Segment
	Duration float64
}

func (p *Profile) Print() float64 {
	total := 0.0
	for _, s := range p.Segment {
		total += s.Time
		fmt.Printf("%v %v\n", s.Depth, s.Time)
	}
	return total
}

func Feet2ATM(feet float64) float64 {
	return (feet + 33) / 33
}

func ATM2Feet(atm float64) float64 {
	return atm*33 - 33
}

func Meter2ATM(meter float64) float64 {
	return (meter + 10) / 10
}

func ATM2Meter(atm float64) float64 {
	return atm*10 - 10
}

// Model interface is used to access a family of decompression model that
// can track gas absortion and elimination in different ways, and are capable
// of calculating a decompression profiles that should be followed after a series
// of exposure segment. ZHL16 and XVal_He_9_040 are examples of such models.
// A model is usually attached to a Dive struct which stores dive parameters,
// as for example, bottom and deco gases, decompression styles, feet/meters and others.
// Depth values are passes to model's function in feet or meters and will be threated depending
// on UseMeters function. 
type Model interface {
	Print(total bool, label string)
	UseMeters()

	Descend(to float64, mix *Mix) float64
	LevelOff(time float64, to float64, mix *Mix)
	Ascend(from, to float64, mix *Mix) float64

	Ceiling() float64

	Decompress(from float64, g *Mix, halfdepth bool) ([]Segment, float64)
	DecompressOCTech(from float64, g *Mix, dg []*Mix,
		laststopdeep bool, suggested_inc float64) ([]Segment, float64)
	DecompressCCR(from float64, g *Mix, ccrPPO2 float64,
		laststopdeep bool, suggested_inc float64) ([]Segment, float64)
}

// Dive struct is use to managed a dive profile and its decompression
type Dive struct {
	m      Model
	mix    *Mix
	bottom *Profile
	deco   *Profile

	lastdepth float64

	usenitrox    bool
	nitrox       int
	useO2        bool
	laststopdeep bool
	inc          float64

	useccr      bool
	ccrPP02     float64
	decoccrPP02 float64

	usehalfdepth  bool
	usesimpledeco bool
	usemeter      bool
}

func NewDive(m Model, mix *Mix, prof *Profile, usenitrox bool, nitrox int, useO2 bool, laststopdeep bool, inc float64) *Dive {
	d := &Dive{}

	d.m = m
	d.mix = mix
	d.bottom = prof
	d.useO2 = useO2
	d.laststopdeep = laststopdeep
	d.inc = inc

	d.usenitrox = true
	d.nitrox = 50
	if usenitrox {
		d.nitrox = nitrox
	}

	d.useccr = false
	d.ccrPP02 = 0.0
	d.decoccrPP02 = 0.0

	d.usehalfdepth = false
	d.usesimpledeco = false
	d.usemeter = false

	d.lastdepth = 0.0
	d.deco = nil

	return d
}

func (d *Dive) SetCCR(ppo2 float64, decoppo2 float64) {
	d.useccr = true
	d.ccrPP02 = ppo2
	d.decoccrPP02 = decoppo2
}

func (d *Dive) UseHalfDepth() {
	d.usehalfdepth = true
}

func (d *Dive) UseSimpleDeco() {
	d.usesimpledeco = true
}
func (d *Dive) UseMeter() {
	d.usemeter = true
	d.m.UseMeters()
}

func (d *Dive) Print() {
	if d.m != nil || d.mix != nil {
		fmt.Printf("%v %v\n", d.m, d.mix)
	}
	if d.bottom != nil {
		fmt.Println("Profile")
		d.bottom.Print()
	}
	fmt.Println()
	if d.deco != nil {
		if d.useccr {
			fmt.Printf("CCR deco %v\n", d.decoccrPP02)

		} else if d.usesimpledeco == false {
			if d.useO2 {
				fmt.Printf("Deco using %v and 100 inc %v\n", d.nitrox, d.inc)

			} else {
				fmt.Printf("Deco using %v inc %v\n", d.nitrox, d.inc)
			}
		}

		d.deco.Print()
		fmt.Printf("total deco %v runtime %v\n",
			math.Round(d.deco.Duration),
			math.Round(d.bottom.Duration+d.deco.Duration))
	} else {
		fmt.Printf("no deco runtime %v\n",
			math.Round(d.bottom.Duration))
	}
}

func (d *Dive) Dive() error {

	if d.useccr {
		return d.DiveCCR()
	}

	d.lastdepth = 0.0
	total := 0.0
	m := d.m

	for i, s := range d.bottom.Segment {
		depth := s.Depth
		time := s.Time

		if i == 0 {
			dt := m.Descend(depth, d.mix)
			m.LevelOff(time-dt, depth, d.mix)
		} else {
			ceil := m.Ceiling()
			if stop := ceil; depth < stop {
				return errors.New("model: ceiling reached before next segment")
			}

			at := m.Ascend(d.lastdepth, depth, d.mix)
			m.LevelOff(time-at, depth, d.mix)
		}
		total += s.Time
		d.lastdepth = depth
	}
	d.bottom.Duration = total
	return nil
}

func CurrentCCRMix(mix *Mix, atm float64, ccrPPO2 float64) *Mix {
	fO2 := ccrPPO2 / atm
	if fO2 > 1.0 {
		fO2 = 1.0
	}
	fDil := 1.0 - fO2
	fHe := fDil * mix.gas["He"]
	if fHe == 0 {
		return NewNitrox(math.Round(fO2 * 100))
	} else {
		return NewTrimix(math.Round(fO2*100), math.Round(fHe*100))
	}
}

func (d *Dive) CurrentCCRMix(depth float64) *Mix {

	if !d.useccr {
		return d.mix
	}

	atm := Feet2ATM(depth)
	if d.usemeter {
		atm = Meter2ATM(depth)
	}
	ppO2 := d.ccrPP02
	fO2 := ppO2 / atm
	if fO2 > 1.0 {
		fO2 = 1.0
	}
	fDil := 1.0 - fO2
	fHe := fDil * d.mix.gas["He"]
	if fHe == 0 {
		return NewNitrox(math.Round(fO2 * 100))
	} else {
		return NewTrimix(math.Round(fO2*100), math.Round(fHe*100))
	}
}

func (d *Dive) DiveCCR() error {

	d.lastdepth = 0.0
	total := 0.0
	m := d.m

	for i, s := range d.bottom.Segment {
		depth := s.Depth
		time := s.Time

		nmix := d.CurrentCCRMix(depth)
		if i == 0 {
			// TODO: Should we have a CCR Descend function
			dt := m.Descend(depth, nmix)
			m.LevelOff(time-dt, depth, nmix)
		} else {
			ceil := m.Ceiling()
			if stop := ceil; depth < stop {
				return errors.New("model: ceiling reached before next segment")
			}
			// TODO: Should we have a CCR Ascend function
			at := m.Ascend(d.lastdepth, depth, nmix)
			m.LevelOff(time-at, depth, nmix)
		}
		total += s.Time
		d.lastdepth = depth
	}
	d.bottom.Duration = total
	return nil
}

func (d *Dive) DecompressCCR() *Profile {
	var seg []Segment
	total := 0.0
	seg, total = d.m.DecompressCCR(d.lastdepth, d.mix, d.decoccrPP02, d.laststopdeep, d.inc)
	if seg == nil {
		return nil
	}
	p := &Profile{seg, total}
	d.deco = p
	return p

}

func (d *Dive) Decompress() *Profile {
	if d.useccr {
		return d.DecompressCCR()
	}

	var seg []Segment
	total := 0.0
	if d.usesimpledeco {
		seg, total = d.m.Decompress(d.lastdepth, d.mix, d.usehalfdepth)
	} else {
		dg := make([]*Mix, 0, 5)
		if d.usenitrox {
			dg = append(dg, NewNitrox(float64(d.nitrox)))
			if d.useO2 {
				dg = append(dg, NewNitrox(100))
			}
		} else {
			dg = append(dg, NewNitrox(50))
		}
		seg, total = d.m.DecompressOCTech(d.lastdepth, d.mix, dg, d.laststopdeep, d.inc)
	}
	if seg == nil {
		return nil
	}
	p := &Profile{seg, total}
	d.deco = p
	return p
}

var xdebug bool
var dwriter io.Writer     

var pmodel bool
var mwriter io.Writer  

func EnablePmodel(w io.Writer) {
	pmodel = true
	mwriter = w
}

func DisablePmodel() {
	pmodel = true
}

func EnableDebug(w io.Writer) {
	xdebug = true
	dwriter = w
}

func DisableDebug() {
	xdebug = false
}

func debug(format string, args ...interface{}) {
	if xdebug {
		fmt.Fprintf(dwriter,format,args...)
	}
}
