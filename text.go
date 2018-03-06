package main

import (
	"image"
	"math"
	"os"

	"github.com/rjkroege/acme/frame"
	"9fans.net/go/draw"
)

const (
	Ldot   = "."
	TABDIR = 3
)

var (
	left1  = []rune{'{', '[', '(', '<', 0xab}
	right1 = []rune{'}', ']', ')', '>', 0xbb}
	left2  = []rune{'\n'}
	left3  = []rune{'\'', '"', '`'}

	left = [][]rune{
		left1,
		left2,
		left3,
	}

	right = [][]rune{
		right1,
		left2,
		left3,
	}
)

type TextKind byte
const (
	Columntag = iota
	Rowtag
	Tag
	Body
)

type Text struct {
	file *File
	fr 	*frame.Frame
	font *draw.Font
	org     uint
	q0      uint
	q1      uint
	what    TextKind
	tabstop int
	w       *Window
	scrollr image.Rectangle
	lastsr  image.Rectangle
	all     image.Rectangle
	row     *Row
	col     *Column

	iq1         uint
	eq0         uint
	cq0         uint
	ncache      int
	ncachealloc int
	cache       []rune
	nofill      int
	needundo    int
}

func (t *Text)Init(f *File, r image.Rectangle, rf *draw.Font, cols [frame.NumColours]*draw.Image) *Text {
	if t == nil {
		t = new(Text)
	}
	t.file = f
	t.all = r
	t.scrollr = r
	t.scrollr.Max.X = r.Min.X + display.ScaleSize(Scrollwid)
	t.lastsr = nullrect
	r.Min.X += display.ScaleSize(Scrollwid) + display.ScaleSize(Scrollgap)
	t.eq0 = math.MaxUint64
	t.ncache = 0
	t.font = rf
	t.tabstop = int(maxtab)
	t.fr = frame.NewFrame( r, rf, display.ScreenImage, cols)
	t.Redraw(r, rf, display.ScreenImage, -1)
	return t
}

func (t *Text) Redraw(r image.Rectangle, f *draw.Font, b *draw.Image, odx int) {
	t.fr.Init(r, f, b, t.fr.Cols)
	rr := t.fr.Rect
	rr.Min.X -= display.ScaleSize(Scrollwid) + display.ScaleSize(Scrollgap)
//	if !t.fr.noredraw {
		display.ScreenImage.Draw(rr, t.fr.Cols[frame.ColBack], nil, image.ZP)
//	}
	// TODO(flux): Draw the text!
}

func (t *Text) Resize(r image.Rectangle, keepextra bool) int {
	if r.Dy() <= 0 {
		r.Max.Y = r.Min.Y
	} else {
		if(!keepextra) {
			r.Max.Y -= r.Dy()%t.fr.Font.DefaultHeight();
		}
	}
	odx := t.all.Dx()
	t.all = r;
	t.scrollr = r;
	t.scrollr.Max.X = r.Min.X+Scrollwid
	t.lastsr = image.ZR
	r.Min.X += display.ScaleSize(Scrollwid+Scrollgap)
	t.fr.Clear(false)
	t.Redraw(r, t.fr.Font.Impl(), t.fr.Background, odx)
	if keepextra && t.fr.Rect.Max.Y < t.all.Max.Y /* && !t.fr.noredraw */ {
		/* draw background in bottom fringe of window */
		r.Min.X -= display.ScaleSize(Scrollgap)
		r.Min.Y = t.fr.Rect.Max.Y
		r.Max.Y = t.all.Max.Y
		display.ScreenImage.Draw(r, t.fr.Cols[frame.ColBack], nil, image.ZP)
	}
	return t.all.Max.Y
}

func (t *Text) Close() {

}

func (t *Text) Columnate(dlp **Dirlist, ndl int) {

}

