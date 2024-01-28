package opengraphsvc

import (
	"context"
	"net/http"
	"strings"

	"github.com/GDGVIT/opengraph-thumbnail-backend/api/pkg/routes"
	"github.com/PuerkitoBio/goquery"
)

func (svc *OpenGraphSvcImpl) GetMetadata(ctx context.Context, params routes.GetMetadataParams) (routes.Metadata, error) {
	res, err := http.Get(params.Url)
	if err != nil {
		svc.logger.Debugf("error in get request", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		svc.logger.Debugf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		svc.logger.Debugf("error in goquery", err)
	}

	metaData := make(map[string]string)

	// Parse Open Graph and Twitter data
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		property, _ := s.Attr("property")
		name, _ := s.Attr("name")
		content, _ := s.Attr("content")

		if strings.HasPrefix(property, "og:") {
			metaData[property] = content
		}
		if strings.HasPrefix(name, "twitter:") {
			metaData[name] = content
		}
	})

	var response routes.Metadata
	response.Title = metaData["og:title"]
	if response.Title == "" {
		// get title from <title> tag
		response.Title = doc.Find("title").Text()
	}
	response.Description = metaData["og:description"]
	if response.Description == "" {
		// get description from <meta name="description"> tag
		response.Description = doc.Find("meta[name=description]").AttrOr("content", "")
	}
	response.Image = metaData["og:image"]
	response.Url = params.Url

	return response, nil
}
