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
	"math/rand"
	"os"
	"reflect"
	"runtime/pprof"
	"time"

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

	//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file` ")
	//var memprofile = flag.String("memprofile", "", "write memory profile `file` to ")

	cpuprofile := "./cpu.prof"
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)

		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)

		}
		defer pprof.StopCPUProfile()

	}

	wasm.SetDebugMode(*verbose)

	run(os.Stdout, flag.Arg(0), *verify)
}

type Runtime struct {
	Input      []byte
	Output     []byte
	CallOutPut []byte
}

func (self *Runtime) block_height(proc *exec.Process) uint32 {
	//fmt.Printf("outputlength: %d\n", uint32(len(self.Output)))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return uint32(r.Uint32())
}

func (self *Runtime) call_output_length(proc *exec.Process) uint32 {
	//fmt.Printf("outputlength: %d\n", uint32(len(self.Output)))
	return uint32(len(self.Output))
}

func (self *Runtime) ret(proc *exec.Process, ptr uint32, len uint32) {
	//self := proc.HostData().(*Runtime)
	bs := make([]byte, len)
	_, err := proc.ReadAt(bs, int64(ptr))
	if err != nil {
		panic(err)
	}

	self.Output = bs
	//fmt.Printf("ret bytes %x\n", self.Output)
	//proc.Terminate()
}

func (self *Runtime) get_call_output(proc *exec.Process, dst uint32) {
	//self := proc.HostData().(*Runtime)
	_, err := proc.WriteAt(self.Output, int64(dst))
	//fmt.Printf("output bytes %x\n", self.Output)
	if err != nil {
		panic(err)
	}
}

func (self *Runtime) debug(proc *exec.Process, ptr uint32, len uint32) {
	bs := make([]byte, len)
	_, err := proc.ReadAt(bs, int64(ptr))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", bs)
}

func (self *Runtime) abort(proc *exec.Process) {
	fmt.Printf("abort called\n")
	panic(nil)
}

func (self *Runtime) timestamp(proc *exec.Process) uint64 {
	//res := rand.Intn()
	res := 2
	return uint64(res)
}

func (self *Runtime) save_input_arg(proc *exec.Process, ptr uint32, len uint32) {
	bs := make([]byte, len)
	_, err := proc.ReadAt(bs, int64(ptr))
	if err != nil {
		panic(err)
	}

	self.Input = make([]byte, len)
	copy(self.Input, bs)

	//panic(nil)
}

func (self *Runtime) get_input(proc *exec.Process, dst uint32) {
	_, err := proc.WriteAt(self.Input, int64(dst))
	if err != nil {
		panic(err)
	}
}

func (self *Runtime) input_length(proc *exec.Process) uint32 {
	return uint32(len(self.Input))
}

func (self *Runtime) ont_assert(proc *exec.Process, istrue uint32, msg uint32, len uint32) {
	//fmt.Printf("ont_assert called\n")
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
			{
				Form:       0,
				ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			{
				Form:        0,
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			{
				Form:        0,
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
			},
		},
	}

	m.FunctionIndexSpace = []wasm.Function{
		{
			Sig:  &m.Types.Entries[0],
			Host: reflect.ValueOf(host.debug),
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
		{
			Sig:  &m.Types.Entries[0],
			Host: reflect.ValueOf(host.save_input_arg),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[3],
			Host: reflect.ValueOf(host.get_input),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[4],
			Host: reflect.ValueOf(host.input_length),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[5],
			Host: reflect.ValueOf(host.timestamp),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[0],
			Host: reflect.ValueOf(host.ret),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[3],
			Host: reflect.ValueOf(host.get_call_output),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[4],
			Host: reflect.ValueOf(host.call_output_length),
			Body: &wasm.FunctionBody{},
		},
		{
			Sig:  &m.Types.Entries[4],
			Host: reflect.ValueOf(host.block_height),
			Body: &wasm.FunctionBody{},
		},
	}

	m.Export = &wasm.SectionExports{
		Entries: map[string]wasm.ExportEntry{
			"debug": {
				FieldStr: "debug",
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
			"save_input_arg": {
				FieldStr: "save_input_arg",
				Kind:     wasm.ExternalFunction,
				Index:    3,
			},
			"get_input": {
				FieldStr: "get_input",
				Kind:     wasm.ExternalFunction,
				Index:    4,
			},
			"input_length": {
				FieldStr: "input_length",
				Kind:     wasm.ExternalFunction,
				Index:    5,
			},
			"timestamp": {
				FieldStr: "timestamp",
				Kind:     wasm.ExternalFunction,
				Index:    6,
			},
			"ret": {
				FieldStr: "ret",
				Kind:     wasm.ExternalFunction,
				Index:    7,
			},
			"get_call_output": {
				FieldStr: "get_call_output",
				Kind:     wasm.ExternalFunction,
				Index:    8,
			},
			"call_output_length": {
				FieldStr: "call_output_length",
				Kind:     wasm.ExternalFunction,
				Index:    9,
			},
			"block_height": {
				FieldStr: "block_height",
				Kind:     wasm.ExternalFunction,
				Index:    10,
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

	vm.AvaliableGas = &exec.Gas{GasPrice: 500, GasLimit: 1000000000000000000}

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
