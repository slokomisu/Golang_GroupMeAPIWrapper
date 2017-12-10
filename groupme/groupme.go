package groupme

type service struct {
	token   string
	baseUrl string
}

type Client struct {
	service *service
	Groups  *groupsService
	Members *membersService
}

func NewClient(accessToken string) *Client {
	service := &service{
		token:   accessToken,
		baseUrl: "https://api.groupme.com/v3",
	}

	client := &Client{
		service: service,
	}

	client.Groups = (*groupsService)(service)
	client.Members = (*membersService)(service)

	return client
}
