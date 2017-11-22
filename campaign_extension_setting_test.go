package gads

import (
	"strconv"
	"testing"
)

func testCampaignExtensionService(t *testing.T) (service *CampaignExtensionSettingService) {
	return &CampaignExtensionSettingService{Auth: testAuthSetup(t)}
}

func TestCampaignExtension(t *testing.T) {
	campaign, cleanupCampaign := testCampaign(t)
	defer cleanupCampaign()

	ces := testCampaignExtensionService(t)

	campaignExtensionSettings, err := ces.Mutate(
		CampaignExtensionSettingOperations{
			"ADD": {
				CampaignExtensionSetting{
					CampaignID:    campaign.Id,
					ExtensionType: "SITELINK",
					ExtensionSetting: &ExtensionSetting{
						Extensions: []ExtensionFeedItem{
							SitelinkFeedItem{
								CommonExtensionFeedItem: &CommonExtensionFeedItem{
									Type: "SitelinkFeedItem",
								},
								Text:      "hello",
								FinalUrls: &UrlList{Urls: []string{"http://example.com"}},
							},
						},
						PlateformRestrictions: "NONE",
					},
				},
			},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, err = ces.Mutate(CampaignExtensionSettingOperations{"REMOVE": campaignExtensionSettings})
		if err != nil {
			t.Error(err)
		}
	}()

	_, _, err = ces.Get(
		Selector{
			Fields: []string{
				"Extensions",
			},
			Predicates: []Predicate{
				{
					"ExtensionType",
					"EQUALS",
					[]string{"SITELINK"},
				},
				{
					"CampaignId",
					"EQUALS",
					[]string{strconv.FormatInt(campaign.Id, 10)},
				},
			},
			Ordering: []OrderBy{},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
}
