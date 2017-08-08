package utils
type ICachedData interface {
	isReturned()bool;
	onPop();
	onReturn();
	Return();
}
type cached_data struct{
	_cached_data_free bool;
	_cached_data_data ICachedData;
	_cached_data_pool *MemoryPool;
}
func (me *cached_data)isReturned()bool{
	return !me._cached_data_free;
}
func (me *cached_data)Return(){
	if !me.isReturned(){
		me.onReturn();
		me._cached_data_pool.cache<-me._cached_data_data;
	}
}
func (me *cached_data)onPop(){
	me._cached_data_free=true;
}
func (me *cached_data)onReturn(){
	me._cached_data_free=false;
}
type MemoryPool struct{
	cache chan ICachedData;
}
func (pool *MemoryPool)Pop()ICachedData{
	o:=<-pool.cache;
	o.onPop();
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