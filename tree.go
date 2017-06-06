package html_tree_v3

type Visitor interface {
	Visit(node *TreeNode, bStart bool)
}
type TreeNode struct {
	name     string
	data     interface{}
	children []*TreeNode
	parent   *TreeNode
}

func (t *TreeNode) GetName() string {
	return t.name
}

func (t *TreeNode) AddChild(name string, data interface{}) *TreeNode {
	n := TreeNode{name: name}
	n.data = data
	n.parent = t
	t.children = append(t.children, &n)
	return &n
}

func (t *TreeNode) GetChild(idx int) *TreeNode {
	return t.children[idx]
}

func (t *TreeNode) GetParent() *TreeNode {
	return t.parent
}

func (t *TreeNode) GetChildren() []*TreeNode {
	return t.children
}

func (t *TreeNode) GetData() interface{} {
	return t.data
}

func (t *TreeNode) ChildrenCount() int {
	return len(t.children)
}

func (t *TreeNode) DeleteChild(idx int) *TreeNode {
	node := t.children[idx]
	t.children = append(t.children[:idx], t.children[idx+1:]...)
	node.parent = nil
	return node
}
func (t *TreeNode) Visit(v Visitor) {
	v.Visit(t, true)
	for _, child := range t.children {
		child.Visit(v)
	}
	v.Visit(t, false)
}

func NewTreeNode(name string, data interface{}) *TreeNode {
	n := TreeNode{name: name, data: data}
	return &n
}
