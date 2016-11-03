package bencode

import (
	"fmt"
	"bytes"
	"strconv"
)

type encode struct {
	b bytes.Buffer
}

// TODO error check
func (e *encode) encodeString(s string) {
	e.b.WriteString(strconv.Itoa(len(s))+":"+s)
}

func (e *encode) encodeInt(i int64) {
	e.b.WriteString("i")
	e.b.WriteString(strconv.FormatInt(i, 10))
	e.b.WriteString("e")
}

func (e *encode) encodeList(l []interface{}) {
	e.b.WriteString("l")
	e.getData(l)
	e.b.WriteString("e")
}

func (e *encode) encodeDict(d map[string]interface{}) {
	e.b.WriteString("d")
	for key, dta := range d {
		e.b.WriteString(strconv.Itoa(len(key))+":"+key)
		e.getData(dta)
	}
	e.b.WriteString("e")
}

func (e *encode) getData(dict interface{}) {
	switch dta := dict.(type) {
	case map[string]interface{}:
		e.encodeDict(dta)
	case []interface{}:
		e.encodeList(dta)
	case int64:
		e.encodeInt(dta)
	case string:
		e.encodeString(dta)
	}
}

func Encode(dta interface{}) {
	var e encode

	e.getData(dta)
	fmt.Println(e.b.String())
}
