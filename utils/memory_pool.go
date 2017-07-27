package utils

type IUserData interface {
	Clear();
}
type ICachedData interface {
	IsReturned()bool;
	GetUserData()IUserData;
	Return();
}
type cached_data struct{
	free bool;
	data IUserData;
	owner *MemoryPool;
}
func (d *cached_data)IsReturned()bool{
	return !d.free;
}
func (d *cached_data)GetUserData()IUserData{
	return d.data;
}
func (d *cached_data)Return(){
	d.owner.put(d);
}

type MemoryPool struct{
	size int;
	cache chan *cached_data;
}
func (pool *MemoryPool)Get()(ICachedData){
	d:=<-pool.cache;
	d.free=true;
	return d;
}
func (pool *MemoryPool)put(d *cached_data){
	if(d.free){
		d.free=false;
		d.GetUserData().Clear();
		pool.cache<-d;
	}
}

func NewMemoryPool(size int,builder func(ICachedData)IUserData)(*MemoryPool){
	p:=&MemoryPool{
		size,
		make(chan *cached_data,size),
	};
	for i:=0;i<size;i++{
		c:=&cached_data{false,nil,p}
		c.data=builder(c);
		p.cache<-c;
	}
	return p;
}
