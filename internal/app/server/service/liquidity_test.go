package service

import "testing"

func TestLiquidityService_QueryPool(t *testing.T) {
	ret, err := new(LiquidityService).QueryPool("ubnb")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("Pass")
	t.Log(ret)
}
