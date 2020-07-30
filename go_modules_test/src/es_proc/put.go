package es_proc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

func PutData(ctx context.Context, esStuInfo *ESStudentInfo) error {
	if esStuInfo == nil || esStuInfo.EsData == nil {
		fmt.Println("student info is nil")
		return errors.New("student info is nil")
	}

	putResp, err := GEsClient.Index().
		Index(esStuInfo.EsIndex).
		Type(esStuInfo.EsType).
		Id(strconv.Itoa(esStuInfo.EsData.Id)).
		BodyJson(esStuInfo.EsData).
		Do(ctx)

	if err != nil {
		fmt.Println("put fail")
		return errors.New("put fail" + err.Error())
	}

	fmt.Printf("put success id=%s index=%s type=%s ver=%d\n", putResp.Id, putResp.Index, putResp.Type, putResp.Version)
	return nil
}