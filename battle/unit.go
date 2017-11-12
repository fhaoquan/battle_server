package battle
type Unit struct{
	Owner uint32;
	ID uint16;
	Type uint32;
	Card uint32;
	Level uint8;
	Ranks uint8;
	X uint16;
	Y uint16;
	AttackPower uint16;
	HP uint16;
	Direction uint16;
	AimingFace uint16;
	Speed uint16;
	Status uint16;
	Score uint16;
	Killer uint32;
}
func NewUnit(id uint16)*Unit{
	return &Unit{
		ID:id,
	};
}
func (me *Unit)Death()(bool){
	return me.HP==0;
}
func (me *Unit)SetAll(u *Unit){
	me.Owner=u.Owner;
	me.Type =u.Type;
	me.Card =u.Card;
	me.Level =u.Level;
	me.Ranks =u.Ranks;
	me.X =u.X;
	me.Y =u.Y;
	me.AttackPower =u.AttackPower;
	me.HP =u.HP;
	me.Direction =u.Direction;
	me.AimingFace =u.AimingFace;
	me.Speed =u.Speed;
	me.Status =u.Status;
	me.Score=u.Score;
	me.Killer=u.Killer;
}
func (me *Unit)WriteToBuf(wtr *packet_encoder){
	wtr.write_uint16(me.ID).
		write_uint32(me.Owner).
		write_uint8(me.Ranks).
		write_uint8(me.Level).
		write_uint32(me.Type).
		write_uint32(me.Card).
		write_uint16(me.HP).
		write_uint16(me.X).
		write_uint16(me.Y).
		write_uint16(me.Speed).
		write_uint16(me.Direction).
		write_uint16(me.AimingFace);
}