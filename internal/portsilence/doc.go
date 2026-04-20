// Package portsilence provides a time-bounded silence registry for ports.
// Silenced ports are suppressed from downstream pipeline stages for the
// configured duration, helping reduce noise from frequently toggling ports.
package portsilence
