package ikimonoaz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/xerrors"
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

func GetUserInfo(userID string) (User, error) {
	path := fmt.Sprintf("users/get?userId=%s", userID)
	body, err := requestToIkimonoAZ(path)
	if err != nil {
		return User{}, err
	}

	type Response struct {
		Data User `json:"data"`
	}
	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return User{}, err
	}

	return resp.Data, nil
}

func CollectArticles(userID string, page int) ([]Article, error) {
	path := fmt.Sprintf("articles/list?mypageUserId=%s&page=%d&stripHtml=0", userID, page)
	body, err := requestToIkimonoAZ(path)
	if err != nil {
		return nil, err
	}

	type Result struct {
		Articles []Article `json:"articles"`
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

func CollectComments(articleID string) ([]Comment, error) {
	path := fmt.Sprintf("comments/list?articleId=%s&offset=0&isFuture=false", articleID)
	body, err := requestToIkimonoAZ(path)
	if err != nil {
		return nil, err
	}

	type Result struct {
		Comments []Comment `json:"comments"`
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

func CollectAllArticles(userID string) ([]Article, error) {
	page := 1
	allArticles := make([]Article, 0)

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

func CollectAllUserData(userID string) (UserData, error) {
	user, err := GetUserInfo(userID)
	if err != nil {
		return UserData{}, err
	}

	articles, err := CollectAllArticles(userID)
	if err != nil {
		return UserData{}, err
	}

	for _, article := range articles {
		comments, err := CollectComments(string(article.ID))
		if err != nil {
			return UserData{}, err
		}

		article.CommentList = comments
	}

	userdata := UserData{
		User:     user,
		Articles: articles,
	}

	return userdata, nil
}
