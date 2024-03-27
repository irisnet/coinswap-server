package task

import (
	ctx "context"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	"github.com/pkg/errors"
	"github.com/qiniu/qmgo"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestAddLocalFarmToWhiteList(t *testing.T) {

	Convey("set farm in db and make it  invisible", t, func() {

		var farmIndb = model.Farms{
			Id:      primitive.NewObjectID(),
			PoolId:  "farm_in_db",
			Visible: false,
			//IsLocalFarm: true,
		}

		defer func() {
			err := model.NewCollection(model.Farms{}).Remove(ctx.Background(), bson.M{"pool_id": "farm_in_db"})
			if err != nil {
				t.Fatal("clean farm in db err,please clean manual", err)
			}
		}()

		_, err := model.NewCollection(model.Farms{}).InsertOne(ctx.Background(), farmIndb)
		if err != nil {
			t.Fatal("set farm in db, occurs  err", err)
		}
		Convey("add the farm to  white list", func() {

			var whiteList = model.WhiteList{PoolId: "farm_in_db",
				PoolName: "farm_in_db_name",
				OrderBy:  1,
			}
			defer func() {
				err := model.NewCollection(model.WhiteList{}).Remove(ctx.Background(), bson.M{"pool_id": "farm_in_db"})
				if err != nil {
					t.Fatal("clean white list err,please clean manual", err)
				}
			}()
			_, err := model.NewCollection(model.WhiteList{}).InsertOne(ctx.Background(), whiteList)
			if err != nil {
				t.Fatal("set white list occurs  err", err)
			}

			Convey("execute sync farm pool task", func() {

				So(NewFarmPoolSyncer().Execute(), ShouldBeNil)

				So(func() bool {
					var farm model.Farms
					err := model.NewCollection(model.Farms{}).Find(ctx.Background(), bson.M{"pool_id": "farm_in_db"}).One(&farm)
					if err != nil {
						t.Fatal("get farm in db to check result,occurs error  err", err)
					}

					return farm.Visible
				}(), ShouldBeTrue)

			})

		})

	})

}

func TestRemoveLocalFarmFromWhiteList(t *testing.T) {

	Convey("set farm in db and make it visible", t, func() {

		var farmIndb = model.Farms{
			Id:      primitive.NewObjectID(),
			PoolId:  "farm_in_db",
			Visible: true,
			//IsLocalFarm: true,
		}

		defer func() {
			err := model.NewCollection(model.Farms{}).Remove(ctx.Background(), bson.M{"pool_id": "farm_in_db"})
			if err != nil {
				t.Fatal("clean farm in db err,please clean manual", err)
			}
		}()

		_, err := model.NewCollection(model.Farms{}).InsertOne(ctx.Background(), farmIndb)
		if err != nil {
			t.Fatal("set farm in db, occurs  err", err)
		}

		Convey("remove the farm from white list", func() {

			err := model.NewCollection(model.WhiteList{}).Remove(ctx.Background(), bson.M{"pool_id": "farm_in_db"})

			if err != nil && !errors.Is(err, qmgo.ErrNoSuchDocuments) {
				t.Fatal("remove farm from white list, occurs  err", err)
			}

			Convey("execute sync farm pool task", func() {

				So(NewFarmPoolSyncer().Execute(), ShouldBeNil)

				So(func() bool {
					var farm model.Farms
					err := model.NewCollection(model.Farms{}).Find(ctx.Background(), bson.M{"pool_id": "farm_in_db"}).One(&farm)
					if err != nil {
						t.Fatal("get farm in db to check result,occurs error  err", err)
					}

					return farm.Visible
				}(), ShouldBeFalse)

			})

		})

	})
}

//func TestRemoveChainFarmFromWhiteList(t *testing.T) {
//
//
//	irishub.Farm().FarmPools= func(ctx ctx.Context, req *farm.QueryFarmPoolsRequest) (*farm.QueryFarmPoolsResponse, error) {
//
//		return nil,nil
//	}
//
//		Convey("set farm in db",t, func() {
//
//			farmPoolResp, err := irishub.Farm().FarmPools(ctx.Background(), &farm.QueryFarmPoolsRequest{})
//			if err != nil {
//				t.Fatal("get farm from chain ,occurs err", err)
//			}
//
//			if len(farmPoolResp.Pools) == 0 {
//				t.Fatal("no farm in chain, can not  test", err)
//				return
//			}
//			//farmFromChain:=model.FarmFromFarmPoolEntry(farmPoolResp.Pools[0])
//
//			Convey("set white list")
//			defer func() {
//				err := model.NewCollection(model.Farms{}).Remove(ctx.Background(), bson.M{"pool_id": "farm_in_db"})
//				if err != nil {
//					t.Fatal("clean farm in db err,please clean manual", err)
//				}
//			}()
//
//			_, err := model.NewCollection(model.Farms{}).InsertOne(ctx.Background(), farmIndb)
//			if err != nil {
//				t.Fatal("set farm in db, occurs  err", err)
//			}
//
//			Convey("execute sync farm pool task", func() {
//
//				So(NewFarmPoolSyncer().Execute(), ShouldBeNil)
//
//				So(func() bool {
//					var farm model.Farms
//					err := model.NewCollection(model.Farms{}).Find(ctx.Background(), bson.M{"pool_id": "farm_in_db"}).One(&farm)
//					if err != nil {
//						t.Fatal("get farm in db to check result,occurs error  err", err)
//					}
//
//					return farm.Visible
//				}(), ShouldBeFalse)
//
//			})
//
//		})
//
//
//}
