package types

import (
	"github.com/indaco/teseo/opengraph"
	"github.com/indaco/teseo/schemaorg"
	"github.com/indaco/teseo/twittercard"
)

type SEOItems struct {
	WebPage        *schemaorg.WebPage
	SiteNavElement *schemaorg.SiteNavigationElement
	Breadcrumb     *schemaorg.BreadcrumbList
	Profile        *opengraph.Profile
	TwitterCard    *twittercard.TwitterCard
}
