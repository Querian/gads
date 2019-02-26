package gads

// CommonExtensionFeedItem is the Parent type for
// all ExtensionFeedItem
type CommonExtensionFeedItem struct {
	Type       string  `xml:"xsi:type,attr,omitempty"`
	FeedID     int64   `xml:"feedId,omitempty"`
	FeedItemID int64   `xml:"feedItemId,omitempty"`
	Status     string  `xml:"status,omitempty"`
	FeedType   string  `xml:"feedType,omitempty"`
	StartTime  string  `xml:"startTime,omitempty"`
	EndTime    *string `xml:"endTime"`
}

// ExtensionFeedItem is an interface all extension feed item implements
type ExtensionFeedItem interface {
	GetFeedID() int64
	GetFeedItemID() int64
	GetType() string
}

// GetFeedID returns the ID of the feed this item belongs to
func (e *CommonExtensionFeedItem) GetFeedID() int64 {
	return e.FeedID
}

// GetFeedItemID returns the ID of this feed item
func (e *CommonExtensionFeedItem) GetFeedItemID() int64 {
	return e.FeedItemID
}

// GetType returns the type
func (e *CommonExtensionFeedItem) GetType() string {
	return e.Type
}

// SitelinkFeedItem represents a sitelink extension.
//
// see https://developers.google.com/adwords/api/docs/reference/v201809/CampaignExtensionSettingService.SitelinkFeedItem
type SitelinkFeedItem struct {
	*CommonExtensionFeedItem
	Text                string   `xml:"sitelinkText"`
	URL                 *string  `xml:"sitelinkUrl"`
	Line2               *string  `xml:"sitelinkLine2"`
	Line3               *string  `xml:"sitelinkLine3"`
	FinalUrls           *UrlList `xml:"sitelinkFinalUrls"`
	TrackingURLTemplate *string  `xml:"sitelinkTrackingUrlTemplate"`
}
