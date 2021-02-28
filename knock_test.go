package knock

import (
    "fmt"
    "testing"
)

func TestCase1(t *testing.T) {
    Register(0, 1, callback1)
    resp, ok := Call(0, 1, []uint8{2, 3, 4})
    fmt.Println(ok, resp)
}

// CallbackFunc 注册服务回调函数
// 返回值是应答数据和应答标志.应答标志为false表示不需要应答
func callback1(req []uint8) ([]uint8, bool) {
    fmt.Println(req)
    return []uint8{1, 2, 3}, true
}

func TestCase2(t *testing.T) {
    Register(0, 1, callback1)
    resp, ok := Call(1, 1, []uint8{2, 3, 4})
    fmt.Println(ok, resp)
}

func TestCase3(t *testing.T) {
    Register(0, 1, callback1)
    resp := CallAsync(0, 1, []uint8{2, 3, 4})
    <-resp.Done
    fmt.Println(resp.IsNeedResp, resp.Bytes)
}

func TestCase4(t *testing.T) {
    Register(0, 1, callback1)
    resp := CallAsync(1, 1, []uint8{2, 3, 4})
    <-resp.Done
    fmt.Println(resp.IsNeedResp, resp.Bytes)
}

func TestCase5(t *testing.T) {
    Register(0, 2, callback2)
    for i := 0; i < 1000000; i++ {
        resp := CallAsync(1, 1, []uint8{2, 3, 4})
        <-resp.Done
    }
}

// CallbackFunc 注册服务回调函数
// 返回值是应答数据和应答标志.应答标志为false表示不需要应答
func callback2(req []uint8) ([]uint8, bool) {
    for i := 0; i < 1000; i++ {

    }
    return req, true
}

func TestCase6(t *testing.T) {
    for i := 0; i < 10000000; i++ {
        callback2([]uint8{1, 2, 3})
    }
}

func TestCase7(t *testing.T) {
    Register(0, 2, callback2)
    for i := 0; i < 10000000; i++ {
        Call(0, 2, []uint8{2, 3, 4})
    }
}

func TestCase8(t *testing.T) {
    Register(0, 2, callback2)
    for i := 0; i < 10000000; i++ {
        resp := CallAsync(0, 2, []uint8{2, 3, 4})
        <-resp.Done
    }
}
