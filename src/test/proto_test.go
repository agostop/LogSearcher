package test

import (
	"searchEngine/src/server/transport/protocol"
	"testing"
)

func TestProtoOneLineInput(t *testing.T)  {
    var delim protocol.DelimiterProto = "\n"
    i := delim.Input([]byte("testae asefae wfasdlk lafejlsd\nladjflsdkfpwejfkladsf\nlasdjfljwief{}wlfjd\"kwefjiweif\n"))
    t.Logf("input i: %d", i+1)
}

func TestProtoOneLinePackage0(t *testing.T)  {
    // 分隔符在中间
    var delim protocol.DelimiterProto = "\n"
    b, b2, err := delim.Package([]byte("1111111 111111 111111\n22222 22222 22222 2222\n333333 33333 33333\"44444 44444\n"))
    if err != nil {
        t.Error("got error.")
    }

    if len(b) == 0 {
        t.Errorf("result should not be 0. but: %v", b)
    }

    if len(b2) > 0 {
        t.Errorf("remain data should be 0. but: %v", string(b2))
    }

}

func TestProtoOneLinePackage1(t *testing.T)  {
    // 分隔符在末尾
    var delim protocol.DelimiterProto = "\n"
    result, remain, err := delim.Package([]byte("11111 11111 11111\n22222 22222 22222\n33333 33333 33333\"\n44444 44444 44444"))
    if err != nil {
        t.Error("got error.")
    }

    if len(result) == 0 {
        t.Errorf("result should not be 0. but: %v", result)
    }

    if len(remain) == 0 {
        t.Errorf("remain data should not be 0. but: %v", string(remain))
    }
}

func TestProtoOneLinePackage2(t *testing.T)  {
    // 分隔符在最前面
    var delim protocol.DelimiterProto = "\n"
    result, remain, err := delim.Package([]byte("\n11111 11111 11111\n22222 22222 22222\n33333 33333 33333\"\n44444 44444 44444"))
    if err != nil {
        t.Error("got error.")
    }

    if len(result) == 0 {
        t.Errorf("result should not be 0. but: %v", result)
    }

    if len(remain) == 0 {
        t.Errorf("remain data should not be 0. but: %v", string(remain))
    }
}

func TestProtoOneLinePackage3(t *testing.T)  {
    // 没有分隔符
    var delim protocol.DelimiterProto = "\n"
    result, remain, err := delim.Package([]byte("11111111"))
    if err != nil {
        t.Error("got error.")
    }

    if len(result) > 0 {
        t.Errorf("result len should be 0. but: %v", result)
    }

    if len(remain) == 0 {
        t.Errorf("remain data len should not be 0. but: %v", string(remain))
    }
}

func TestProtoOneLinePackage4(t *testing.T)  {
    // 只有分隔符
    var delim protocol.DelimiterProto = "\n"
    result, remain, err := delim.Package([]byte("\n"))
    if err != nil {
        t.Error("got error.")
    }

    if len(result) > 0 {
        t.Errorf("result len should be 0. but: %v", result)
    }

    if len(remain) == 0 {
        t.Errorf("remain data len should not be 0. but: %v", string(remain))
    }
}

func TestProtoOneLinePackage5(t *testing.T)  {
    // 空字符串
    var delim protocol.DelimiterProto = "\n"
    result, remain, err := delim.Package([]byte(""))
    if err != nil {
        t.Error("got error.")
    }

    if len(result) > 0 {
        t.Errorf("result len should be 0. but: %v", result)
    }

    if len(remain) > 0 {
        t.Errorf("remain data len should not be 0. but: %v", string(remain))
    }
}

func TestProtoOneLinePackage6(t *testing.T)  {
    // 单个字符
    var delim protocol.DelimiterProto = "\n"
    result, remain, err := delim.Package([]byte("n"))
    if err != nil {
        t.Error("got error.")
    }

    if len(result) > 0 {
        t.Errorf("result len should be 0. but: %v", result)
    }

    if len(remain) == 0 {
        t.Errorf("remain data len should not be 0. but: %v", string(remain))
    }
}