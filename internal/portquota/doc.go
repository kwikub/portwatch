// Package portquota flags ports that open more than a configured number of
// times within a rolling time window, helping surface noisy or flapping
// services before they saturate downstream alerting.
package portquota
