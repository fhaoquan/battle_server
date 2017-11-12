package utils

const MaxPktSize  				=1024*1;
const MaxRoomSize 				=1024;
const UdpListenStart			=10000;

const CMD_pingpong =uint8(0);
const CMD_battle_wating_start =uint8(1);
const CMD_battle_all =uint8(2);
const CMD_unit_movment =uint8(3);
const CMD_create_unit =uint8(4);
const CMD_attack_done =uint8(5);
const CMD_attack_start =uint8(6);
const CMD_battle_start =uint8(7);
const CMD_unit_destory =uint8(8);
const CMD_battle_end =uint8(9);
const CMD_battle_remaining_time =uint8(10);
