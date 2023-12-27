package service

import (
	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/example/gameserver/util"
	"github.com/cat3306/gnetrpc/protocol"
)

type RoomMgr struct {
	rooms map[string]*Room
}

func (r *RoomMgr) Init(v ...interface{}) gnetrpc.IService {
	r.rooms = make(map[string]*Room)
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
func (r *RoomMgr) CreateRoom(ctx *protocol.Context, req *CreateRoomReq, rsp *ApiRsp) *gnetrpc.CallMode {
	_, exists := ctx.Conn.GetProperty(RoomIdKey)
	if exists {
		rsp.Err("already create room")
		return gnetrpc.CallSelf()
	}
	id := util.GenId(6)
	room := &Room{
		maxNum:     req.MaxNum,
		pwd:        req.Pwd,
		joinState:  req.JoinState,
		gameState:  false,
		scene:      0,
		id:         id,
		connMatrix: gnetrpc.NewConnMatrix(false),
	}
	room.connMatrix.Add(ctx.Conn)
	r.AddRoom(room)
	rsp.Ok(id)
	ctx.Conn.SetProperty(RoomIdKey, id)
	return gnetrpc.CallSelf()
}

func (r *RoomMgr) LeaveRoom(ctx *protocol.Context, req *struct{}, rsp *ApiRsp) *gnetrpc.CallMode {
	roomId, exists := ctx.Conn.GetProperty(RoomIdKey)
	if !exists {
		rsp.Err("not join room yet")
		return gnetrpc.CallSelf()
	}
	room, ok := r.GetRoom(roomId.(string))
	if !ok {
		rsp.Err("not found room")
		return gnetrpc.CallSelf()
	}
	room.connMatrix.Remove(ctx.Conn.Id())
	if room.connMatrix.Len() == 0 {
		r.DelRoom(room.id)
	}
	ctx.Conn.DelProperty(RoomIdKey)
	rsp.Ok(nil)
	return gnetrpc.CallSelf()
}

type RoomInfosRsp struct {
	Id     string `json:"id"`
	MaxNum int    `json:"max_num"` //人数
	Pwd    string `json:"pwd"`     //密码
}

func (r *RoomMgr) RoomInfos(ctx *protocol.Context, req *struct{}, rsp *ApiRsp) *gnetrpc.CallMode {
	list := make([]RoomInfosRsp, 0)
	for _, v := range r.rooms {
		list = append(list, RoomInfosRsp{
			Id:     v.id,
			MaxNum: v.maxNum,
			Pwd:    v.pwd,
		})
	}
	rsp.Ok(list)
	return gnetrpc.CallSelf()
}
