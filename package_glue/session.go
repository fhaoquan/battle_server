package package_glue
import (
	"../session"
)
type Session session.S_session
func (s *Session)GetUserID()uint32{
	return s.V_user_id;
}