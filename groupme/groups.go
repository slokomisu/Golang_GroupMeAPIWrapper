package groupme

import (
	"github.com/levigross/grequests"
	"strconv"
	"errors"
)

type groupsService service

type Group struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url"`
	CreatorUserID string `json:"creator_user_id"`
	CreatedAt     int    `json:"created_at"`
	UpdatedAt     int    `json:"updated_at"`
	Members []Member `json:"members"`
	ShareURL string `json:"share_url"`
	Messages struct {
		Count                int    `json:"count"`
		LastMessageID        string `json:"last_message_id"`
		LastMessageCreatedAt int    `json:"last_message_created_at"`
		Preview struct {
			Nickname string `json:"nickname"`
			Text     string `json:"text"`
			ImageURL string `json:"image_url"`
			Attachments []struct {
				Type        string  `json:"type"`
				URL         string  `json:"url,omitempty"`
				Lat         string  `json:"lat,omitempty"`
				Lng         string  `json:"lng,omitempty"`
				Name        string  `json:"name,omitempty"`
				Token       string  `json:"token,omitempty"`
				Placeholder string  `json:"placeholder,omitempty"`
				Charmap     [][]int `json:"charmap,omitempty"`
			} `json:"attachments"`
		} `json:"preview"`
	} `json:"messages"`
}

type IndexOptions struct {
	Page    int
	PerPage int
	Omit    string
}

type (
	groupArrayResponse struct {
		Response []Group `json:"response"`
	}

	groupSingleResponse struct {
		Response *Group `json:"response"`
	}

	GroupParams struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ImageUrl    string `json:"image_url"`
		Share       bool   `json:"share"`
		OfficeMode  bool   `json:"office_mode"`
	}

	ChangeOwnerRequest struct {
		GroupId string `json:"group_id"`
		OwnerId string `json:"owner_id"`
	}

	ChangeOwnerResponse struct {
		Results []struct {
			GroupId string `json:"group_id"`
			OwnerId string `json:"owner_id"`
			Status  string `json:"status"`
		} `json:"results"`
	}
)

// Index gets all of the authenticated user's active groups
// Response paginated with a default of 10 groups per page
func (s *groupsService) Index(params *IndexOptions) ([]Group, error) {

	// Default parameters
	if params == nil {
		params = &IndexOptions{
			Page:    1,
			PerPage: 10,
			Omit:    "",
		}
	}

	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"page":     strconv.Itoa(params.Page),
			"per_page": strconv.Itoa(params.PerPage),
			"omit":     params.Omit,
			"token":    s.token,
		},
	}

	var response groupArrayResponse

	resp, requestErr := grequests.Get(s.baseUrl+"/groups", ro)
	if requestErr != nil {
		return nil, requestErr
	}

	if resp.StatusCode == 401 {
		return nil, errors.New("401 unauthorized")
	}

	jsonErr := resp.JSON(&response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return response.Response, nil
}

// Load a specific group by group id
func (s *groupsService) Show(id int) (*Group, error) {



	var response groupSingleResponse

	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"token": s.token,
		},
	}

	resp, requestErr := grequests.Get(s.baseUrl+"/groups/"+strconv.Itoa(id), ro)
	if requestErr != nil {
		return nil, requestErr
	}

	if resp.StatusCode == 401 {
		return nil, errors.New("401 unauthorized")
	}

	jsonErr := resp.JSON(&response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return response.Response, nil

}

// Shows all the groups of the authenticated user that they have left and can rejoin
func (s *groupsService) Former() ([]Group, error) {
	var response groupArrayResponse

	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"token": s.token,
		},
	}

	resp, requestErr := grequests.Get(s.baseUrl+"/groups/former", ro)
	if requestErr != nil {
		return nil, requestErr
	}

	if resp.StatusCode == 401 {
		return nil, errors.New("401 unauthorized")
	}

	jsonErr := resp.JSON(&response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return response.Response, nil

}

// Creates a new group
func (s *groupsService) Create(createOptions *GroupParams) (*Group, error) {
	var response groupSingleResponse

	ro := &grequests.RequestOptions{
		JSON: createOptions,
		Params: map[string]string{
			"token": s.token,
		},
	}

	resp, responseErr := grequests.Post(s.baseUrl+"/groups", ro)
	if responseErr != nil {
		return nil, responseErr
	}

	jsonErr := resp.JSON(&response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return response.Response, nil
}

// Updates a group that the authorized user has access to
func (s *groupsService) Update(id int, updateOptions *GroupParams) (*Group, error) {
	var response groupSingleResponse

	ro := &grequests.RequestOptions{
		JSON: updateOptions,
		Params: map[string]string{
			"token": s.token,
		},
	}

	resp, responseErr := grequests.Post(s.baseUrl+"/groups/"+strconv.Itoa(id)+"/update", ro)
	if responseErr != nil {
		return nil, responseErr
	}

	jsonErr := resp.JSON(&response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return response.Response, nil
}


// Deletes a group the authorized user is the owner of
func (s *groupsService) Destroy(id int) error {
	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"token": s.token,
		},
	}

	resp, responseErr := grequests.Post(s.baseUrl+"/groups/"+strconv.Itoa(id)+"/destroy", ro)
	if responseErr != nil {
		return responseErr
	}

	if resp.StatusCode == 200 {
		return nil
	}

	return nil
}

// Joins a group with a given id and share token
func (s *groupsService) Join(id int, shareToken string) (*Group, error) {
	var response groupSingleResponse

	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"token": s.token,
		},
	}

	resp, responseErr := grequests.Post(s.baseUrl+"/groups/"+strconv.Itoa(id)+"/join/"+shareToken, ro)
	if responseErr != nil {
		return nil, responseErr
	}

	jsonErr := resp.JSON(&response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return response.Response, nil
}

// Rejoin a group that the authorized user can rejoin
func (s *groupsService) Rejoin(id int) (*Group, error) {
	var response groupSingleResponse

	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"token":    s.token,
			"group_id": strconv.Itoa(id),
		},
	}

	resp, responseErr := grequests.Post(s.baseUrl+"/groups/join/", ro)
	if responseErr != nil {
		return nil, responseErr
	}

	jsonErr := resp.JSON(&response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return response.Response, nil
}

// Change the owners of the group
func (s *groupsService) ChangeOwners(changeRequests []ChangeOwnerRequest) (*ChangeOwnerResponse, error) {
	var response ChangeOwnerResponse

	ro := &grequests.RequestOptions{
		JSON: changeRequests,
		Params: map[string]string{
			"token": s.token,
		},
	}

	resp, responseErr := grequests.Get(s.baseUrl+"/groups/change_owners", ro)
	if responseErr != nil {
		return nil, responseErr
	}
	jsonErr := resp.JSON(&response)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &response, nil
}
