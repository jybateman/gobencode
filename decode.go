package bencode

import (
	"bytes"
	"strconv"
)

type decode struct {
	r *bytes.Buffer
}

// TODO error check
func (d *decode) decodeString() string {
	ns, _ := d.r.ReadBytes(':')
	i, _ := strconv.Atoi(string(ns[:len(ns)-1]))
	key := make([]byte, i)
	d.r.Read(key)	
	return string(key)
}

func (d *decode) decodeInt() int {
	n, _ := d.r.ReadBytes('e')
	i, _ := strconv.Atoi(string(n[:len(n)-1]))
	return i
}

func (d *decode) decodeList() []interface{} {
	var l []interface{}
	for {
		t, _ := d.r.ReadByte()
		e := d.getType(t)
		l = append(l, e)
		t, _ = d.r.ReadByte()
		if t == 'e' {
			break
		}
		d.r.UnreadByte()
	}
	return l
}

func (d *decode) decodeDict() map[string]interface{} {
	dict := make(map[string]interface{})
	for {
		key := d.decodeString()
		t, _ := d.r.ReadByte()
		dict[string(key)] = d.getType(t)
		t, _ = d.r.ReadByte()
		if t == 'e' {
			break
		}
		d.r.UnreadByte()
	}
	return dict
}

func (d *decode) getType(t byte) interface{} {
	switch {
	case t == 'd':
		return d.decodeDict()
	case t == 'l':
		return d.decodeList()
	case t == 'i':
		return d.decodeInt()
	default:
		d.r.UnreadByte()
		return d.decodeString()
	}
	return nil
}

func Decode(dta string) map[string]interface{}{
	var d decode

	d.r = bytes.NewBufferString(dta)
	dct, _ := d.r.ReadByte()
	if dct == 'd' {
		return d.decodeDict()
	}
	return nil
}
