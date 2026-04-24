// Package portscope provides a configurable gate that restricts port scanning
// to explicitly defined protocol/range pairs. An unconfigured Scope passes all
// ports through, making it safe to embed unconditionally in the scan pipeline.
package portscope
