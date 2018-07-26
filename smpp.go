package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const host = ""
const port = "3000"
const proto = "tcp"

type sms struct {
	size_hex string
	size uint64
	command_id_hex string
	command_id string
	command_status string
	sequence_number string
	pdu []byte
	body string
}

func main() {

	l, err := net.Listen(proto, host+":"+port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not listen TCP: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()
	fmt.Fprintf(os.Stderr, "TCP Listening\n")
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in Accept: %v\n", err)
			os.Exit(1)
		}
		go receive(conn)
	}

}

func hex2int(hexStr string) uint64 {
	cleaned := strings.Replace(hexStr, "0x", "", -1)
	result, _ := strconv.ParseUint(cleaned, 16, 64)
	return uint64(result)
}

func getCommandName (command_id string) string {
	c := make(map[uint64]string)
	c[hex2int("0x80000000")] = "generic_nack"
	c[hex2int("0x00000001")] = "bind_receiver"
	c[hex2int("0x80000001")] = "bind_receiver_resp"
	c[hex2int("0x00000002")] = "bind_transmitter"
	c[hex2int("0x80000002")] = "bind_transmitter_resp"
	c[hex2int("0x00000003")] = "query_sm"
	c[hex2int("0x80000003")] = "query_sm_resp"
	c[hex2int("0x00000009")] = "bind_transceiver"
	c[hex2int("0x80000009")] = "bind_transceiver_resp"
	c[hex2int("0x00000015")] = "enquire_link"
	c[hex2int("0x80000015")] = "enquire_link_resp"
	c[hex2int("0x00000003")] = "query_sm"
	c[hex2int("0x80000003")] = "query_sm_resp"
	c[hex2int("0x00000004")] = "submit_sm"
	c[hex2int("0x80000004")] = "submit_sm_resp"
	c[hex2int("0x00000005")] = "deliver_sm"
	c[hex2int("0x80000005")] = "deliver_sm_resp"
	c[hex2int("0x00000006")] = "unbind"
	c[hex2int("0x80000006")] = "unbind_resp"
	c[hex2int("0x00000007")] = "replace_sm"
	c[hex2int("0x80000007")] = "replace_sm_resp"
	c[hex2int("0x00000008")] = "cancel_sm"
	c[hex2int("0x80000008")] = "cancel_sm_resp"
	c[hex2int("0x0000000B")] = "outbind"
	c[hex2int("0x00000021")] = "submit_multi"
	c[hex2int("0x80000021")] = "submit_multi_resp"
	c[hex2int("0x00000102")] = "alert_notification"
	c[hex2int("0x00000103")] = "data_sm"
	c[hex2int("0x80000103")] = "data_sm_resp"
	
	if _, ok := c[hex2int(command_id)] ; ok {
		return c[hex2int(command_id)]
	}
	return ""
}

func getCommandHex (commandName string) uint64 {

	c := make(map[string]string)

	c["generic_nack"] = "0x80000000"
	c["bind_receiver"] = "0x00000001"
	c["bind_receiver_resp"] = "0x80000001"
	c["bind_transmitter"] = "0x00000002"
	c["bind_transmitter_resp"] = "0x80000002"
	c["query_sm"] = "0x00000003"
	c["query_sm_resp"] = "0x80000003"
	c["bind_transceiver"] = "0x00000009"
	c["bind_transceiver_resp"] = "0x80000009"
	c["enquire_link"] = "0x00000015"
	c["enquire_link_resp"] = "0x80000015"
	c["query_sm"] = "0x00000003"
	c["query_sm_resp"] = "0x80000003"
	c["submit_sm"] = "0x00000004"
	c["submit_sm_resp"] = "0x80000004"
	c["deliver_sm"] = "0x00000005"
	c["deliver_sm_resp"] = "0x80000005"
	c["unbind"] = "0x00000006"
	c["unbind_resp"] = "0x80000006"
	c["replace_sm"] = "0x00000007"
	c["replace_sm_resp"] = "0x80000007"
	c["cancel_sm"] = "0x00000008"
	c["cancel_sm_resp"] = "0x80000008"
	c["outbind"] = "0x0000000B"
	c["submit_multi"] = "0x00000021"
	c["submit_multi_resp"] = "0x80000021"
	c["alert_notification"] = "0x00000102"
	c["data_sm"] = "0x00000103"
	c["data_sm_resp"] = "0x80000103"
	
	if _, ok := c[commandName]; ok {
		return hex2int(c[commandName])
	}
	return 0
}

func enquire_link(conn net.Conn, sms sms) {
	var response string;
	response = fmt.Sprintf("%08X%08X%08X%s", 16, getCommandHex("enquire_link_resp"), 0, sms.sequence_number)
	
	fmt.Fprintf(os.Stderr, "%s\n", response)
	brsp := make([]byte, hex.DecodedLen(len(response)))
	hex.Decode(brsp,[]byte(response))
	fmt.Fprintf(os.Stderr, "-%s-\n", brsp)
	conn.Write(brsp)
}

func bind_transceiver(conn net.Conn, sms sms){
	
	var response string
	response = fmt.Sprintf("%08X%08X%08X%s", 16, getCommandHex("bind_transceiver_resp"), 0, sms.sequence_number)
	
	fmt.Fprintf(os.Stderr, "%s\n", response)
	brsp := make([]byte, hex.DecodedLen(len(response)))
	hex.Decode(brsp,[]byte(response))
	fmt.Fprintf(os.Stderr, "-%s-\n", brsp)
	conn.Write(brsp)
}

func submit_sm(conn net.Conn, sms sms) {

	var response string
	response = fmt.Sprintf("%08X%08X%08X%s", 16, getCommandHex("submit_sm_resp"), 0, sms.sequence_number)
	
	fmt.Fprintf(os.Stderr, "%s\n", response)
	brsp := make([]byte, hex.DecodedLen(len(response)))
	hex.Decode(brsp,[]byte(response))
	fmt.Fprintf(os.Stderr, "-%s-\n", brsp)
	conn.Write(brsp)
}

func receive(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		fmt.Fprintf(os.Stderr,"Next input\n")
		rlen, err := conn.Read(buf)
		dst := make([]byte, hex.EncodedLen(len(buf)))

		fmt.Fprintf(os.Stderr, "%s\n", string(dst))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in read: %v\n", err)
			break
		}
		fmt.Fprintf(os.Stderr, " Received: %s (%d)\n", string(buf), rlen)
		hex.Encode(dst, buf)
		fmt.Fprintf(os.Stderr, "Hex value: %s\n", string(dst))

		sms := sms{sequence_number : string(dst)[24:32]}
		sms.size_hex = string(dst)[0:8]
		sms.command_id_hex = string(dst)[8:16]
		sms.command_status = string(dst)[16:24]
		sms.size = hex2int(sms.size_hex)
		sms.command_id = getCommandName(sms.command_id_hex)
		sms.pdu = buf
		sms.body = string(buf[16:sms.size])

		fmt.Fprintf(os.Stderr, "Size: %s (%d)\n", sms.size_hex, sms.size)
		fmt.Fprintf(os.Stderr, "Command: %s (%s)\n", sms.command_id_hex, sms.command_id)
		fmt.Fprintf(os.Stderr, "Status: %s\n", sms.command_status)
		fmt.Fprintf(os.Stderr, "Sequence: %s\n", sms.sequence_number)

		fmt.Fprintf(os.Stderr, "Message: -%s-\n", sms.body)

		if sms.command_id == "bind_transceiver" {
			go bind_transceiver(conn,sms)
		}

		if sms.command_id == "enquire_link" {
			go enquire_link(conn, sms)
		}

		
	}

}

