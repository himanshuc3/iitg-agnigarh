package callbacks

// TODO: Benchmark between adding 1000 elements after initiating
// and appending an element to vector

import (
	"fmt"
	"net/url"
	"github.com/gocolly/colly/v2"
	"net/http"
	"github.com/iitg-agnigarh/utils"
	"github.com/spf13/cobra"
)

func scrapeMagicNumber(){
	c := colly.NewCollector()

	c.OnHTML("input[name='magic']", func(e *colly.HTMLElement){
		magicNumber := e.Attr("value")
		loginToFirewall(magicNumber)
	})

	c.OnError(func(_ *colly.Response, err error) {
		
		fmt.Println(err)
	})

	c.Visit(utils.LoginUrl)

}

func loginToFirewall(magicNumber string){

	data := url.Values{
		"username": {"c.himanshu"},
		"password": {"amphion&Z3thus"},
		"magic": {magicNumber},
		"4Tredir": {utils.LoginUrl}}


	response, err := http.PostForm(utils.BaseUrl, data)

	if err != nil {
		fmt.Println("Failed to login: ", err)
		return
	}

	// Defer until return is hit
	defer response.Body.Close()

	fmt.Println(response)
	return

}

func Connect(cmd *cobra.Command, args []string){
	fmt.Println("Connecting to ", utils.BaseUrl)
	scrapeMagicNumber()
}