// Copyright 2017 The go-interpreter Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/go-interpreter/wagon/exec"
	"github.com/go-interpreter/wagon/validate"
	"github.com/go-interpreter/wagon/wasm"
	"math"
)

func main() {
	log.SetPrefix("wasm-run: ")
	log.SetFlags(0)

	verbose := flag.Bool("v", false, "enable/disable verbose mode")
	verify := flag.Bool("verify-module", false, "run module verification")

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	wasm.SetDebugMode(*verbose)

	run(os.Stdout, flag.Arg(0), *verify)
}

type Runtime struct {
	Input      []byte
	Output     []byte
	CallOutPut []byte
}

func (self *Runtime) prints_l(proc *exec.Process, ptr uint32, len uint32) {
	//fmt.Printf("prints_l called\n")
	bs := make([]byte, len)
	_, err := proc.ReadAt(bs, int64(ptr))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", bs)
}

func (self *Runtime) abort(proc *exec.Process) {
	fmt.Printf("abort called\n")
}

func (self *Runtime) ont_assert(proc *exec.Process, istrue uint32, msg uint32, len uint32) {
	fmt.Printf("ont_assert called\n")
	if istrue != 0 {
		bs := make([]byte, len)
		_, err := proc.ReadAt(bs, int64(msg))
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", bs)
	}
}

func NewHostModule(host *Runtime) *wasm.Module {
	m := wasm.NewModule()
	m.Types = &wasm.SectionTypes{
		Entries: []wasm.FunctionSig{
			{
				Form:       0,
				ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
			},
			{
				Form: 0,
			},
			{
				Form:       0,
				ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			},
		},
	}

	m.FunctionIndexSpace = []wasm.Function{
		{
			Sig:  &m.Types.Entries[0],
			Host: reflect.ValueOf(host.prints_l),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[1],
			Host: reflect.ValueOf(host.abort),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[2],
			Host: reflect.ValueOf(host.ont_assert),
			Body: &wasm.FunctionBody{},
		},
	}

	m.Export = &wasm.SectionExports{
		Entries: map[string]wasm.ExportEntry{
			"prints_l": {
				FieldStr: "prints_l",
				Kind:     wasm.ExternalFunction,
				Index:    0,
			},
			"abort": {
				FieldStr: "abort",
				Kind:     wasm.ExternalFunction,
				Index:    1,
			},
			"ont_assert": {
				FieldStr: "ont_assert",
				Kind:     wasm.ExternalFunction,
				Index:    2,
			},
		},
	}
	return m
}

func run(w io.Writer, fname string, verify bool) {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	host := &Runtime{}
	raw, err := ioutil.ReadFile(fname)
	m, err := wasm.ReadModule(bytes.NewReader(raw), func(name string) (*wasm.Module, error) {
		switch name {
		case "env":
			return NewHostModule(host), nil
		}
		return nil, fmt.Errorf("module %q unknown", name)
	})

	//m, err := wasm.ReadModule(f, importer)
	if err != nil {
		log.Fatalf("could not read module: %v", err)
	}

	if verify {
		err = validate.VerifyModule(m)
		if err != nil {
			log.Fatalf("could not verify module: %v", err)
		}
	}

	if m.Export == nil {
		log.Fatalf("module has no export section")
	}

	vm, err := exec.NewVM(m, math.MaxUint64)
	if err != nil {
		log.Fatalf("could not create VM: %v", err)
	}

	vm.AvaliableGas = &exec.Gas{GasPrice: 500, GasLimit: 1000000}

	entryname := "invoke"
	entry, ok := m.Export.Entries[entryname]
	if !ok {
		log.Fatalf("method: " + entryname + " do not exist")
	}
	index := int64(entry.Index)
	params := make([]uint64, 0)

	res, err := vm.ExecCode(index, params...)

	fmt.Printf("exec res : %d\n", res)
}

func importer(name string) (*wasm.Module, error) {
	f, err := os.Open(name + ".wasm")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m, err := wasm.ReadModule(f, nil)
	if err != nil {
		return nil, err
	}
	err = validate.VerifyModule(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
