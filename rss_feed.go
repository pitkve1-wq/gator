package main
import "net/http"
import "io"
import "encoding/xml"
import "html"
import "context"
import "fmt"

func handlerAgg(s *state, cmd command) error {
feed, err := fetchFeed(context.TODO(),"https://www.wagslane.dev/index.xml")
if err != nil {
	return err
}
fmt.Printf("%+v",feed)
return nil
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}
type RSSItem struct {
	Title string		`xml:"title"`
	Link string			`xml:"link"`
	Description string	`xml:"description"`
	PubDate string		`xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL,nil)
	if err != nil {
		return nil , err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	c, err2 := client.Do(req)
	if err2 != nil {
		return nil, err2
	}
	defer c.Body.Close()
	data, err3 := io.ReadAll(c.Body)
	if err3 != nil {
		return nil, err3
	}
	var ab RSSFeed
	err4 := xml.Unmarshal(data, &ab)
	if err4 != nil {
		return nil, err4
	}
	for i, item := range ab.Channel.Item {
	item.Title = html.UnescapeString(item.Title)
    item.Description = html.UnescapeString(item.Description)
	ab.Channel.Item[i] = item
	}
	ab.Channel.Title = html.UnescapeString(ab.Channel.Title)
	ab.Channel.Description = html.UnescapeString(ab.Channel.Description)
	return &ab, nil
}