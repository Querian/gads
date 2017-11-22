package gads

import (
	"encoding/xml"
)

// DataService struct for service
type DataService struct {
	Auth
}

// CriterionBidLandscape struct for CriterionBidLandscape
type CriterionBidLandscape struct {
	CampaignID      uint64           `xml:"campaignId"`
	AdGroupID       uint64           `xml:"adGroupId"`
	CriterionID     uint64           `xml:"criterionId"`
	StartDate       string           `xml:"startDate"`
	EndDate         string           `xml:"endDate"`
	LandscapePoints []LandscapePoint `xml:"landscapePoints"`
}

// LandscapePoint struct for LandscapePoint
type LandscapePoint struct {
	Bid                           uint64  `xml:"bid>microAmount"`
	Clicks                        uint64  `xml:"clicks"`
	Cost                          uint64  `xml:"cost>microAmount"`
	Impressions                   uint64  `xml:"impressions"`
	PromotedImpressions           uint64  `xml:"promotedImpressions"`
	RequiredBudget                uint64  `xml:"requiredBudget>microAmount"`
	BidModifier                   float64 `xml:"bidModifier"`
	TotalImpressions              uint64  `xml:"totalLocalImpressions"`
	TotalLocalClicks              uint64  `xml:"totalLocalClicks"`
	TotalLocalCost                uint64  `xml:"totalLocalCost>microAmount"`
	TotalLocalPromotedImpressions uint64  `xml:"totalLocalPromotedImpressions"`
}

// NewDataService returns new instance of DataService
func NewDataService(auth *Auth) *DataService {
	return &DataService{Auth: *auth}
}

// GetCriterionBidLandscape returns CriterionBidLandscape
func (s *DataService) GetCriterionBidLandscape(selector Selector) ([]CriterionBidLandscape, int64, error) {
	selector.XMLName = xml.Name{Space: "", Local: "serviceSelector"}
	respBody, err := s.Auth.request(
		dataServiceUrl,
		"getCriterionBidLandscape",
		struct {
			XMLName xml.Name
			Sel     Selector
		}{
			XMLName: xml.Name{Space: baseUrl, Local: "getCriterionBidLandscape"},
			Sel:     selector,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	getResp := struct {
		Size                   int64                   `xml:"rval>totalNumEntries"`
		CriterionBidLandscapes []CriterionBidLandscape `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal(respBody, &getResp)
	if err != nil {
		return nil, 0, err
	}
	return getResp.CriterionBidLandscapes, getResp.Size, err
}

// GetCampaignCriterionBidLandscape returns CriterionBidLandscape
func (s *DataService) GetCampaignCriterionBidLandscape(selector Selector) ([]CriterionBidLandscape, int64, error) {
	selector.XMLName = xml.Name{Space: "", Local: "serviceSelector"}
	respBody, err := s.Auth.request(
		dataServiceUrl,
		"getCampaignCriterionBidLandscape",
		struct {
			XMLName xml.Name
			Sel     Selector
		}{
			XMLName: xml.Name{Space: baseUrl, Local: "getCampaignCriterionBidLandscape"},
			Sel:     selector,
		},
	)
	if err != nil {
		return nil, 0, err
	}
	getResp := struct {
		Size                   int64                   `xml:"rval>totalNumEntries"`
		CriterionBidLandscapes []CriterionBidLandscape `xml:"rval>entries"`
	}{}
	err = xml.Unmarshal(respBody, &getResp)
	if err != nil {
		return nil, 0, err
	}
	return getResp.CriterionBidLandscapes, getResp.Size, err
}
