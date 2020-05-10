package handler

import "github/com/yuuki80code/game-server/ws"

func Draw(c *ws.Context) {
	roomBoradcast := ws.RoomBroadcast{
		RoomID: c.Client.RoomID,
		Data: ws.Result{
			CMD:  c.CMD,
			Data: c.Data,
			Msg:  "",
		},
		ClientID: c.Client.ID,
	}
	c.Client.RoomBroadcast(roomBoradcast)
}

func Answer(c *ws.Context) {
	roomBoradcast := ws.RoomBroadcast{
		RoomID: c.Client.RoomID,
		Data: ws.Result{
			CMD:  c.CMD,
			Data: c.Data,
			Msg:  "",
		},
		ClientID: c.Client.ID,
	}
	c.Client.RoomBroadcast(roomBoradcast)
}
