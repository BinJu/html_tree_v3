package html_tree_v3

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	file, err := os.Open("log.html")
	if nil != err {
		t.Error("Failed to open test file")
	}
	defer file.Close()

	data, errRead := ioutil.ReadAll(file)
	if nil != errRead {
		t.Error("Read data from file failed")
	}
	//r := HtmlReader{data: "<!--this is a comment --><html style='good' aaa bbb>this is content<body><div  class='inner'/></body></html>", pos: 0, lastT: 0, lastV: ""}
	r := HtmlReader{data: string(data), pos: 0, lastT: 0}
	for {
		elm, _ := r.Next()
		if nil == elm {
			break
		}
		/*t.Logf("#V: %s T:%c", elm.Value, elm.Type)
		for k, v := range elm.Props {
			t.Logf("\t\tK: %s, V: %s", k, v)
		}*/
	}
}
