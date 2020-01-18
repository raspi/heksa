package units

import "testing"

func Test100(t *testing.T) {
	actual := `100`
	expected := int64(100)

	got, err := Parse(actual)

	if err != nil {
		t.Fail()
	}

	if got != expected {
		t.Fail()
	}
}

func TestMinus100(t *testing.T) {
	actual := `-100`
	expected := int64(-100)

	got, err := Parse(actual)

	if err != nil {
		t.Fail()
	}

	if got != expected {
		t.Fail()
	}
}

func Test2KB(t *testing.T) {
	actual := `2KB`
	expected := int64(2000)

	got, err := Parse(actual)

	if err != nil {
		t.Fail()
	}

	if got != expected {
		t.Fail()
	}
}

func TestMinus2KB(t *testing.T) {
	actual := `-2KB`
	expected := int64(-2000)

	got, err := Parse(actual)

	if err != nil {
		t.Fail()
	}

	if got != expected {
		t.Fail()
	}
}

func TestHex2KB(t *testing.T) {
	actual := `0xAKB`
	expected := int64(10000)

	got, err := Parse(actual)

	if err != nil {
		t.Fail()
	}

	if got != expected {
		t.Fail()
	}
}

func Test2KiB(t *testing.T) {
	actual := `2KiB`
	expected := int64(2048)

	got, err := Parse(actual)

	if err != nil {
		t.Fail()
	}

	if got != expected {
		t.Fail()
	}
}

func TestMinus2KiB(t *testing.T) {
	actual := `-2KiB`
	expected := int64(-2048)

	got, err := Parse(actual)

	if err != nil {
		t.Fail()
	}

	if got != expected {
		t.Fail()
	}
}
