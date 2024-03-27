package irishub

//
//import (
//	"fmt"
//	"testing"
//
//	sdktype "github.com/irisnet/irishub-sdk-go/types"
//	"github.com/stretchr/testify/require"
//)
//
//func TestTxSearch(t *testing.T) {
//	Connect("tcp://192.168.150.60:26657", "192.168.150.60:29090", "test", "4iris")
//
//	to := "iaa1r04vmlycmhw6yw8sp4jka42w3uqf0w409nmdq3"
//	height := int64(1261950)
//
//	builder := sdktype.NewEventQueryBuilder().AddCondition(
//		sdktype.NewCond("transfer", "recipient").EQ(to),
//	).AddCondition(
//		sdktype.NewCond("message", "module").EQ("bank"),
//	).AddCondition(
//		sdktype.NewCond("tx", "height").EQ(height),
//	)
//	//require.Equal(t, "C827FA1E1ACBF007ECF19513F09E7C93F182A38658879D8EE5EFE0A81BE4039A", txs[0].Hash)
//}

// func TestMarshalJSON(t *testing.T) {
// 	Connect("tcp://localhost:26657", "localhost:9090", "test", "4iris")
// 	msg := bank.MsgSend{
// 		FromAddress: "iaa1ekm8qfqcl54z5l4pm9d4v7th72vd5qfu5k2642",
// 		ToAddress:   "iaa1ekm8qfqcl54z5l4pm9d4v7th72vd5qfu5k2642",
// 		Amount:      sdktype.NewCoins(sdktype.NewCoin("uiris", sdktype.NewInt(100))),
// 	}

// 	//msgBz, err := codec.ProtoMarshalJSON(&msg)
// 	msgBz, err := MarshalJSON(msg)
// 	require.NoError(t, err)
// 	//msgBz, err := irishub.EncodingConfig().Marshaler.MarshalJSON(o proto.Message)
// 	var msg1 sdktype.Msg
// 	err = UnmarshalJSON(msgBz, &msg1)
// 	require.NoError(t, err)
// 	fmt.Println(string(msgBz))
// }
