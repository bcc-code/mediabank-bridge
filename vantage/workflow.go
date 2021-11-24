package vantage

import "github.com/ansel1/merry/v2"

type workflowName string
type workflowParamName string
type workflowID string

// WorkflowMap for mapping the name of the worflow to it's ID
type WorkflowMap map[workflowName]workflowID

// Sentinel errors
var (
	ErrWorkflowIDMissing = merry.Sentinel("Missing workflow ID")
)

const (
	workflowCreateSubclip workflowName = "create_subclip"
)

var workflowNamesList = []workflowName{
	workflowCreateSubclip,
}

var (
	paramEndSamples   workflowParamName = "VS Subclip TC End Samples INT"
	paramStartSamples workflowParamName = "VS Subclip TC Start Samples INT"
	paramSubclipTitle workflowParamName = "VS API Subclip Title"
	paramAPIAssetID   workflowParamName = "VS API Asset ID"
)

var workflowParamNames = map[workflowName](workflowParamNameList){
	workflowCreateSubclip: workflowParamNameList{
		paramEndSamples,
		paramStartSamples,
		paramSubclipTitle,
		paramAPIAssetID,
	},
}

// Validate that all workflow ids are provided
func (wm WorkflowMap) Validate() error {
	for _, name := range workflowNamesList {
		if _, ok := wm[name]; !ok {
			return merry.Wrap(ErrWorkflowIDMissing, merry.WithMessagef("Missing workflow id for %s", name))
		}
	}
	return nil
}

type workflowParamNameList []workflowParamName

func (l workflowParamNameList) Contains(name string) bool {
	for _, n := range l {
		if string(n) == name {
			return true
		}
	}

	return false
}
