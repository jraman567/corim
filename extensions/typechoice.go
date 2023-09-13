package extensions

type ITypeChoiceValue interface {
	// String returns the string representation of the ITypeChoiceValue.
	String() string
	// Valid returns an error if validation of the ITypeChoiceValue fails,
	// or nil if it succeeds.
	Valid() error
	// Type returns the type name of this ITypeChoiceValue implementation.
	Type() string
}
