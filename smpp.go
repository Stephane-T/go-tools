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

var sequence_number int

type pdu struct {
	size uint64
	command_id string
	command_status string
	sequence_number string
	pdu []byte
	body string
	// SMS
	service_type string
	source_addr_ton string
	source_addr_npi string
	source_addr string
	dest_addr_ton string
	dest_addr_npi string
	dest_addr string
	esm_class string
	protocol_id string
	priority_flag string
	delivery_time string
	validity_period string
	registered_delivery string
	replace_if_present_flag string
	data_coding string
	sm_default_msg_id string
	sm_length uint64
	sm string
}

func get_sequence_number () string {
	sequence_number++
	return fmt.Sprintf("%08d",sequence_number)
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

func enquire_link(conn net.Conn, sms pdu) {
	var response string;
	response = fmt.Sprintf("%08X%08X%08X%s", 16, getCommandHex("enquire_link_resp"), 0, sms.sequence_number)
	
	fmt.Fprintf(os.Stderr, "%s\n", response)
	brsp := make([]byte, hex.DecodedLen(len(response)))
	hex.Decode(brsp,[]byte(response))
	conn.Write(brsp)
}

func bind_transceiver(conn net.Conn, sms pdu){
	
	var response string
	response = fmt.Sprintf("%08X%08X%08X%s", 16, getCommandHex("bind_transceiver_resp"), 0, sms.sequence_number)
	
	fmt.Fprintf(os.Stderr, "%s\n", response)
	brsp := make([]byte, hex.DecodedLen(len(response)))
	hex.Decode(brsp,[]byte(response))
	conn.Write(brsp)
}

func get_addr_ton(value string) string {
	c := make(map[string]string)
	c["00"] = "Unknown"
	c["01"] = "International"
	c["02"] = "National"
	c["03"] = "Network Specific"
	c["04"] = "Subscriber Number"
	c["05"] = "Alphanumeric"
	c["06"] = "Abbreviated"
	if _, ok := c[value]; ok {
		return c[value]
	}
	return "Error"
}

func get_addr_npi(value string) string {
	c := make(map[string]string)
	c["00"] = "Unknown"
	c["01"] = "ISDN"
	c["02"] = "Data"
	c["04"] = "Telex"
	c["06"] = "Land Mobile"
	c["08"] = "National"
	c["09"] = "Private"
	c["10"] = "ERMES"
	c["14"] = "Internet"
	c["18"] = "WAP"
	if _, ok := c[value]; ok {
		return c[value]
	}
	return "Error"

}


func HexToBin(hex string) (string, error) {
	ui, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		return "", err
	}

	// %016b indicates base 2, zero padded, with 16 characters
	return fmt.Sprintf("%08b", ui), nil
}

func get_data_coding(dc string) string {
	c := make(map[string]string)
	c["00"] = "Default"
	c["01"] = "IA5"
	c["02"] = "8b"
	c["03"] = "Latin1"
	c["04"] = "8b"
	c["05"] = "JIS"
	c["06"] = "Cyrillic"
	c["07"] = "Hebrew"
	c["08"] = "UCS2"
	c["09"] = "Pictogram"
	c["10"] = "ISO-2022-JP"
	c["13"] = "Kanji"
	c["14"] = "KS C 5601"
	if _, ok := c[dc]; ok {
		return c[dc]
	}
	return "Reserved"
	
	

}

