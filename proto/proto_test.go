package example

import (
	"fmt"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"testing"
)

func TestProtoAny(t *testing.T) {
	before := &TestAny{
		Id:      1,
		Title:   "标题",
		Content: "内容",
	}
	any, err := anypb.New(before)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(any) // [type.googleapis.com/TestAny]:{Id:1  Title:"标题"  Content:"内容"}

	resp := &AnyResponse{
		Code: 0,
		Msg:  "success",
		Data: any,
	}

	msgName := resp.Data.MessageName()
	fmt.Println(msgName)
	after := &TestAny{}
	if err := resp.Data.UnmarshalTo(after); err != nil {
		log.Fatal(err)
	}

	fmt.Println(after) // Id:1  Title:"标题"  Content:"内容"
}

func TestProtoOneof(t *testing.T) {
	resp := OneofResponse{
		Result: &OneofResponse_Correct{
			Correct: &Correct{
				Rank: 233,
			}},
	}

	//resp := OneofResponse{
	//	Result: &OneofResponse_Wrong{
	//		Wrong: &Wrong{
	//			Code: 99,
	//			Msg:  "wrong answer",
	//		},
	//	},
	//}

	switch resp.Result.(type) {
	case *OneofResponse_Correct:
		fmt.Println("正确！！！", resp.GetCorrect())
	case *OneofResponse_Wrong:
		fmt.Println("错误:( ", resp.GetWrong())
	}
}
