package feed

import (
	"context"
	"net/url"
	"time"

	"github.com/mmcdole/gofeed"

	"MrRSS/backend"
)

type RssOriginContentsInfo struct {
	Title string
	Image string
	Item  gofeed.Item
}

func getRssContent(rssLinks []string) []backend.FeedContentsInfo {
	var rssOriginContent []RssOriginContentsInfo

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	rssParser := gofeed.NewParser()

	for _, link := range rssLinks {
		rss, err := rssParser.ParseURLWithContext(link, ctx)
		if err != nil {
			continue
		}

		for _, item := range rss.Items {
			// Get the URL of the rss
			u, err := url.Parse(rss.Link)
			if err != nil {
				panic(err)
			}

			// Get the image URL for the rss
			imageUrl := "https://www.google.com/s2/favicons?sz=16&domain=" + u.Host
			if rss.Image != nil {
				imageUrl = rss.Image.URL
			}

			// Append the item to the rssOriginContent
			rssOriginContent = append(rssOriginContent, RssOriginContentsInfo{
				Title: rss.Title,
				Image: imageUrl,
				Item:  *item,
			})
		}
	}

	var result []backend.FeedContentsInfo

	for _, item := range rssOriginContent {
		// Get the image URL
		imageURL := ""
		filterImageUrl := filterImage(item.Item.Content)
		if item.Item.Image != nil {
			imageURL = item.Item.Image.URL
		} else if filterImageUrl != nil {
			imageURL = *filterImageUrl
		}

		// Get the time since the item was published
		timeSinceStr := getTimeSince(item.Item.PublishedParsed)

		rssContentItem := backend.FeedContentsInfo{
			FeedTitle: item.Title,
			FeedImage: item.Image,
			Title:     item.Item.Title,
			Link:      item.Item.Link,
			TimeSince: timeSinceStr,
			Time:      item.Item.PublishedParsed.Format("2006-01-02 15:04"),
			Image:     imageURL,
			Content:   item.Item.Content,
			Readed:    false,
		}

		// Append the item to the rssContent
		result = append(result, rssContentItem)
	}

	return result
}
