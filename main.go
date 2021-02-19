package main

import (
	"fmt"

	"github.com/bytecodealliance/wasmtime-go"

	_ "github.com/bytecodealliance/wasmtime-go/build/include"        // to include the C headers.
	_ "github.com/bytecodealliance/wasmtime-go/build/linux-x86_64"   // to include the static lib for linking.
	_ "github.com/bytecodealliance/wasmtime-go/build/macos-x86_64"   // to include the static lib for linking.
	_ "github.com/bytecodealliance/wasmtime-go/build/windows-x86_64" // to include the static lib for linking.
)

func main() {

	bs, err := wasmtime.Wat2Wasm(`
	(module
		(func (export "f1")
			(unreachable)))
	`)
	check(err)

	for i := 0; ; i++ {
		fmt.Println("i:", i)
		store := wasmtime.NewStore(wasmtime.NewEngine())

		module, err := wasmtime.NewModule(store.Engine, bs)
		check(err)

		instance, err := wasmtime.NewInstance(store, module, nil)
		check(err)

		_, err = instance.GetExport("f1").Func().Call()
		if err == nil {
			panic("expected error")
		}
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
