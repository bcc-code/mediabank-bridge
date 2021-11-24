package vantage

type typeCode string

// TypeCode contants
const (
	Int32  typeCode = "Int32"
	String typeCode = "String"
)

type subclipRequest struct {
	Attachments []string
	JobName     string
	Labels      []string
	Medias      []string
	Priority    uint32
	Variables   []vantageWorkflowVariable
}

type vantageWorkflowVariable struct {
	Identifier   string
	DefaultValue string
	Description  string
	Name         workflowParamName
	TypeCode     typeCode
	Value        string
}