func submit_sm(conn net.Conn, sms pdu) {

	var response string
	response = fmt.Sprintf("%08X%08X%08X%s", 16, getCommandHex("submit_sm_resp"), 0, sms.sequence_number)

	dst := make([]byte, hex.EncodedLen(len([]byte(sms.body))))
	hex.Encode(dst, []byte(sms.body))

	var tmp_sum string
	var tmp string
	var tmp_b []byte
	i := 0
	
	for ; i < 10; i = i+2 {
		tmp = string(dst[i:i+2])
		if tmp == "00" {
			break
		}
		tmp_sum += tmp
	}
	tmp_b = make([]byte, hex.DecodedLen(len([]byte(tmp_sum))))
	hex.Decode(tmp_b,[]byte(tmp_sum))
	sms.service_type = string(tmp_b)
	tmp_sum = ""
	i = i+2
	
	sms.source_addr_ton = string(dst[i:i+2])
	i = i+2

	sms.source_addr_npi = string(dst[i:i+2])
	i = i+2
	
	for ; i < 50; i = i+2 {
		tmp = string(dst[i:i+2])
		if tmp == "00" {
			break
		}
		tmp_sum += tmp
	}
	tmp_b = make([]byte, hex.DecodedLen(len([]byte(tmp_sum))))
	hex.Decode(tmp_b, []byte(tmp_sum))
	fmt.Fprintf(os.Stderr,"!!!%s\n",tmp_sum)
	sms.source_addr = string(tmp_b)
	tmp_sum = ""
	i = i+2

	sms.dest_addr_ton = string(dst[i:i+2])
	i = i+2

	sms.dest_addr_npi = string(dst[i:i+2])
	i = i+2
	
	for ; i < 50; i = i+2 {
		tmp = string(dst[i:i+2])
		if tmp == "00" {
			break
		}
		tmp_sum += tmp
	}
	tmp_b = make([]byte, hex.DecodedLen(len([]byte(tmp_sum))))
	hex.Decode(tmp_b, []byte(tmp_sum))
	sms.dest_addr = string(tmp_b)
	tmp_sum = ""
	i = i+2

	sms.esm_class = string(dst[i:i+2]) 
	i = i+2

	sms.protocol_id = string(dst[i:i+2])
	i = i+2

	sms.priority_flag = string(dst[i:i+2])
	i = i+2
	
	for ; i < 50; i = i+2 {
		tmp = string(dst[i:i+2])
		if tmp == "00" {
			break
		}
		tmp_sum += tmp
	}
	tmp_b = make([]byte, hex.DecodedLen(len([]byte(tmp_sum))))
	hex.Decode(tmp_b, []byte(tmp_sum))
	sms.delivery_time = string(tmp_b)
	tmp_sum = ""
	i = i+2
	
	for ; i < 50; i = i+2 {
		tmp = string(dst[i:i+2])
		if tmp == "00" {
			break
		}
		tmp_sum += tmp
	}
	tmp_b = make([]byte, hex.DecodedLen(len([]byte(tmp_sum))))
	hex.Decode(tmp_b, []byte(tmp_sum))
	sms.validity_period = string(tmp_b)
	tmp_sum = ""
	i = i+2

	sms.registered_delivery = string(dst[i:i+2])
	i = i+2

	sms.replace_if_present_flag = string(dst[i:i+2])
	i = i+2

	sms.data_coding = string(dst[i:i+2])
	i = i+2

	sms.sm_default_msg_id = string(dst[i:i+2])
	i = i+2

	sms.sm_length = hex2int(string(dst[i:i+2]))
	fmt.Fprintf(os.Stderr,"---%s\n", string(dst[i:i+2]))

	sms.sm = string(sms.body[i/2:(i/2)+int(sms.sm_length)+1])

	fmt.Fprintf(os.Stderr,"Service type: %s\n", sms.service_type)
	fmt.Fprintf(os.Stderr,"Source Addr Ton: %s (%s)\n", get_addr_ton(sms.source_addr_ton), sms.source_addr_ton)
	fmt.Fprintf(os.Stderr,"Source Addr NPI : %s (%s)\n", get_addr_npi(sms.source_addr_npi), sms.source_addr_npi)
	fmt.Fprintf(os.Stderr,"Source Addr: %s\n", sms.source_addr)
	fmt.Fprintf(os.Stderr,"Dest Addr Ton: %s (%s)\n", sms.dest_addr_ton, get_addr_ton(sms.dest_addr_ton))
	fmt.Fprintf(os.Stderr,"Dest Addr NPI: %s (%s)\n", sms.dest_addr_npi, get_addr_npi(sms.dest_addr_npi))
	fmt.Fprintf(os.Stderr,"Dest Addr: %s\n", sms.dest_addr)
	fmt.Fprintf(os.Stderr,"ESM Class: %s\n", sms.esm_class)
	fmt.Fprintf(os.Stderr,"Protocol ID: %s\n", sms.protocol_id)
	fmt.Fprintf(os.Stderr,"Priority flag: %s\n", sms.priority_flag)
	fmt.Fprintf(os.Stderr,"Delivery time: %s\n", sms.delivery_time)
	fmt.Fprintf(os.Stderr,"Validity Period: %s\n", sms.validity_period)
	fmt.Fprintf(os.Stderr,"Registered Delivery: %s\n", sms.registered_delivery)
	fmt.Fprintf(os.Stderr,"Replace if present: %s\n", sms.replace_if_present_flag)
	fmt.Fprintf(os.Stderr,"Data coding: %s (%s)\n", sms.data_coding, get_data_coding(sms.data_coding))
	fmt.Fprintf(os.Stderr,"SM Default MSG ID: %s\n", sms.sm_default_msg_id)
	fmt.Fprintf(os.Stderr,"SM Length: %d\n", sms.sm_length)
	fmt.Fprintf(os.Stderr,"POS %d , %d\n", i, i/2)
	fmt.Fprintf(os.Stderr,"SM : -%s-\n", sms.sm)
	
	fmt.Fprintf(os.Stderr, "%s\n", response)
	brsp := make([]byte, hex.DecodedLen(len(response)))
	hex.Decode(brsp,[]byte(response))
	fmt.Fprintf(os.Stderr, "-%s-\n", brsp)
	conn.Write(brsp)
	send_delivery_sm(conn, sms)
}

