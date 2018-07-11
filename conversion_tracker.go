package gads

import (
	"encoding/xml"
	"fmt"
)

type CommonConversionTracker struct {
	ID                            string  `xml:"id,omitempty"`
	OriginalConversionTypeID      int64   `xml:"originalConversionTypeId,omitempty"`
	Name                          string  `xml:"name,omitempty"`
	Status                        string  `xml:"status,omitempty"`
	Category                      string  `xml:"category,omitempty"`
	ConversionTypeOwnerCustomerID int64   `xml:"conversionTypeOwnerCustomerId,omitempty"`
	ViewthroughLookbackWindow     int64   `xml:"viewthroughLookbackWindow,omitempty"`
	CtcLookbackWindow             int64   `xml:"ctcLookbackWindow,omitempty"`
	CountingType                  string  `xml:"countingType,omitempty"`
	DefaultRevenueValue           float64 `xml:"defaultRevenueValue,omitempty"`
	DefaultRevenueCurrencyCode    string  `xml:"defaultRevenueCurrencyCode,omitempty"`
	AlwaysUseDefaultRevenueValue  bool    `xml:"alwaysUseDefaultRevenueValue,omitempty"`
	MostRecentConversionDate      string  `xml:"mostRecentConversionDate,omitempty"`
	LastReceivedRequestTime       string  `xml:"lastReceivedRequestTime,omitempty"`
}

// ConversionTracker is an interface for conversion tackers
type ConversionTracker interface {
	GetStatus() string
	SetStatus(string)
	GetID() string
}

// SetStatus setter for conversion tracker status
func (c *CommonConversionTracker) SetStatus(status string) {
	c.Status = status
}

// GetID getter for conversion tracker ID
func (c *CommonConversionTracker) GetID() string {
	return c.ID
}

// GetStatus getter for conversion tracker status
func (c *CommonConversionTracker) GetStatus() string {
	return c.Status
}

type AdWordsConversionTracker struct {
	*CommonConversionTracker
	Type                   string `xml:"xsi:type,attr,omitempty"`
	TextFormat             string `xml:"textFormat,omitempty"`
	ConversionPageLanguage string `xml:"ConversionPageLanguage,omitempty"`
	BackgroundColor        string `xml:"backgroundColor,omitempty"`
	TrackingCodeType       string `xml:"trackingCodeType,omitempty"`
}

type ConversionTrackers []ConversionTracker
type ConversionTrackerOperations map[string][]ConversionTracker

func (cts *ConversionTrackers) UnmarshalXML(
	dec *xml.Decoder,
	start xml.StartElement,
) error {
	typeName := xml.Name{
		Space: "http://www.w3.org/2001/XMLSchema-instance",
		Local: "type",
	}

	typeValue, err := findAttr(start.Attr, typeName)
	if err != nil {
		return err
	}

	switch typeValue {
	case "AdWordsConversionTracker":
		item := AdWordsConversionTracker{}
		err := dec.DecodeElement(&item, &start)
		if err != nil {
			return err
		}
		*cts = append(*cts, item)
	default:
		return fmt.Errorf("unknown Conversion Tracker -> %#v", typeName)
	}

	return nil
}

type ConversionTrackerService struct {
	Auth
}

func NewConversionTrackerService(auth *Auth) *ConversionTrackerService {
	return &ConversionTrackerService{Auth: *auth}
}

func (s *ConversionTrackerService) Mutate(
	conversionTrackerOperations ConversionTrackerOperations,
) (conversionTrackers ConversionTrackers, err error) {

	//TODO: there should be a way to factorize things so that one
	// should only have to do a call
	// func (
	//     auth Auth,
	//	   itemsOperations map[string][]interface{},
	//	   mutateResp *interface{},
	// ) (err error) {
	// in all Mutate

	type operation struct {
		Action string      `xml:"operator"`
		Item   interface{} `xml:"operand"`
	}
	operations := []operation{}
	for action, items := range conversionTrackerOperations {
		for _, item := range items {
			operations = append(
				operations,
				operation{
					Action: action,
					Item:   item,
				},
			)
		}
	}
	mutation := struct {
		XMLName xml.Name
		Ops     []operation `xml:"operations"`
	}{
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "mutate",
		},
		Ops: operations,
	}
	respBody, err := s.Auth.request(conversionTrackerServiceUrl, "mutate", mutation)
	if err != nil {
		return conversionTrackers, err
	}
	mutateResp := struct {
		BaseResponse
		ConversionTrackers ConversionTrackers `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return conversionTrackers, err
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.ConversionTrackers, err
}

func (s *ConversionTrackerService) Get(selector Selector) (
	items ConversionTrackers,
	totalCount int64,
	err error,
) {
	selector.XMLName = xml.Name{"", "serviceSelector"}
	respBody, err := s.Auth.request(
		conversionTrackerServiceUrl,
		"get",
		struct {
			XMLName xml.Name
			Sel     Selector
		}{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "get",
			},
			Sel: selector,
		},
	)
	if err != nil {
		return items, totalCount, err
	}
	getResp := struct {
		Size  int64              `xml:"rval>totalNumEntries"`
		Items ConversionTrackers `xml:"rval>entries"`
	}{}

	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return items, totalCount, err
	}
	return getResp.Items, getResp.Size, err
}
