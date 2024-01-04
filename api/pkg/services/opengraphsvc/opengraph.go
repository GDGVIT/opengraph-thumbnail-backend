package opengraphsvc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/GDGVIT/opengraph-thumbnail-backend/api/pkg/routes"
	"github.com/PuerkitoBio/goquery"
)

func (svc *OpenGraphSvcImpl) OpenGraphEditor(c context.Context, params routes.OpenGraphParams) (string, error) {
	// Get the metadata
	metaData, err := getMetadata(params.Url, params.Title, params.Description, params.Image)
	if err != nil {
		log.Fatal(err)
	}

	// Generate a temporary HTML page with the metadata
	html := generateTemporaryHTML(metaData, params.Url)

	return html, nil
}

// Helper functions
func getMetadata(url string, customTitle *string, customDescription *string, customImage *string) (map[string]string, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
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

	// Modify the title with custom title
	if customTitle != nil {
		metaData["og:title"] = *customTitle
		metaData["twitter:title"] = *customTitle
	}

	// Modify the description with custom description if provided
	if customDescription != nil {
		metaData["og:description"] = *customDescription
		metaData["twitter:description"] = *customDescription
	}

	if customImage != nil {
		metaData["og:image"] = *customImage
		metaData["twitter:image"] = *customImage
	}

	return metaData, nil
}

func generateTemporaryHTML(metaData map[string]string, originalURL string) string {
	var builder strings.Builder

	builder.WriteString("<html><head>")
	for key, value := range metaData {
		builder.WriteString(fmt.Sprintf("<meta name=\"%s\" content=\"%s\">", key, value))
	}
	builder.WriteString("</head><body>")
	builder.WriteString("<p></p>")
	builder.WriteString("<script>setTimeout(function() { window.location.href = '" + originalURL + "'; }, 200);</script>")
	builder.WriteString("</body></html>")

	return builder.String()
}
