package dynamodb

type AttributeUpdate struct {
	Action string
	Value  Attribute
}

type Attribute struct {
	S  string   `json:",omitempty"`
	SS []string `json:",omitempty"`

	N  string   `json:",omitempty"`
	NS []string `json:",omitempty"`

	B  string   `json:",omitempty"`
	BS []string `json:",omitempty"`
}

type UpdateItemRequest struct {
	Key                         map[string]Attribute
	AttributeUpdates            map[string]AttributeUpdate
	ReturnConsumedCapacity      string `json:",omitempty"`
	ReturnItemCollectionMetrics string `json:",omitempty"`
	ReturnValues                string `json:",omitempty"`
	TableName                   string
}
