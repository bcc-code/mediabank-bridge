package vantage

import (
	"fmt"

	"github.com/ansel1/merry/v2"
	"github.com/davecgh/go-spew/spew"
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
	req.Post(uri, req.BodyJSON(data))
	return nil
}

func (c Client) getVariablesForWOrkflow(wName workflowName) error {
	// http://XXXX/Rest/Workflows/{ID}/JobInputs
	uri := c.makeURI("Workflows", fmt.Sprintf("%s/JobInputs", c.workflows[wName]))
	println(uri)
	res, err := req.Get(uri)
	if err != nil {
		return merry.Wrap(err)
	}

	data := subclipRequest{}
	if err := res.ToJSON(&data); err != nil {
		return err
	}

	params := []vantageWorkflowVariable{}

	for _, n := range workflowParamNames[wName] {
		found := false
		for _, v := range data.Variables {
			if v.Name != n {
				continue
			}
			params = append(params, v)
			found = true
			break
		}

		if !found {
			return merry.Errorf("Unable to find variable %s", n)
		}
	}

	spew.Dump(params)
	return nil
}

func (c Client) CreateSubclip() error {

	return c.getVariablesForWOrkflow(workflowCreateSubclip)

	/*
		data := subclipRequest{
			Attachments: []string{},
			JobName:     "Test job 1",
			Labels:      []string{},
			Medias:      []string{},
			Variables:   []vantageWorkflowVariable{},
		}

		c.StartWorkflow("asd", data)

		return nil

		/*
			{
			    "Attachments": [],
			    "JobName": "Populate with your job name",
			    "Labels": [],
			    "Medias": [],
			    "Priority": 0,
			    "Variables": [
			        {
			            "Identifier": "8eeac8c9-b3fd-4437-9f2c-1e342d0f14b8",
			            "DefaultValue": "0",
			            "Description": "",
			            "Name": "VS Subclip TC End Samples INT",
			            "TypeCode": "Int32",
			            "Value": "0"
			        },
			        {
			            "Identifier": "e6c77eb8-8c46-435d-b471-66ddf6be9efe",
			            "DefaultValue": "0",
			            "Description": "",
			            "Name": "VS Subclip TC Start Samples INT",
			            "TypeCode": "Int32",
			            "Value": "0"
			        },
			        {
			            "Identifier": "ce3aac9f-0f3b-4731-ab5b-0da9f6e0dacd",
			            "DefaultValue": "false",
			            "Description": "",
			            "Name": "VS API Subclip Title",
			            "TypeCode": "String",
			            "Value": "false"
			        },
			        {
			            "Identifier": "82f0f6fb-07da-45d1-912c-36ff3afb4878",
			            "DefaultValue": "False",
			            "Description": "",
			            "Name": "VS API Asset ID",
			            "TypeCode": "String",
			            "Value": "PLACEHOLDER1"
			        }
			    ]
			}*/
}
