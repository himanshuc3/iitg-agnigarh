package callbacks

// TODO: Benchmark between adding 1000 elements after initiating
// and appending an element to vector

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gocolly/colly"
	"github.com/iitg-agnigarh/utils"
	"github.com/spf13/cobra"
)

const MAX_ELAPSED_TIME = 2 * time.Minute

func scrapeDomain() {

	scraper := colly.NewCollector()
	backoffHandler := utils.NewExponentialBackoff()

	scraping := func() {

		scraper.OnHTML("input[name='magic']", func(e *colly.HTMLElement) {
			magicNumber := e.Attr("value")
			loginToFirewall(magicNumber)
		})

		scraper.OnError(func(r *colly.Response, err error) {
			fmt.Println("Failed attempt")
			fmt.Println(err)
			fmt.Println("Before sleep")
			time.Sleep(backoffHandler.NextBackoff())
			fmt.Println("After sleep")
			r.Request.Retry()
		})

		scraper.Visit(utils.LoginUrl)

	}
	scraping()
}

func loginToFirewall(magicNumber string) {

	data := url.Values{
		"username": {"c.himanshu"},
		"password": {"amphion&Z3thus"},
		"magic":    {magicNumber},
		"4Tredir":  {utils.LoginUrl}}

	response, err := http.PostForm(utils.BaseUrl, data)

	if err != nil {
		fmt.Println("Failed to login: ", err)
		return
	}

	// Defer until return is hit
	defer response.Body.Close()

	fmt.Println(response)
	fmt.Println("Successfully connected to domain")
}

func Connect(cmd *cobra.Command, args []string) {
	fmt.Println("Connecting to ", utils.BaseUrl)

	scrapeDomain()
}
