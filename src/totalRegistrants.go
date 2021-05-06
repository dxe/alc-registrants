package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const eventbriteURL = "https://www.eventbriteapi.com/v3/events/149821691713/attendees/"

type EventbriteResponse struct {
	Pagination struct {
		ObjectCount int `json:"object_count"`
	} `json:"pagination"`
}

func getAPI(path string, resp interface{}) error {
	response, err := http.Get(path)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return errors.New("GET request failed. Status: " + strconv.Itoa(response.StatusCode))
	}
	return json.NewDecoder(response.Body).Decode(&resp)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Content-Type": "text/html",
	}

	url := eventbriteURL + "?token=" + os.Getenv("EVENTBRITE_TOKEN")

	var ebRes EventbriteResponse
	err := getAPI(url, &ebRes)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "<p>Error. Please try again.</p>", StatusCode: http.StatusInternalServerError, Headers: headers}, nil
	}

	totalRegistrants := strconv.Itoa(ebRes.Pagination.ObjectCount)

	body := `
<!DOCTYPE html>
<html style="background-color: black;">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>ALC 2021 Registrants</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.2/css/bulma.min.css">
  </head>
  <body>

  <section class="section">

    <div class="columns">
      <div class="column is-one-quarter"></div>
      <div class="column is-half">
        <figure class="image is-square">
          <img src="https://images.squarespace-cdn.com/content/v1/597571c4197aea04d1d21b52/1619561658354-JTEDPD04ROZ21VB9C4YH/ke17ZwdGBToddI8pDm48kNiEM88mrzHRsd1mQ3bxVct7gQa3H78H3Y0txjaiv_0fDoOvxcdMmMKkDsyUqMSsMWxHk725yiiHCCLfrh8O1z4YTzHvnKhyp6Da-NYroOW3ZGjoBKy3azqku80C789l0topjEaZcWjtmMYdCWL4dkGbxs35J-ZjFa9s1e3LsxrX8g4qcOj2k2AL08mW_Htcgg/alc+2021+clear+no+background.png?format=1500w">
        </figure>
      </div>
      <div class="column is-one-quarter"></div>
    </div>

    <div class="columns">
      <div class="column">

        <nav class="level">
          <div class="level-item has-text-centered">
            <div>
              <p class="heading has-text-light">Total Registrants</p>
              <p class="title has-text-light">` + totalRegistrants + `</p>
            </div>
          </div>
        </nav>

      </div>
    </div>

  </section>

  </body>
</html>
`

	return events.APIGatewayProxyResponse{Body: body, StatusCode: http.StatusOK, Headers: headers}, nil
}

func main() {
	lambda.Start(Handler)
}
