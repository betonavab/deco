package deco

import (
	"testing"
)

func cmpDeco(t *testing.T, deco *Profile, r []Segment) {
	if len(deco.Segment) != len(r) {
		t.Errorf("different total len deco %v; want some %v", deco.Segment, r)
		return
	}

	for i, s := range deco.Segment {
		if d := r[i].Depth; s.Depth != d {
			t.Errorf("different  depth at entry %v  %v; want some %v", i, s, r[i])
		}
		if tm := r[i].Time; s.Time != tm {
			t.Errorf("different  time at entry %v  %v; want some %v", i, s, r[i])
		}
	}

}

//
// NewDive(m, trimix, &Profile{segment, descent_speed, ascend_speed,
//							usenitrox, nitrox, useO2
//							laststop20, inc)

func Test_rec32(t *testing.T) {
	nitrox := NewNitrox(32)
	segment := []Segment{
		{100, 30},
	}

	m := ZHL16C(0.30, 0.95)
	d := NewDive(m, nitrox, &Profile{segment, 0.0}, true, 32, false, true, 1)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{70, 1},
		{60, 1},
		{50, 1},
		{40, 1},
		{30, 1},
		{20, 1},
	})
}

func Test_rec32_GUE100_30(t *testing.T) {
	nitrox := NewNitrox(32)
	segment := []Segment{
		{100, 30},
	}

	m := ZHL16C(0.30, 0.95)
	d := NewDive(m, nitrox, &Profile{segment, 0.0}, true, 32, false, true, 1)
	d.UseSimpleDeco()
	d.UseHalfDepth()

	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{50, 1},
		{40, 1},
		{30, 1},
		{20, 1},
		{10, 1},
	})
}

func Test_rec32_GUE90_40(t *testing.T) {
	nitrox := NewNitrox(32)
	segment := []Segment{
		{90, 40},
	}

	m := ZHL16C(0.30, 0.95)
	d := NewDive(m, nitrox, &Profile{segment, 0.0}, true, 32, false, true, 1)
	d.UseSimpleDeco()
	d.UseHalfDepth()

	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{50, 1},
		{40, 1},
		{30, 1},
		{20, 1},
		{10, 1},
	})
}

func Test_rec32_GUE80_50(t *testing.T) {
	nitrox := NewNitrox(32)
	segment := []Segment{
		{80, 50},
	}

	m := ZHL16C(0.30, 0.95)
	d := NewDive(m, nitrox, &Profile{segment, 0.0}, true, 32, false, true, 1)
	d.UseSimpleDeco()
	d.UseHalfDepth()

	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{40, 1},
		{30, 1},
		{20, 1},
		{10, 1},
	})
}

func Test_rec32_GUE70_60(t *testing.T) {
	nitrox := NewNitrox(32)
	segment := []Segment{
		{70, 60},
	}

	m := ZHL16C(0.30, 0.95)
	d := NewDive(m, nitrox, &Profile{segment, 0.0}, true, 32, false, true, 1)
	d.UseSimpleDeco()
	d.UseHalfDepth()

	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{40, 1},
		{30, 1},
		{20, 1},
		{10, 1},
	})
}

func Test_rec32_GUE60_70(t *testing.T) {
	nitrox := NewNitrox(32)
	segment := []Segment{
		{60, 70},
	}

	m := ZHL16C(0.30, 0.95)
	d := NewDive(m, nitrox, &Profile{segment, 0.0}, true, 32, false, true, 1)
	d.UseSimpleDeco()
	d.UseHalfDepth()

	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{30, 1},
		{20, 1},
		{10, 1},
	})
}

func Test_rec32_GUE50_80(t *testing.T) {
	nitrox := NewNitrox(32)
	segment := []Segment{
		{50, 80},
	}

	m := ZHL16C(0.30, 0.95)
	d := NewDive(m, nitrox, &Profile{segment, 0.0}, true, 32, false, true, 1)
	d.UseSimpleDeco()
	d.UseHalfDepth()

	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{30, 1},
		{20, 1},
	})
}

