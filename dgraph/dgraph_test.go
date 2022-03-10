package dgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"google.golang.org/grpc"
	"log"
	"testing"
)

const (
	IsPartOf  = "is_part_of"
	LocatedIn = "located_in"
)

var landNum = 2
var buildingNum = 5
var levelNum = 5
var roomNum = 5
var deviceNum = 100
var capNum = 300

type Twin struct {
	Uid          string
	ID           string
	Capabilities map[string]interface{}
	DType        []string
}

type TwinSlice []*Twin

func (ts TwinSlice) Marshal() ([]byte, error) {
	var twins []map[string]interface{}
	for _, twin := range ts {
		kv := make(map[string]interface{})
		kv["uid"] = twin.Uid
		kv["id"] = twin.ID
		kv["dgraph.type"] = twin.DType
		for k, v := range twin.Capabilities {
			kv[k] = v
		}
		twins = append(twins, kv)
	}
	return json.Marshal(twins)
}

type CancelFunc func()

func getDgraphClient() (*dgo.Dgraph, CancelFunc) {
	conn, err := grpc.Dial("192.168.30.58:29080", grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}

	dc := api.NewDgraphClient(conn)
	dg := dgo.NewDgraphClient(dc)
	//ctx := context.Background()

	// Perform login call. If the Dgraph cluster does not have ACL and
	// enterprise features enabled, this call should be skipped.
	//for {
	//	// Keep retrying until we succeed or receive a non-retriable error.
	//	err = dg.Login(ctx, "groot", "password")
	//	if err == nil || !strings.Contains(err.Error(), "Please retry") {
	//		break
	//	}
	//	time.Sleep(time.Second)
	//}
	if err != nil {
		log.Fatalf("While trying to login %v", err.Error())
	}

	return dg, func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing connection:%v", err)
		}
	}
}

func TestDropAll(t *testing.T) {
	dg, cancel := getDgraphClient()
	defer cancel()

	if err := dg.Alter(context.Background(), &api.Operation{DropAll: true}); err != nil {
		t.Fatal(err)
	}
}

func TestDropType(t *testing.T) {
	dg, cancel := getDgraphClient()
	defer cancel()

	if err := dg.Alter(context.Background(), &api.Operation{
		DropOp:    api.Operation_TYPE,
		DropValue: "user"}); err != nil {
		t.Fatal(err)
	}
}

func TestAlter(t *testing.T) {
	dg, cancel := getDgraphClient()
	defer cancel()

	op := &api.Operation{}
	op.Schema = `
		id: string @index(hash).
		is_part_of: uid @reverse .
		located_in: uid @reverse .
		type Land {}
		type Building {}
		type Level {}
		type Room {}
	`
	typeDevice := `
		type Device {
	`
	for i := 1; i <= capNum; i++ {
		op.Schema += fmt.Sprintf("spot%d: int .\n", i)
		typeDevice += fmt.Sprintf("spot%d\n", i)
	}
	op.Schema += typeDevice + `}`
	fmt.Println(op.Schema)

	ctx := context.Background()
	if err := dg.Alter(ctx, op); err != nil {
		t.Fatal(err)
	}
}

func TestBatchMutate(t *testing.T) {
	dg, cancel := getDgraphClient()
	defer cancel()

	mu := &api.Mutation{
		CommitNow: true,
	}
	ctx := context.Background()

	var twins TwinSlice
	// land
	for landIndex := 1; landIndex <= landNum; landIndex++ {
		land := &Twin{
			Uid:          fmt.Sprintf("_:land%d", landIndex),
			ID:           fmt.Sprintf("land%d", landIndex),
			Capabilities: nil,
			DType:        []string{"Land"},
		}
		twins = append(twins, land)

		// building
		for buildingIndex := 1; buildingIndex <= buildingNum; buildingIndex++ {
			building := &Twin{
				Uid:          fmt.Sprintf("%s_building%d", land.Uid, buildingIndex),
				ID:           fmt.Sprintf("%s_building%d", land.ID, buildingIndex),
				Capabilities: nil,
				DType:        []string{"Building"},
			}
			twins = append(twins, building)
			mu.Set = append(mu.Set, &api.NQuad{
				Subject:   building.Uid,
				Predicate: IsPartOf,
				ObjectId:  land.Uid,
			})

			// level
			for levelIndex := 1; levelIndex <= levelNum; levelIndex++ {
				level := &Twin{
					Uid:          fmt.Sprintf("%s_level%d", building.Uid, levelIndex),
					ID:           fmt.Sprintf("%s_level%d", building.ID, levelIndex),
					Capabilities: nil,
					DType:        []string{"Level"},
				}
				twins = append(twins, level)
				mu.Set = append(mu.Set, &api.NQuad{
					Subject:   level.Uid,
					Predicate: IsPartOf,
					ObjectId:  building.Uid,
				})

				// room
				for roomIndex := 1; roomIndex <= roomNum; roomIndex++ {
					room := &Twin{
						Uid:          fmt.Sprintf("%s_room%d", level.Uid, roomIndex),
						ID:           fmt.Sprintf("%s_room%d", level.ID, roomIndex),
						Capabilities: nil,
						DType:        []string{"Room"},
					}
					twins = append(twins, room)
					mu.Set = append(mu.Set, &api.NQuad{
						Subject:   room.Uid,
						Predicate: IsPartOf,
						ObjectId:  level.Uid,
					})

					// device
					for deviceIndex := 1; deviceIndex <= deviceNum; deviceIndex++ {
						device := &Twin{
							Uid:          fmt.Sprintf("%s_device%d", room.Uid, deviceIndex),
							ID:           fmt.Sprintf("%s_device%d", room.ID, deviceIndex),
							Capabilities: make(map[string]interface{}),
							DType:        []string{"Device"},
						}
						// Case 1: use capability as predicate
						for capIndex := 1; capIndex <= capNum; capIndex++ {
							device.Capabilities[fmt.Sprintf("spot%d", capIndex)] = 233
						}
						twins = append(twins, device)
						mu.Set = append(mu.Set, &api.NQuad{
							Subject:   device.Uid,
							Predicate: LocatedIn,
							ObjectId:  room.Uid,
						})
					}
				}
			}
		}
	}

	mu.SetJson, _ = twins.Marshal()
	if _, err := dg.NewTxn().Mutate(ctx, mu); err != nil {
		t.Fatal(err)
	}
}

func Test_Query(t *testing.T) {

}
