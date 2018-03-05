package main

import (
	"image"
	"fmt"

	"github.com/rjkroege/acme/frame"
)

var (
	Lheader = []string{
		"New ",
		"Cut ",
		"Paste ",
		"Snarf ",
		"Sort ",
		"Zerox ",
		"Delcol ",
	}
)

type Column struct {
	r    image.Rectangle
	tag  Text
	row  *Row
	w    []*Window
	safe bool
}

func (c *Column)nw() int {
	return len(c.w)
}

func (c *Column)Init(r image.Rectangle) *Column {
	if c == nil {
		c = &Column{}
	}
	c.w = []*Window{}
fmt.Println("init", r)
	display.ScreenImage.Draw(r, display.White, nil, image.ZP)
	c.r = r
	tagfile := NewTagFile()
	tagr := r
	tagr.Max.Y = tagr.Min.Y + tagfont.Height
	c.tag.Init(tagfile, tagr, tagfont, tagcolors)
	c.tag.what = Columntag
	c.tag.col = c
	return c
}

/*
func (c *Column) AddFile(f *File) *Window {
	w := NewWindow(f)
	c.Add(w, nil, 0)
}
*/

func (c *Column) Add(w, clone *Window, y int) *Window {
	// Figure out new window placement
	var v *Window
	var ymax int
	r := c.r
	r.Min.Y = c.tag.fr.Rect.Max.Y + display.ScaleSize(Border)
	if y < r.Min.Y && c.nw() > 0 { // Steal half the last window
		v = c.w[c.nw() - 1]
		y = v.body.fr.Rect.Min.Y + v.body.fr.Rect.Dx()/2
	}
	// Which window will we land on?
	var windex int
	for windex := range c.w {
		v = c.w[windex]
		if y < v.r.Max.Y {
			break
		}
	}
	buggered := false 	// historical variable name
	if c.nw() > 0 {
		if windex < c.nw() {
			windex++
		}
		/*
		 * if landing window (v) is too small, grow it first.
		 */
		minht := v.tag.fr.Font.DefaultHeight() + display.ScaleSize(Border) + 1
		j := 0
		for !c.safe || v.body.fr.MaxLines < 3 || v.body.all.Dy() <= minht {
			j++
			if j > 10 {
				buggered = true // Too many windows in column
				break
			}
			c.Grow(v, 1)
		}
		
		/*
		 * figure out where to split v to make room for w
		 */
	
		/* new window stops where next window begins */
		if windex < c.nw() {
			ymax = c.w[windex].r.Min.Y-display.ScaleSize(Border)
		} else {
			ymax = c.r.Max.Y
		}

		/* new window must start after v's tag ends */
		y = max(y, v.tagtop.Max.Y+display.ScaleSize(Border))

		/* new window must start early enough to end before ymax */
		y = min(y, ymax - minht)

		/* if y is too small, too many windows in column */
		if y < v.tagtop.Max.Y+display.ScaleSize(Border) {
			buggered = true
		}

		// Resize & redraw v
		r = v.r
		r.Max.Y = ymax
	fmt.Println("A", r)
		display.ScreenImage.Draw(r, textcolors[frame.ColBack], nil, image.ZP)
		r1 := r
		y = min(y, ymax-(v.tag.fr.Font.DefaultHeight()*v.taglines+v.body.fr.Font.DefaultHeight()+display.ScaleSize(Border)+1))
		r1.Max.Y = min(y, v.body.fr.Rect.Min.Y+v.body.fr.Nlines*v.body.fr.Font.DefaultHeight())
		r1.Min.Y = v.Resize(r1, false, false)
		r1.Max.Y = r1.Min.Y+display.ScaleSize(Border)
	fmt.Println("B", r1)
		display.ScreenImage.Draw(r1, display.Black, nil, image.ZP)
		
		/*
		 * leave r with w's coordinates
		 */
		r.Min.Y = r1.Max.Y
	}
	if w == nil {
		w = NewWindow()
		w.col = c
	fmt.Println("C", r)
		display.ScreenImage.Draw(r, textcolors[frame.ColBack], nil, image.ZP)
		w.Init(clone, r)
	} else {
		w.col = c
		w.Resize(r, false, true)
	}
	w.tag.col = c
	w.tag.row = c.row
	w.body.col = c
	w.body.row = c.row
	c.w = append(c.w, w)
	c.safe = true
	if buggered {
		c.Resize(c.r)
	}
	// savemouse(w) // TODO(flux): Mouse?
	// mousectl.Moveto(w.tag.Scrollr.Max.Add(image.Point(3,3)))
	barttext = &w.body
	return w
}

