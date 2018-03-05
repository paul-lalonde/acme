package main

import (
	"image"

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

func (c *Column)nw() uint {
	return uint(len(c.w))
}

func (c *Column)Init(r image.Rectangle) *Column {
	if c == nil {
		c = &Column{}
	}
	c.w = []*Window{}
	display.ScreenImage.Draw(r, display.White, nil, image.ZP)
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
	var winindex uint
	for winindex := range c.w {
		v = c.w[winindex]
		if y < v.r.Max.Y {
			break
		}
	}
	buggered := false 	// historical variable name
	if c.nw() > 0 {
		if winindex < c.nw() {
			winindex++
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
		if winindex < c.nw() {
			ymax = c.w[winindex].r.Min.Y-display.ScaleSize(Border)
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
		display.ScreenImage.Draw(r, textcolors[frame.ColBack], nil, image.ZP)
		r1 := r
		y = min(y, ymax-(v.tag.fr.Font.DefaultHeight()*v.taglines+v.body.fr.Font.DefaultHeight()+display.ScaleSize(Border)+1))
		r1.Max.Y = min(y, v.body.fr.Rect.Min.Y+v.body.fr.Nlines*v.body.fr.Font.DefaultHeight())
		r1.Min.Y = v.Resize(r1, false, false)
		r1.Max.Y = r1.Min.Y+display.ScaleSize(Border)
		display.ScreenImage.Draw(r1, display.Black, nil, image.ZP)
		
		/*
		 * leave r with w's coordinates
		 */
		r.Min.Y = r1.Max.Y
	}
	if w == nil {
		w = NewWindow()
		w.col = c
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
	c.tag.Resize(r1, true) 
	// TODO(flux): Column button
	//display.ScreenImage.Draw(c.tag.scrollr, colbutton, nil, colbutton->r.Min)
	r1.Min.Y = r1.Max.Y
	r1.Max.Y += r.Max.Y
	display.ScreenImage.Draw(r1, display.Black, nil, image.ZP)
	r1.Max.Y = r.Max.Y
	for i, win := range c.w {
		win.maxlines = 0
		if i == len(c.w)-1 {
			r1.Max.Y = r.Max.Y
		} else {
			r1.Max.Y = r1.Min.Y + (win.r.Dx() + display.ScaleSize(Border))*r.Dy()/c.r.Dy()
		}
		r1.Max.Y = max(r1.Max.Y, r1.Min.Y + display.ScaleSize(Border) + tagfont.Height)
		r2 := r1
		r2.Max.Y = r2.Max.Y + display.ScaleSize(Border)
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
//	var i, j, k, l, y1, y2, *nl, *ny, tot, nnl, onl, dnl, h int
//	Window *v;
//
//	var winindex uint
//
//	for winindex := range c.w {
//		if(&c.w[winindex] == w) {
//			break
//		}
//	}
//	if winindex = c.nw() {
//		coerror("can't find window");
//	}
//
// 	cr := c.r
//	if but < 0 {	/* make sure window fills its own space properly */
//		r: = w.r
//		if i == c.nw()-1 || !c.safe {
//			r.Max.Y = cr.Max.Y
//		} else {
//			r.Max.Y = c.w[winindex+1]->r.Min.Y - display.ScaleSize(Border)
//		}
//		w.Resize(r, false, true)
//		return
//	}
//panic("unimplemented") XXXXXXXXXXXXXXXXXXXXXXXXx
//	cr.min.y = c->w[0]->r.min.y;
//	if(but == 3){	/* full size */
//		if(i != 0){
//			v = c->w[0];
//			c->w[0] = w;
//			c->w[i] = v;
//		}
//		draw(screen, cr, textcols[BACK], nil, ZP);
//		winresize(w, cr, FALSE, TRUE);
//		for(i=1; i<c->nw; i++)
//			c->w[i]->body.fr.maxlines = 0;
//		c->safe = FALSE;
//		return;
//	}
//	/* store old #lines for each window */
//	onl = w->body.fr.maxlines;
//	nl = emalloc(c->nw * sizeof(int));
//	ny = emalloc(c->nw * sizeof(int));
//	tot = 0;
//	for(j=0; j<c->nw; j++){
//		l = c->w[j]->taglines-1 + c->w[j]->body.fr.maxlines;
//		nl[j] = l;
//		tot += l;
//	}
//	/* approximate new #lines for this window */
//	if(but == 2){	/* as big as can be */
//		memset(nl, 0, c->nw * sizeof(int));
//		goto Pack;
//	}
//	nnl = min(onl + max(min(5, w->taglines-1+w->maxlines), onl/2), tot);
//	if(nnl < w->taglines-1+w->maxlines)
//		nnl = (w->taglines-1+w->maxlines + nnl)/2;
//	if(nnl == 0)
//		nnl = 2;
//	dnl = nnl - onl;
//	/* compute new #lines for each window */
//	for(k=1; k<c->nw; k++){
//		/* prune from later window */
//		j = i+k;
//		if(j<c->nw && nl[j]){
//			l = min(dnl, max(1, nl[j]/2));
//			nl[j] -= l;
//			nl[i] += l;
//			dnl -= l;
//		}
//		/* prune from earlier window */
//		j = i-k;
//		if(j>=0 && nl[j]){
//			l = min(dnl, max(1, nl[j]/2));
//			nl[j] -= l;
//			nl[i] += l;
//			dnl -= l;
//		}
//	}
//    Pack:
//	/* pack everyone above */
//	y1 = cr.min.y;
//	for(j=0; j<i; j++){
//		v = c->w[j];
//		r = v->r;
//		r.min.y = y1;
//		r.max.y = y1+Dy(v->tagtop);
//		if(nl[j])
//			r.max.y += 1 + nl[j]*v->body.fr.font->height;
//		r.min.y = winresize(v, r, c->safe, FALSE);
//		r.max.y += Border;
//		draw(screen, r, display->black, nil, ZP);
//		y1 = r.max.y;
//	}
//	/* scan to see new size of everyone below */
//	y2 = c->r.max.y;
//	for(j=c->nw-1; j>i; j--){
//		v = c->w[j];
//		r = v->r;
//		r.min.y = y2-Dy(v->tagtop);
//		if(nl[j])
//			r.min.y -= 1 + nl[j]*v->body.fr.font->height;
//		r.min.y -= Border;
//		ny[j] = r.min.y;
//		y2 = r.min.y;
//	}
//	/* compute new size of window */
//	r = w->r;
//	r.min.y = y1;
//	r.max.y = y2;
//	h = w->body.fr.font->height;
//	if(Dy(r) < Dy(w->tagtop)+1+h+Border)
//		r.max.y = r.min.y + Dy(w->tagtop)+1+h+Border;
//	/* draw window */
//	r.max.y = winresize(w, r, c->safe, TRUE);
//	if(i < c->nw-1){
//		r.min.y = r.max.y;
//		r.max.y += Border;
//		draw(screen, r, display->black, nil, ZP);
//		for(j=i+1; j<c->nw; j++)
//			ny[j] -= (y2-r.max.y);
//	}
//	/* pack everyone below */
//	y1 = r.max.y;
//	for(j=i+1; j<c->nw; j++){
//		v = c->w[j];
//		r = v->r;
//		r.min.y = y1;
//		r.max.y = y1+Dy(v->tagtop);
//		if(nl[j])
//			r.max.y += 1 + nl[j]*v->body.fr.font->height;
//		y1 = winresize(v, r, c->safe, j==c->nw-1);
//		if(j < c->nw-1){	/* no border on last window */
//			r.min.y = y1;
//			r.max.y += Border;
//			draw(screen, r, display->black, nil, ZP);
//			y1 = r.max.y;
//		}
//	}
//	free(nl);
//	free(ny);
//	c->safe = TRUE;
//	winmousebut(w);
}

func (c *Column) DragWin(w *Window, but int) {

}

func (c *Column) Which(p image.Point) *Text {
	return nil
}

func (c *Column) Clean() int {
	return 0
}
