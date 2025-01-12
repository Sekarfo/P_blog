package articles

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Sekarfo/P_blog/models"
)

const (
	apiKey  = "ec79b26e36ca4f3ba4f320ceec94a351"
	baseURL = "https://newsapi.org/v2/everything"
)

type ArticleService interface {
	GetArticles(params models.ArticleSearch) ([]models.Article, int, error)
}

type articleSearcher struct {
	httpCli *http.Client
}

func NewArticleGetter() ArticleService {
	return &articleSearcher{httpCli: http.DefaultClient}
}

func (a *articleSearcher) GetArticles(
	params models.ArticleSearch,
) ([]models.Article, int, error) {
	toReq := toNewsAPIReq{
		ArticleSearch: params,
		ApiKey:        apiKey,
	}
	url := baseURL + toReq.toQueryParams()
	fmt.Println("URL", url)
	resp, err := a.httpCli.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	var newsResponse newsResponse
	err = json.Unmarshal(body, &newsResponse)
	if err != nil {
		return nil, 0, err
	}

	return newsResponse.Articles, newsResponse.TotalResults, nil
}
