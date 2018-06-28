package gads

import "encoding/xml"

// CommonAd define the parent type Ad type as defined
// https://developers.google.com/adwords/api/docs/reference/v201806/AdGroupAdService.Ad
type CommonAd struct {
	Type                string            `xml:"xsi:type,attr,omitempty"`
	ID                  int64             `xml:"id,omitempty"`
	URL                 string            `xml:"url,omitempty"`
	DisplayURL          string            `xml:"displayUrl,omitempty"`
	FinalURLs           []string          `xml:"finalUrls,omitempty"`
	FinalMobileURLs     []string          `xml:"finalMobileUrls,omitempty"`
	FinalAppURLs        []AppUrl          `xml:"finalAppUrls,omitempty"`
	TrackingURLTemplate *string           `xml:"trackingUrlTemplate"`
	URLCustomParameters *CustomParameters `xml:"urlCustomParameters,omitempty"`
	DevicePreference    int64             `xml:"devicePreference,omitempty"`
}

// Ad is the common interface for accessing properties shared by all
// child types of ads
type Ad interface {
	GetID() int64
	GetURL() string
	GetType() string
	GetTrackingURLTemplate() *string
	GetFinalURLs() []string

	CloneForTemplate([]string, *string) Ad
}

// GetID returns the Ad's id
func (c CommonAd) GetID() int64 {
	return c.ID
}

// Get Type returns the type of the Ad
func (c CommonAd) GetType() string {
	return c.Type
}

// GetURL returns the Ad's url (a.k.a destination url)
func (c CommonAd) GetURL() string {
	return c.URL
}

// GetTrackingURLTemplate returns the  tracking url template
func (c CommonAd) GetTrackingURLTemplate() *string {
	return c.TrackingURLTemplate
}

// GetFinalURLs returns the list of Final urls (can be empty array)
func (c CommonAd) GetFinalURLs() []string {
	return c.FinalURLs
}

// CloneForTemplate create a clone of an Ad, to recreate it for changing the tracking Url Template (as Ad are immutable)
func (c CommonAd) CloneForTemplate(finalURLs []string, trackingURLTemplate *string) Ad {
	c.ID = 0   // value used by go for omitempty
	c.URL = "" // template needs an empty destination url (as it deprecates this field)
	c.FinalURLs = finalURLs
	c.TrackingURLTemplate = trackingURLTemplate
	return c
}

// CloneForTemplate create a clone of an Ad, to recreate it for changing the tracking Url Template
func (c TextAd) CloneForTemplate(finalURLs []string, trackingURLTemplate *string) Ad {
	c.ID = 0   // value used by go for omitempty
	c.URL = "" // template needs an empty destination url (as it deprecates this field)
	c.FinalURLs = finalURLs
	c.TrackingURLTemplate = trackingURLTemplate

	return c
}

// CloneForTemplate create a clone of an Ad, to recreate it for changing the tracking Url Template
func (c ImageAd) CloneForTemplate(finalURLs []string, trackingURLTemplate *string) Ad {
	c.ID = 0   // value used by go for omitempty
	c.URL = "" // template needs an empty destination url (as it deprecates this field)
	c.FinalURLs = finalURLs
	c.TrackingURLTemplate = trackingURLTemplate
	return c
}

// CloneForTemplate create a clone of an Ad, to recreate it for changing the tracking Url Template
func (c TemplateAd) CloneForTemplate(finalURLs []string, trackingURLTemplate *string) Ad {
	c.ID = 0   // value used by go for omitempty
	c.URL = "" // template needs an empty destination url (as it deprecates this field)
	c.FinalURLs = finalURLs
	c.TrackingURLTemplate = trackingURLTemplate
	return c
}

// CloneForTemplate create a clone of an Ad, to recreate it for changing the tracking Url Template
func (c DynamicSearchAd) CloneForTemplate(finalURLs []string, trackingURLTemplate *string) Ad {
	c.ID = 0   // value used by go for omitempty
	c.URL = "" // template needs an empty destination url (as it deprecates this field)
	c.FinalURLs = finalURLs
	c.TrackingURLTemplate = trackingURLTemplate
	return c
}

// CloneForTemplate create a clone of an Ad, to recreate it for changing the tracking Url Template
func (c ExpandedTextAd) CloneForTemplate(finalURLs []string, trackingURLTemplate *string) Ad {
	c.ID = 0   // value used by go for omitempty
	c.URL = "" // template needs an empty destination url (as it deprecates this field)
	c.FinalURLs = finalURLs
	c.TrackingURLTemplate = trackingURLTemplate
	return c
}

// TextAd represents the TextAd object as documented here
// https://developers.google.com/adwords/api/docs/reference/v201806/AdGroupAdService.TextAd
type TextAd struct {
	CommonAd
	Headline     string `xml:"headline"`
	Description1 string `xml:"description1"`
	Description2 string `xml:"description2"`
}

type ProductAd struct {
	CommonAd
}

// DynamicSearchAd represents the equivalent object documented here
// https://developers.google.com/adwords/api/docs/reference/v201806/AdGroupAdService.DynamicSearchAd
type DynamicSearchAd struct {
	CommonAd
	Description1 string `xml:"description1"`
	Description2 string `xml:"description2"`
}

// ImageAd represents the equivalent object documented here
// https://developers.google.com/adwords/api/docs/reference/v201806/AdGroupAdService.ImageAd
type ImageAd struct {
	CommonAd
	Image             int64  `xml:"imageId"` //TODO should actually be Image object, not just an int
	Name              string `xml:"name"`
	AdToCopyImageFrom int64  `xml:"adToCopyImageFrom"`
}

// TemplateAd represents the equivalent object documented here
// https://developers.google.com/adwords/api/docs/reference/v201806/AdGroupAdService.TemplateAd
type TemplateAd struct {
	CommonAd
	TemplateID       int64             `xml:"templateId"`
	AdUnionID        int64             `xml:"adUnionId>id"`
	TemplateElements []TemplateElement `xml:"templateElements"`
	Dimensions       []Dimensions      `xml:"dimensions"`
	Name             string            `xml:"name"`
	Duration         int64             `xml:"duration"`
	originAdID       *int64            `xml:"originAdId"`
}

// ExpandedTextAd epresents the equivalent object documented here
// https://developers.google.com/adwords/api/docs/reference/v201806/AdGroupAdService.ExpandedTextAd
type ExpandedTextAd struct {
	CommonAd
	HeadlinePart1 string `xml:"headlinePart1"`
	HeadlinePart2 string `xml:"headlinePart2"`
	Description   string `xml:"description"`
	Path1         string `xml:"path1"`
	Path2         string `xml:"path2"`
}

// MarshalXML returns unimplemented error as the structure does not
// match yet 100% of the field required by google api
func (c ImageAd) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return ERROR_NOT_YET_IMPLEMENTED
}

// MarshalXML returns unimplemented error as the structure does not
// match yet 100% of the field required by google api
func (c TemplateAd) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return ERROR_NOT_YET_IMPLEMENTED
}