func send_delivery_sm (conn net.Conn, sms pdu) {
	sequence_number := get_sequence_number()
	command_id := fmt.Sprintf("%08X",getCommandHex("deliver_sm"))
	command_status := fmt.Sprintf("%08X",0)
	
	tmp_b := make([]byte, hex.EncodedLen(len([]byte(sms.service_type))))
	hex.Encode(tmp_b, []byte(sms.service_type))
	service_type := string(tmp_b) + "00"

	tmp_b = make([]byte, hex.EncodedLen(len([]byte(sms.source_addr))))
	hex.Encode(tmp_b, []byte(sms.source_addr))
	source_addr := string(tmp_b) + "00"

	tmp_b = make([]byte, hex.EncodedLen(len([]byte(sms.dest_addr))))
	hex.Encode(tmp_b, []byte(sms.dest_addr))
	dest_addr := string(tmp_b) + "00"

	delivery_time := "00"
	validity_period := "00"
	replace_if_present := "00"
	sm_default_msg_id := "00"
	sm_length := fmt.Sprintf("%02X",sms.sm_length)

	tmp_b = make([]byte, hex.EncodedLen(len([]byte(sms.sm))))
	hex.Encode(tmp_b, []byte(sms.sm))
	sm := string(tmp_b)

	var fr string
	
	fmt.Fprintf(os.Stderr,"Send Delivery SM\n")
	fmt.Fprintf(os.Stderr,"Sequence number: %s\n", sequence_number)
	fr = sequence_number
	fmt.Fprintf(os.Stderr,"Command ID: %s\n", command_id)
	fr += command_id
	fmt.Fprintf(os.Stderr,"Command status: %s\n", command_status)
	fr += command_status
	fmt.Fprintf(os.Stderr,"Service type: %s\n", service_type)
	fr += service_type
	fmt.Fprintf(os.Stderr,"Source Addr Ton: %s\n", sms.source_addr_ton)
	fr += sms.source_addr_ton
	fmt.Fprintf(os.Stderr,"Source Addr NPI: %s\n", sms.source_addr_npi)
	fr += sms.source_addr_npi
	fmt.Fprintf(os.Stderr,"Source Addr: %s\n", source_addr)
	fr += source_addr
	fmt.Fprintf(os.Stderr,"Dest Addr Ton: %s\n", sms.dest_addr_ton)
	fr += sms.dest_addr_ton
	fmt.Fprintf(os.Stderr,"Dest Addr NPI: %s\n", sms.dest_addr_npi)
	fr += sms.dest_addr_npi
	fmt.Fprintf(os.Stderr,"Dest Addr: %s\n", dest_addr)
	fr += dest_addr
	fmt.Fprintf(os.Stderr,"ESM Class: %s\n", sms.esm_class)
	fr += sms.esm_class
	fmt.Fprintf(os.Stderr,"Protocol ID: %s\n", sms.protocol_id)
	fr += sms.protocol_id
	fmt.Fprintf(os.Stderr,"Priority Flag: %s\n", sms.priority_flag)
	fr += sms.priority_flag
	fmt.Fprintf(os.Stderr,"Delivery Time: %s\n", delivery_time)
	fr += delivery_time
	fmt.Fprintf(os.Stderr,"Validity Period: %s\n", validity_period)
	fr += validity_period
	fmt.Fprintf(os.Stderr,"Registered Delivery: %s\n", sms.registered_delivery)
	fr += sms.registered_delivery
	fmt.Fprintf(os.Stderr,"Replace if present: %s\n", replace_if_present)
	fr += replace_if_present
	fmt.Fprintf(os.Stderr,"Data Coding: %s\n", sms.data_coding)
	fr += sms.data_coding
	fmt.Fprintf(os.Stderr,"SM Default Msg Id: %s\n", sm_default_msg_id)
	fr += sm_default_msg_id
	fmt.Fprintf(os.Stderr,"SM Length: %s\n", sm_length)
	fr += sm_length
	fmt.Fprintf(os.Stderr,"SM: %s\n", sm)
	fr += sm
	fmt.Fprintf(os.Stderr,"Frame: %s\n", fr)
	fr_length := fmt.Sprintf("%08X", (len(fr) + 8) /2)
	fmt.Fprintf(os.Stderr,"Frame length: %s\n", fr_length)
	fr = fr_length + fr

	brsp := make([]byte, hex.DecodedLen(len(fr)))
	hex.Decode(brsp,[]byte(fr))
	conn.Write(brsp)

}

