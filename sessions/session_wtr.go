package sessions

type message_getter func()[]byte
type send_loop_context struct {
	s *session
}

func (me *send_loop_context)Do(){

}
type i_WithMsgGetterRtn interface{Do()}
func (me *send_loop_context)WithMsgGetter(getter message_getter)(i_WithMsgGetterRtn){
	return me;
}

type i_send_loop_context_WithSessionRtn interface{WithMsgGetter(getter message_getter)(i_WithMsgGetterRtn)}
func (me *send_loop_context)WithSession(s *session)i_send_loop_context_WithSessionRtn{
	return me;
}
type i_SdlpRtn interface {WithSession(*session)(i_send_loop_context_WithSessionRtn)}
func NewSendLoop()i_SdlpRtn  {
	return &send_loop_context{}
}