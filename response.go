package gads

// BaseResponse is the common structure of every response received from
// Google Adwords in the most common case
type BaseResponse struct {
	PartialFailureErrors PartialFailureErrors `xml:"rval>partialFailureErrors,omitempty"`
}
