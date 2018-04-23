package rand

import (
	"github.com/ReconfigureIO/fixed"
)

const (
	c    = 256
	mask = c - 1
	r    = 220
	rInv = 18
)

// Modified Ziggurat Algorithm based upon the paper
// "Hardware-Optimized Ziggurat Algorithm for High-Speed Gaussian Random Number Generators"

// 2.6 width integer
type I6 int8

func I6F(a int8) I6 {
	return I6(a & 0x3f)
}

func (a I6) Mul(b I6) I6 {
	return I6((int16(a) * int16(b)) >> 6)
}

// restricted ln from [0, 1) using a 32 lookup table
func log(x I6) I6 {
	return [32]I6{-0, -1, -3, -5, -7, -9, -10, -12, -14, -15, -17, -18, -20, -21, -23, -24, -25, -27, -28, -29, -31, -32, -33, -34, -35, -36, -38, -39, -40, -41, -42, -43}[(x>>1)&0x1f]
}

type param struct {
	X     I6
	F     I6
	M     I6
	tMask uint8
}

// Normals writes a stream of Int26_6, normally distributed
func (rand Rand) Normals(output chan<- fixed.Int26_6) {
	a := rand.Iteration()

	uint32s := make(chan uint32, 1)
	tailUs := make(chan uint32, 1)

	tailRand := New(a)

	go rand.Uint32s(uint32s)

	go tailRand.Uint32s(tailUs)

	for {
		u := int32(<-uint32s)
		// the index we'll use
		i := uint8(u & mask)
		negate := u < 0

		// see cmd/tables/main.go for how to generate it
		tmp := [c]uint32{0, 219611136, 303166144, 353365722, 403565797, 437054187, 470542574, 487254001, 520808179, 537519604, 554296566, 587850743, 604562167, 621339128, 638116344, 654893305, 671604985, 688381946, 705159162, 721936122, 738713339, 738713339, 755490299, 772201979, 788979195, 788978940, 805756156, 822533372, 839310588, 839310588, 856087548, 872864764, 872864764, 889641980, 906353661, 906353661, 923130877, 923130621, 939907837, 956685053, 956685053, 973462269, 973462269, 990239485, 990239485, 1007016701, 1023793917, 1023793917, 1040571133, 1040571133, 1057348349, 1057348349, 1074125565, 1074060029, 1090837245, 1090837246, 1107614462, 1107614462, 1124391678, 1124391678, 1141168894, 1141168894, 1157946110, 1157946110, 1174723326, 1174723326, 1191500542, 1191500542, 1208277758, 1208277758, 1208277758, 1225054974, 1225054974, 1241832446, 1241832446, 1258609662, 1258609662, 1275386878, 1275386878, 1292164094, 1292164094, 1292164094, 1308941310, 1308941566, 1325718782, 1325718782, 1342495998, 1342495998, 1359207678, 1359207678, 1359207678, 1375984894, 1375985150, 1392762366, 1392762366, 1409539582, 1409539582, 1409539582, 1426316798, 1426316798, 1443094270, 1443094270, 1459871486, 1459871486, 1459871486, 1476648702, 1476648702, 1493426174, 1493426174, 1510203390, 1510203390, 1510203390, 1526980606, 1526980862, 1543758078, 1543758078, 1560535294, 1560535294, 1560535294, 1577312766, 1577312766, 1594089982, 1594089982, 1610867198, 1610867198, 1610867454, 1627644670, 1627644670, 1644421886, 1644421886, 1661199102, 1661199358, 1677976574, 1677976574, 1677976574, 1694753790, 1694754046, 1711531262, 1711531262, 1728308478, 1728308478, 1728308734, 1745085950, 1745085950, 1761863166, 1761863166, 1778640638, 1778640638, 1795417854, 1795417854, 1812129534, 1812129790, 1812129790, 1828907006, 1828907006, 1845684222, 1845684478, 1862461694, 1862461694, 1879238910, 1879238910, 1896016382, 1896016382, 1912793598, 1912793598, 1929571070, 1929571070, 1946348286, 1946348286, 1963125502, 1963125758, 1979902974, 1979902974, 1996680190, 1996680446, 2013457662, 2013457662, 2030234878, 2030234878, 2047012350, 2047012350, 2063789566, 2063789566, 2080567038, 2097344254, 2097344254, 2114121470, 2114121470, 2130898942, 2130898942, 2147676158, 2164453374, 2164453630, 2181230846, 2181230846, 2198008062, 2214785534, 2214785534, 2231562750, 2231562750, 2248340222, 2265117438, 2265117438, 2281894654, 2298672126, 2298672126, 2315449342, 2332226558, 2349004030, 2349004030, 2365781246, 2382558462, 2399335934, 2399335934, 2416113150, 2432890366, 2449667838, 2466445054, 2466445054, 2483222270, 2499999742, 2516776958, 2533554174, 2550331390, 2567108862, 2583886078, 2600663294, 2617440510, 2634217982, 2650995198, 2667772414, 2684549630, 2701261566, 2734815998, 2751593214, 2768370685, 2801925117, 2818702333, 2852256765, 2869034237, 2902588669, 2919365885, 2952920317, 2986475005, 3020029437, 3053583868, 3103915516, 3137470204, 3187801852, 3238133499, 3305242619, 3372351482, 3456237561, 3556900856, 3691053302, 3909157105}[i]
		// choose our param param and unpack them into a struct
		// Unpacking should be nearly free
		p := param{I6(tmp >> 24), I6(tmp >> 16), I6(tmp >> 8), uint8(tmp)}

		// use u as an 0.8 fixed point from [0..1)
		t := I6(uint32(u) >> 6)

		// 2.6
		x := p.X

		z := I6((int16(t) * int16(x)) >> 8)

		if uint8(t) < p.tMask {
			// Bulk, this path should happen very frequently
			if negate {
				output <- fixed.Int26_6(-z)
			} else {
				output <- fixed.Int26_6(z)
			}
		} else if i == 0 {
			// Tail
			var x2 I6
			k := true
			for k {
				var vars [2]I6
				rand := <-tailUs
				for i := 0; i < 2; i++ {
					vars[i] = log(I6F(int8(rand)))
					rand = rand >> 6
				}
				t := vars[0]
				x2 = t.Mul(rInv)

				y := vars[1] << 1
				k = y < x2.Mul(x2)
			}
			if negate {
				output <- fixed.Int26_6(-r - int16(x2))
			} else {
				output <- fixed.Int26_6(r + int16(x2))
			}
		} else {
			// wedge
			f := p.F

			// 0.8
			t := I6(<-tailUs)

			// 2.14
			m := p.M
			check := m.Mul(z - x)

			// 2.14
			y := f.Mul(t)

			if y < check {
				if u < 0 {
					output <- fixed.Int26_6(-y)
				} else {
					output <- fixed.Int26_6(y)
				}
			}
		}

	}

}