func (c *Column) Close(w *Window, dofree bool) {

}

func (c *Column) CloseAll() {

}

func (c *Column) MouseBut() {

}

func (c *Column) Resize(r image.Rectangle) {
	// clearmouse() // TODO(flux) Mouse?
	r1 := r
	r1.Max.Y = r1.Min.Y + c.tag.fr.Font.DefaultHeight()
	c.tag.Resize(r1, true) // And draw the tag
	// TODO(flux): Column button
	//display.ScreenImage.Draw(c.tag.scrollr, colbutton, nil, colbutton->r.Min)
	r1.Min.Y = r1.Max.Y // Walk past the tag
	r1.Max.Y += display.ScaleSize(Border)
fmt.Println("D", r1)
	display.ScreenImage.Draw(r1, display.Black, nil, image.ZP)
	r1.Max.Y = r.Max.Y
	for i, win := range c.w {
		win.maxlines = 0
		if i == len(c.w)-1 {
			r1.Max.Y = r.Max.Y
		} else {
fmt.Println(c.r)
			r1.Max.Y = r1.Min.Y + (win.r.Dx() + display.ScaleSize(Border))*r.Dy()/c.r.Dy()
		}
		r1.Max.Y = max(r1.Max.Y, r1.Min.Y + display.ScaleSize(Border) + tagfont.Height)
		r2 := r1
		r2.Max.Y = r2.Max.Y + display.ScaleSize(Border)
fmt.Println("E", r2)
		display.ScreenImage.Draw(r2, display.Black, nil, image.ZP)
		r1.Min.Y = r2.Max.Y
		r1.Min.Y = win.Resize(r1, false, i == len(c.w)-1)
	}
}

func cmp(a, b interface{}) int {
	return 0
}

func (c *Column) Sort() {

}

