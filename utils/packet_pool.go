package utils

type I_cached_data interface {
	Clear();
	ReturnToPool();
}

type PacketPool struct{
	size int;
	cache chan I_cached_data;
}
type inner_cached_data struct{
	pool *PacketPool;
	owner I_cached_data;
}
func (b *inner_cached_data)Clear(){
}
func (b *inner_cached_data)ReturnToPool(){
	b.owner.Clear();
	b.pool.cache<-b.owner;
}

func (pool *PacketPool)GetEmptyPkt()(I_cached_data){
	return <-pool.cache;
}

func NewPacketPool(size int,builder func(I_cached_data)I_cached_data)(*PacketPool){
	p:=&PacketPool{
		size,
		make(chan I_cached_data,size),
	};
	for i:=0;i<size;i++{
		d:=&inner_cached_data{p,nil};
		i:=builder(d);
		p.cache<-i;
		d.owner=i;
	}
	return p;
}