// TODO(flux): turn return int into error check.
func (t *Text) Load(q0 uint, filename string, setquid bool) int {
	if t.ncache!=0 || t.file.b.nc() > 0 || t.w==nil || t!=&t.w.body {
		panic("text.load")
	}
	if t.w.isdir && t.file.nname==0{
		warning(nil, "empty directory name")
		return -1
	}
	if ismtpt(file){
		warning(nil, "will not open self mount point %s\n", file)
		return -1
	}
	fd, err := os.Open(file);
	if err != nil{
		warning(nil, "can't open %s: %v\n", file, err)
		return -1
	}
	d, err = fd.Stat()
	if err != nil{
		warning(nil, "can't fstat %s: %r\n", file, err)
		goto Rescue
	}
	nulls = false;
	h = nil;
	if d.IsDir() {
		/* this is checked in get() but it's possible the file changed underfoot */
		if len(t.file.text) > 1{
			warning(nil, "%s is a directory; can't read with multiple windows on it\n", file)
			goto Rescue
		}
		t.w.isdir = true;
		t.w.filemenu = false;
		// TODO(flux): Find all '/' and replace with filepath.Separator properly
		if len(t.file.name) > 0 && !strings.HasSuffix(t.file.name, "/") {
			t.file.name = t.file.name + "/"
			t.w.SetName(t.file.name);
		}
		dirNames, err := fd.Readdirnames(0)
		if err != nil {
			warning(nil, "failed to Readdirnames: %s\n", file)
			goto Rescue
		}
		for i, dn := range dirNames {
			f. err := os.Open(dn)
			if err != nil{
				warning(nil, "can't open %s: %v\n", dn, err)
				return -1
			}
			s, err  := f.Stat()
			if err != nil{
				warning(nil, "can't fstat %s: %r\n", dn, err)
			} else {
				if s.IsDir() {
					dirNames[i] = dn+"/"
				}
			}
		}
		sort.Strings(dirNames)
		widths := make([]int, len(dirNames))
		for i, s := range dirNames {
			widths[i] = t.fr.Font.StringWidth(s)
		}
		t.Columnate(dirNames, widths)
		q1 = t.file.b.nc()
	}else{ //////////HERE
		t.w.isdir = false
		t.w.filemenu = true
		if q0 == 0 {
			h = sha1(nil, 0, nil, nil)
		}
		q1 = q0 + fileload(t.file, q0, fd, &nulls, h);
	}
	if setqid{
		if h != nil {
			sha1(nil, 0, t.file.sha1, h);
			h = nil;
		} else {
			memset(t.file.sha1, 0, sizeof t.file.sha1);
		}
		t.file.dev = d.dev;
		t.file.mtime = d.mtime;
		t.file.qidpath = d.qid.path;
	}
	close(fd);
	rp = fbufalloc();
	for(q=q0; q<q1; q+=n){
		n = q1-q;
		if n > RBUFSIZE
			n = RBUFSIZE;
		bufread(&t.file.b, q, rp, n);
		if q < t.org
			t.org += n;
		else if q <= t.org+t.fr.nchars
			frinsert(&t.fr, rp, rp+n, q-t.org);
		if t.fr.lastlinefull
			break;
	}
	fbuffree(rp);
	for(i=0; i<t.file.ntext; i++){
		u = t.file.text[i];
		if u != t{
			if u.org > u.file.b.nc)	/* will be 0 because of reset(, but safety first */
				u.org = 0;
			textresize(u, u.all, true);
			textbacknl(u, u.org, 0);	/* go to beginning of line */
		}
		textsetselect(u, q0, q0);
	}
	if nulls
		warning(nil, "%s: NUL bytes elided\n", file);
	free(d);
	return q1-q0;

    Rescue:
	close(fd);
	return -1;
}

func (t *Text) Backnl(p, n uint) uint {
	return 0
}

func (t *Text) BsInsert(q0 uint, r []rune, n uint, tofile bool, nrp *int) uint {
	return 0
}

func (t *Text) Insert(q0 uint, r []rune, tofile bool) {

}

func (t *Text) TypeCommit() {

}

func (t *Text) Fill() {

}

func (t *Text) Delete(q0, q1 uint, tofile bool) {

}

func (t *Text) Constrain(q0, q1 uint, p0, p1 *uint) {
	*p0 = minu(q0, t.file.b.nc())
	*p1 = minu(q1, t.file.b.nc())
}

func (t *Text) ReadRune(q uint) rune {
	return ' '
}

func (t *Text) BsWidth(c rune) int {
	return 0
}

func (t *Text) FileWidth(q0 uint, oneelement int) int {
	return 0
}

func (t *Text) Complete() []rune {
	return nil
}

func (t *Text) Type(r rune) {

}

func (t *Text) Commit(tofile bool) {

}

func (t *Text) FrameScroll(dl int) {

}

func (t *Text) Select() {

}

func (t *Text) Show(q0, q1 uint, doselect bool) {

}

func (t *Text) SetSelect(q0, q1 uint) {

}

func (t *Text) Select23(q0, q1 *uint, high *draw.Image, mask int) int {
	return 0
}

func (t *Text) Select3(q0, q1 *uint) int {
	return 0
}

func (t *Text) DoubleClick(q0, q1 *uint) {

}

func (t *Text) ClickMatch(cl, cr, dir int, q *uint) int {
	return 0
}

func (t *Text) ishtmlstart(q uint, q1 *uint) bool {
	return false
}

func (t *Text) ishtmlend(q uint, q0 *uint) bool {
	return false
}

func (t *Text) ClickHTMLMatch(q0, q1 *uint) int {
	return 0
}

func (t *Text) BackNL(p, n uint) uint {
	return 0
}

func (t *Text) SetOrigin(org uint, exact int) {

}

func (t *Text) Reset() {

}
