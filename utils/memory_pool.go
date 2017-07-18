package utils
/*
type user_data_builder func()UserData;

type UserData interface {
}

type CachedData struct{
	data UserData;
	owner *MemoryPool;
}

func (d *CachedData)GetUserData()UserData{
	return d.data;
}
func (d *CachedData)Return(){
	d.owner.Put(d);
}

type MemoryPool struct{
	size int;
	cache chan *CachedData;
}
func (pool *MemoryPool)Chan()chan *CachedData{
	return pool.cache;
}
func (pool *MemoryPool)Get()(*CachedData){
	return <-pool.cache;
}
func (pool *MemoryPool)Put(d *CachedData)(bool){
	if(d.owner==pool){
		pool.cache<-d;
		return true;
	}else{
		return false;
	}
}

func NewPool(size int,builder user_data_builder)(*MemoryPool){
	p:=&MemoryPool{
		size,
		make(chan *CachedData,size),
	};
	for i:=0;i<size;i++{
		p.cache<-&CachedData{builder(),p};
	}
	return p;
}
*/