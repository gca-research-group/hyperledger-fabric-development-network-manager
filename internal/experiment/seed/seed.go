// Package seed provides the complete example used to start experiments.
package seed

import _ "embed"

// YAML contains every configuration field, including unused optional settings.
//
//go:embed seed.yaml
var YAML string
