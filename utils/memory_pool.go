package utils
type IDataOwner interface {
	IsReturned()bool;
	GetUserData()interface{};
	Return();
	GetUseOneTime()IUseOneTime;
}
type IUseOneTime interface{
	UseOneTime(func(interface{}));
}
type cached_data struct{
	free bool;
	data interface{};
	owner *MemoryPool;
}
func (me *cached_data)IsReturned()bool{
	return !me.free;
}
func (me *cached_data)GetUserData()interface{}{
	return me.data;
}
func (me *cached_data)Return(){
	if !me.IsReturned(){
		me.owner.cache<-me;
		me.free=false;
	}
}
func (me *cached_data)UseOneTime(f func(interface{})){
	if(me.IsReturned()){
		return ;
	}
	if(f!=nil){
		f(me.GetUserData());
		me.Return();
	}
}
func (me *cached_data)GetUseOneTime()IUseOneTime{
	return me;
}
type MemoryPool struct{
	cache chan *cached_data;
}
func (pool *MemoryPool)PopOne()IDataOwner{
	return <-pool.cache;
}
func NewMemoryPool(size int,builder func()interface{})(*MemoryPool){
	p:=&MemoryPool{
		make(chan *cached_data,size),
	};
	for i:=0;i<size;i++{
		c:=&cached_data{false,nil,p}
		c.data=builder();
		p.cache<-c;
	}
	return p;
}