func Test_rec32_GUE40_90(t *testing.T) {
	nitrox := NewNitrox(32)
	segment := []Segment{
		{40, 90},
	}

	m := ZHL16C(0.30, 0.95)
	d := NewDive(m, nitrox, &Profile{segment, 0.0}, true, 32, false, true, 1)
	d.UseSimpleDeco()
	d.UseHalfDepth()

	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{20, 1},
	})
}
func Test_rec32_GUE30_100(t *testing.T) {
	nitrox := NewNitrox(32)
	segment := []Segment{
		{30, 100},
	}

	m := ZHL16C(0.30, 0.95)
	d := NewDive(m, nitrox, &Profile{segment, 0.0}, true, 32, false, true, 1)
	d.UseSimpleDeco()
	d.UseHalfDepth()

	if d.Dive() != nil {
		t.Error("Dive failed")
	}

	deco := d.Decompress()
	if deco != nil {
		t.Error("has deco; want none")
	}
}
func Test_rec32_Multi(t *testing.T) {
	nitrox := NewNitrox(32)
	segment := []Segment{
		{100, 20},
		{60, 20},
		{40, 20},
	}

	m := ZHL16C(0.30, 0.95)
	d := NewDive(m, nitrox, &Profile{segment, 0.0}, true, 32, false, true, 1)
	d.UseSimpleDeco()
	d.UseHalfDepth()

	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{20, 1},
		{10, 1},
	})
}

func Test_rec3_1(t *testing.T) {
	trimix := NewTrimix(30, 30)
	segment := []Segment{
		{100, 40},
	}

	m := ZHL16C(0.20, 0.85)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 32, false, true, 2)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{70, 2},
		{60, 2},
		{50, 2},
		{40, 2},
		{30, 2},
		{20, 10},
	})
}

func Test_rec3_2(t *testing.T) {
	trimix := NewTrimix(21, 35)
	segment := []Segment{
		{120, 20},
	}

	m := ZHL16C(0.20, 0.85)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 32, false, true, 1)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{70, 1},
		{60, 1},
		{50, 1},
		{40, 1},
		{30, 1},
		{20, 7},
	})
}

func Test_tech1_class1(t *testing.T) {
	trimix := NewTrimix(21, 35)
	segment := []Segment{
		{150, 20},
	}

	m := ZHL16C(0.20, 0.85)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 50, false, true, segment[0].Time/10)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{70, 2},
		{60, 2},
		{50, 2},
		{40, 2},
		{30, 2},
		{20, 8},
	})
}

func Test_tech1_class2(t *testing.T) {
	trimix := NewTrimix(18, 45)
	segment := []Segment{
		{170, 20},
	}

	m := ZHL16C(0.20, 0.85)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 50, false, true, segment[0].Time/10)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{90, 1},
		{80, 1},
		{70, 2},
		{60, 2},
		{50, 2},
		{40, 2},
		{30, 2},
		{20, 18},
	})
}

func Test_tech1_class2_meter(t *testing.T) {
	trimix := NewTrimix(18, 45)
	segment := []Segment{
		{51, 20},
	}

	m := ZHL16C(0.20, 0.85)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 50, false, true, segment[0].Time/10)
	d.UseMeter()
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{27, 1},
		{24, 1},
		{21, 2},
		{18, 2},
		{15, 2},
		{12, 2},
		{9, 2},
		{6, 18},
	})
}
func Test_K1(t *testing.T) {
	trimix := NewTrimix(21, 40)
	segment := []Segment{
		{150, 15},
		{130, 15},
	}

	m := ZHL16C(0.20, 0.90)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 50, false, true, 3)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{70, 3},
		{60, 3},
		{50, 3},
		{40, 3},
		{30, 3},
		{20, 15},
	})
}

func Test_K1_deep(t *testing.T) {
	trimix := NewTrimix(18, 45)
	segment := []Segment{
		{170, 15},
		{130, 15},
	}

	m := ZHL16C(0.20, 0.90)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 50, false, true, 3)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{80, 1},
		{70, 3},
		{60, 3},
		{50, 3},
		{40, 3},
		{30, 3},
		{20, 21},
	})
}

func Test_tech2_class1(t *testing.T) {
	trimix := NewTrimix(15, 55)
	segment := []Segment{
		{220, 20},
	}

	m := ZHL16C(0.20, 0.85)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 50, true, true, 5)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{130, 1},
		{120, 1},
		{110, 1},
		{100, 1},
		{90, 2},
		{80, 2},
		{70, 5},
		{60, 5},
		{50, 5},
		{40, 5},
		{30, 5},
		{20, 25},
	})
}

func Test_tech2_class1_meter(t *testing.T) {
	trimix := NewTrimix(15, 55)
	segment := []Segment{
		{66, 20},
	}

	m := ZHL16C(0.20, 0.85)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 50, true, true, 5)
	d.UseMeter()
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{39, 1},
		{36, 1},
		{33, 1},
		{30, 1},
		{27, 2},
		{24, 2},
		{21, 5},
		{18, 5},
		{15, 5},
		{12, 5},
		{9, 5},
		{6, 25},
	})
}

