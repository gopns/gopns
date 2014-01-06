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

type AttributeDefinition struct {
	AttributeName string
	AttributeType string
}
type KeySchema struct {
	AttributeName string
	KeyType       string
}

type ProvisionedThroughput struct {
	ReadCapacityUnits  int
	WriteCapacityUnits int
}

type UpdateItemRequest struct {
	Key                         map[string]Attribute
	AttributeUpdates            map[string]AttributeUpdate
	ReturnConsumedCapacity      string `json:",omitempty"`
	ReturnItemCollectionMetrics string `json:",omitempty"`
	ReturnValues                string `json:",omitempty"`
	TableName                   string
}

type CreateTableRequest struct {
	AttributeDefinitions  []AttributeDefinition
	TableName             string
	KeySchema             []KeySchema
	ProvisionedThroughput ProvisionedThroughput
}

type GetItemRequest struct {
	Key       map[string]Attribute
	TableName string
}

type ScanRequest struct {
	ExclusiveStartKey map[string]Attribute `json:",omitempty"`
	Limit             int
	TableName         string
}

type ScanResponse struct {
	Items            []map[string]Attribute
	LastEvaluatedKey map[string]Attribute
}
