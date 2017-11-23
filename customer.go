package gads

import "encoding/xml"

// CustomerService fetches or modify Customer properties
type CustomerService struct {
	Auth
}

// Customer represents a customer in the Google Adwords API
type Customer struct {
	ID                  uint64  `xml:"customerId,omitempty"`
	CurrencyCode        *string `xml:"currencyCode,omitempty"`
	DateTimeZone        *string `xml:"dateTimeZone,omitempty"`
	DescriptiveName     string  `xml:"descriptiveName,omitempty"`
	CanManageClients    bool    `xml:"canManageClients,omitempty"`
	TestAccount         bool    `xml:"testAccount,omitempty"`
	AutoTaggingEnabled  *bool   `xml:"autoTaggingEnabled"`
	TrackingURLTemplate *string `xml:"trackingUrlTemplate"`
	// TODO: implements this - not needed for the moment
	//ConversionTrackingSettings *ConversionTrackingSettings `xml:"conversionTrackingSettings,omitempty"`
	//RemarketingSettings        *RemarketingSettings        `xml:"remarketingSettings.omitempty"`
}

// NewCustomerService creates a CustomerService
func NewCustomerService(auth *Auth) *CustomerService {
	return &CustomerService{Auth: *auth}
}

// Get fetches the managed customers of the current account
func (m *CustomerService) GetCustomers(s *Selector) (
	customers []Customer,
	err error,
) {
	if s != nil {
		s.XMLName = xml.Name{"", "serviceSelector"}
	}

	var respBody []byte
	respBody, err = m.request(
		customerServiceUrl,
		"getCustomers",
		struct {
			XMLName xml.Name
			Sel     *Selector
		}{
			XMLName: xml.Name{
				Space: managedCustomerUrl,
				Local: "getCustomers",
			},
			Sel: s,
		},
	)

	if err != nil {
		return
	}

	getResp := struct {
		Customers []Customer `xml:"rval"`
	}{}

	err = xml.Unmarshal(respBody, &getResp)
	if err != nil {
		return
	}

	return getResp.Customers, err
}

// Mutate performs modifications of one or many customer
func (m *CustomerService) Mutate(c Customer) (customer Customer, err error) {

	mutation := struct {
		XMLName  xml.Name
		Customer Customer `xml:"customer"`
	}{
		XMLName: xml.Name{
			Space: managedCustomerUrl,
			Local: "mutate",
		},
		Customer: c,
	}
	respBody, err := m.request(customerServiceUrl, "mutate", mutation)
	if err != nil {
		return customer, err
	}

	mutateResp := struct {
		BaseResponse
		RespCustomer Customer `xml:"rval"`
	}{}

	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return customer, err
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.RespCustomer, err
}

// ServiceLink struct for customer service link
type ServiceLink struct {
	Type   string `xml:"serviceType,omitempty"`
	ID     int64  `xml:"serviceLinkId,omitempty"`
	Status string `xml:"linkStatus,omitempty"`
	Name   string `xml:"name,omitempty"`
}

// GetServiceLinks fetches the service links
func (m *CustomerService) GetServiceLinks(s *Selector) (serviceLinks []ServiceLink, err error) {
	if s != nil {
		s.XMLName = xml.Name{"", "selector"}
	}

	var respBody []byte
	respBody, err = m.request(
		customerServiceUrl,
		"getServiceLinks",
		struct {
			XMLName xml.Name
			Sel     *Selector
		}{
			XMLName: xml.Name{
				Space: managedCustomerUrl,
				Local: "getServiceLinks",
			},
			Sel: s,
		},
	)

	if err != nil {
		return
	}

	getResp := struct {
		ServiceLinks []ServiceLink `xml:"rval"`
	}{}

	err = xml.Unmarshal(respBody, &getResp)
	if err != nil {
		return
	}

	return getResp.ServiceLinks, err
}

type ServiceLinkOperations map[string][]ServiceLink

func (s *CustomerService) MutateServiceLinks(ops ServiceLinkOperations) (links []ServiceLink, err error) {
	type linkOperation struct {
		Action      string      `xml:"https://adwords.google.com/api/adwords/cm/v201710 operator"`
		ServiceLink ServiceLink `xml:"operand"`
	}
	operations := []linkOperation{}
	for action, links := range ops {
		for _, link := range links {
			operations = append(operations, linkOperation{Action: action, ServiceLink: link})
		}
	}
	respBody, err := s.Auth.request(
		customerServiceUrl,
		"mutateServiceLinks",
		struct {
			XMLName xml.Name
			Ops     []linkOperation `xml:"operations"`
		}{
			XMLName: xml.Name{
				Space: managedCustomerUrl,
				Local: "mutateServiceLinks",
			},
			Ops: operations,
		},
	)
	if err != nil {
		return links, err
	}
	mutateResp := struct {
		BaseResponse
		Links []ServiceLink `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return links, err
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.Links, err
}
