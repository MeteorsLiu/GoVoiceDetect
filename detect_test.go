package govoicedetect

import "testing"

func TestDetect(t *testing.T) {
	v, err := NewVad("/home/nfs/py/GHKP-50/GHKP-50_02.mp4")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	t.Log(v.Detect())
}
