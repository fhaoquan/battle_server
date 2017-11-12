package battle

func (context *Battle)ComputeResultScore(pid uint32)(score int){
	score=0;
	context.ForEachUnitDo(func(u *Unit) bool {
		if u.Death()&&u.Killer==pid{
			score+=(int)(u.Score);
		}
		return true;
	});
	return score;
}