func get_pdu_size (conn net.Conn) uint64 {
	buf := make([]byte, 4)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading on socket: %v\n", err)
		return 0
	}
	dst := make([]byte, hex.EncodedLen(len(buf)))
	hex.Encode(dst,buf)
	size_hex := string(dst)[0:8]
	return hex2int(size_hex)
}

func get_pdu (conn net.Conn, size uint64) ([]byte,[]byte) {
	buf := make([]byte, size -4)

	_, err := conn.Read(buf)
	
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from socket: %v\n", err)
		return []byte("Error"), []byte("Error")
	}

	dst := make([]byte, hex.EncodedLen(len(buf)))
	hex.Encode(dst, buf)

	return dst, buf
}
	

func receive(conn net.Conn) {
	defer conn.Close()
	count_err := 0

	for {
		var sms pdu

		fmt.Fprintf(os.Stderr,"Next input\n")

		sms.size = get_pdu_size(conn)
		
		fmt.Fprintf(os.Stderr, "Size: %d\n", sms.size)

		if sms.size == 0 || sms.size > 1024 {
			count_err++
			if count_err > 20 {
				break
			}
			continue
		}
		count_err = 0

		dst, buf := get_pdu(conn, sms.size)
		if string(dst) == "Error" {
			break
		}
		
		sms.sequence_number = string(dst)[16:24]
		sms.command_status = string(dst)[8:16]
		sms.command_id = getCommandName(string(dst)[0:8])
		sms.pdu = buf
		sms.body = string(buf[12:])

		fmt.Fprintf(os.Stderr, "Command: %s\n", sms.command_id)
		fmt.Fprintf(os.Stderr, "Status: %s\n", sms.command_status)
		fmt.Fprintf(os.Stderr, "Sequence: %s\n", sms.sequence_number)

		if sms.command_id == "bind_transceiver" {
			go bind_transceiver(conn,sms)
		}

		if sms.command_id == "enquire_link" {
			go enquire_link(conn, sms)
		}

		if sms.command_id == "submit_sm" {
			go submit_sm(conn, sms)
		}

	}
}



