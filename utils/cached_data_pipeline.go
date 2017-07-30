package utils

type CachedDataAdapter struct {
	u interface{Clear()}
	o ICachedData;
}
func(c *CachedDataAdapter)UseUntilTrue(f func(interface{})bool){
	if(c.o.IsReturned()){
		return ;
	}
	if(f(c.u)){
		c.o.Return();
	}
}
func(c *CachedDataAdapter)Clear(){
	c.u.Clear();
}
type CachedDataPipeline struct {
	ch chan func(f func(interface{})bool);
}
func (c *CachedDataPipeline)In(d *CachedDataAdapter){
	c.ch<-d.UseUntilTrue;
}
func (c *CachedDataPipeline)Ch()chan func(f func(interface{})bool){
	return c.ch;
}
type MemoryPoolWithCachedDataPipeline struct{
	*MemoryPool
}
func (pool *MemoryPoolWithCachedDataPipeline)Get()(*CachedDataAdapter){
	return pool.MemoryPool.Get().GetUserData().(*CachedDataAdapter)
}
func NewMemoryPoolWithCachedDataPipeline(size int,f func()IUserData)(*MemoryPoolWithCachedDataPipeline){
	return &MemoryPoolWithCachedDataPipeline{NewMemoryPool(size,func(data ICachedData)IUserData{
		return &CachedDataAdapter{
			f(),data,
		}
	})}
}

