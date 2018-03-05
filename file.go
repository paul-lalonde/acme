package main

type File struct {
	b         Buffer
	delta     Buffer
	epsilon   Buffer
	elogbuf   *Buffer
	elog      Elog
	name      []rune
	qidpath   uint64
	mtime     uint64
	dev       int
	unread    bool
	editclean bool
	seq       int
	mod       bool

	curtext *Text
	text    []*Text
	dumpid  int
}

func (f *File) AddText(t *Text) *File {
	f.text = append(f.text, t)
	f.curtext = t
	return f
}

func (f *File) DelText(t *Text) {

}

func (f *File) Insert(q0 uint, s []rune, ns uint) {

}

func (f *File) Uninsert(delta *Buffer, q0, ns uint) {

}

func (f *File) Delete(p0, p1 uint) {

}

func (f *File) Undelete(delta *Buffer, p0, p1 uint) {

}

func (f *File) SetName(name string, n int) {

}

func (f *File) UnsetName(delta *Buffer) {

}

func NewFile(filename string) *File {
	return &File{
	b:        NewBuffer(),
/*	delta     Buffer
	epsilon   Buffer
*/
	elog: MakeElog(),
	name:      []rune(filename),
//	qidpath   uint64
//	mtime     uint64
//	dev       int
	unread: true,
	editclean: true,
//	seq       int
	mod:      false,

	curtext: nil,
	text: []*Text{},
//	ntext   int
//	dumpid  int
	}
}

func NewTagFile() *File {

	return &File{
	b:        NewBuffer(),
/*	delta     Buffer
	epsilon   Buffer
*/
	elog: MakeElog(),
	name:      []rune{},
//	qidpath   uint64
//	mtime     uint64
//	dev       int
	unread: true,
	editclean: true,
//	seq       int
	mod:      false,

//	curtext *Text
//	text    **Text
//	ntext   int
//	dumpid  int
	}
}

func (f *File) RedoSeq() uint {
	return 0
}

func (f *File) Undo(isundo bool, q0p, q1p *uint) {

}

func (f *File) Reset() {

}

func (f *File) Close() {

}

func (f *File) Mark() {

}
