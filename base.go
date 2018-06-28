package gads

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

const (
	// https://developers.google.com/adwords/api/docs/reference/
	apiVersion         = "v201806"
	baseUrl            = "https://adwords.google.com/api/adwords/cm/" + apiVersion
	rmktgBaseUrl       = "https://adwords.google.com/api/adwords/rm/" + apiVersion
	managedCustomerUrl = "https://adwords.google.com/api/adwords/mcm/" + apiVersion
	reportAPIURL       = "https://adwords.google.com/api/adwords/reportdownload/" + apiVersion
	// used for developpement, if true all unknown field will raise an error
	StrictMode = false
)

type HasXSIType interface {
	GetType() string
}

type ServiceUrl struct {
	Url  string
	Name string
}

type AWQLQuery struct {
	XMLName xml.Name
	Query   string `xml:"query"`
}

var (
	// exceptions
	ERROR_NOT_YET_IMPLEMENTED = fmt.Errorf("Not yet implemented")
	hasXSIType                = reflect.TypeOf((*HasXSIType)(nil)).Elem()
)

var (
	configJson = flag.String("config_json", "./config.json", "API credentials")

	// service urls
	adGroupAdServiceUrl                = ServiceUrl{baseUrl, "AdGroupAdService"}
	adGroupBidModifierServiceUrl       = ServiceUrl{baseUrl, "AdGroupBidModifierService"}
	adGroupCriterionServiceUrl         = ServiceUrl{baseUrl, "AdGroupCriterionService"}
	adGroupFeedServiceUrl              = ServiceUrl{baseUrl, "AdGroupFeedService"}
	adGroupServiceUrl                  = ServiceUrl{baseUrl, "AdGroupService"}
	adParamServiceUrl                  = ServiceUrl{baseUrl, "AdParamService"}
	adwordsUserListServiceUrl          = ServiceUrl{rmktgBaseUrl, "AdwordsUserListService"}
	biddingStrategyServiceUrl          = ServiceUrl{baseUrl, "BiddingStrategyService"}
	budgetOrderServiceUrl              = ServiceUrl{baseUrl, "BudgetOrderService"}
	budgetServiceUrl                   = ServiceUrl{baseUrl, "BudgetService"}
	campaignAdExtensionServiceUrl      = ServiceUrl{baseUrl, "CampaignAdExtensionService"}
	campaignCriterionServiceUrl        = ServiceUrl{baseUrl, "CampaignCriterionService"}
	campaignFeedServiceUrl             = ServiceUrl{baseUrl, "CampaignFeedService"}
	campaignServiceUrl                 = ServiceUrl{baseUrl, "CampaignService"}
	campaignExtensionSettingServiceUrl = ServiceUrl{baseUrl, "CampaignExtensionSettingService"}
	campaignSharedSetServiceUrl        = ServiceUrl{baseUrl, "CampaignSharedSetService"}
	constantDataServiceUrl             = ServiceUrl{baseUrl, "ConstantDataService"}
	conversionTrackerServiceUrl        = ServiceUrl{baseUrl, "ConversionTrackerService"}
	customerFeedServiceUrl             = ServiceUrl{baseUrl, "CustomerFeedService"}
	customerServiceUrl                 = ServiceUrl{managedCustomerUrl, "CustomerService"}
	customerSyncServiceUrl             = ServiceUrl{baseUrl, "CustomerSyncService"}
	dataServiceUrl                     = ServiceUrl{baseUrl, "DataService"}
	experimentServiceUrl               = ServiceUrl{baseUrl, "ExperimentService"}
	feedItemServiceUrl                 = ServiceUrl{baseUrl, "FeedItemService"}
	feedMappingServiceUrl              = ServiceUrl{baseUrl, "FeedMappingService"}
	feedServiceUrl                     = ServiceUrl{baseUrl, "FeedService"}
	geoLocationServiceUrl              = ServiceUrl{baseUrl, "GeoLocationService"}
	labelServiceUrl                    = ServiceUrl{baseUrl, "LabelService"}
	locationCriterionServiceUrl        = ServiceUrl{baseUrl, "LocationCriterionService"}
	managedCustomerServiceUrl          = ServiceUrl{managedCustomerUrl, "ManagedCustomerService"}
	mediaServiceUrl                    = ServiceUrl{baseUrl, "MediaService"}
	mutateJobServiceUrl                = ServiceUrl{baseUrl, "Mutate_JOB_Service"}
	offlineConversionFeedServiceUrl    = ServiceUrl{baseUrl, "OfflineConversionFeedService"}
	reportDefinitionServiceUrl         = ServiceUrl{baseUrl, "ReportDefinitionService"}
	sharedCriterionServiceUrl          = ServiceUrl{baseUrl, "SharedCriterionService"}
	sharedSetServiceUrl                = ServiceUrl{baseUrl, "SharedSetService"}
	targetingIdeaServiceUrl            = ServiceUrl{baseUrl, "TargetingIdeaService"}
	trafficEstimatorServiceUrl         = ServiceUrl{baseUrl, "TrafficEstimatorService"}
)

