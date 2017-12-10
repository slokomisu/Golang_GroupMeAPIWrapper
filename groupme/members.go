package groupme

import (
	"github.com/levigross/grequests"
	"strconv"
	"errors"
)

type membersService service

type Member struct {
	Nickname string `json:"nickname"`
	ID string `json:"id"`
	UserID string `json:"user_id"`
	Muted bool `json:"muted"`
	ImageURL string `json:"image_url"`
	Autokicked bool `json:"autokicked"`
	AppInstalled bool `json:"app_installed"`
	GUID string `json:"guid"`
	PhoneNumber string `json:"phone_number"`
	Email string `json:"email"`
}

func (s *membersService) Add(groupId int, newMembers []Member) (*string, error) {
	type (
		addMemberResponse struct {
			Response struct {
				ResultsId string `json:"results_id"`
			} `json:"response"`
		}

		newMembersRequest struct {
			Members []Member `json:"members"`
		}
	)

	request := &newMembersRequest{
		Members: newMembers,
	}

	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"token": s.token,
		},
		JSON: request,
	}

	resp, responseErr := grequests.Post(s.baseUrl+"/groups/"+strconv.Itoa(groupId)+"/members/add", ro)
	if responseErr != nil {
		return nil, responseErr
	}

	var response addMemberResponse

	jsonErr := resp.JSON(&response)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return &response.Response.ResultsId, nil
}

func (s *membersService) Results(groupId int , resultsId string) ([]Member, error) {
	type (
		membersResults struct {
			Response struct {
				Members []Member `json:"members"`
			} `json:"response"`
		}
	)


	ro := &grequests.RequestOptions{
		Params: map[string]string {
			"token": s.token,
			"results_id": resultsId,
		},
	}

	resp, responseErr := grequests.Get(s.baseUrl+"/groups/"+strconv.Itoa(groupId)+"/"+resultsId, ro)
	if responseErr != nil {
		return nil, responseErr
	}


	if resp.StatusCode == 404 {
		return nil, errors.New("results are no longer available")
	}

	if resp.StatusCode == 503 {
		return nil, errors.New("results aren't ready yet. try again later")
	}

	var response membersResults
	jsonErr := resp.JSON(&response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return response.Response.Members, nil
}

func (s *membersService) Remove(groupId, membershipId int) error {
	ro := &grequests.RequestOptions{
		Params: map[string]string {
			"token": s.token,
		},
	}

	resp, responseErr := grequests.Post(s.baseUrl+"/groups/"+strconv.Itoa(groupId)+"/members/"+strconv.Itoa(membershipId)+"/remove", ro)

	if responseErr != nil {
		return responseErr
	}

	if resp.StatusCode != 200 {
		return errors.New("could not remove member from group")
	}

	return nil

}

func (s *membersService) Update(groupId int, newNickname string) (*Member, error) {
	type (
		updateRequest struct {
			Membership struct{
				Nickname string `json:"nickname"`
			} `json:"membership"`
		}

		memberResults struct {
			Response *Member `json:"response"`
		}
	)

	update := &updateRequest{}

	update.Membership.Nickname = newNickname

	ro := &grequests.RequestOptions{
		Params: map[string]string {
			"token": s.token,
		},
		JSON: update,
	}

	resp, responseErr := grequests.Post(s.baseUrl+"/groups/"+strconv.Itoa(groupId)+"/memberships/update", ro)
	if responseErr != nil {
		return nil, responseErr
	}



	var response memberResults

	jsonErr := resp.JSON(&response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return response.Response, nil
}




