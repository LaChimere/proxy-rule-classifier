package errors

var (
	InsufficientFields = NewError(10000001, "insufficient fields to resolve a Rule entity")
	TooManyFields      = NewError(10000002, "too many fields to resolve a Rule entity")
	InvalidPolicy      = NewError(10000003, "invalid policy type")
)
