package battle
type Unit struct{
	ID uint16;
	Camps uint16;
	X uint16;
	Y uint16;
	AttackPower uint16;
	HP uint16;
	Direction uint16;
	AimingFace uint16;
	Speed uint16;
	Status uint16;
}
func NewUnit(id uint16)*Unit{
	return &Unit{
		ID:id,
	};
}