package encoding

// TypeAndValue stores a JSON object with two attributes: a string "type"
// and a generic "value" (string) defined by type.  This type is used in
// a few places to implement the choice types that CBOR handles using tags.
type TypeAndValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
