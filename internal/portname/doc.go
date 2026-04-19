// Package portname maps port/protocol pairs to human-readable service names.
// It ships with a built-in table of common well-known ports and supports
// custom registrations at runtime.
//
// Basic usage:
//
//	name, ok := portname.Lookup(80, "tcp")
//	if ok {
//		fmt.Println(name) // "http"
//	}
//
// Custom entries can be registered at startup:
//
//	portname.Register(8080, "tcp", "my-service")
package portname
