// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package validate provides functions for validating WebAssembly modules.
package validate

import (
	"bytes"
	"errors"
	"io"

	"github.com/go-interpreter/wagon/wasm"
	ops "github.com/go-interpreter/wagon/wasm/operators"
)

func verifyBodyWithSpec(fn *wasm.FunctionSig, body *wasm.FunctionBody, module *wasm.Module) (*mockSpecVM, error) {
	vm := &mockSpecVM{
		opdStack:   []wasm.ValueType{},
		ctrlStack:  []ctrlFrame{},
		code:       bytes.NewReader(body.Code),
		origLength: uint32(len(body.Code)),
		curFunc:    fn,
	}

	localVariables := []wasm.ValueType{}
	for _, entry := range fn.ParamTypes {
		localVariables = append(localVariables, entry)
	}

	for _, entry := range body.Locals {
		vars := make([]wasm.ValueType, entry.Count)
		for i := uint32(0); i < entry.Count; i++ {
			vars[i] = entry.Type
		}
		localVariables = append(localVariables, vars...)
	}

	var fnsig wasm.ValueType
	if len(fn.ReturnTypes) == 0 {
		fnsig = wasm.ValueType(wasm.BlockTypeEmpty)
	} else if len(fn.ReturnTypes) == 1 {
		fnsig = fn.ReturnTypes[0]
	} else {
		return vm, errors.New("MVP only support single return value")
	}

	err := vm.pushCtrl(wasm.ValueType(fnsig), wasm.ValueType(fnsig), false)
	if err != nil {
		return vm, err
	}

	for {
		op, err := vm.code.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return vm, err
		}
		opStruct, err := ops.New(op)
		if err != nil {
			return vm, err
		}

		if !opStruct.Polymorphic {
			if err := vm.adjustStack(opStruct); err != nil {
				return vm, err
			}
		}

		switch op {
		case ops.Drop:
			_, err := vm.popOpd()
			if err != nil {
				return vm, err
			}
		case ops.Unreachable:
			err := vm.unreachable()
			if err != nil {
				return vm, err
			}
		case ops.Block:
			sig, err := vm.fetchByte()
			if err != nil {
				return vm, err
			}
			err = vm.pushCtrl(wasm.ValueType(sig), wasm.ValueType(sig), false)
			if err != nil {
				return vm, err
			}
		case ops.If:
			sig, err := vm.fetchByte()
			if err != nil {
				return vm, err
			}
			/*If is not PolymorphicOp. handle operand stack already.*/
			err = vm.pushCtrl(wasm.ValueType(sig), wasm.ValueType(sig), true)
			if err != nil {
				return vm, err
			}
		case ops.Loop:
			sig, err := vm.fetchByte()
			if err != nil {
				return vm, err
			}
			err = vm.pushCtrl(wasm.ValueType(wasm.BlockTypeEmpty), wasm.ValueType(sig), false)
			if err != nil {
				return vm, err
			}
		case ops.Else:
			_, err = vm.matchElse()
			if err != nil {
				return vm, err
			}

			typ, err := vm.popCtrl()
			if err != nil {
				return vm, err
			}
			err = vm.pushCtrl(typ, typ, false)
			if err != nil {
				return vm, err
			}
		case ops.End:
			cFrame, err := vm.topCtrl()
			if err != nil {
				return vm, err
			}

			if cFrame.ifType && cFrame.endType != wasm.ValueType(wasm.BlockTypeEmpty) {
				return vm, errors.New("type mismatch in if false branch")
			}

			typ, err := vm.popCtrl()
			if err != nil {
				return vm, err
			}
			err = vm.pushOpd(typ)
			if err != nil {
				return vm, err
			}
		case ops.Br:
			depth, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
			/*due to the function block frame.*/
			frame, err := vm.pickCtrl(depth)
			if err != nil {
				return vm, err
			}
			_, err = vm.popOpdExpect(frame.labelTypes)
			if err != nil {
				return vm, err
			}
			err = vm.unreachable()
			if err != nil {
				return vm, err
			}
		case ops.BrIf:
			depth, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
			/*sames like if. is not PolymorphicOp. handle operand stack before.*/
			frame, err := vm.pickCtrl(depth)
			if err != nil {
				return vm, err
			}
			_, err = vm.popOpdExpect(frame.labelTypes)
			if err != nil {
				return vm, err
			}
			err = vm.pushOpd(frame.labelTypes)
			if err != nil {
				return vm, err
			}
		case ops.BrTable:
			targetCount, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}

			var targetTable []uint32
			for i := uint32(0); i < targetCount; i++ {
				entry, err := vm.fetchVarUint()
				if err != nil {
					return vm, err
				}
				targetTable = append(targetTable, entry)
			}
			defaultTarget, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}

			/*pickCtrl will check the max depth*/
			defaultFrame, err := vm.pickCtrl(defaultTarget)
			if err != nil {
				return vm, err
			}

			for _, target := range targetTable {
				targetFrame, err := vm.pickCtrl(target)
				if err != nil {
					return vm, err
				}
				if targetFrame.labelTypes != defaultFrame.labelTypes {
					return vm, errors.New("matched error labelTypes")
				}
			}
			_, err = vm.popOpdExpect(wasm.ValueTypeI32)
			if err != nil {
				return vm, err
			}
			_, err = vm.popOpdExpect(defaultFrame.labelTypes)
			if err != nil {
				return vm, err
			}
			vm.unreachable()
			if err != nil {
				return vm, err
			}
		case ops.Return:
			if vm.ctrlSize() < 1 {
				return vm, errors.New("no function block can return")
			}
			funcFrame := vm.ctrlStack[0]

			_, err = vm.popOpdExpect(funcFrame.labelTypes)
			if err != nil {
				return vm, err
			}

			vm.unreachable()
			if err != nil {
				return vm, err
			}
		case ops.I32Const:
			_, err := vm.fetchVarInt()
			if err != nil {
				return vm, err
			}
		case ops.I64Const:
			_, err := vm.fetchVarInt64()
			if err != nil {
				return vm, err
			}
		case ops.F32Const:
			_, err := vm.fetchUint32()
			if err != nil {
				return vm, err
			}
		case ops.F64Const:
			_, err := vm.fetchUint64()
			if err != nil {
				return vm, err
			}
		case ops.GetLocal, ops.SetLocal, ops.TeeLocal:
			i, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
			if i >= uint32(len(localVariables)) {
				return vm, InvalidLocalIndexError(i)
			}
			typ := localVariables[i]
			if op == ops.GetLocal {
				err = vm.pushOpd(typ)
				if err != nil {
					return vm, err
				}
			} else if op == ops.SetLocal {
				_, err := vm.popOpdExpect(typ)
				if err != nil {
					return vm, err
				}
			} else {
				_, err := vm.popOpdExpect(typ)
				if err != nil {
					return vm, err
				}
				err = vm.pushOpd(typ)
				if err != nil {
					return vm, err
				}
			}
		case ops.GetGlobal, ops.SetGlobal:
			index, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
			gv := module.GetGlobal(int(index))
			if gv == nil {
				return vm, wasm.InvalidGlobalIndexError(index)
			}
			if op == ops.GetGlobal {
				err = vm.pushOpd(gv.Type.Type)
				if err != nil {
					return vm, err
				}
			} else {
				expectType := gv.Type.Type
				if gv.Type.Mutable != true {
					return vm, errors.New("try to set immutable global var.")
				}
				_, err := vm.popOpdExpect(expectType)
				if err != nil {
					return vm, err
				}
			}
		case ops.I32Load, ops.I64Load, ops.F32Load, ops.F64Load, ops.I32Load8s, ops.I32Load8u, ops.I32Load16s, ops.I32Load16u, ops.I64Load8s, ops.I64Load8u, ops.I64Load16s, ops.I64Load16u, ops.I64Load32s, ops.I64Load32u, ops.I32Store, ops.I64Store, ops.F32Store, ops.F64Store, ops.I32Store8, ops.I32Store16, ops.I64Store8, ops.I64Store16, ops.I64Store32:
			_, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
			// offset
			_, err = vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
		case ops.CurrentMemory, ops.GrowMemory:
			memIndex, err := vm.fetchByte()
			if err != nil {
				return vm, err
			}

			if memIndex != 0x00 {
				return vm, errors.New("validate: memory index must be 0")
			}
		case ops.Call:
			index, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}
			fn := module.GetFunction(int(index))
			if fn == nil {
				return vm, wasm.InvalidFunctionIndexError(index)
			}

			for index := range fn.Sig.ParamTypes {
				argType := fn.Sig.ParamTypes[len(fn.Sig.ParamTypes)-index-1]
				_, err = vm.popOpdExpect(argType)
				if err != nil {
					return vm, err
				}
			}

			if len(fn.Sig.ReturnTypes) > 1 {
				return vm, errors.New("validate: MVP not support multiple return types")
			}

			if len(fn.Sig.ReturnTypes) > 0 {
				err = vm.pushOpd(fn.Sig.ReturnTypes[0])
				if err != nil {
					return vm, err
				}
			}
		case ops.CallIndirect:
			if module.Table == nil || len(module.Table.Entries) == 0 {
				return vm, NoSectionError(wasm.SectionIDTable)
			}
			index, err := vm.fetchVarUint()
			if err != nil {
				return vm, err
			}

			tableIndex, err := vm.fetchByte()
			if err != nil {
				return vm, err
			}

			if tableIndex != 0x00 {
				return vm, errors.New("validate: table index in call_indirect must be 0")
			}

			if index > uint32(len(module.Types.Entries)) {
				return vm, errors.New("validate: type index out of bound")
			}

			fnExpectSig := module.Types.Entries[index]

			_, err = vm.popOpdExpect(wasm.ValueTypeI32)
			if err != nil {
				return vm, err
			}
			for index := range fnExpectSig.ParamTypes {
				argType := fnExpectSig.ParamTypes[len(fnExpectSig.ParamTypes)-index-1]
				_, err = vm.popOpdExpect(argType)
				if err != nil {
					return vm, err
				}
			}

			if len(fnExpectSig.ReturnTypes) > 1 {
				return vm, errors.New("validate: MVP not support multiple return types")
			}

			if len(fnExpectSig.ReturnTypes) > 0 {
				err = vm.pushOpd(fnExpectSig.ReturnTypes[0])
				if err != nil {
					return vm, err
				}
			}
		case ops.Select:
			_, err := vm.popOpdExpect(wasm.ValueTypeI32)
			if err != nil {
				return vm, err
			}
			typ1, err := vm.popOpd()
			if err != nil {
				return vm, err
			}
			typ2, err := vm.popOpdExpect(typ1)
			if err != nil {
				return vm, err
			}
			err = vm.pushOpd(typ2)
			if err != nil {
				return vm, err
			}
		}
	}

	if vm.ctrlSize() != 0 {
		return vm, errors.New("function frame mismatched")
	}

	if body.Code[len(body.Code)-1] != ops.End {
		return vm, wasm.ErrFunctionNoEnd
	}

	return vm, nil
}
