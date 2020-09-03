package xdp

const (
	svrCmdHandshake = 1
)

type HandshakeRequest struct {
	AppID string
	Key   string
}

func (r *HandshakeRequest) Marshal() ([]byte, error) {
	// bb := bytes.NewBuffer(nil)
	// if err := writeString(bb, r.Key); err != nil {
	// 	return nil, err
	// }
	// return bb.Bytes(), nil
	return nil, nil
}

func (r *HandshakeRequest) Unmarshal(data []byte) error {
	// var err error
	// buf := bytes.NewBuffer(data)

	// if r.Key, err = readString(buf); err != nil {
	// 	return err
	// }
	return nil
}
