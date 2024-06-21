package wrapper

import (
	"io"
	"net"
	"time"

	"ferxes.uz/bcp/protocol"
)

type Response struct {
    conn *net.Conn
	CacheTime time.Duration     `json:"cacheTime"`
	ServedBy  string            `json:"servedBy"`
	Additions map[string]string `json:"additions"`
	Data      *io.Reader        `json:"data"`
}

func (r *Response) Reply(){
    response := protocol.Response{
        CacheTime: r.CacheTime,
        ServedBy: r.ServedBy,
        Additions: r.Additions,
        Data: r.Data,
    }

    response.ProtocolFormatWrite(r.conn)
}
