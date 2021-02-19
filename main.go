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
		(import "" "f2" (func $f2))
		(func (export "f1")
			(unreachable))
	`)
	check(err)

	for {

		store := wasmtime.NewStore(wasmtime.NewEngine())

		imports := []*wasmtime.Extern{
			wasmtime.NewFunc(store, wasmtime.NewFuncType(nil, nil), func(c *wasmtime.Caller, args []wasmtime.Val) ([]wasmtime.Val, *wasmtime.Trap) {
				panic("xxx")
				return nil, nil
			}).AsExtern(),
		}

		module, err := wasmtime.NewModule(store.Engine, bs)
		check(err)

		instance, err := wasmtime.NewInstance(store, module, imports)
		check(err)

		err = func() error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()

			_, err = instance.GetExport("f1").Func().Call()
			return err
		}()
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
