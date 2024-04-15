package external

import "net/http"

func GetExternalResource(link string) (*http.Response, error) {
	client := http.DefaultClient
	request, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
