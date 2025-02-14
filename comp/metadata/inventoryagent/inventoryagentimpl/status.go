// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package inventoryagentimpl

import (
	"embed"
	"io"

	"github.com/DataDog/datadog-agent/comp/core/status"
)

//go:embed status_templates
var templatesFS embed.FS

// Name renders the name
func (ia *inventoryagent) Name() string {
	return "metadata"
}

// Index renders the index
func (ia *inventoryagent) Index() int {
	return 3
}

// JSON populates the status map
func (ia *inventoryagent) JSON(_ bool, stats map[string]interface{}) error {
	for k, v := range ia.Get() {
		stats[k] = v
	}

	return nil
}

// Text renders the text output
func (ia *inventoryagent) Text(_ bool, buffer io.Writer) error {
	return status.RenderText(templatesFS, "inventory.tmpl", buffer, ia.Get())
}

// HTML renders the html output
func (ia *inventoryagent) HTML(_ bool, buffer io.Writer) error {
	return status.RenderHTML(templatesFS, "inventoryHTML.tmpl", buffer, ia.Get())
}