func (s ServiceUrl) String() string {
	return s.Url + "/" + s.Name
}

type Auth struct {
	CustomerId     string
	DeveloperToken string
	UserAgent      string
	ValidateOnly   bool         `json:"-"`
	PartialFailure bool         `json:"-"`
	Testing        *testing.T   `json:"-"`
	Client         *http.Client `json:"-"`
}

// Date is a google date, a simple type inference with methods
type Date time.Time

// String gives the google representation of a date
func (d Date) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(time.Time(d).Format("20060102"), start)
}

//
// Selector structs
//
type DateRange struct {
	Min Date `xml:"min"`
	Max Date `xml:"max"`
}

type Predicate struct {
	Field    string   `xml:"field"`
	Operator string   `xml:"operator"`
	Values   []string `xml:"values"`
}

type OrderBy struct {
	Field     string `xml:"field"`
	SortOrder string `xml:"sortOrder"`
}

type Paging struct {
	Offset int64 `xml:"startIndex"`
	Limit  int64 `xml:"numberResults"`
}

type Selector struct {
	XMLName    xml.Name
	Fields     []string    `xml:"fields,omitempty"`
	Predicates []Predicate `xml:"predicates"`
	DateRange  *DateRange  `xml:"dateRange,omitempty"`
	Ordering   []OrderBy   `xml:"ordering"`
	Paging     *Paging     `xml:"paging,omitempty"`
}

type UrlList struct {
	Urls []string `xml:"urls"`
}

// error parsers
func selectorError() (err error) {
	return err
}

