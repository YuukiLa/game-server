package cmd_router

import "github/com/yuuki80code/game-server/ws"

func initTest() {
	ws.AddHandler("1", func(c *ws.Context) {
		c.Client.RoomID = "123"
		err := c.Client.EnterRoom()
		if err != nil {
			c.SendString(err.Error())
			return
		}
		c.SendString("您进入了123房间")
	})
	ws.AddHandler("2", func(c *ws.Context) {
		c.Client.UnEnterRoom()
		c.SendString("您退出了123房间")
	})
	ws.AddHandler("3", func(c *ws.Context) {
		if c.Client.RoomID == "" {
			c.SendString("您当前不在任何房间")
			return
		}
		roomBoradcast := ws.RoomBroadcast{
			RoomID: c.Client.RoomID,
			//Data:   c.Data,
			ClientID: c.Client.ID,
		}
		c.Client.RoomBroadcastAll(roomBoradcast)
	})
	//ws.AddHandler("4", func(c *ws.Context) {
	//	if c.Client.RoomID == "" {
	//		c.SendString("您当前不在任何房间")
	//		return
	//	}
	//	roomBoradcast := ws.RoomBroadcast{
	//		RoomID: c.Client.RoomID,
	//		Data:   c.Data,
	//		ClientID: c.Client.ID,
	//	}
	//	c.Client.RoomBroadcast(roomBoradcast)
	//})

}
