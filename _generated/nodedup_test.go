package _generated

import (
	"bytes"
	"fmt"
	"github.com/glycerine/greenpack2/msgp"
	"testing"
)

//  --no-dedup Test
func Test3444NoDedupOfSamePointer(t *testing.T) {

	ptr := &Target2{ID: 3}
	greet := &Greeter2{Style: 7}
	nd := &NoDedup{
		MyPtr0:   ptr,
		MyPtr1:   ptr,
		MyIface0: greet,
		MyIface1: greet,
		Slice:    []Hello2{greet, greet},
	}

	var buf bytes.Buffer
	wr := msgp.NewWriter(&buf)
	panicOn(nd.EncodeMsg(wr))
	wr.Flush()

	//fmt.Printf("\nAFTER EncodeMsg WRITE: PointerCount=%v. buf='%#v'\n buf as string='%s'\n", wr.DedupPointerCount(), buf.Bytes(), string(buf.Bytes()))

	rd := msgp.NewReader(&buf)
	var nd2 NoDedup
	panicOn(nd2.DecodeMsg(rd))
	if nd2.MyPtr0 == nd2.MyPtr1 {
		panic(fmt.Sprintf("expected pointers to be different"))
	}
	if nd2.MyIface0.(*Greeter2) == nd2.MyIface1.(*Greeter2) {
		panic(fmt.Sprintf("expected pointers behind interfaces to be different"))
	}
	// check slices of interfaces for dedup too.
	if nd2.Slice[0].(*Greeter2) == nd2.Slice[1].(*Greeter2) {
		panic(fmt.Sprintf("expected pointers behind interfaces to be different"))
	}
	if nd2.Slice[0].(*Greeter2) == nd2.MyIface0.(*Greeter2) {
		panic(fmt.Sprintf("expected pointers behind interfaces to be different"))
	}
}

// -no-dedup Test2
func Test3445NoDedupOfSamePointer(t *testing.T) {

	// dedup within interface slice alone?

	//ptr := &Target2{ID: 3}
	iface := &Greeter2{Style: 7}
	nd := &NoDedup{
		//MyPtr0: ptr,
		//MyPtr1: ptr,
		//		MyIface0: iface,
		//		MyIface1: iface,
		Slice: []Hello2{iface, iface},
	}

	var buf bytes.Buffer
	wr := msgp.NewWriter(&buf)
	//wr.ResetDedup()
	panicOn(nd.EncodeMsg(wr))
	wr.Flush()

	//fmt.Printf("\nAFTER EncodeMsg WRITE: PointerCount=%v. buf='%#v'\n buf as string='%s'\n", wr.PointerCount(), buf.Bytes(), string(buf.Bytes()))

	rd := msgp.NewReader(&buf)
	//rd.ResetDedup()
	var nd2 NoDedup
	panicOn(nd2.DecodeMsg(rd))
	// check slices of interfaces for dedup too.
	if nd2.Slice[0].(*Greeter2) == nd2.Slice[1].(*Greeter2) {
		panic(fmt.Sprintf("expected pointers behind interfaces to be different"))
	}
}

// -no-dedup Test3
func Test3446DedupOfSamePointerWorks(t *testing.T) {

	// dedup within slices of pointers

	//ptr := &Target2{ID: 3}
	iface := &Greeter2{Style: 7}
	nd := &NoDedup{
		//MyPtr0: ptr,
		//MyPtr1: ptr,
		//		MyIface0: iface,
		//		MyIface1: iface,
		Slice:    []Hello2{iface},
		SlicePtr: []*Greeter2{iface},
	}

	var buf bytes.Buffer
	wr := msgp.NewWriter(&buf)
	//wr.ResetDedup()
	panicOn(nd.EncodeMsg(wr))
	wr.Flush()

	//fmt.Printf("\nAFTER EncodeMsg WRITE: PointerCount=%v. buf='%#v'\n buf as string='%s'\n", wr.PointerCount(), buf.Bytes(), string(buf.Bytes()))

	rd := msgp.NewReader(&buf)
	//rd.ResetDedup()
	var nd2 NoDedup
	panicOn(nd2.DecodeMsg(rd))
	// check across slices of interfaces and slices of pointers for dedup.
	if nd2.Slice[0].(*Greeter2) == nd2.SlicePtr[0] {
		panic(fmt.Sprintf("expected pointers/interfaces to be different"))
	}
}

// -no-dedup Test4
func Test3500NestedDedup(*testing.T) {

	// slices of interfaces within slices of interfaces
	// should still dedup correctly.

	inner := &Inner2{}
	mid := &Middle2{
		Children: []Shouter2{inner, inner},
	}
	outer := &Outer2{
		Slc: []Imid2{mid},
	}

	var buf bytes.Buffer
	wr := msgp.NewWriter(&buf)
	panicOn(outer.EncodeMsg(wr))
	wr.Flush()

	//fmt.Printf("\nAFTER EncodeMsg WRITE: PointerCount=%v. buf='%#v'\n buf as string='%s'\n", wr.DedupPointerCount(), buf.Bytes(), string(buf.Bytes()))

	rd := msgp.NewReader(&buf)
	var o2 Outer2
	panicOn(o2.DecodeMsg(rd))
	// check dedup of the inner and inner
	if o2.Slc[0].(*Middle2).Children[0].(*Inner2) == o2.Slc[0].(*Middle2).Children[1].(*Inner2) {
		panic(fmt.Sprintf("expected pointers to be different"))
	}
}
