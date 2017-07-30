package room

type kcp_channel struct {
	c chan func(func(uid uint32,rid uint32,bdy []byte));
}
