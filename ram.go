package main

type Ram struct {
    all []byte
}

func (ram *Ram) Init() {
    ram.all = make([]byte, 65536, 65536)
}

func (ram *Ram) Read(loc uint16) (byte) {
    return ram.all[loc]
}

func (ram *Ram) Write(loc uint16, val byte) {
    ram.all[loc] = val
}

func (ram *Ram) WriteBlock(loc uint16, vals []byte) {
    copy(ram.all[loc:], vals)
}
