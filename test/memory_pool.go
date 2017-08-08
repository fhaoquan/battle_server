package test

type ICachedData interface {
	IsReturned()bool;
	OnPop();
	OnReturn();
	Return();
}
type cached_data struct{
	_cached_data_free bool;
	_cached_data_data ICachedData;
	_cached_data_pool *MemoryPool;
}
func (me *cached_data)IsReturned()bool{
	return !me._cached_data_free;
}
func (me *cached_data)Return(){
	if !me.IsReturned(){
		me._cached_data_pool.cache<-me._cached_data_data;
		me.OnReturn();
	}
}
func (me *cached_data)OnPop(){
	me._cached_data_free=true;
}
func (me *cached_data)OnReturn(){
	me._cached_data_free=false;
}
type MemoryPool struct{
	cache chan ICachedData;
}
func (pool *MemoryPool)PullOne()ICachedData{
	o:=<-pool.cache;
	o.OnPop();
	return o;
}
func NewMemoryPool(size int,builder func(ICachedData)ICachedData)(*MemoryPool){
	p:=&MemoryPool{
		make(chan ICachedData,size),
	};
	for i:=0;i<size;i++{
		c:=&cached_data{true,nil,p};
		c._cached_data_data=builder(c);
		c.Return();
	}
	return p;
}