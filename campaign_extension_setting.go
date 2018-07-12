package gads

import (
	"encoding/xml"
	"fmt"
)

// CampaignExtensionSettingService (v201806)
// Service used to manage extensions at the campaign level.
// The extensions are managed by AdWords using existing feed services,
// including creating and modifying feeds, feed items,
// and campaign feeds for the user.
//
// see https://developers.google.com/adwords/api/docs/reference/v201806/CampaignExtensionSettingService
type CampaignExtensionSettingService struct {
	Auth
}

// CampaignExtensionSetting is used to add or
// modify extensions being served for the specified campaign.
//
// see https://developers.google.com/adwords/api/docs/reference/v201806/CampaignExtensionSettingService.CampaignExtensionSetting
type CampaignExtensionSetting struct {
	CampaignID       int64             `xml:"campaignId,omitempty"`
	ExtensionType    string            `xml:"extensionType"`
	ExtensionSetting *ExtensionSetting `xml:"extensionSetting"`
}

// GetSitelinks is a conveniency methods to get the list of sitelinks
// out of the extension
func (c *CampaignExtensionSetting) GetSitelinks() (sitelinks []SitelinkFeedItem) {
	if c.ExtensionSetting == nil {
		return
	}
	for _, extension := range c.ExtensionSetting.Extensions {
		sitelink, ok := extension.(SitelinkFeedItem)
		if !ok {
			continue
		}
		sitelinks = append(sitelinks, sitelink)
	}
	return
}

type ExtensionFeedItems []ExtensionFeedItem

func (ex *ExtensionFeedItems) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	feedItemType, err := findAttr(start.Attr, xml.Name{
		Space: "http://www.w3.org/2001/XMLSchema-instance", Local: "type"})
	if err != nil {
		return err
	}
	switch feedItemType {
	case "SitelinkFeedItem":
		slfi := SitelinkFeedItem{}
		err := dec.DecodeElement(&slfi, &start)
		if err != nil {
			return err
		}
		slfi.Type = "SitelinkFeedItem"
		*ex = append(*ex, slfi)
	default:
		if StrictMode {
			return fmt.Errorf("unknown feed item type -> %#v", feedItemType)
		}
	}
	return nil
}

// ExtensionSetting specifies when and which extensions should serve
// at a given level (customer, campaign, or ad group).
//
// see https://developers.google.com/adwords/api/docs/reference/v201806/CampaignExtensionSettingService.ExtensionSetting
type ExtensionSetting struct {
	Extensions            ExtensionFeedItems `xml:"extensions"`
	PlateformRestrictions string             `xml:"platformRestrictions"`
}

// CampaignExtensionSettingOperations is a conveniency map on CampaignExtensionSetting
// to manipulate
// it can have the 3 following keys (case sensitive)
//
// ADD
// SET
// REMOVE
type CampaignExtensionSettingOperations map[string][]CampaignExtensionSetting

// NewCampaignExtensionSettingService is a constructor for CampaignExtensionSettingService
func NewCampaignExtensionSettingService(auth *Auth) *CampaignExtensionSettingService {
	return &CampaignExtensionSettingService{Auth: *auth}
}

// Get returns an array of CampaignExtensionSettings' and
// the total number of CampaignExtensionSettings' matching the selector.
//
// see https://developers.google.com/adwords/api/docs/reference/v201806/CampaignExtensionSettingService#get
func (s *CampaignExtensionSettingService) Get(selector Selector) (
	extensionSettings []CampaignExtensionSetting,
	totalCount int64,
	err error,
) {
	selector.XMLName = xml.Name{"", "selector"}
	respBody, err := s.Auth.request(
		campaignExtensionSettingServiceUrl,
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
		return extensionSettings, totalCount, err
	}
	getResp := struct {
		Size              int64                      `xml:"rval>totalNumEntries"`
		ExtensionSettings []CampaignExtensionSetting `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return extensionSettings, totalCount, err
	}
	return getResp.ExtensionSettings, getResp.Size, err

}

// Mutate allows you to add, modify and remove CampaignExtensionSetting, returning the
// modified ones.
//
// see https://developers.google.com/adwords/api/docs/reference/v201806/CampaignExtensionSettingService#mutate
func (s *CampaignExtensionSettingService) Mutate(
	campaignExtensionSettingOperations CampaignExtensionSettingOperations,
) (campaignExtensionSettings []CampaignExtensionSetting, err error) {
	type operation struct {
		Action                   string                   `xml:"operator"`
		CampaignExtensionSetting CampaignExtensionSetting `xml:"operand"`
	}
	operations := []operation{}
	for action, campaignExtensionSettings := range campaignExtensionSettingOperations {
		for _, campaignExtensionSetting := range campaignExtensionSettings {
			operations = append(
				operations,
				operation{
					Action: action,
					CampaignExtensionSetting: campaignExtensionSetting,
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
	respBody, err := s.Auth.request(
		campaignExtensionSettingServiceUrl,
		"mutate",
		mutation,
	)
	if err != nil {
		return campaignExtensionSettings, err
	}
	mutateResp := struct {
		BaseResponse
		CampaignExtensionSettings []CampaignExtensionSetting `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return campaignExtensionSettings, err
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.CampaignExtensionSettings, err
}

// Query allows to use AWQL to Get CampaignExtensionSettings matching
// the query
//
// see https://developers.google.com/adwords/api/docs/reference/v201806/CampaignExtensionSettingService#query
func (s *CampaignExtensionSettingService) Query(query string) (campaignExtensionSettings []CampaignExtensionSetting, totalCount int64, err error) {

	respBody, err := s.Auth.request(
		adGroupServiceUrl,
		"query",
		AWQLQuery{
			XMLName: xml.Name{
				Space: baseUrl,
				Local: "query",
			},
			Query: query,
		},
	)

	if err != nil {
		return campaignExtensionSettings, totalCount, err
	}
	getResp := struct {
		Size                      int64                      `xml:"rval>totalNumEntries"`
		CampaignExtensionSettings []CampaignExtensionSetting `xml:"rval>entries"`
	}{}
	fmt.Printf("%s\n", respBody)
	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return campaignExtensionSettings, totalCount, err
	}
	return getResp.CampaignExtensionSettings, getResp.Size, err

}
