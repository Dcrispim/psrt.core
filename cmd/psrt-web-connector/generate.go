//go:build windows

package main

//go:generate go run github.com/tc-hib/go-winres@v0.3.3 make --in winres/winres.json --out resource.syso --arch amd64
