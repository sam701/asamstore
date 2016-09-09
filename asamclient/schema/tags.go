package schema

type AttributeChange struct {
	AttributeName string              `json:"attributeName,omitempty"`
	ChangeType    AttributeChangeType `json:"type,omitempty"`
	Values        []string            `json:"values,omitempty"`
}

type AttributeChangeType string

const (
	AttributeChangeTypeSet    AttributeChangeType = "set"
	AttributeChangeTypeAdd                        = "add"
	AttributeChangeTypeDelete                     = "del"
)
