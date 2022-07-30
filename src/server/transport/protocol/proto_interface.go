package protocol

type ProtoInf interface {
    Input([]byte) (int)
    Package([]byte) ([][]byte, []byte, error)
}