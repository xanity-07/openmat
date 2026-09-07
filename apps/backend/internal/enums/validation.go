// Package enums contains all the enums data types for our application
package enums

type BindingSource string

const (
	BindingJSON  BindingSource = "json"
	BindingQuery BindingSource = "query"
	BindingURI   BindingSource = "uri"
)
