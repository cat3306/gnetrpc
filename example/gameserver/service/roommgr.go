package service

import (
	"errors"

	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/common"
	"github.com/cat3306/gnetrpc/example/gameserver/util"
	"github.com/cat3306/gnetrpc/protocol"
	rpcutil "github.com/cat3306/gnetrpc/util"
	"github.com/panjf2000/gnet/v2"
)

var (
	roomMgr *RoomMgr
)

type RoomMgr struct {
	rooms map[string]*Room
}

func (r *RoomMgr) Init(v ...interface{}) gnetrpc.IService {
	r.rooms = make(map[string]*Room)
	roomMgr = r
	return r
}
func (r *RoomMgr) Alias() string {
	return ""
}

type CreateRoomReq struct {
	Pwd       string `json:"pwd"`
	MaxNum    int    `json:"max_num"` //
	JoinState bool   `json:"join_state"`
}
type CreateRoomRsp struct {
	Id string `json:"id"`
}

func (r *RoomMgr) AddRoom(room *Room) {
	r.rooms[room.id] = room
}
func (r *RoomMgr) DelRoom(id string) {
	delete(r.rooms, id)
}
func (r *RoomMgr) GetRoom(id string) (*Room, bool) {
	room, ok := r.rooms[id]
	return room, ok
}

// {"pwd":"123","max_num":10,"join_state":true}
func (r *RoomMgr) Create(ctx *protocol.Context, req *CreateRoomReq, rsp *ApiRsp) *gnetrpc.CallMode {
	// _, exists := ctx.Conn.GetProperty(RoomIdKey)
	// if exists {
	// 	rsp.Err("already in room")
	// 	return gnetrpc.CallSelf()
	// }
	id := util.GenId(6)
	room := &Room{
		maxNum:     req.MaxNum,
		pwd:        req.Pwd,
		joinState:  req.JoinState,
		gameState:  false,
		scene:      0,
		id:         id,
		connMatrix: common.NewConnMatrix(false),
	}
	room.connMatrix.Add(ctx.Conn)
	r.AddRoom(room)
	rsp.Ok(id)
	//ctx.Conn.SetProperty(RoomIdKey, id)
	return gnetrpc.CallBroadcast()
}

func (r *RoomMgr) ConnOnClose(conn gnet.Conn) error {
	roomId := ""
	room, ok := r.GetRoom(roomId)
	if !ok {
		return errors.New("not found room")
	}

	room.connMatrix.Remove("")
	if room.connMatrix.Len() == 0 {
		r.DelRoom(room.id)
	} else {
		ctx := protocol.GetCtx()
		ctx.H.Fill(protocol.Json)
		ctx.ServiceMethod = "RoomMgr"
		ctx.ServicePath = "Leave"
		room.connMatrix.Broadcast(protocol.Encode(ctx, new(ApiRsp).Ok("哈哈")))
	}
	return nil
}
func (r *RoomMgr) Leave(ctx *protocol.Context, req *struct{}, rsp *ApiRsp) *gnetrpc.CallMode {
	//roomId, exists := ctx.Conn.GetProperty(RoomIdKey)
	// if !exists {
	// 	rsp.Err("not join room yet")
	// 	return gnetrpc.CallSelf()
	// }
	roomId := ""
	room, ok := r.GetRoom(roomId)
	if !ok {
		rsp.Err("not found room")
		return gnetrpc.CallSelf()
	}
	room.connMatrix.Remove("")
	if room.connMatrix.Len() == 0 {
		r.DelRoom(room.id)
	} else {
		room.connMatrix.Broadcast(protocol.Encode(ctx, rsp.Ok("leave")))
	}
	//ctx.Conn.DelProperty(RoomIdKey)
	rsp.Ok(nil)
	return gnetrpc.CallNone()
}

type RoomInfosRsp struct {
	Id     string `json:"id"`
	MaxNum int    `json:"max_num"` //人数
	Pwd    string `json:"pwd"`     //密码
	Cnt    int    `json:"cnt"`
}

func (r *RoomMgr) RoomsInfo(ctx *protocol.Context, req *struct{}, rsp *ApiRsp) *gnetrpc.CallMode {
	list := make([]RoomInfosRsp, 0)
	for _, v := range r.rooms {
		list = append(list, RoomInfosRsp{
			Id:     v.id,
			MaxNum: v.maxNum,
			Pwd:    v.pwd,
			Cnt:    v.connMatrix.Len(),
		})
	}
	rsp.Ok(list)
	return gnetrpc.CallSelf()
}

func (r *RoomMgr) Join(ctx *protocol.Context, id *string, rsp *ApiRsp) *gnetrpc.CallMode {
	ctx.H.SerializeType = byte(protocol.Json)
	// _, exists := ctx.Conn.GetProperty(RoomIdKey)
	// if exists {
	// 	rsp.Err("already in room")
	// 	return gnetrpc.CallSelf()
	// }
	room, ok := r.GetRoom(*id)
	if !ok {
		return gnetrpc.CallSelf()
	}
	room.connMatrix.Add(ctx.Conn)
	//ctx.Conn.SetProperty(RoomIdKey, *id)
	rsp.Ok(nil)
	room.connMatrix.Broadcast(protocol.Encode(ctx, rsp))
	return gnetrpc.CallNone()
}

func (r *RoomMgr) Chat(ctx *protocol.Context, txt *string, rsp *ApiRsp) *gnetrpc.CallMode {
	ctx.H.SerializeType = byte(protocol.Json)
	// v, ok := ctx.Conn.GetProperty(RoomIdKey)
	// if !ok {
	// 	rsp.Err("not join room")
	// 	return gnetrpc.CallNone()
	// }
	roomId := ""
	room, ok := r.GetRoom(roomId)
	if !ok {
		rsp.Err("not found room")
		return gnetrpc.CallSelf()
	}
	rsp.Ok(txt)
	room.connMatrix.BroadcastExceptOne(protocol.Encode(ctx, rsp), rpcutil.GetConnId(ctx.Conn))
	return gnetrpc.CallNone()
}

func (r *RoomMgr) StartGame(ctx *protocol.Context, id *string, rsp *ApiRsp) *gnetrpc.CallMode {
	return gnetrpc.CallNone()
}
