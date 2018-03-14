package main

/*

enum
{
	None = 0,
	Fore = '+',
	Back = '-'
};

enum
{
	Char,
	Line
};

func isaddrc (r  int) (int) {
	if r && utfrune("0123456789+-/$.#,;?", r)!=nil 
		return true;
	return false;
}


 //* quite hard: could be almost anything but white space, but we are a little conservative,
 //* aiming for regular expressions of alphanumerics and no white space

func isregexc (r  int) (int) {
	if r == 0 
		return false;
	if isalnum(r) 
		return true;
	if utfrune("^+-.*?#,;[]()$", r)!=nil 
		return true;
	return false;
}

// nlcounttopos starts at q0 and advances nl lines,
// being careful not to walk past the end of the text,
// and then nr chars, being careful not to walk past
// the end of the current line.
// It returns the final position.
func nlcounttopos (t * Text, q0  long, nl  long, nr  long) (long) {
	while(nl > 0 && q0 < t.file.b.nc) {
		if textreadc(t, q0++) == '\n' 
			nl--;
	}
	if nl > 0 
		return q0;
	while(nr > 0 && q0 < t.file.b.nc && textreadc(t, q0) != '\n') {
		q0++;
		nr--;
	}
	return q0;
}

func number (showerr  uint, t * Text, r  Range, line  int, dir  int, size  int, evalp * int) (Range) {
	uint q0, q1;

	if size == Char {
		if dir == Fore 
			line = r.q1+line;
		else if dir == Back {
			if r.q0==0 && line>0 
				r.q0 = t.file.b.nc;
			line = r.q0 - line;
		}
		if line<0 || line>t.file.b.nc 
			goto Rescue;
		*evalp = true;
		return range(line, line);
	}
	q0 = r.q0;
	q1 = r.q1;
	switch(dir){
	case None:
		q0 = 0;
		q1 = 0;
	Forward:
		while(line>0 && q1<t.file.b.nc)
			if textreadc(t, q1++) == '\n' || q1==t.file.b.nc 
				if --line > 0 
					q0 = q1;
		if line==1 && q1==t.file.b.nc  // 6 goes to end of 5-line file
			break;
		if line > 0 
			goto Rescue;
		break;
	case Fore:
		if q1 > 0 
			while(q1<t.file.b.nc && textreadc(t, q1-1) != '\n')
				q1++;
		q0 = q1;
		goto Forward;
	case Back:
		if q0 < t.file.b.nc 
			while(q0>0 && textreadc(t, q0-1)!='\n')
				q0--;
		q1 = q0;
		while(line>0 && q0>0){
			if textreadc(t, q0-1) == '\n' {
				if --line >= 0 
					q1 = q0;
			}
			--q0;
		}
		// :1-1 is :0 = #0, but :1-2 is an error 
		if line > 1 
			goto Rescue;
		while(q0>0 && textreadc(t, q0-1)!='\n')
			--q0;
	}
	*evalp = true;
	return range(q0, q1);

    Rescue:
	if showerr 
		warning(nil, "address out of range\n");
	*evalp = false;
	return r;
}


func regexp (showerr  uint, t * Text, lim  Range, r  Range, pat * Rune, dir  int, foundp * int) (Range) {
	int found;
	Rangeset sel;
	int q;

	if pat[0] == '\0' && rxnull() {
		if showerr 
			warning(nil, "no previous regular expression\n");
		*foundp = false;
		return r;
	}
	if pat[0] && rxcompile(pat) == false){
		*foundp = false;
		return r;
	}
	if dir == Back 
		found = rxbexecute(t, r.q0, &sel);
	else{
		if lim.q0 < 0 
			q = Infinity;
		else
			q = lim.q1;
		found = rxexecute(t, nil, r.q1, q, &sel);
	}
	if !found && showerr 
		warning(nil, "no match for regexp\n");
	*foundp = found;
	return sel.r[0];
}
*/
func address (showerr  bool, t * Text, lim  Range, ar  Range, a []rune, q0  uint, q1  uint) (r Range, evalp bool , qp uint) {
Unimpl()
	return Range{0,0}, false, 0
}
/*
	int dir, size, npat;
	int prevc, c, nc, n;
	uint q;
	Rune *pat;
	Range r, nr;

	func getc(r[]rune, q uint) int {
		return r[q]
	}

	r = ar;
	q = q0;
	dir = None;
	size = Line;
	c = 0;
	while(q < q1){
		prevc = c;
		c = (*getc)(a, q++);
		switch(c){
		default:
			*qp = q-1;
			return r;
		case ';':
			ar = r;
			// fall through 
		case ',':
			if prevc == 0 	// lhs defaults to 0 
				r.q0 = 0;
			if q>=q1 && t!=nil && t.file!=nil 	// rhs defaults to $ 
				r.q1 = t.file.b.nc;
			else{
				nr = address(showerr, t, lim, ar, a, q, q1, getc, evalp, &q);
				r.q1 = nr.q1;
			}
			*qp = q;
			return r;
		case '+':
		case '-':
			if *evalp && (prevc=='+' || prevc=='-') 
				if (nc=(*getc)(a, q))!='#' && nc!='/' && nc!='?' 
					r = number(showerr, t, r, 1, prevc, Line, evalp);	// do previous one 
			dir = c;
			break;
		case '.':
		case '$':
			if q != q0+1 {
				*qp = q-1;
				return r;
			}
			if *evalp 
				if c == '.' 
					r = ar;
				else
					r = range(t.file.b.nc, t.file.b.nc);
			if q < q1 
				dir = Fore;
			else
				dir = None;
			break;
		case '#':
			if q==q1 || (c=(*getc)(a, q++))<'0' || '9'<c {
				*qp = q-1;
				return r;
			}
			size = Char;
			// fall through 
		case '0': case '1': case '2': case '3': case '4':
		case '5': case '6': case '7': case '8': case '9':
			n = c -'0';
			while(q<q1){
				c = (*getc)(a, q++);
				if c<'0' || '9'<c {
					q--;
					break;
				}
				n = n*10+(c-'0');
			}
			if *evalp 
				r = number(showerr, t, r, n, dir, size, evalp);
			dir = None;
			size = Line;
			break;
		case '?':
			dir = Back;
			// fall through 
		case '/':
			npat = 0;
			pat = nil;
			while(q<q1){
				c = (*getc)(a, q++);
				switch(c){
				case '\n':
					--q;
					goto out;
				case '\\':
					pat = runerealloc(pat, npat+1);
					pat[npat++] = c;
					if q == q1 
						goto out;
					c = (*getc)(a, q++);
					break;
				case '/':
					goto out;
				}
				pat = runerealloc(pat, npat+1);
				pat[npat++] = c;
			}
		    out:
			pat = runerealloc(pat, npat+1);
			pat[npat] = 0;
			if *evalp 
				r = regexp(showerr, t, lim, r, pat, dir, evalp);
			free(pat);
			dir = None;
			size = Line;
			break;
		}
	}
	if *evalp && dir != None 
		r = number(showerr, t, r, 1, dir, Line, evalp);	// do previous one 
	*qp = q;
	return r;
}
*/