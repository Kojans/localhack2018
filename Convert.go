package main

import (
	"math"
	"encoding/binary"
)

func Float32bytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func GetByte(c *Client) []byte {
	l := byte(len(c.name))
	lMain_Gun := uint8(len(c.Main_Guns))
	lAdditionals_Gun := uint8(len(c.Additional_Guns))
	forSend := make([]byte, l+lMain_Gun*2+lAdditionals_Gun*2+38)
	forSend[0] = 1
	forSend[1] = l
	var i uint8 = 2
	for ; i < l+2; i++ {
		forSend[i] = c.name[i-2]
	}
	a := make([]byte, 2)
	f := Float64bytes(c.x)
	forSend[l+2] = f[0]
	forSend[l+3] = f[1]
	forSend[l+4] = f[2]
	forSend[l+5] = f[3]
	forSend[l+6] = f[4]
	forSend[l+7] = f[5]
	forSend[l+8] = f[6]
	forSend[l+9] = f[7]
	f = Float64bytes(c.y)
	forSend[l+10] = f[0]
	forSend[l+11] = f[1]
	forSend[l+12] = f[2]
	forSend[l+13] = f[3]
	forSend[l+14] = f[4]
	forSend[l+15] = f[5]
	forSend[l+16] = f[6]
	forSend[l+17] = f[7]
	a = make([]byte, 2)
	binary.LittleEndian.PutUint16(a, c.hp)
	forSend[l+18] = a[1]
	forSend[l+19] = a[0]
	binary.LittleEndian.PutUint16(a, c.max_hp)
	forSend[l+20] = a[1]
	forSend[l+21] = a[0]
	binary.LittleEndian.PutUint16(a, c.shield)
	forSend[l+22] = a[1]
	forSend[l+23] = a[0]
	binary.LittleEndian.PutUint16(a, c.max_shield)
	forSend[l+24] = a[1]
	forSend[l+25] = a[0]

	forSend[l+26] = c.Ship.Type

	forSend[l+27] = lMain_Gun
	i = l + 28
	for q, o := range c.Main_Guns {
		forSend[i] = q
		forSend[i+1] = o.t
		i += 2
	}
	i = l + lMain_Gun*2 + 28
	forSend[i] = lAdditionals_Gun //[+28]
	i++
	for q, o := range c.Additional_Guns {
		forSend[i] = q
		forSend[i+1] = o.t
		i += 2
	}

	i = 0
	for ; i < 4; i++ {
		forSend[l+lAdditionals_Gun*2+lMain_Gun*2+29+i] = c.angleB[i]
	}
	if (c.base) {
		forSend[l+lAdditionals_Gun*2+lMain_Gun*2+29+i] = 1
	} else {
		forSend[l+lAdditionals_Gun*2+lMain_Gun*2+29+i] = 0
	}
	i = 0
	for ; i < 4; i++ {
		forSend[l+lAdditionals_Gun*2+lMain_Gun*2+34+i] = c.a_gun_angleB[i]
	}
	//log.Println(forSend)
	return forSend[:]
}
