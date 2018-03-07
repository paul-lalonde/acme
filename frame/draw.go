package frame

import (
	"image"
//	"log"

	"9fans.net/go/draw"
)

func (f *Frame) DrawText(pt image.Point, text *draw.Image, back *draw.Image) {
//	log.Println("DrawText at", pt, "NoRedraw", f.NoRedraw, text)
	for nb := 0; nb < f.nbox; nb++ {
		b := f.box[nb]
		f.cklinewrap(&pt, b)
//		log.Printf("box [%d] %#v pt %v NoRedraw %v nrune %d\n",  nb, string(b.Ptr), pt, f.NoRedraw, b.Nrune)

		if !f.NoRedraw && b.Nrune >= 0 {
			f.Background.Bytes(pt, text, image.ZP, f.Font.Impl(), b.Ptr)
		}
		pt.X += b.Wid
	}
}

func (f *Frame) DrawSel(pt image.Point, p0, p1 int, issel bool) {
//	log.Println("DrawSel")
	var back, text *draw.Image

	if f.Ticked {
		f.Tick(f.Ptofchar(f.p0), false)
	}

	if p0 == p1 {
		f.Tick(pt, issel)
		return
	}

	if issel {
		back = f.Cols[ColHigh]
		text = f.Cols[ColHText]
	} else {
		back = f.Cols[ColBack]
		text = f.Cols[ColText]
	}

	f.drawsel0(pt, p0, p1, back, text)
}

func (f *Frame) drawsel0(pt image.Point, p0, p1 int, back *draw.Image, text *draw.Image) image.Point {
//	log.Println("drawsel0")
	p := 0
	bi := 0
	b := f.box[bi]
	trim := false
	i := 0
	x := 0
	var w int

	for nb := 0; nb < f.nbox && p < p1; nb++ {
		nr := b.Nrune
		if nr < 0 {
			nr = 1
		}
		if p+nr <= p0 {
			goto Continue
		}
		if p >= p0 {
			qt := pt
			f.cklinewrap(&pt, b)
			// fill in the end of a wrapped line
			if pt.Y > qt.Y {
				f.Background.Draw(image.Rect(qt.X, qt.Y, f.Rect.Max.X, pt.Y), back, nil, qt)
			}
		}
		i = 0
		if p < p0 {
			// beginning of region: advance into box
			i += len(b.Ptr[:int(p0)-p])
			nr -= int(p0) - p
			p = int(p0)
		}
		trim = false
		if p+nr > p1 {
			// end of region: trim box
			nr -= (p + nr) - int(p1)
			trim = true
		}

		if b.Nrune < 0 || nr == b.Nrune {
			w = b.Wid
		} else {
			// Corresponds to the native code but does the wrong thing if frbox.Nrune is
			// is actually the number of runes (as opposed to the number of bytes)
			// In that case, this code and the code below would fail on UTF code points
			// that are more than one byte each.
			//
			// Given that the native code in frdraw.c also has this issue, I'll revisit this
			// problem later.
			w = f.Font.StringWidth(string(b.Ptr[i : i+nr]))
		}
		x = pt.X + w
		if x > f.Rect.Max.X {
			x = f.Rect.Max.X
		}
		f.Background.Draw(image.Rect(pt.X, pt.Y, x, pt.Y+f.Font.DefaultHeight()), back, nil, pt)
		if b.Nrune >= 0 {
			// See comment above. Same issue applies.
			f.Background.String(pt, text, image.ZP, f.Font.Impl(), string(b.Ptr[i:i+nr]))
		}
		pt.X += w
	Continue:
		bi++
		b = f.box[bi]
		p += nr
	}

	if p1 > p0 && bi > 0 && bi < f.nbox && f.box[bi-1].Nrune > 0 && !trim {
		qt := pt
		f.cklinewrap(&pt, b)
		if pt.Y > qt.Y {
			f.Background.Draw(image.Rect(qt.X, qt.Y, f.Rect.Max.X, pt.Y), back, nil, qt)
		}
	}
	return pt
}

func (f *Frame) Redraw() {
//	log.Println("Redraw")
	Ticked := false
	var pt image.Point

	if f.p0 == f.p1 {
		Ticked = f.Ticked
		if Ticked {
			f.Tick(f.Ptofchar(f.p0), false)
		}
		f.drawsel0(f.Ptofchar(0), 0, f.NChars, f.Cols[ColBack], f.Cols[ColText])
		if Ticked {
			f.Tick(f.Ptofchar(f.p0), true)
		}
	}

	pt = f.Ptofchar(0)
	pt = f.drawsel0(pt, 0, f.p0, f.Cols[ColBack], f.Cols[ColText])
	pt = f.drawsel0(pt, f.p0, f.p1, f.Cols[ColHigh], f.Cols[ColHText])
	pt = f.drawsel0(pt, f.p1, f.NChars, f.Cols[ColBack], f.Cols[ColText])

}

func (f *Frame) tick(pt image.Point, Ticked bool) {
//	log.Println("_tick")
	if f.Ticked == Ticked || f.TickImage == nil || !pt.In(f.Rect) {
		return
	}

	pt.X -= f.TickScale
	r := image.Rect(pt.X, pt.Y, pt.X+frtickw*f.TickScale, pt.Y+f.Font.DefaultHeight())

	if r.Max.X > f.Rect.Max.X {
		r.Max.X = f.Rect.Max.X
	}
	if Ticked {
		f.TickBack.Draw(f.TickBack.R, f.Background, nil, pt)
		f.Background.Draw(r, f.TickImage, nil, image.ZP)
	} else {
		f.Background.Draw(r, f.TickBack, nil, image.ZP)
	}
	f.Ticked = Ticked
}

func (f *Frame) Tick(pt image.Point, Ticked bool) {
//	log.Println("Tick")
	if f.TickScale != f.Display.ScaleSize(1) {
		if f.Ticked {
			f.tick(pt, false)
		}
		f.InitTick()
	}
	f.tick(pt, Ticked)
}

func (f *Frame) _draw(pt image.Point) image.Point {
//	log.Println("_draw")
	for nb := 0; nb < f.nbox; nb++ {
		b := f.box[nb]
		f.cklinewrap0(&pt, b)
		if pt.Y == f.Rect.Max.Y {
			f.NChars -= f.strlen(nb)
			f.delbox(nb, f.nbox-1)
			break
		}

		if b.Nrune > 0 {
			n, fits := f.canfit(pt, b)
			if !fits {
				break
			}
			if n != b.Nrune {
				f.splitbox(nb, n)
				b = f.box[nb]
			}
			pt.X += b.Wid
		} else {
			if b.Bc == '\n' {
				pt.X = f.Rect.Min.X
				pt.Y += f.Font.DefaultHeight()
			} else {
				pt.X += f.newwid(pt, b)
			}
		}
	}
	return pt
}

func (f *Frame) strlen(nb int) int {
//	log.Println("strlen")
	var n int
	for n = 0; nb < f.nbox; nb++ {
		n += nrune(f.box[nb])
	}
	return n
}
