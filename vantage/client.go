package vantage

import (
	"fmt"

	"github.com/ansel1/merry/v2"
	"github.com/bcc-code/mediabank-bridge/log"
	"github.com/imroc/req"
)

// Sentinel Errors
var (
	ErrAddressEmpty = merry.Sentinel("Address can not be empty")
)

// Client for vantage API
type Client struct {
	baseAddress string
	dryRun      bool
	workflows   WorkflowMap
}

// NewClient with specified settings
func NewClient(s ClientSettings) (*Client, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}

	c := Client{
		baseAddress: s.Address,
		dryRun:      s.DryRun,
		workflows:   s.Workflows,
	}

	return &c, nil
}

func (c Client) makeURI(action, path string) string {
	// http://XXXXXXX/Rest/Workflows/YYYYYYY/Submit
	return fmt.Sprintf("%s/Rest/%s/%s", c.baseAddress, action, path)
}

// StartWorkflow with specified params
func (c Client) StartWorkflow(wID workflowID, data interface{}) error {
	uri := c.makeURI("Workflows", fmt.Sprintf("%s/Submit", wID))
	res, err := req.Post(uri, req.BodyJSON(data))
	log.L.Debug().Str("return body", res.String()).Msg("Return value from vantage")
	return err
}

func (c Client) getVariablesForWOrkflow(wName workflowName) (map[workflowParamName]*vantageWorkflowVariable, error) {
	// http://XXXX/Rest/Workflows/{ID}/JobInputs
	uri := c.makeURI("Workflows", fmt.Sprintf("%s/JobInputs", c.workflows[wName]))
	println(uri)
	res, err := req.Get(uri)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	data := subclipRequest{}
	if err := res.ToJSON(&data); err != nil {
		return nil, merry.Wrap(err)
	}

	params := map[workflowParamName]*vantageWorkflowVariable{}

	for _, n := range workflowParamNames[wName] {
		found := false
		for _, v := range data.Variables {
			if v.Name != n {
				continue
			}
			params[n] = &v
			found = true
			break
		}

		if !found {
			return nil, merry.Errorf("Unable to find variable %s", n)
		}
	}

	return params, nil
}

// CreateSubclipParams to be passed to the workflow
type CreateSubclipParams struct {
	In      string
	Out     string
	AssetID string
	Title   string
}

// ToWorkflowParams fills the Vantage style params
func (p CreateSubclipParams) ToWorkflowParams(tpl map[workflowParamName]*vantageWorkflowVariable) []vantageWorkflowVariable {
	tpl[paramStartSamples].Value = p.In
	tpl[paramEndSamples].Value = p.Out
	tpl[paramSubclipTitle].Value = p.Title
	tpl[paramAPIAssetID].Value = p.AssetID
	return mapToList(tpl)
}

// CreateSubclip on the specified asset ID
func (c Client) CreateSubclip(in CreateSubclipParams) error {
	paramsTemplate, err := c.getVariablesForWOrkflow(workflowCreateSubclip)
	if err != nil {
		return err
	}

	data := subclipRequest{
		Attachments: []string{},
		JobName:     fmt.Sprintf("Create subclip - %s - %s", in.AssetID, in.Title),
		Labels:      []string{},
		Medias:      []string{},
		Variables:   in.ToWorkflowParams(paramsTemplate),
	}

	return c.StartWorkflow(c.workflows[workflowCreateSubclip], data)
}
