/*=============================================================================
#     FileName: datagram_test.go
#         Desc: Datagram pack/unpack
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-05-09 16:07:54
#      History:
=============================================================================*/
package net

import (
    //"encoding/binary"
    "testing"
    "bytes"
)

/*

//datagram and datapacket define
type IDatagram interface {
    Fetch(c *Transport) (n int, dps []*DataPacket)
    Pack(dp *DataPacket) []byte
}
*/

func Test_Pack(t *testing.T) {
    datagram := &Datagram{}

    data := []byte("1234567890")
    dp := &DataPacket{Type:1,Data:data}
    data2 := datagram.Pack(dp)

    _data := [17]byte{0x59,0x7a,1,0,0,0,10}
    copy(_data[7:],data)

    if !bytes.Equal(_data[:],data2) {
        t.Error("pack return is equal:",data2)
    }

}

func Test_Fetch(t *testing.T) {
    /*
  username, _ := getAdmin(1)
    if (username != "admin") {
         t.Error("getAdmin get data error")
    }
    */

    datagram := &Datagram{}
    trans := NewTransport(1,nil,nil)

    trans.Buff = []byte{0x59,0x7a,1,0,0,0,10}

    data0 := []byte("1234567890")
    trans.BuffAppend(data0)
    data := trans.Buff
    trans.BuffAppend(data)
    trans.BuffAppend(data)
    trans.BuffAppend(data)
    trans.BuffAppend(data)

    Log("buff:",trans.Buff)

    n,dps := datagram.Fetch(trans)
    if n != 5 || len(dps)!= 5 {
        t.Error("fetch dps len is error:",n,len(dps))
    }

    dp := dps[0]
    if dp.Type != 1 {
        t.Error("fetch dps data is error")
    }

    if !bytes.Equal(dp.Data,data0) {
        t.Error("fetch dps data is error" + string(dp.Data))
    }

    dp = dps[4]
    if dp.Type != 1 {
        t.Error("fetch dps data is error",4)
    }

    if !bytes.Equal(dp.Data,data0) {
        t.Error("fetch dps data is error" ,4, string(dp.Data))
    }
}
