package es_proc

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
	"reflect"
)

func GetData(ctx context.Context, queryStr string, esStuInfo *ESStudentInfo) ([]*StudentInfo, error) {
	if esStuInfo == nil || queryStr == "" || ctx == nil {
		fmt.Println("student info is nil or query string is empty")
		return nil, errors.New("student info is nil or query string is empty")
	}

	rawQuery := elastic.NewRawStringQuery(queryStr)
	fmt.Println(rawQuery.Source())

	// Query 接口函数默认带有 query 查询，故rawQuery放查询子句
	getResult, err := GEsClient.Search(esStuInfo.EsIndex).
		Type(esStuInfo.EsType).
		Query(rawQuery).
		From(0).		// 读取时，索引从0开始
		Size(30).		// 读取的大小是1
		Do(ctx)
	
	if err != nil {
		fmt.Println("get es fail", err.Error())
		return nil, err
	}
	fmt.Printf("get es success %d milliseconds. total: %d num:%d\n", getResult.TookInMillis, getResult.Hits.TotalHits, len(getResult.Each(reflect.TypeOf(ESStudentInfo{}))))

	// 解析读取的数据，读取和写入时的类型保持一致
	stuInfoArray := make([]*StudentInfo, 0)
	for _, item := range getResult.Each(reflect.TypeOf(StudentInfo{})) {
		if ele, ok := item.(StudentInfo); ok {
			stuInfoArray = append(stuInfoArray, &ele)
		}
	}

	return stuInfoArray, nil
}