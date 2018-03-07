package main

import (
	"image"
	"math"
	"sync"

	"github.com/paul-lalonde/acme/frame"
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
	autoindent bool
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
	var rf *draw.Font
	//	var rp []rune
	//	var nc int

	bodyFile := NewFile("")
	w.body.Init(bodyFile, r, tagfont, textcolors)

	w.tag.w = w
	w.taglines = 1
	w.tagsafe = true
	w.tagexpand = true
	w.body.w = w

	WinId++
	w.id = WinId

	w.ctlfid = math.MaxUint64
	w.utflastqid = -1
	r1 := r

	w.tagtop = r
	w.tagtop.Max.Y = r.Min.Y + tagfont.Height
	r1.Max.Y = r1.Min.Y + w.taglines*tagfont.Height

	f := &File{}
	f.AddText(&w.tag)
	w.tag.Init(f, r1, tagfont, tagcolors);
	w.tag.what = Tag;

	/* tag is a copy of the contents, not a tracked image */
/* TODO(flux): Unimplemented Clone
	if(clone){
		w.tag.Delete(&, 0, w.tag.file.b.nc, true);
		w.tag.Insert(0, clone.tag.file.b, true);
		w.tag.file.Reset();
		w.tag.Setselect(len(w.tag.file.b), len(w.tag.file.b))
	}
*/
	r1 = r
	r1.Min.Y += w.taglines*tagfont.Height + 1
	if r1.Max.Y < r1.Min.Y {
		r1.Max.Y = r1.Min.Y
	}
/* TODO(flux): Unimplemented Clone
	f = nil;
	if clone {
		f = clone.body.file;
		w.body.org = clone.body.org;
		w.isscratch = clone.isscratch;
		rf = rfget(false, false, false, clone.body.reffont.f.name);
	} else {
*/
		rf = fontget(0, false, false, "");
//	}
	f = f.AddText(&w.body)
	w.body.what = Body
	w.body.Init(f, r1, rf, textcolors)
	r1.Min.Y -= 1
	r1.Max.Y = r1.Min.Y +1
	display.ScreenImage.Draw(r1, tagcolors[frame.ColBord], nil, image.ZP)
	// TODO(flux) w.body.Scrdraw()
	w.r = r
	var br image.Rectangle
	br.Min = w.tag.scrollr.Min
	br.Max.X = br.Min.X + button.R.Dx()
	br.Max.Y = br.Min.Y + button.R.Dy()
	display.ScreenImage.Draw(br, button, nil, button.R.Min)
	w.filemenu = true
	w.maxlines = w.body.fr.MaxLines
	w.autoindent = globalautoindent
/* TODO(flux): Unimplemented Clone
	if clone != nil{
		w.dirty = clone.dirty
		w.autoindent = clone.autoindent
		w.body.Setselect(clone.body.q0, clone.body.q1)
		w.Settag(w)
	}
*/
}

func (w *Window) DrawButton() {

	b := button;
	if !w.isdir && !w.isscratch && w.body.file.mod { // TODO(flux) validate text cache stuff goes away
		b = modbutton
	}
	var br image.Rectangle
	br.Min = w.tag.scrollr.Min;
	br.Max.X = br.Min.X + b.R.Dx()
	br.Max.Y = br.Min.Y + b.R.Dy()
	display.ScreenImage.Draw(br, b, nil, b.R.Min)

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
	if !safe || !w.tagsafe || !w.tag.all.Eq(r1) {
		w.tag.Resize(r1, true)
		y = w.tag.fr.Rect.Max.Y
		w.DrawButton()
		w.tagsafe = true;

		// If mouse is in tag, pull up as tag closes. 
/* TODO(flux): Mouse
		if(mouseintag && !ptinrect(mouse.xy, w.tag.all)){
			p = mouse.xy;
			p.y = w.tag.all.Max.Y-3;
			moveto(mousectl, p);
		}
		// If mouse is in body, push down as tag expands. 
		if(mouseinbody && ptinrect(mouse.xy, w.tag.all)){
			p = mouse.xy;
			p.y = w.tag.all.Max.Y+3;
			moveto(mousectl, p);
		}
*/
	}
	// Redraw body
	r1 = r
	r1.Min.Y = y
	if !safe || !w.body.all.Eq(r1) {
		oy := y
		if y+1+w.body.fr.Font.DefaultHeight() <= r.Max.Y { /* room for one line */
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
	w.maxlines = min(w.body.fr.NLines, max(w.maxlines, w.body.fr.MaxLines))
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
