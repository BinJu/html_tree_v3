package html_tree_v3

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

func BuildTree(r *HtmlReader) (*TreeNode, error) {
	root := NewTreeNode("ROOT", nil)
	node := root
	for {
		elm, err := r.Next()
		if err != nil {
			return nil, err
		}
		if elm == nil {
			return root, nil
		}
		if data := node.GetData(); data != nil && elm.Type != ELM_CLOSE && nodeAutoClose(data.(*Element).Value) {
			node = node.GetParent()
		}
		/*p := ""
		if node.GetData() != nil {
			p = node.GetData().(*Element).Value
		}
		fmt.Printf("@Creating Node: %s\t Type: %c\t Parent:%s\n", elm.Value, elm.Type, p)*/
		switch elm.Type {
		case ELM_OPEN:
			node = node.AddChild(elm.Value, elm)
		case ELM_CLOSE:
			if data := node.GetData(); data != nil && nodeAutoClose(data.(*Element).Value) &&
				!strings.EqualFold(data.(*Element).Value, elm.Value) {
				node = node.GetParent()
			}
			node = node.GetParent()
		case ELM_SELFCLOSE:
			node.AddChild(elm.Value, elm)
		case ELM_TEXT:
			node.AddChild(elm.Value, elm)
		case ELM_COMMENT:
			node.AddChild(elm.Value, elm)
		default:
			return nil, errors.New("Unexpected ELEMENT: " + string(elm.Type))
		}
	}
}

func nodeAutoClose(name string) bool {
	for _, ac := range []string{"!doctype", "link", "input", "meta", "br", "param", "img"} {
		if strings.EqualFold(ac, name) {
			return true
		}
	}
	return false
}

type HtmlVisitor struct {
	Buff bytes.Buffer
}

func (v *HtmlVisitor) Visit(node *TreeNode, bStart bool) {
	if node.GetData() == nil {
		return
	}
	data := node.GetData().(*Element)
	if bStart {
		switch data.Type {
		case ELM_OPEN:
			v.writeLeadingSpaces(node)
			v.Buff.WriteString(fmt.Sprintf("<%s", data.Value))
			for k, val := range data.Props {
				if len(val) > 0 {
					v.Buff.WriteString(fmt.Sprintf(" %s=\"%s\"", k, val))
				} else {
					v.Buff.WriteString(fmt.Sprintf(" %s", k))
				}
			}
			if node.ChildrenCount() > 0 {
				if d := node.GetChild(0).GetData(); d != nil && d.(*Element).Type == ELM_TEXT {
					v.Buff.WriteString(">")
				} else {
					v.Buff.WriteString(">\n")
				}
			} else {
				v.Buff.WriteString(">")
			}
		case ELM_SELFCLOSE:
			v.writeLeadingSpaces(node)
			v.Buff.WriteString(fmt.Sprintf("<%s", data.Value))
			for k, val := range data.Props {
				if len(val) > 0 {
					v.Buff.WriteString(fmt.Sprintf(" %s=\"%s\"", k, val))
				} else {
					v.Buff.WriteString(fmt.Sprintf(" %s", k))
				}
			}

			v.Buff.WriteString("/>\n")
		case ELM_TEXT:
			v.Buff.WriteString(data.Value)
		case ELM_COMMENT:
			v.writeLeadingSpaces(node)
			v.Buff.WriteString(fmt.Sprintf("<!--%s-->\n", data.Value))
		}
	} else {
		if data.Type == ELM_OPEN {
			if node.ChildrenCount() > 0 {
				if d := node.GetChild(0).GetData(); d != nil && d.(*Element).Type != ELM_TEXT {
					v.writeLeadingSpaces(node)
				}
			}
			v.Buff.WriteString("</" + data.Value + ">\n")
		}
	}
}

func (v *HtmlVisitor) writeLeadingSpaces(node *TreeNode) {
	for tmp := node.GetParent(); tmp.GetParent() != nil; tmp = tmp.GetParent() {
		v.Buff.WriteString("  ")
	}
}
