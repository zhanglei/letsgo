/*=============================================================================
#     FileName: defaultdispatcher.go
#         Desc: default dispatcher
#       Author: sunminghong
#        Email: allen.fantasy@gmail.com
#     HomePage: http://weibo.com/5d13
#      Version: 0.0.1
#   LastChange: 2013-06-06 14:37:57
#      History:
=============================================================================*/
package gate

import (
    "strconv"
    "strings"
    "github.com/sunminghong/letsgo/log"
)

type DefaultDispatcher struct {
    messageCodemaps map[int]int
}

func (r *DefaultDispatcher) Init() {
    r.messageCodemaps = make(map[int]int)
}

func (r *DefaultDispatcher) Add(gridID int, messageCodes string) {
    codes:= strings.Split(messageCodes,",")
    for _,p_ := range codes {
        p := strings.Trim(p_," ")
        if len(p) == 0 {
            continue
        }
        pmessageCode, err := strconv.Atoi(p)
        if err ==nil {
            r.messageCodemaps[pmessageCode] = gridID
        }
    }
}

func (r *DefaultDispatcher) Dispatch(messageCode int) (gridID int,ok bool) {
    if dp.Type == DATAPACKET_TYPE_DELAY {
        return 0,false
    }

    gcode := r.GroupCode(messageCode)

    gridID,ok = r.messageCodemaps[gcode]

    log.Trace("dispatcher Handler func messageCode,messageCode,gridID:",messageCode,gcode,gridID)

    return gridID,ok
}

//将协议编号分组以供Dispatch决策用那个Grid 来处理
func (r *DefaultDispatcher) GroupCode(messageCode uint16) int {
    return int(messageCode / 100)
}