func (a *Auth) request(
	serviceUrl ServiceUrl,
	action string,
	body interface{},
) (respBody []byte, err error) {

	type devToken struct {
		XMLName xml.Name
	}
	type soapReqHeader struct {
		XMLName          xml.Name
		UserAgent        string `xml:"userAgent"`
		DeveloperToken   string `xml:"developerToken"`
		ClientCustomerId string `xml:"clientCustomerId,omitempty"`
		PartialFailure   bool   `xml:"partialFailure,omitempty"`
		ValidateOnly     bool   `xml:"validateOnly,omitempty"`
	}

	type SoapReqBody struct {
		Body interface{}
	}

	type soapReqEnvelope struct {
		XMLName      xml.Name
		Header       soapReqHeader `xml:"Header>RequestHeader"`
		XSINamespace string        `xml:"xmlns:xsi,attr"`
		Body         SoapReqBody   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	}

	// we make  sure the xsi:type attribute is here for every single
	// structs in the body
	body = addXSIType(body)
	soapReqBody := SoapReqBody{body}

	reqBody, err := xml.MarshalIndent(
		soapReqEnvelope{
			XMLName: xml.Name{"http://schemas.xmlsoap.org/soap/envelope/", "Envelope"},
			Header: soapReqHeader{
				XMLName:          xml.Name{serviceUrl.Url, "RequestHeader"},
				UserAgent:        a.UserAgent,
				DeveloperToken:   a.DeveloperToken,
				ClientCustomerId: a.CustomerId,
				PartialFailure:   a.PartialFailure,
				ValidateOnly:     a.ValidateOnly,
			},
			XSINamespace: "http://www.w3.org/2001/XMLSchema-instance",
			Body:         soapReqBody,
		},
		"  ",
		"  ",
	)
	if err != nil {
		return []byte{}, err
	}

	//fmt.Println(string(reqBody[:]))

	req, err := http.NewRequest("POST", serviceUrl.String(), bytes.NewReader(reqBody))
	req.Header.Add("Accept", "text/xml")
	req.Header.Add("Accept", "multipart/*")
	req.Header.Add("Content-Type", "text/xml;charset=UTF-8")
	contentLength := fmt.Sprintf("%d", len(reqBody))
	req.Header.Add("Content-length", contentLength)
	req.Header.Add("SOAPAction", action)
	if a.Testing != nil {
		a.Testing.Logf("request ->\n%s\n%#v\n%s\n", req.URL.String(), req.Header, string(reqBody))
	}
	resp, err := a.Client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	respBody, err = ioutil.ReadAll(resp.Body)
	if a.Testing != nil {
		a.Testing.Logf("respBody ->\n%s\n%s\n", string(respBody), resp.Status)
	}

	//fmt.Println(string(respBody[:]))

	type soapRespHeader struct {
		RequestId    string `xml:"requestId"`
		ServiceName  string `xml:"serviceName"`
		MethodName   string `xml:"methodName"`
		Operations   int64  `xml:"operations"`
		ResponseTime int64  `xml:"responseTime"`
	}

	type soapRespBody struct {
		Response []byte `xml:",innerxml"`
	}

	soapResp := struct {
		XMLName xml.Name       `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
		Header  soapRespHeader `xml:"Header>RequestHeader"`
		Body    soapRespBody   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	}{}

	err = xml.Unmarshal([]byte(respBody), &soapResp)
	if err != nil {
		return respBody, err
	}
	if resp.StatusCode == 400 || resp.StatusCode == 401 || resp.StatusCode == 403 || resp.StatusCode == 405 || resp.StatusCode == 500 {
		fault := Fault{}
		// fmt.Printf("unknown error ->\n%s\n", string(soapResp.Body.Response))
		err = xml.Unmarshal(soapResp.Body.Response, &fault)
		if err != nil {
			return respBody, err
		}
		return soapResp.Body.Response, &fault.Errors
	}
	return soapResp.Body.Response, err
}

// walk through all the fields recursively
// and set the XSI types if necessary
func addXSIType(obj interface{}) interface{} {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)

	copy := reflect.New(original.Type()).Elem()
	addXSITypeRecursive(copy, original)

	// Remove the reflection wrapper
	return copy.Interface()
}

func addXSITypeRecursive(copy, original reflect.Value) {
	switch original.Kind() {
	// The first cases handle nested structures recursively

	// If it is a pointer we need to unwrap and call once again
	case reflect.Ptr:
		// To get the actual value of the original we have to call Elem()
		// At the same time this unwraps the pointer so we don't end up in
		// an infinite recursion
		originalValue := original.Elem()
		// Check if the pointer is nil
		if !originalValue.IsValid() {
			return
		}
		// Allocate a new object and set the pointer to it
		copy.Set(reflect.New(originalValue.Type()))
		// Unwrap the newly created pointer
		addXSITypeRecursive(copy.Elem(), originalValue)

	// If it is an interface (which is very similar to a pointer), do basically the
	// same as for the pointer. Though a pointer is not the same as an interface so
	// note that we have to call Elem() after creating a new object because otherwise
	// we would end up with an actual pointer
	case reflect.Interface:
		// Get rid of the wrapping interface
		originalValue := original.Elem()
		if !originalValue.IsValid() {
			return
		}

		// Create a new object. Now new gives us a pointer, but we want the value it
		// points to, so we have to call Elem() to unwrap it
		copyValue := reflect.New(originalValue.Type()).Elem()
		addXSITypeRecursive(copyValue, originalValue)
		copy.Set(copyValue)

	// If it is a struct we look for each field
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			f := original.Field(i)
			structField := original.Type().Field(i)

			// If the struct has an empty string field "Type"
			// then
			//    if the struct implements HasXSIType, we take the value from GetType
			//    else we take the name of the struct by reflection
			//    and we set it as the value of "Type"
			if structField.Name == "Type" && f.Kind() == reflect.String && f.String() == "" {
				typeName := original.Type().Name()
				if original.Type().Implements(hasXSIType) {
					v, _ := original.Interface().(HasXSIType)
					typeName = v.GetType()
				}
				copy.Field(i).SetString(typeName)
			} else {
				addXSITypeRecursive(copy.Field(i), original.Field(i))
			}
		}

	// If it is a slice we create a new slice and check each element
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(original.Type(), original.Len(), original.Cap()))
		for i := 0; i < original.Len(); i += 1 {
			addXSITypeRecursive(copy.Index(i), original.Index(i))
		}

	// If it is a map we create a new map and check each value
	case reflect.Map:
		copy.Set(reflect.MakeMap(original.Type()))
		for _, key := range original.MapKeys() {
			originalValue := original.MapIndex(key)
			// New gives us a pointer, but again we want the value
			copyValue := reflect.New(originalValue.Type()).Elem()
			addXSITypeRecursive(copyValue, originalValue)
			copy.SetMapIndex(key, copyValue)
		}

	// Otherwise we cannot traverse anywhere so this finishes the the recursion
	// And everything else will simply be taken from the original
	default:
		copy.Set(original)
	}

}
