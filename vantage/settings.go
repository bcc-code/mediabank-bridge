package vantage

import "github.com/ansel1/merry/v2"

// ClientSettings to be passed into NewClient()
type ClientSettings struct {
	Address   string
	DryRun    bool
	Workflows WorkflowMap
}

// Validate that the settings make some basic sense
func (cs ClientSettings) Validate() error {
	if cs.Address == "" {
		return merry.Wrap(ErrAddressEmpty)
	}

	return cs.Workflows.Validate()
}
