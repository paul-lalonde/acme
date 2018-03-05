package main

import (
	"image"
	"math"
	"sync"

	"github.com/rjkroege/acme/frame"
	"9fans.net/go/draw"
)

type Window struct {
	lk   *sync.Mutex
	ref  Ref
	tag  Text
	body Text
	r    image.Rectangle

	isdir     bool
	isscratch bool
	filemenu  bool
	dirty     bool
	autoident bool
	showdel   bool

	id    int
	addr  Range
	limit Range

	nopen     [QMAX]bool
	nomark    bool
	wselrange Range
	rdselfd   int

	col    *Column
	eventx Xfid
	events string

	nevents     int
	owner       int
	maxlines    int
	dlp         **Dirlist
	ndl         int
	putseq      int
	nincl       int
	incl        []string
	reffont     *draw.Font
	ctrllock    *sync.Mutex
	ctlfid      uint
	dumpstr     string
	dumpdir     string
	dumpid      int
	utflastqid  int
	utflastboff int
	utflastq    int
	tagsafe     bool
	tagexpand   bool
	taglines    int
	tagtop      image.Rectangle
	editoutlk   *sync.Mutex
}

func NewWindow() *Window { 
	return &Window{
	}
}

func (w *Window) Init(clone *Window, r image.Rectangle) {

	//	var r1, br image.Rectangle
	//	var f *File
	//	var rf *Reffont
	//	var rp []rune
	//	var nc int

	bodyFile := NewFile("")
	w.body.Init(bodyFile, r, tagfont, textcolors)
	tagFile := NewFile("")
	w.tag.Init(tagFile, r, tagfont, textcolors)

	w.tag.w = w
	w.taglines = 1
	w.tagsafe = true
	w.tagexpand = true
	w.body.w = w

	//	WinId++
	//	w.id = WinId

	w.ref.Inc()
	if globalincref {
		w.ref.Inc()
	}
	w.ctlfid = math.MaxUint64
	w.utflastqid = -1
	//	r1 = r

	w.tagtop = r
	w.tagtop.Max.Y = r.Min.Y + tagfont.Height
}

func (w *Window) DrawButton() {

}

func (w *Window) RunePos() int {
	return 0
}

func (w *Window) ToDel() {

}

func (w *Window) TagLines(r image.Rectangle) int {
	return 1
}

func (w *Window) Resize(r image.Rectangle, safe, keepextra bool) int {
	// mouseintag := mouse.xy.In(w.tag.all)
	// mouseinbody := mouse.xy.In(w.body.all) // TODO(flux): Mouse

	w.tagtop = r
	w.tagtop.Max.Y = r.Min.Y + tagfont.Height

	r1 := r
	r1.Max.Y = min(r.Max.Y, r1.Min.Y + w.taglines * tagfont.Height)
	if !safe || !w.tagsafe || w.tag.all.Eq(r1) {
		w.taglines = w.TagLines(r)
		r1.Max.Y = min(r.Max.Y, r1.Min.Y + w.taglines * tagfont.Height)
	}

	y := r1.Max.Y;

	// Resize/redraw tag TODO(flux)
	/*
	if(!safe || !w->tagsafe || !eqrect(w->tag.all, r1)){
		textresize(&w->tag, r1, TRUE);
		y = w->tag.fr.r.max.y;
		windrawbutton(w);
		w->tagsafe = TRUE;

		// If mouse is in tag, pull up as tag closes. 
		if(mouseintag && !ptinrect(mouse->xy, w->tag.all)){
			p = mouse->xy;
			p.y = w->tag.all.max.y-3;
			moveto(mousectl, p);
		}

		// If mouse is in body, push down as tag expands. 
		if(mouseinbody && ptinrect(mouse->xy, w->tag.all)){
			p = mouse->xy;
			p.y = w->tag.all.max.y+3;
			moveto(mousectl, p);
		}
	}
	*/
	// Redraw body
	r1 = r
	r1.Min.Y = y
	if !safe || !w.body.all.Eq(r1) {
		oy := y
		if y+1+w.body.fr.Font.DefaultHeight() < r.Max.Y { /* no room for one line */
			r1.Min.Y = y
			r1.Max.Y = y+1
			display.ScreenImage.Draw(r1, tagcolors[frame.ColBord], nil, image.ZP)
			y++
			r1.Min.Y = min(y, r.Max.Y)
			r1.Max.Y = r.Max.Y
		} else {
			r1.Min.Y = y
			r1.Max.Y = y
		}
		y = w.body.Resize(r1, keepextra)
		w.r = r
		w.r.Max.Y = y
		// w.body.Scrdraw()  // TODO(flux) scrollbars
		w.body.all.Min.Y = oy
	}
	w.maxlines = min(w.body.fr.Nlines, max(w.maxlines, w.body.fr.MaxLines))
	return w.r.Max.Y
}

func (w *Window) Lock1(owner int) {

}

func (w *Window) Lock(owner int) {

}

func (w *Window) Unlock() {

}

func (w *Window) MouseBut() {

}

func (w *Window) DirFree() {

}

func (w *Window) Close() {

}

func (w *Window) Delete() {

}

func (w *Window) Undo(isundo bool) {

}

func (w *Window) SetName(name string) {

}

func (w *Window) Type(t *Text, r rune) {

}

func (w *Window) ClearTag() {

}

func (w *Window) SetTag1() {

}

func (w *Window) SetTag() {

}

func (w *Window) Commit(t *Text) {

}

func (w *Window) AddIncl(r string, n int) {

}

func (w *Window) Clean(conservative bool) int {
	return 0
}

func (w *Window) CtlPrint(buf string, fonts int) string {
	return ""
}

func (w *Window) Event(fmt string, args ...interface{}) {

}
