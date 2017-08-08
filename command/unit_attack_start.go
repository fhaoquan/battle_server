package command

import (
	"fmt"
	"errors"
)

func (cmd *Commamd)UnitAttackStart(data []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	return nil;
}