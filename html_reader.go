package html_tree_v3

import (
	"strings"
)

const ELM_OPEN = 'O'
const ELM_CLOSE = 'C'
const ELM_SELFCLOSE = 'B'
const ELM_COMMENT = 'c'
const ELM_TEXT = 'T'
const LAST_SCRIPT = 'S'

type Element struct {
	Value string
	Props map[string]string
	Type  byte
}

type HtmlReader struct {
	data  string
	pos   int
	lastT byte
}

func (r *HtmlReader) Next() (*Element, error) {
	if r.eof() {
		return nil, nil
	}
	r.skip(" \r\n\t")

	elm := Element{}
	if r.is("<!--") {
		elm.Type = ELM_COMMENT
		r.move(4)
		elm.Value = r.to("-->")
		r.move(3)
	} else if r.is("</") {
		elm.Type = ELM_CLOSE
		r.move(2)
		elm.Value = r.to(">")
		elm.Value = strings.Trim(elm.Value, " \r\n\t")
		r.move(1)
	} else if r.is("<") {
		if r.is_nocase("<script") {
			r.lastT = LAST_SCRIPT
		} else {
			r.lastT = 0
		}
		elm.Type = ELM_OPEN
		r.move(1)
		elm.Value = r.to(">", "/>", " ")
		elm.Value = strings.Trim(elm.Value, " \r\n\t")
		if r.is(">") {
			r.move(1)
		} else if r.is("/>") {
			elm.Type = ELM_SELFCLOSE
			r.move(2)
		} else if r.is(" ") {
			lastKey := ""
			for {
				r.skip(" \r\n\t")
				if len(lastKey) > 0 {
					if elm.Props == nil {
						elm.Props = make(map[string]string)
					}
					elm.Props[lastKey] = ""
				}
				lastKey = strings.Trim(r.to("=", ">", "/>", " "), " \r\n\t")
				if r.is(" ") {
					r.move(1)
				} else if r.is("=") {
					r.move(1)
					r.skip(" \r\n\t")
					val := r.to(" ", "'", "\"", ">", "/>")
					if r.is(" ") {
						if len(lastKey) > 0 {
							if elm.Props == nil {
								elm.Props = make(map[string]string)
							}
							elm.Props[lastKey] = val
						}
						r.move(1)
						lastKey = ""
					} else if r.is("'") {
						r.move(1)
						val := r.to("'")
						if len(lastKey) > 0 {
							if elm.Props == nil {
								elm.Props = make(map[string]string)
							}
							elm.Props[lastKey] = val
						}
						r.move(1)
						lastKey = ""
					} else if r.is("\"") {
						r.move(1)
						val := r.to("\"")
						if len(lastKey) > 0 {
							if elm.Props == nil {
								elm.Props = make(map[string]string)
							}
							elm.Props[lastKey] = val
						}
						r.move(1)
						lastKey = ""
					} else if r.is(">") {
						if len(lastKey) > 0 {
							if elm.Props == nil {
								elm.Props = make(map[string]string)
							}
							elm.Props[lastKey] = ""
						}
						r.move(1)
						lastKey = ""
						break
					} else if r.is("/>") {
						elm.Type = ELM_SELFCLOSE
						if len(lastKey) > 0 {
							if elm.Props == nil {
								elm.Props = make(map[string]string)
							}
							elm.Props[lastKey] = ""
						}
						r.move(2)
						lastKey = ""
						break
					}
				} else if r.is(">") {
					if len(lastKey) > 0 {
						if elm.Props == nil {
							elm.Props = make(map[string]string)
						}
						elm.Props[lastKey] = ""
					}
					r.move(1)
					lastKey = ""
					break
				} else if r.is("/>") {
					elm.Type = ELM_SELFCLOSE
					if len(lastKey) > 0 {
						if elm.Props == nil {
							elm.Props = make(map[string]string)
						}
						elm.Props[lastKey] = ""
					}
					r.move(2)
					lastKey = ""
					break
				}
			}
		}
	} else { //text
		elm.Type = ELM_TEXT
		if r.lastT == LAST_SCRIPT {
			elm.Value = r.to("</script")
		} else {
			elm.Value = r.to("<")
		}
		elm.Value = strings.Trim(elm.Value, " \r\n\t")
	}
	return &elm, nil
}

func (r *HtmlReader) to(exps ...string) string {
	p := r.pos
	for ; r.pos < len(r.data); r.pos++ {
		for _, ex := range exps {
			if r.is_nocase(ex) {
				return r.data[p:r.pos]
			}
		}
	}
	return ""
}

func (r *HtmlReader) is(exp string) bool {
	i := 0
	for ; r.pos+i < len(r.data) && i < len(exp) && r.data[r.pos+i] == exp[i]; i++ {
	}
	if i == len(exp) {
		return true
	} else {
		return false
	}
}

func (r *HtmlReader) is_nocase(exp string) bool {
	if r.pos+len(exp) < len(r.data) {
		return strings.EqualFold(exp, r.data[r.pos:r.pos+len(exp)])
	}
	return false
}

func (r *HtmlReader) skip(set string) {
	for ; r.pos < len(r.data); r.pos++ {
		if strings.IndexAny(set, string(r.data[r.pos])) < 0 {
			return
		}
	}
}

func (r *HtmlReader) move(count int) {
	r.pos += count
}

func (r *HtmlReader) eof() bool {
	return r.pos >= len(r.data)
}
