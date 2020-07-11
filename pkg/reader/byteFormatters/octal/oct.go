package octal

import "github.com/raspi/heksa/pkg/reader/byteFormatters/base"

type printer struct {
}

func New() base.ByteFormatter {
	return printer{}
}

func (p printer) Print(b byte) (o string) {
	return octByteToString[b]
}

func (p printer) GetPrintSize() int {
	return 3
}
