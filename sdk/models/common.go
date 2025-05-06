package models

type EmptyResponse struct{}

type Condition struct {
	Column  `json:",inline" yaml:",inline"`
	Filters []Filter `json:"filters,omitempty" yaml:"filters,omitempty"`
}

type Column struct {
	Key    string `json:"key,omitempty" yaml:"key,omitempty"`
	Origin string `json:"origin,omitempty" yaml:"origin,omitempty"`
	Type   string `json:"type,omitempty" yaml:"type,omitempty"`
}

type Filter struct {
	Op    string      `json:"op,omitempty" yaml:"op,omitempty"`
	Value interface{} `json:"value,omitempty" yaml:"value,omitempty"`
}

type Group struct {
	Conditions []Condition `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	Operator   string      `json:"operator,omitempty" yaml:"operator,omitempty"`
	Groups     []Group     `json:"groups,omitempty" yaml:"groups,omitempty"`
}

type Pipeline struct {
	Selectors []Selector `json:"selectors,omitempty" yaml:"selectors,omitempty"`
	Except    []Selector `json:"except,omitempty" yaml:"except,omitempty"`
	From      *Pipeline  `json:"from,omitempty" yaml:"from,omitempty"`
	Filters   *Group     `json:"filters,omitempty" yaml:"filters,omitempty"`
	GroupBy   []Selector `json:"groupBy,omitempty" yaml:"groupBy,omitempty"`
	Having    *Group     `json:"having,omitempty" yaml:"having,omitempty"`
	OrderBy   []OrderBy  `json:"orderBy,omitempty" yaml:"orderBy,omitempty"`
	Limit     uint64     `json:"limit,omitempty" yaml:"limit,omitempty"`
	Offset    uint64     `json:"offset,omitempty" yaml:"offset,omitempty"`
}

type Selector struct {
	Column     `yaml:",inline"`
	Processors []Processor `json:"processors,omitempty" yaml:"processors,omitempty"`
	Alias      string      `json:"alias,omitempty" yaml:"alias,omitempty"`
}

type Processor struct {
	Op   string   `json:"op,omitempty" yaml:"op,omitempty"`
	Args []string `json:"args,omitempty" yaml:"args,omitempty"`
}

type OrderBy struct {
	Selector  Selector `json:"selector"`
	Direction string   `json:"direction"`
}

type PolicyMetadata struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type PromqlPipeline struct {
	Metric     string          `json:"metric,omitempty" yaml:"metric,omitempty"`
	Conditions []Condition     `json:"conditions,omitempty" yaml:"conditions,omitempty"`
	Function   *PromqlFunction `json:"function,omitempty" yaml:"function,omitempty"`
	Template   string          `json:"template,omitempty" yaml:"template,omitempty"`
}

type PromqlFunction struct {
	Name      string           `json:"name,omitempty" yaml:"name,omitempty"`
	Pipelines []PromqlPipeline `json:"pipelines,omitempty" yaml:"pipelines,omitempty"`
	Args      []string         `json:"args,omitempty" yaml:"args,omitempty"`
}

type KnownPipelines map[string]PromqlPipeline

func NewEqualStringCondition(key, value string) Condition {
	return Condition{
		Column: Column{
			Key:    key,
			Type:   ColumnTypeString,
			Origin: ColumnOriginRoot,
		},
		Filters: []Filter{
			{
				Op:    OpEqual,
				Value: value,
			},
		},
	}
}
