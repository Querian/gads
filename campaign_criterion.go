package gads

import (
	//  "strings"
	//  "strconv"
	"encoding/xml"
)

type CampaignCriterion struct {
	Type        string    `xml:"xsi:type,attr,omitempty"`
	CampaignId  int64     `xml:"campaignId"`
	Criterion   Criterion `xml:"criterion"`
	BidModifier *float64  `xml:"bidModifier,omitempty"`
	IsNegative  bool      `xml:"isNegative"`
	Errors      []error   `xml:"-"`
}

func (cc *CampaignCriterion) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		if err != nil {
			return err
		}
		switch start := token.(type) {
		case xml.StartElement:
			switch start.Name.Local {
			case "campaignId":
				if err := dec.DecodeElement(&cc.CampaignId, &start); err != nil {
					return err
				}
			case "criterion":
				criterion, err := criterionUnmarshalXML(dec, start)
				if err != nil {
					return err
				}
				cc.Criterion = criterion
			case "bidModifier":
				if err := dec.DecodeElement(&cc.BidModifier, &start); err != nil {
					return err
				}
			case "isNegative":
				if err := dec.DecodeElement(&cc.IsNegative, &start); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (cc CampaignCriterion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(
		start.Attr,
		xml.Attr{
			xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"},
			"CampaignCriterion",
		},
	)
	e.EncodeToken(start)
	e.EncodeElement(&cc.CampaignId, xml.StartElement{Name: xml.Name{"", "campaignId"}})
	e.EncodeElement(&cc.IsNegative, xml.StartElement{Name: xml.Name{"", "isNegative"}})
	e.EncodeElement(
		&cc.Criterion,
		xml.StartElement{xml.Name{"", "criterion"}, []xml.Attr{}},
	)

	if cc.BidModifier != nil {
		e.EncodeElement(&cc.BidModifier, xml.StartElement{Name: xml.Name{"", "bidModifier"}})
	}

	e.EncodeToken(start.End())
	return nil
}

type NegativeCampaignCriterion struct {
	Type        string    `xml:"xsi:type,attr,omitempty"`
	CampaignId  int64     `xml:"campaignId"`
	IsNegative  bool      `xml:"isNegative"`
	Criterion   Criterion `xml:"criterion"`
	BidModifier float64   `xml:"bidModifier,omitempty"`
	Errors      []error   `xml:"-"`
}

func (cc *NegativeCampaignCriterion) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for token, err := dec.Token(); err == nil; token, err = dec.Token() {
		if err != nil {
			return err
		}
		switch start := token.(type) {
		case xml.StartElement:
			switch start.Name.Local {
			case "campaignId":
				if err := dec.DecodeElement(&cc.CampaignId, &start); err != nil {
					return err
				}
			case "criterion":
				criterion, err := criterionUnmarshalXML(dec, start)
				if err != nil {
					return err
				}
				cc.Criterion = criterion
			case "bidModifier":
				if err := dec.DecodeElement(&cc.BidModifier, &start); err != nil {
					return err
				}
			case "isNegative":
				if err := dec.DecodeElement(&cc.IsNegative, &start); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (ncc NegativeCampaignCriterion) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	isNegative := true
	//fmt.Printf("processing -> %#v\n",ncc)
	start.Attr = append(
		start.Attr,
		xml.Attr{
			xml.Name{"http://www.w3.org/2001/XMLSchema-instance", "type"},
			"NegativeCampaignCriterion",
		},
	)
	e.EncodeToken(start)
	e.EncodeElement(&ncc.CampaignId, xml.StartElement{Name: xml.Name{"", "campaignId"}})
	e.EncodeElement(&isNegative, xml.StartElement{Name: xml.Name{"", "isNegative"}})
	e.EncodeElement(
		&ncc.Criterion,
		xml.StartElement{xml.Name{"", "criterion"}, []xml.Attr{}},
	)
	e.EncodeToken(start.End())
	return nil
}

type CampaignCriterions []interface{}
type CampaignCriterionOperations map[string]CampaignCriterions

func (ccs *CampaignCriterions) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	campaignCriterionType, _ := findAttr(start.Attr, xml.Name{
		Space: "http://www.w3.org/2001/XMLSchema-instance", Local: "type"})

	switch campaignCriterionType {
	case "NegativeCampaignCriterion":
		cc := NegativeCampaignCriterion{}
		err := dec.DecodeElement(&cc, &start)
		if err != nil {
			return err
		}
		*ccs = append(*ccs, cc)
	default:
		cc := CampaignCriterion{}
		err := dec.DecodeElement(&cc, &start)
		if err != nil {
			return err
		}
		*ccs = append(*ccs, cc)
	}
	return nil
}

/*
func NewNegativeCampaignCriterion(campaignId int64, bidModifier float64, criterion interface{}) CampaignCriterion {
  return CampaignCriterion{
    CampaignId: campaignId,
    Criterion: criterion,
    BidModifier: bidModifier
  }
  switch c := criterion.(type) {
  case AdScheduleCriterion:
  case AgeRangeCriterion:
  case ContentLabelCriterion:
  case GenderCriterion:
  case KeywordCriterion:
  case LanguageCriterion:
  case LocationCriterion:
  case MobileAppCategoryCriterion:
  case MobileApplicationCriterion:
  case MobileDeviceCriterion:
  case OperatingSystemVersionCriterion:
  case PlacementCriterion:
  case PlatformCriterion:
  case ProductCriterion:
  case ProximityCriterion:
  case UserInterestCriterion:
    cc.Criterion = criterion
  case UserListCriterion:
    cc.Criterion = criterion
  case VerticalCriterion:
  }
}
*/

type CampaignCriterionService struct {
	Auth
}

func NewCampaignCriterionService(auth *Auth) *CampaignCriterionService {
	return &CampaignCriterionService{Auth: *auth}
}

func (s *CampaignCriterionService) Get(selector Selector) (campaignCriterions CampaignCriterions, totalCount int64, err error) {
	selector.XMLName = xml.Name{"", "serviceSelector"}
	respBody, err := s.Auth.request(
		campaignCriterionServiceUrl,
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
		return campaignCriterions, totalCount, err
	}
	getResp := struct {
		Size               int64              `xml:"rval>totalNumEntries"`
		CampaignCriterions CampaignCriterions `xml:"rval>entries"`
	}{}

	err = xml.Unmarshal([]byte(respBody), &getResp)
	if err != nil {
		return campaignCriterions, totalCount, err
	}
	return getResp.CampaignCriterions, getResp.Size, err
}

func (s *CampaignCriterionService) Mutate(campaignCriterionOperations CampaignCriterionOperations) (campaignCriterions CampaignCriterions, err error) {
	type campaignCriterionOperation struct {
		Action            string      `xml:"operator"`
		CampaignCriterion interface{} `xml:"operand"`
	}
	operations := []campaignCriterionOperation{}
	for action, campaignCriterions := range campaignCriterionOperations {
		for _, campaignCriterion := range campaignCriterions {
			operations = append(operations,
				campaignCriterionOperation{
					Action:            action,
					CampaignCriterion: campaignCriterion,
				},
			)
		}
	}
	mutation := struct {
		XMLName xml.Name
		Ops     []campaignCriterionOperation `xml:"operations"`
	}{
		XMLName: xml.Name{
			Space: baseUrl,
			Local: "mutate",
		},
		Ops: operations,
	}
	respBody, err := s.Auth.request(campaignCriterionServiceUrl, "mutate", mutation)
	if err != nil {
		/*
			    switch t := err.(type) {
			    case *ErrorsType:
				    for action, campaignCriterions := range campaignCriterionOperations {
					    for _, campaignCriterion := range campaignCriterions {
			          campaignCriterions = append(campaignCriterions,campaignCriterion)
			        }
			      }
			      for _, aef := range t.ApiExceptionFaults {
			        for _,e := range aef.Errors {
			          switch et := e.(type) {
			          case CriterionError:
			            offset, err := strconv.ParseInt(strings.Trim(et.FieldPath,"abcdefghijklmnop.]["),10,64)
			            if err != nil {
			              return CampaignCriterions{}, err
			            }
			            cc := campaignCriterions[offset]
			            switch c := cc.(type) {
			            case CampaignCriterion:
			              CampaignCriterion(campaignCriterions[offset]).Errors = append(campaignCriterions[offset].(CampaignCriterion).Errors,fmt.Errorf(et.Reason))
			            case NegativeCampaignCriterion:
			              NegativeCampaignCriterion(campaignCriterions[offset]).Errors = append(NegativeCampaignCriterion(campaignCriterions[offset].Errors),fmt.Errorf(et.Reason))
			            }
			          }
			        }
			      }
			    default:
		*/
		return campaignCriterions, err
		//}
	}
	mutateResp := struct {
		BaseResponse
		CampaignCriterions CampaignCriterions `xml:"rval>value"`
	}{}
	err = xml.Unmarshal([]byte(respBody), &mutateResp)
	if err != nil {
		return campaignCriterions, err
	}

	if len(mutateResp.PartialFailureErrors) > 0 {
		err = mutateResp.PartialFailureErrors
	}

	return mutateResp.CampaignCriterions, err
}

func (s *CampaignCriterionService) Query(query string) (campaignCriterions CampaignCriterions, err error) {
	return campaignCriterions, ERROR_NOT_YET_IMPLEMENTED
}
