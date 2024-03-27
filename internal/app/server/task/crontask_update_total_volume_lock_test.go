package task

import "testing"

func TestSyncUpdateTotalVolumnLockTask_calculateTotalVolumeLock(t *testing.T) {
	res, err := new(SyncUpdateTotalVolumnLockTask).calculateTotalVolumeLock()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(res)
}

func TestSyncUpdateTotalVolumnLockTask_caculateFarmLPAmt(t *testing.T) {
	res, err := new(SyncUpdateTotalVolumnLockTask).caculateFarmLPAmt()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(res)
}
