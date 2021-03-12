// Copyright 2021-2021 The jdh99 Authors. All rights reserved.
// knock包用于通信时帧头和帧正文处理解耦
// 帧正文处理函数注册成knock包的服务,帧头处理调用服务,调用knock服务获取响应数据
// Authors: jdh99 <jdh821@163.com>

package knock

// Resp 异步调用应答
type Resp struct {
	IsNeedResp bool
	Bytes      []uint8
	Done       chan *Resp
}

// done 结果返回.框架内调用
func (resp *Resp) done() {
	select {
	case resp.Done <- resp:
	default:
	}
}

// CallbackFunc 注册服务回调函数
// 返回值是应答数据和应答标志.应答标志为false表示不需要应答
type CallbackFunc func(req []uint8, params ...interface{}) ([]uint8, bool)

var services = make(map[int]CallbackFunc)

// Call 同步调用
// 返回值是应答字节流和是否需要应答标志
func Call(protocol uint16, cmd uint16, req []uint8, params ...interface{}) ([]uint8, bool) {
	return callback(protocol, cmd, req, params...)
}

// CallAsync 异步调用
// 返回值是应答字节流和是否需要应答标志
func CallAsync(protocol uint16, cmd uint16, req []uint8, params ...interface{}) *Resp {
	var resp Resp
	resp.Done = make(chan *Resp, 10)

	go func() {
		resp.Bytes, resp.IsNeedResp = callback(protocol, cmd, req, params...)
		resp.done()
	}()

	return &resp
}

// Register 注册服务回调函数
func Register(protocol uint16, cmd uint16, callback CallbackFunc) {
	rid := int(cmd) + (int(protocol) << 16)
	services[rid] = callback
}

// callback 回调命令字对应的函数
func callback(protocol uint16, cmd uint16, req []uint8, params ...interface{}) ([]uint8, bool) {
	rid := int(cmd) + (int(protocol) << 16)
	v, ok := services[rid]
	if ok == false {
		return nil, false
	}
	return v(req, params...)
}