func Test_K2(t *testing.T) {
	trimix := NewTrimix(18, 45)
	segment := []Segment{
		{220, 20},
		{150, 15},
	}

	m := ZHL16C(0.20, 0.90)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 50, true, true, 5)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{100, 1},
		{90, 2},
		{80, 3},
		{70, 5},
		{60, 5},
		{50, 5},
		{40, 5},
		{30, 5},
		{20, 30},
	})
}

func Test_HN_OC(t *testing.T) {
	trimix := NewTrimix(18, 45)
	segment := []Segment{
		{150, 90},
	}

	m := ZHL16C(0.30, 0.90)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 50, true, true, 8)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{90, 5},
		{80, 6},
		{70, 8},
		{60, 8},
		{50, 8},
		{40, 8},
		{30, 16},
		{20, 72},
	})
}

func Test_HN_CCR(t *testing.T) {
	trimix := NewTrimix(18, 45)
	segment := []Segment{
		{150, 90},
	}

	m := ZHL16C(0.30, 0.90)
	d := NewDive(m, trimix, &Profile{segment, 0.0}, true, 50, true, true, 8)
	d.SetCCR(1.2, 1.4)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{80, 4},
		{70, 8},
		{60, 8},
		{50, 8},
		{40, 8},
		{30, 16},
		{20, 64},
	})
}

func Test_rec32_GUE30_30(t *testing.T) {
	nitrox := NewNitrox(32)
	segment := []Segment{
		{30, 30},
	}

	m := ZHL16C(0.30, 0.95)
	d := NewDive(m, nitrox, &Profile{segment, 0.0}, true, 32, false, true, 1)
	d.UseSimpleDeco()
	d.UseHalfDepth()
	d.UseMeter()

	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	cmpDeco(t, deco, []Segment{
		{15, 1},
		{12, 1},
		{9, 1},
		{6, 1},
		{3, 1},
	})
}

func Test_LEM_150_60(t *testing.T) {
	heliox := NewHeliox(18)
	segment := []Segment{
		{150, 60},
	}

	m := XVal_He_9_040_fsw()
	d := NewDive(m, heliox, &Profile{segment, 0.0}, false, 50, false, true, 8)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	d.SetCCR(1.3, 1.3)
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	// Difference from LEM report after //
	cmpDeco(t, deco, []Segment{
		{60, 3},	// 2
		{50, 3},	// 4
		{40, 5},	// 4
		{30, 10},	// 11
		{20, 57},	// 58
	})
}

func Test_LEM_150_90(t *testing.T) {
	heliox := NewHeliox(18)
	segment := []Segment{
		{150, 90},
	}

	m := XVal_He_9_040_fsw()
	d := NewDive(m, heliox, &Profile{segment, 0.0}, false, 50, false, true, 8)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	d.SetCCR(1.3, 1.3)
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	// Difference from LEM report after //
	cmpDeco(t, deco, []Segment{
		{70, 3},	// 2
		{60, 3},	// 4
		{50, 4},
		{40, 11},	// 10
		{30, 10},	// 11
		{20, 113},	// 115
	})
}

func Test_LEM_250_40(t *testing.T) {
	heliox := NewHeliox(15)
	segment := []Segment{
		{250, 40},
	}

	m := XVal_He_9_040_fsw()
	d := NewDive(m, heliox, &Profile{segment, 0.0}, false, 50, false, true, 8)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	d.SetCCR(1.3, 1.3)
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	// Difference from LEM report after //
	cmpDeco(t, deco, []Segment{
		{140, 1},	
		{130, 4},
		{120, 3},
		{110, 3},	//4
		{100, 3},
		{90, 5},
		{80, 10},	//11
		{70, 10},	//11
		{60, 11},	//10
		{50, 10},	//11
		{40, 11},
		{30, 10},	//11
		{20, 117},	//118
	})
}
func Test_LEM_300_30(t *testing.T) {
	heliox := NewHeliox(12)
	segment := []Segment{
		{300, 30},
	}

	m := XVal_He_9_040_fsw()
	d := NewDive(m, heliox, &Profile{segment, 0.0}, false, 50, false, true, 8)
	if d.Dive() != nil {
		t.Error("Dive failed")
	}
	d.SetCCR(1.3, 1.3)
	deco := d.Decompress()
	if deco == nil {
		t.Error("has no deco; want some")
	}
	// Difference from LEM report after //
	cmpDeco(t, deco, []Segment{
		{160, 4},	//3	
		{150, 3},	//4	
		{140, 3},	
		{130, 4},
		{120, 3},	//4
		{110, 3},	
		{100, 4},
		{90, 3},	//5
		{80, 10},	//11
		{70, 11},
		{60, 10},
		{50, 11},	
		{40, 10},	//11
		{30, 10},	//11
		{20, 125},	//126
	})
}


