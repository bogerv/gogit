package infoquery

import (
	"context"
	"encoding/json"
	"gitshell/utilpkg/es"
	"gopkg.in/olivere/elastic.v6"
	"log"
)

type InfoQuery struct {
	Id              string
	Index           string
	Query           string
	Language        string
	From            int
	Size            int
	PhrasePrefix    bool // true:按词拆分搜索 false:按字拆分进行搜索
	SearchFields    []string
	HighlightFields []string
}

func New(index string) *InfoQuery {
	query := &InfoQuery{
		Index: index,
	}
	if query.Size == 0 {
		query.Size = 50
	}
	return query
}

func (slf *InfoQuery) SetQuery(param string) *InfoQuery {
	slf.Query = param
	return slf
}

func (slf *InfoQuery) SetLanguage(param string) *InfoQuery {
	slf.Language = param
	return slf
}

func (slf *InfoQuery) SetFrom(param int) *InfoQuery {
	if param == 0 {
		return slf
	}
	slf.From = param
	return slf
}

func (slf *InfoQuery) SetSize(param int) *InfoQuery {
	if param == 0 {
		return slf
	}
	slf.Size = param
	return slf
}

func (slf *InfoQuery) SetSearchFields(param []string) *InfoQuery {
	slf.SearchFields = param
	return slf
}

func (slf *InfoQuery) SetHighlightFields(param []string) *InfoQuery {
	slf.HighlightFields = param
	return slf
}

func (slf *InfoQuery) SetPhrasePrefix(param bool) *InfoQuery {
	slf.PhrasePrefix = param
	return slf
}

func (slf *InfoQuery) CreateIndex(ctx context.Context) error {
	res, err := es.ES().IndexExists(slf.Index).Do(ctx)
	if !res {
		_, err = es.ES().CreateIndex(slf.Index).Body(InfoMapping).Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (slf *InfoQuery) Create(ctx context.Context, body interface{}) error {
	_, err := es.ES().Index().
		Index(slf.Index).
		Type("_doc").
		Id(slf.Id).
		BodyJson(body).
		Do(ctx)
	return err
}

func (slf *InfoQuery) Update(ctx context.Context, body map[string]interface{}) error {
	_, err := es.ES().Update().
		Index(slf.Index).
		Type("_doc").
		Id(slf.Id).
		Doc(body).
		DetectNoop(true).
		Do(ctx)
	return err
}

//condition := make(map[string]interface{})
//	condition["uuid"] = uuid
//	infoes := infoquery.New("info")
//	if err := infoes.UpdateByQuery(context.Background(), "ctx._source.isDel=1", condition); err != nil {
//		log.Printf("[es] delete err: %s; uuid: %s\n", err.Error(), uuid)
//	}
func (slf *InfoQuery) UpdateByQuery(ctx context.Context, scriptStr string, condition map[string]interface{}) error {
	conditions := make([]elastic.Query, 0)
	for key, value := range condition {
		conditions = append(conditions, elastic.NewMatchQuery(key, value))
	}
	boolQuery := elastic.NewBoolQuery().Must(conditions...)
	PrintQueryDSL(boolQuery)
	script := elastic.NewScript(scriptStr)
	_, err := es.ES().UpdateByQuery().
		Index(slf.Index).
		Query(boolQuery).
		Script(script).
		Do(ctx)
	return err
}

func (slf *InfoQuery) MultiMatchQuery(ctx context.Context) (*elastic.SearchResult, error) {
	searchResult, err := slf.buildCondition().
		Sort("checkTime", false).      // sort by "user" field, ascending
		From(slf.From).Size(slf.Size). // take documents 0-9
		//Pretty(true).       // pretty print request and response JSON
		Do(ctx)

	return searchResult, err
}

func (slf *InfoQuery) buildCondition() *elastic.SearchService {
	var multiMatchQuery *elastic.MultiMatchQuery
	if slf.PhrasePrefix {
		multiMatchQuery = elastic.NewMultiMatchQuery(slf.Query, slf.SearchFields...).Type("phrase_prefix")
	} else {
		multiMatchQuery = elastic.NewMultiMatchQuery(slf.Query, slf.SearchFields...)
	}
	boolQuery := elastic.NewBoolQuery().
		Must(elastic.NewMatchQuery("isOnline", 1)).
		MustNot(elastic.NewMatchQuery("isDel", 1), elastic.NewMatchQuery("checkTime", 0))
	if len(slf.Language) > 0 {
		boolQuery = boolQuery.Must(elastic.NewMatchQuery("language", slf.Language))
	}
	query := elastic.NewBoolQuery().Must(multiMatchQuery, boolQuery)

	highlighterFields := make([]*elastic.HighlighterField, 0)
	if len(slf.HighlightFields) > 0 {
		for _, field := range slf.HighlightFields {
			highlight := elastic.NewHighlighterField(field).PreTags("<tag1>").PostTags("</tag1>")
			highlighterFields = append(highlighterFields, highlight)
		}
	}
	highLight := elastic.NewHighlight().Fields(highlighterFields...)

	return es.ES().
		Search().
		Index(slf.Index).
		Query(query).
		Highlight(highLight)
}

func PrintQueryDSL(query *elastic.BoolQuery) {
	src, _ := query.Source()
	data, err := json.MarshalIndent(src, "", "  ")
	log.Println("PrintSource data:", string(data), err)
}