func (c *Column) Grow(w *Window, but int) {
	//var nl, ny *int
	var v *Window

	var windex int

	for windex := range c.w {
		if(c.w[windex] == w) {
			break
		}
	}
	if windex == c.nw() {
		panic("can't find window") // TODO(flux): implement counterpart to the C version's error()
	}

 	cr := c.r
	if but < 0 {	/* make sure window fills its own space properly */
		r := w.r
		if windex == int(c.nw()-1) || !c.safe {
			r.Max.Y = cr.Max.Y
		} else {
			r.Max.Y = c.w[windex+1].r.Min.Y - display.ScaleSize(Border)
		}
		w.Resize(r, false, true)
		return
	}
	cr.Min.Y = c.w[0].r.Min.Y
	if but == 3 {	/* Switch to full size window */
		if windex != 0 {
			v = c.w[0]
			c.w[0] = w
			c.w[windex] = v
		}
fmt.Println("F", cr)
		display.ScreenImage.Draw(cr, textcolors[frame.ColBack], nil, image.ZP);
		w.Resize(cr, false, true);
		for i:=1; i<c.nw(); i++ {
			c.w[i].body.fr.MaxLines = 0;
		}
		c.safe = false
		return
	}
	/* store old #lines for each window */
	onl := w.body.fr.MaxLines
	nl := make([]int, c.nw())
	ny := make([]int, c.nw())
	tot := 0
	for j:=0; j<c.nw(); j++{
		l := c.w[j].taglines-1 + c.w[j].body.fr.MaxLines; // TODO(flux): This taglines subtraction (for scrolling tags) assumes tags take the same number of pixels height as the body lines.  This is clearly false.
		nl[j] = l
		tot += l
	}
	/* approximate new #lines for this window */
	if but == 2 {	/* as big as can be */
		for i:=range nl { nl[i] = 0 }
		goto Pack
	}
	{	// Scope for nnl & dln
		nnl := min(onl + max(min(5, w.taglines-1+w.maxlines), onl/2), tot) // TODO(flux) more bad taglines use
		if nnl < w.taglines-1+w.maxlines {
			nnl = (w.taglines-1+w.maxlines + nnl)/2
		}
		if nnl == 0 {
			nnl = 2
		}
		dnl := nnl - onl;
		/* compute new #lines for each window */
		for k:=1; k<c.nw(); k++{
			/* prune from later window */
			j := windex+k
			if j<c.nw() && nl[j] != 0{
				l := min(dnl, max(1, nl[j]/2))
				nl[j] -= l
				nl[windex] += l
				dnl -= l
			}
			/* prune from earlier window */
			j = windex-k
			if j>=0 && nl[j] != 0{
				l := min(dnl, max(1, nl[j]/2))
				nl[j] -= l
				nl[windex] += l
				dnl -= l
			}
		}
	}
    Pack:
	/* pack everyone above */
	y1 := cr.Min.Y
	for j:=0; j<windex; j++{
		v = c.w[j]
		r := v.r
		r.Min.Y = y1
		r.Max.Y = y1+ v.tagtop.Dy()
		if nl[j] != 0 {
			r.Max.Y += 1 + nl[j]*v.body.fr.Font.DefaultHeight()
		}
		r.Min.Y = v.Resize(r, c.safe, false)
		r.Max.Y += display.ScaleSize(Border)
fmt.Println("G", r)
		display.ScreenImage.Draw(r, display.Black, nil, image.ZP)
		y1 = r.Max.Y
	}
	/* scan to see new size of everyone below */
	y2 := c.r.Max.Y
	for j:=c.nw()-1; j>windex; j-- {
		v = c.w[j]
		r := v.r
		r.Min.Y = y2-v.tagtop.Dy()
		if nl[j] != 0{
			r.Min.Y -= 1 + nl[j]*v.body.fr.Font.DefaultHeight()
		}
		r.Min.Y -= display.ScaleSize(Border)
		ny[j] = r.Min.Y
		y2 = r.Min.Y
	}
	/* compute new size of window */
	r := w.r
	r.Min.Y = y1
	r.Max.Y = y2
	h := w.body.fr.Font.DefaultHeight() // TODO(flux) Is this the right frame font height to use?
	if r.Dy() < w.tagtop.Dy()+1+h+display.ScaleSize(Border) {
		r.Max.Y = r.Min.Y + w.tagtop.Dy()+1+h+display.ScaleSize(Border)
	}
	/* draw window */
	r.Max.Y = w.Resize(r, c.safe, true)
	if windex < c.nw()-1 {
		r.Min.Y = r.Max.Y
		r.Max.Y += display.ScaleSize(Border)
fmt.Println("H", r)
		display.ScreenImage.Draw(r, display.Black, nil, image.ZP)
		for j:=windex+1; j<c.nw(); j++ {
			ny[j] -= (y2-r.Max.Y)
		}
	}
	/* pack everyone below */
	y1 = r.Max.Y
	for j:=windex+1; j<c.nw(); j++ {
		v = c.w[j]
		r = v.r
		r.Min.Y = y1
		r.Max.Y = y1+v.tagtop.Dy()
		if nl[j] != 0 {
			r.Max.Y += 1 + nl[j]*v.body.fr.Font.DefaultHeight()
		}
		y1 = v.Resize(r, c.safe, j==c.nw()-1)
		if j < c.nw()-1 {	/* no border on last window */
			r.Min.Y = y1
			r.Max.Y += display.ScaleSize(Border)
fmt.Println("I", r)
			display.ScreenImage.Draw(r, display.Black, nil, image.ZP)
			y1 = r.Max.Y
		}
	}
	c.safe = true
	// winmousebut(w); TODO(flux): Mouse?
}

func (c *Column) DragWin(w *Window, but int) {

}

func (c *Column) Which(p image.Point) *Text {
	return nil
}

func (c *Column) Clean() int {
	return 0
}
