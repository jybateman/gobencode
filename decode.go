package bencode

import (
	"bytes"
	"errors"
	"strconv"
)

type decode struct {
	r *bytes.Buffer
}

// TODO error check
func (d *decode) decodeString() (string, error) {
	ns, err := d.r.ReadBytes(':')
	if err != nil {
		return "", err
	}
	i, err := strconv.Atoi(string(ns[:len(ns)-1]))
	if err != nil {
		return "", err
	}
	key := make([]byte, i)
	d.r.Read(key)	
	return string(key), nil
}

func (d *decode) decodeInt() (int, error) {
	n, err := d.r.ReadBytes('e')
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(string(n[:len(n)-1]))
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (d *decode) decodeList() ([]interface{}, error) {
	var l []interface{}
	for {
		t, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		e, err := d.getType(t)
		if err != nil {
			return nil, err
		}
		l = append(l, e)
		t, err = d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		if t == 'e' {
			break
		}
		d.r.UnreadByte()
	}
	return l, nil
}

func (d *decode) decodeDict() (map[string]interface{}, error) {
	dict := make(map[string]interface{})
	for {
		key, err := d.decodeString()
		if err != nil {
			return nil, err
		}
		t, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		v, err := d.getType(t)
		if err != nil {
			return nil, err
		}		
		dict[string(key)] = v
		t, err = d.r.ReadByte()
		if err != nil {
			return nil, err
		}
		if t == 'e' {
			break
		}
		d.r.UnreadByte()
	}
	return dict, nil
}

func (d *decode) getType(t byte) (interface{}, error) {
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
	return nil, errors.New("Type does not exist!")
}

func Decode(dta string) (map[string]interface{}, error) {
	var d decode

	d.r = bytes.NewBufferString(dta)
	dct, _ := d.r.ReadByte()
	if dct == 'd' {
		return d.decodeDict()
	}
	return nil, nil
}
