package es_proc

import (
	"context"
	"errors"
	"fmt"
	"strconv"
)

func DeleteData(ctx context.Context, esStuInfo *ESStudentInfo) error {
	if esStuInfo == nil || esStuInfo.EsData == nil {
		fmt.Println("student info is nil")
		return errors.New("student info is nil")
	}

	putResp, err := GEsClient.Delete().
		Index(esStuInfo.EsIndex).
		Type(esStuInfo.EsType).
		Id(strconv.Itoa(esStuInfo.EsData.Id)).
		Do(ctx)

	if err != nil {
		fmt.Println("delete fail")
		return errors.New("delete fail" + err.Error())
	}

	fmt.Printf("delete success id=%s index=%s type=%s ver=%d", putResp.Id, putResp.Index, putResp.Type, putResp.Version)
	return nil
}