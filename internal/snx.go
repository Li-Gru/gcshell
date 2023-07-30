package internal
/*
BSD 3-Clause License

Copyright (c) 2023, Francisco Anderson Rodrigues da Silva

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
// Copy from https://github.com/francisco-anderson/snx-go

import (
"math/big"
"net"
"os/exec"
"strconv"

binary_pack "github.com/roman-kachanovsky/go-binary-pack/binary-pack"
)

type SNX struct {
	Params  map[string]string
	Debug   bool
	SnxPath string
	info    []byte
}

func (extender *SNX) GenerateSNXInfo() {
	params := extender.Params
	gwIP, err := net.LookupHost(params["host_name"])
	Iferr(err)

	bp := new(binary_pack.BinaryPack)

	ip := net.ParseIP(gwIP[0])
	ipv4 := big.NewInt(0)
	ipv4.SetBytes(ip.To4())
	tmp := ip.To4()

	hwData, err := bp.UnPack([]string{"I"}, []byte{tmp[3], tmp[2], tmp[1], tmp[0]})
	Iferr(err)

	gwInt := hwData[0].(int)

	magic := string([]byte{0x13, 0x11, 0x00, 0x00})
	length := 0x3d0

	port, err := strconv.Atoi(params["port"])
	Iferr(err)

	format := []string{"4s", "L", "L", "64s", "L", "6s", "256s", "256s", "128s", "256s", "H"}

	values := []interface{}{
		magic,
		length,
		gwInt,
		params["host_name"],
		port,
		string([]byte{0}),
		params["server_cn"],
		params["user_name"],
		params["password"],
		params["server_fingerprint"],
		1,
	}

	data, err := bp.Pack(format, values)
	Iferr(err)

	extender.info = data
}

func (extender *SNX) CallSNX() {
	snxCmd := exec.Command(extender.SnxPath, "-Z")

	_, err := snxCmd.Output()
	Iferr(err)

	connection, err := net.Dial("tcp", "localhost:7776")
	Iferr(err)

	_, err = connection.Write(extender.info)
	Iferr(err)

	buffer := make([]byte, 4096)

	_, err = connection.Read(buffer)
	Iferr(err)

	connection.Read(buffer) //Block execution

}


