package memory_pool

type ICachedData interface {
	onReturn();
	onPop();
}
type IReturnable interface {
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
		me._cached_data_pool.cache<-me;
	}
}
func (me *cached_data)OnPop(){
}
func (me *cached_data)OnReturn(){
}
func (me *cached_data)onPop(){
	me._cached_data_free=true;
	me._cached_data_data.onPop();
}
func (me *cached_data)onReturn(){
	me._cached_data_free=false;
	me._cached_data_data.onReturn();
}
type MemoryPool struct{
	cache chan *cached_data;
}
func (pool *MemoryPool)Len()int{
	return len(pool.cache);
}
func (pool *MemoryPool)Pop() ICachedData {
	o:=<-pool.cache;
	o.onPop();
	return o._cached_data_data;
}
func NewMemoryPool(size int,builder func(IReturnable) ICachedData)(*MemoryPool){
	p:=&MemoryPool{
		make(chan *cached_data,size),
	};
	for i:=0;i<size;i++{
		c:=&cached_data{true,nil,p};
		c._cached_data_data=builder(c);
		c.Return();
	}
	return p;
}