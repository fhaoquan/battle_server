package builder
import (
	"../session"
)
type SessionGlue session.S_session
func (s *SessionGlue)GetUserID()uint32{
	return s.V_user_id;
}