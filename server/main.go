package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
)

type (
	GetUserReq struct {
		ID string `json:"id"`
	}

	GetUserResp struct {
		Id    string `json:"id"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
)

type User struct {
	Id    string
	Name  string
	Phone string
}

var users = map[string]*User{
	"1": {
		Id:    "1",
		Name:  "木兮",
		Phone: "13800001111",
	},
	"2": {
		Id:    "2",
		Name:  "小慕",
		Phone: "15688880000",
	},
}

type UserServer struct {
}

func (u *UserServer) GetUser(req *GetUserReq, resp *GetUserResp) error {
	if u, ok := users[req.ID]; ok {
		*resp = GetUserResp{
			Id:    u.Id,
			Name:  u.Name,
			Phone: u.Phone,
		}
		return nil
	} else {
		return errors.New("user not found")
	}
}

func main() {
	userServer := new(UserServer)
	rpc.Register(userServer)
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		go rpc.ServeConn(conn)
	}
}
