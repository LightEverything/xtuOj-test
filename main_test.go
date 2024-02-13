package main

import (
	"context"
	"fmt"
	"testing"
	"time"
	"xtuOj/helper"
	"xtuOj/models"
)

var ctx = context.Background()

func TestSetValue(t *testing.T) {
	models.RDB.Set("wanxinnb@outlook.com", "123", time.Second*100)
}

func TestGetValue(t *testing.T) {
	v, err := models.RDB.Get("name").Result()

	if err != nil {
		t.Fail()
		t.Log(err)
	}
	fmt.Println(v)
}

func TestGetToken(t *testing.T) {
	st, err := helper.GetToken("2223", "123123", 1)
	if err != nil {
		t.Fail()
		t.Log(err)
		return
	}
	t.Log("successful : " + st)
}

func TestAnlayToken(t *testing.T) {
	data, err := helper.AnalyseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiMjIyMjIzIiwiaWRlbnRpdHkiOiIxMjMxMjMiLCJpc19hZG1pbiI6MH0.GH22knmNHAQnOI_RhSvrxNSKmFjJkJqxDrXPhZPIwSI")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	print(data.Is_admin)
}

func TestRandomCode(t *testing.T) {
	fmt.Println(helper.GetRandCode())
	fmt.Println(helper.GetRandCode())
	fmt.Println(helper.GetRandCode())
	fmt.Println(helper.GetRandCode())
	fmt.Println(helper.GetRandCode())
}
