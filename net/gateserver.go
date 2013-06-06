/*=============================================================================
#     FileName: gateserver.go
#         Desc: server base 
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-23 15:42:14
#      History:
=============================================================================*/
package net

import (
    "strconv"
    "time"
    goconf "github.com/hgfischer/goconf"
    //"github.com/sunminghong/letsgo/helper"
    "github.com/sunminghong/letsgo/log"
)
/*
[GateServer]
name = gate1
host = :12000
maxConnections = 3000
endian = 1 #BigEndian= 0 ,LittleEndian = 1

[LogicServer1]
name = game1
host = :12001
process = 1,2,3,4,5,6
endian = 1

[LogicServer2]
name = game2
host = :12002
process = 
endian = 1

[LogicServer3]
name = game3
host = :12003
process = 7,8
endian = 1

*/

type GateServer struct {
    Name                string

    datagram            IDatagram
    GateService         *Server
    LogicServicePool    *ClientPool

    RouteHandler IRouter
}

func NewGateServer(Client,makeclient NewClientFunc,
    makeLogicService NewClientFunc, datagram IDatagram,
    routehand IRouter) *GateServer {

    gs := &GateServer{}

    gs.datagram = datagram
    gs.GateService = NewServer(makeclient, datagram)
    gs.LogicServicePool = NewClientPool(makeLogicService,datagram)
    gs.RouteHandler = routehand

    return gs
}

func (gs *GateServer) Start(configfile string) {
    //parse config ini file
    c, err := goconf.ReadConfigFile(configfile)
    if err != nil {
        log.Error(err.Error())
        return
    }

    //make some connection to game server
    for i:=1; i<50; i++ {
        section := "GateServer" + strconv.Itoa(i)
        gname, err := c.GetString(section,"name")
        if err != nil {
            //if err.Reason == goconf.SectionNotFound {
            //    break
            //} else {
                log.Error(err.Error())
            //    continue
            //}
            break
        }

        host, err := c.GetString(section,"host")
        if err != nil {
            continue
        }

        protocols, err := c.GetString(section,"protocols")
        if err != nil {
            protocols = ""
        }

        endian, err := c.GetInt("GateServer","endian")
        if err == nil {
            da := gs.datagram.Clone(endian)
            gs.AddLogicServer(gname, host, protocols,da)
        } else {
            gs.AddLogicServer(gname, host, protocols,nil)
        }

    }

    //start gate service
    gatename, err := c.GetString("GateServer","name")
    if err != nil {
        log.Error(err.Error())
        return
    }

    gatehost, err := c.GetString("GateServer","host")
    if err != nil {
        log.Error(err.Error())
        return
    }

    maxConnections, err := c.GetInt("GateServer","maxConnections")
    if err != nil {
        log.Error(err.Error())
        return
    }

    endian, err := c.GetInt("GateServer","endian")
    if err == nil {
        gs.datagram.SetEndian(endian)
    }

    gs.Name = gatename
    gs.GateService.Start(gatehost,maxConnections)
}

func (gs *GateServer) AddLogicServer(name string,host string,protocols string,datagram IDatagram) {

        pool := gs.LogicServicePool
        go pool.Start(name, host, datagram)
        time.Sleep(1)

        //if Pool don't find it ,then that is no success!
        c := pool.Clients.GetByName(name)
        if c == nil {
            log.Error(host + " can't connect")
            return
        }

        //add route
        gs.RouteHandler.Add(c.GetTransport().Cid,protocols)
}
