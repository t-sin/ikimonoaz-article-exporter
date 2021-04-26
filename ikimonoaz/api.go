package ikimonoaz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/xerrors"

	"ikimonoaz-exporter/userdata"
)

const baseURL string = "https://ikimonoaz.ikimonopal.jp/api"

func requestToIkimonoAZ(path string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", baseURL, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, xerrors.New("Status code is not 200")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetUserInfo(userID string) (userdata.User, error) {
	path := fmt.Sprintf("users/get?userId=%s", userID)
	body, err := requestToIkimonoAZ(path)
	if err != nil {
		return userdata.User{}, err
	}

	type Response struct {
		Data userdata.User `json:"data"`
	}
	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return userdata.User{}, err
	}

	return resp.Data, nil
}

func CollectArticles(userID string, page int) ([]userdata.Article, error) {
	path := fmt.Sprintf("articles/list?mypageUserId=%s&page=%d&stripHtml=0", userID, page)
	body, err := requestToIkimonoAZ(path)
	if err != nil {
		return nil, err
	}

	type Result struct {
		Articles []userdata.Article `json:"articles"`
	}
	type Response struct {
		Data Result `json:"data"`
	}
	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Data.Articles, nil
}

func CollectComments(articleID int) ([]userdata.Comment, error) {
	path := fmt.Sprintf("comments/list?articleId=%d&offset=0&isFuture=false", articleID)
	body, err := requestToIkimonoAZ(path)
	if err != nil {
		return nil, err
	}

	type Result struct {
		Comments []userdata.Comment `json:"comments"`
	}
	type Response struct {
		Data Result `json:"data"`
	}
	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp.Data.Comments, nil
}

func CollectAllArticles(userID string) ([]userdata.Article, error) {
	page := 1
	allArticles := make([]userdata.Article, 0)

	for true {
		articles, err := CollectArticles(userID, page)
		if err != nil {
			return nil, err
		}

		if len(articles) == 0 {
			break
		}

		allArticles = append(allArticles, articles...)
		page += 1
	}

	return allArticles, nil
}

func CollectAllUserData(userID string) (userdata.UserData, error) {
	user, err := GetUserInfo(userID)
	if err != nil {
		return userdata.UserData{}, err
	}

	articles, err := CollectAllArticles(userID)
	if err != nil {
		return userdata.UserData{}, err
	}

	for i, article := range articles {
		comments, err := CollectComments(article.ID)
		if err != nil {
			return userdata.UserData{}, err
		}

		articles[i].CommentList = comments
	}

	userdata := userdata.UserData{
		User:     user,
		Articles: articles,
	}

	return userdata, nil
}
