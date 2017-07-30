package test

import "bytes"

type IUserData interface {
	Clear();
}
type IDataOwner interface {
	IsReturned()bool;
	GetUserData()IUserData;
	Return();
	GetUseOneTime()IUseOneTime;
}
type IUseOneTime interface{
	UseOneTime(func(interface{}));
}
type cached_data struct{
	free bool;
	data IUserData;
	owner *MemoryPool;
}
func (me *cached_data)IsReturned()bool{
	return !me.free;
}
func (me *cached_data)GetUserData()IUserData{
	return me.data;
}
func (me *cached_data)Return(){
	if !me.IsReturned(){
		me.data.Clear();
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
func (pool *MemoryPool)PullOne()IDataOwner{
	return <-pool.cache;
}
func NewMemoryPool(size int,builder func()IUserData)(*MemoryPool){
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
func uuu(){
	bytes.NewBuffer(make([]byte,1024)).WriteRune()
}