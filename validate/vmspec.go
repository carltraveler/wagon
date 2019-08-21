// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package validate

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/go-interpreter/wagon/wasm"
	"github.com/go-interpreter/wagon/wasm/leb128"
	ops "github.com/go-interpreter/wagon/wasm/operators"
)

const (
	ValueTypeUnk wasm.ValueType = 0x7b
)

type mockSpecVM struct {
	opdStack   []wasm.ValueType
	ctrlStack  []ctrlFrame
	origLength uint32
	code       *bytes.Reader
	curFunc    *wasm.FunctionSig
}

type ctrlFrame struct {
	labelTypes  wasm.ValueType
	endType     wasm.ValueType
	height      uint32
	unreachable bool
	ifType      bool
}

func (vm *mockSpecVM) matchElse() (bool, error) {
	cFrame, err := vm.topCtrl()
	if err != nil {
		return false, err
	}

	if cFrame.ifType {
		return true, nil
	}

	return false, errors.New("frame do not match else.")
}

func (vm *mockSpecVM) fetchVarUint() (uint32, error) {
	return leb128.ReadVarUint32(vm.code)
}

func (vm *mockSpecVM) fetchVarInt() (int32, error) {
	return leb128.ReadVarint32(vm.code)
}

func (vm *mockSpecVM) fetchByte() (byte, error) {
	return vm.code.ReadByte()
}

func (vm *mockSpecVM) fetchBlockType() (wasm.ValueType, error) {
	typ, err := vm.code.ReadByte()
	if err != nil {
		return 0, err
	}

	switch wasm.ValueType(typ) {
	case wasm.ValueTypeI32, wasm.ValueTypeI64, wasm.ValueTypeF32, wasm.ValueTypeF64, wasm.ValueType(wasm.BlockTypeEmpty):
	default:
		return 0, errors.New("unknow value type")
	}

	return wasm.ValueType(typ), nil
}

func (vm *mockSpecVM) fetchVarInt64() (int64, error) {
	return leb128.ReadVarint64(vm.code)
}

func (vm *mockSpecVM) fetchUint32() (uint32, error) {
	var buf [4]byte
	_, err := io.ReadFull(vm.code, buf[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

func (vm *mockSpecVM) fetchUint64() (uint64, error) {
	var buf [8]byte
	_, err := io.ReadFull(vm.code, buf[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf[:]), nil
}

func (vm *mockSpecVM) opdSize() uint32 {
	return uint32(len(vm.opdStack))
}

func (vm *mockSpecVM) ctrlSize() uint32 {
	return uint32(len(vm.ctrlStack))
}

func (vm *mockSpecVM) pushOpd(typ wasm.ValueType) {
	if typ != wasm.ValueType(wasm.BlockTypeEmpty) {
		vm.opdStack = append(vm.opdStack, typ)
	}
}

func (vm *mockSpecVM) popOpd() (wasm.ValueType, error) {
	topFrame, err := vm.topCtrl()
	if err != nil {
		return ValueTypeUnk, err
	}

	if topFrame.height == vm.opdSize() && topFrame.unreachable {
		return ValueTypeUnk, nil
	} else if topFrame.height == vm.opdSize() {
		return ValueTypeUnk, errors.New("operand stack underflow")
	}

	if len(vm.opdStack) == 0 {
		panic("code logic error")
	}

	typ := vm.opdStack[len(vm.opdStack)-1]

	vm.opdStack = vm.opdStack[:len(vm.opdStack)-1]
	return typ, nil
}

func (vm *mockSpecVM) topCtrl() (*ctrlFrame, error) {
	if len(vm.ctrlStack) == 0 {
		return nil, errors.New("no block frame in stack")
	}

	return &vm.ctrlStack[len(vm.ctrlStack)-1], nil
}

func (vm *mockSpecVM) pickCtrl(index uint32) (*ctrlFrame, error) {
	if len(vm.ctrlStack) == 0 {
		return nil, errors.New("no block frame in stack")
	}

	/*always use + other then sub to compare on uint32.*/
	if uint32(len(vm.ctrlStack)) < index+1 {
		return nil, errors.New("block frame execeed max")
	}

	return &vm.ctrlStack[uint32(len(vm.ctrlStack))-index-1], nil
}

func (vm *mockSpecVM) popOpdExpect(expect wasm.ValueType) (wasm.ValueType, error) {
	if expect == wasm.ValueType(wasm.BlockTypeEmpty) {
		return wasm.ValueType(wasm.BlockTypeEmpty), nil
	}

	actual, err := vm.popOpd()
	if err != nil {
		return ValueTypeUnk, err
	}

	if actual == ValueTypeUnk {
		return expect, nil
	}

	if expect == ValueTypeUnk {
		return actual, nil
	}

	if actual != expect {
		return ValueTypeUnk, errors.New("not expect type")
	}

	return actual, nil
}

func (vm *mockSpecVM) pushCtrl(labelTypes wasm.ValueType, endType wasm.ValueType, ifType bool) {
	frame := ctrlFrame{
		labelTypes:  labelTypes,
		endType:     endType,
		height:      vm.opdSize(),
		unreachable: false,
		ifType:      ifType,
	}

	vm.ctrlStack = append(vm.ctrlStack, frame)
}

func (vm *mockSpecVM) popCtrl() (wasm.ValueType, error) {
	topFrame, err := vm.topCtrl()
	if err != nil {
		return ValueTypeUnk, err
	}

	_, err = vm.popOpdExpect(topFrame.endType)
	if err != nil {
		return ValueTypeUnk, err
	}

	if vm.opdSize() != topFrame.height {
		return ValueTypeUnk, errors.New("stack overflow")
	}

	vm.ctrlStack = vm.ctrlStack[:len(vm.ctrlStack)-1]
	return topFrame.endType, nil
}

func (vm *mockSpecVM) unreachable() error {
	topFrame, err := vm.topCtrl()
	if err != nil {
		return err
	}

	if topFrame.height > vm.opdSize() {
		return errors.New("code logic error")
	}

	vm.opdStack = vm.opdStack[:topFrame.height]

	topFrame, err = vm.topCtrl()
	topFrame.unreachable = true
	return nil
}

func (vm *mockSpecVM) adjustStack(op ops.Op) error {
	for _, t := range op.Args {
		_, err := vm.popOpdExpect(t)
		if err != nil {
			return err
		}
	}

	if op.Returns != wasm.ValueType(wasm.BlockTypeEmpty) {
		vm.pushOpd(op.Returns)
	}
	return nil
}

func (vm *mockSpecVM) pc() int {
	return int(vm.origLength - uint32(vm.code.Len()))
}
