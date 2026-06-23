//go:build js

package wasmbridge

import (
	"encoding/json"
	"errors"
	"syscall/js"
)

func resultOK(data []byte) js.Value {
	obj := js.Global().Get("Object").New()
	obj.Set("ok", true)
	if len(data) > 0 {
		uint8arr := js.Global().Get("Uint8Array").New(len(data))
		js.CopyBytesToJS(uint8arr, data)
		obj.Set("data", uint8arr)
	}
	return obj
}

func resultErr(err error) js.Value {
	obj := js.Global().Get("Object").New()
	obj.Set("ok", false)
	if err != nil {
		obj.Set("err", err.Error())
	}
	return obj
}

func wrap(fn func(args []js.Value) ([]byte, error)) js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		out, err := fn(args)
		if err != nil {
			return resultErr(err)
		}
		return resultOK(out)
	})
}

func wrapString(fn func(args []js.Value) (string, error)) js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		out, err := fn(args)
		if err != nil {
			return resultErr(err)
		}
		return resultOK([]byte(out))
	})
}

func errMissing(field string) error {
	return errors.New("missing argument: " + field)
}

func bytesArg(args []js.Value, i int) ([]byte, error) {
	if i >= len(args) {
		return nil, errMissing("argument")
	}
	v := args[i]
	if v.Type() == js.TypeString {
		return []byte(v.String()), nil
	}
	if v.Type() != js.TypeObject {
		return nil, errors.New("expected string or Uint8Array")
	}
	if v.Get("byteLength").Type() == js.TypeNumber {
		n := v.Get("byteLength").Int()
		buf := make([]byte, n)
		js.CopyBytesToGo(buf, v)
		return buf, nil
	}
	return []byte(v.String()), nil
}

func stringArg(args []js.Value, i int) (string, error) {
	if i >= len(args) || args[i].Type() != js.TypeString {
		return "", errMissing("string argument")
	}
	return args[i].String(), nil
}

func boolArg(args []js.Value, i int, def bool) bool {
	if i >= len(args) || args[i].Type() != js.TypeBoolean {
		return def
	}
	return args[i].Bool()
}

func intArg(args []js.Value, i int, def int) int {
	if i >= len(args) || args[i].Type() != js.TypeNumber {
		return def
	}
	return args[i].Int()
}

func optionalIndexArg(args []js.Value, i int) int {
	if i >= len(args) || args[i].Type() != js.TypeNumber {
		return -1
	}
	return args[i].Int()
}

func parseJSONArg(args []js.Value, i int, dest any) error {
	b, err := bytesArg(args, i)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return errMissing("json argument")
	}
	return json.Unmarshal(b, dest)
}

func compileOptsFromArg(args []js.Value, i int) compileoptsFromJS {
	var o compileoptsFromJS
	if i < len(args) && args[i].Type() == js.TypeObject {
		v := args[i]
		o.LinksOnly = v.Get("linksOnly").Truthy()
		o.NoScript = v.Get("noScript").Truthy()
	}
	return o
}

type compileoptsFromJS struct {
	LinksOnly bool
	NoScript  bool
}
