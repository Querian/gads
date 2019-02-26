package gads

import "encoding/xml"

// ManagedCustomerService represents the api that handle links between accounts
type ManagedCustomerService struct {
	Auth
}

// LinkStatus represents the state of the link
type LinkStatus string

const (
	LinkStatusActive    LinkStatus = "ACTIVE"
	LinkStatusInactive  LinkStatus = "INACTIVE"
	LinkStatusPending   LinkStatus = "PENDING"
	LinkStatusRefused   LinkStatus = "REFUSED"
	LinkStatusCancelled LinkStatus = "CANCELLED"
	LinkStatusUnkwown   LinkStatus = "UNKNOWN"
)

// NewManagedCustomerService is the ManagedCustomerService constructor
func NewManagedCustomerService(auth *Auth) *ManagedCustomerService {
	return &ManagedCustomerService{Auth: *auth}
}

// ManagedCustomer represents the accounts handled by the current clientCustomerID
// https://developers.google.com/adwords/api/docs/reference/v201809/ManagedCustomerService.ManagedCustomer
type ManagedCustomer struct {
	Name                 string `xml:"name,omitempty"`
	CustomerID           uint   `xml:"customerId,omitempty"`
	CanManageClients     bool   `xml:"canManageClients,omitempty"`
	CurrencyCode         string `xml:"currencyCode,omitempty"`
	TestAccount          bool   `xml:"testAccount,omitempty"`
	DateTimeZone         string `xml:"dateTimeZone,omitempty"`
	ExcludeHiddenAccount bool   `xml:"excludeHiddenAccounts,omitempty"`
}

// ManagedCustomerLinkOperations are used when you change links between mcc and classic adwords account
type ManagedCustomerLinkOperations map[string][]*ManagedCustomerLink

// ManagedCustomerLink represents the status of the link between the current
// account and the account it manages
// https://developers.google.com/adwords/api/docs/reference/v201809/ManagedCustomerService.ManagedCustomerLink
type ManagedCustomerLink struct {
	ManagerCustomerID      uint   `xml:"managerCustomerId,omitempty"`
	ClientCustomerId       uint   `xml:"clientCustomerId,omitempty"`
	LinkStatus             string `xml:"linkStatus,omitempty"`
	PendingDescriptiveName string `xml:"pendingDescriptiveName,omitempty"`
	Hidden                 bool   `xml:"isHidden,omitempty"`
}

// ManagedCustomerLinkResult is the response of the MutateLink service
// @see https://developers.google.com/adwords/api/docs/reference/v201809/ManagedCustomerService.MutateLinkResults
type ManagedCustomerLinkResult struct {
	Links []*ManagedCustomerLink `xml:"links,omitempty"`
}

// ManagedCustomerMoveOperations are performed when you move an account inside the hierarchy
type ManagedCustomerMoveOperations map[string][]ManagedCustomerMoveOperation

// ManagedCustomerMoveOperation is one operation in the move actions
type ManagedCustomerMoveOperation struct {
	OldManagerCustomerId uint
	Link                 ManagedCustomerLink
}

// Get fetches the managed customers of the current account
func (m *ManagedCustomerService) Get(selector Selector) (
	customers []ManagedCustomer,
	managedCustomerLinks []ManagedCustomerLink,
	totalCount int64,
	err error,
) {
	selector.XMLName = xml.Name{"", "serviceSelector"}
	var respBody []byte
	respBody, err = m.Auth.request(
		managedCustomerServiceUrl,
		"get",
		struct {
			XMLName xml.Name
			Sel     Selector
		}{
			XMLName: xml.Name{
				Space: managedCustomerUrl,
				Local: "get",
			},
			Sel: selector,
		},
	)

	if err != nil {
		return
	}

	getResp := struct {
		Size                int64                 `xml:"rval>totalNumEntries"`
		ManagedCustomers    []ManagedCustomer     `xml:"rval>entries"`
		ManagedCustomerLink []ManagedCustomerLink `xml:"rval>links"`
	}{}

	err = xml.Unmarshal(respBody, &getResp)
	if err != nil {
		return
	}

	totalCount = getResp.Size

	return getResp.ManagedCustomers, getResp.ManagedCustomerLink, totalCount, err
}

// MutateManager takes a budgetOperations and creates, modifies or destroys the associated budgets.
func (m *ManagedCustomerService) MutateManager(mcmOps ManagedCustomerMoveOperations) (links []ManagedCustomerLink, err error) {
	type managedCustomerMoveOperation struct {
		Action               string              `xml:"https://adwords.google.com/api/adwords/cm/v201809 operator"`
		Link                 ManagedCustomerLink `xml:"operand"`
		OldManagerCustomerId uint                `xml:"oldManagerCustomerId"`
	}

	operations := []managedCustomerMoveOperation{}
	for action, ops := range mcmOps {
		for _, op := range ops {
			operations = append(
				operations,
				managedCustomerMoveOperation{
					Action:               action,
					Link:                 op.Link,
					OldManagerCustomerId: op.OldManagerCustomerId,
				},
			)
		}
	}
	respBody, err := m.Auth.request(
		managedCustomerServiceUrl,
		"mutateManager",
		struct {
			XMLName xml.Name
			Ops     []managedCustomerMoveOperation `xml:"operations"`
		}{
			XMLName: xml.Name{
				Space: managedCustomerUrl,
				Local: "mutateManager",
			},
			Ops: operations,
		},
	)
	if err != nil {
		return links, err
	}
	mutateResp := struct {
		BaseResponse
		Links []ManagedCustomerLink `xml:"rval>value"`
	}{}
	err = xml.Unmarshal(respBody, &mutateResp)
	if err != nil {
		return links, err
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.Links, err
}

// MutateLink changes the links between mcc and classic adwords account
func (m *ManagedCustomerService) MutateLink(mcl ManagedCustomerLinkOperations) ([]*ManagedCustomerLink, error) {

	type linkOperation struct {
		Action string               `xml:"https://adwords.google.com/api/adwords/cm/v201809 operator"`
		Link   *ManagedCustomerLink `xml:"operand"`
	}

	operations := []*linkOperation{}
	for action, ops := range mcl {
		for _, op := range ops {
			operations = append(
				operations,
				&linkOperation{
					Action: action,
					Link:   op,
				},
			)
		}
	}
	respBody, err := m.Auth.request(
		managedCustomerServiceUrl,
		"mutateLink",
		struct {
			XMLName xml.Name
			Ops     []*linkOperation `xml:"operations"`
		}{
			XMLName: xml.Name{
				Space: managedCustomerUrl,
				Local: "mutateLink",
			},
			Ops: operations,
		},
	)
	if err != nil {
		return nil, err
	}
	mutateResp := struct {
		BaseResponse
		Links []*ManagedCustomerLink `xml:"rval>value>links"`
	}{}
	err = xml.Unmarshal(respBody, &mutateResp)
	if err != nil {
		return nil, err
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.Links, err
}
