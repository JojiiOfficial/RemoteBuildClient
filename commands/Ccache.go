package commands

import (
	"fmt"
)

// ClearCcache clear ccache on server
func (cData *CommandData) ClearCcache() {
	// Do request
	resp, err := cData.Librb.ClearCcache()
	if err != nil {
		printResponseError(err, "clearing cache")
		return
	}

	// Print result if no
	// error occured
	fmt.Println(resp)
}

// QueryCcache get ccache info
func (cData *CommandData) QueryCcache() {
	// Do request
	resp, err := cData.Librb.QueryCcache()
	if err != nil {
		printResponseError(err, "querying cache")
		return
	}

	// Print result if no
	// error occured
	fmt.Println(resp.String)
}
