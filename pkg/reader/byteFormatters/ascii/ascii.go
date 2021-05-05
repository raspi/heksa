package ascii

import (
	"fmt"
	"github.com/raspi/heksa/pkg/reader/byteFormatters/base"
)

// Check implementation
var _ base.ByteFormatter = AsciiPrinter{}

func PrintSpecial(specialBreak, hilightBreak string, b byte) string {
	return fmt.Sprintf(`%[1]s[%[2]s%[3]c%[1]s]`, specialBreak, hilightBreak, AsciiByteToChar[b])
}

type AsciiPrinter struct {
}

func New() AsciiPrinter {
	return AsciiPrinter{}
}

func (p AsciiPrinter) Print(b byte) (o string) {
	return string(AsciiByteToChar[b])
}

func (p AsciiPrinter) GetPrintSize() int {
	return 1
}

func (p AsciiPrinter) UseSplitter() bool {
	return true
}
