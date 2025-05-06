package models

const (
	ColumnTypeString = "string"
	ColumnOriginRoot = "root"
)

// Filter Operators
const (
	OperatorEqual                 = "eq"
	OperatorNotEqual              = "ne"
	OperatorContains              = "contains"
	OperatorNotContains           = "notcontains"
	OperatorContainsIgnoreCase    = "icontains"
	OperatorNotContainsIgnoreCase = "inotcontains"
	OperatorStartsWith            = "startswith"
	OperatorStartsWithIgnoreCase  = "istartswith"
)

// Filter keys
const (
	Workload  = "workload"
	Cluster   = "cluster"
	Env       = "env"
	Namespace = "namespace"
	Instance  = "instance"
)
