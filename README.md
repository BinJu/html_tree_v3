## What?
This is tool which can parse html source code to a tree. 
Later you can:
###1. Work on this tree with your logic.
###2. Output the tree as **XML qulified** HTML code.(That means The output is suitable to any other XML tools).
## How?
```
package html_tree_v3

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestTreeBuild(t *testing.T) {
	file, err := os.Open("log.html")
	if nil != err {
		t.Error("Failed to open test file")
	}
	defer file.Close()

	data, errRead := ioutil.ReadAll(file)
	if nil != errRead {
		t.Error("Read data from file failed")
	}
	r := HtmlReader{data: string(data), pos: 0, lastT: 0}
	tree, buildErr := BuildTree(&r)
	if nil != buildErr {
		t.Error(buildErr.Error())
	}
	if nil != tree {
		v := &HtmlVisitor{}
		tree.Visit(v)
		t.Log(v.Buff.String())
	}
}
```