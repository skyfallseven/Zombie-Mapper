/*
File: 		zm.c
Author: 	Jeremiah Faison

Descript: 	Live packet capture analysis/report generation
			Also analyzes pcap files post capture
			Detects SSH Bruteforcing 
			Logs IP addresses that attempt SYN attacks
			Reports to AbuseIPDB
			DNS info and geolocation on all attackers

Info:		compile with gcc zm.c -lpcap
Version:	0.001-Alpha (ain't do nuffin)
*/

#include <stdio.h>
#include <time.h>
#include <pcap/pcap.h>
#include <netinet/in.h>
#include <netinet/if_ether.h>
#include <string.h>
void print_pkt_info(const u_char *packet, struct pcap_pkthdr pkt_header);

int main(int argc, char *argv[]) {
	char *device; //Name of device; wlan0, eth0
	char ip[13]; 
	bpf_u_int32 ip_raw; //IP addr as integer
	bpf_u_int32 subnet_raw;
	int lookup_return_code;
	char err_buff[PCAP_ERRBUF_SIZE]; //Size defined in pcap.h
	struct in_addr address; //Used for IP

	pcap_t *handle;
	const u_char *packet;
	struct pcap_pkthdr packet_header;
	int pkt_count = 1;
	int timeout = 10000; //In milliseconds

	//Find a device
	device = pcap_lookupdev(err_buff);
	if (device == NULL) {
		printf("Error opening device: %s\n", err_buff);
		return 1;
	}

	//Get device info
	lookup_return_code = pcap_lookupnet(device, &ip_raw, &subnet_raw, err_buff);
	if (lookup_return_code == -1) {
		printf("%s\n", err_buff);
		return 1;
	}

	//Get IP in human form
	address.s_addr = ip_raw;
	strncpy(ip, inet_ntoa(address), 13);
	if (ip == NULL) {
		perror("inet_ntoa");
		return 1;
	}
	
	//Start live capture
	handle = pcap_open_live(device, BUFSIZ, pkt_count, timeout, err_buff);
	packet = pcap_next(handle, &packet_header);
	if (packet == NULL) {
		printf("No packet found.\n");
		return 2;
	}	
	printf("Zombie Mapper -- 0.001 Alpha\n");
	printf("Interface: %s\n", device);
	printf("IP address: %s\n", ip);

	//Output some info
	print_pkt_info(packet, packet_header);
	
	return 0;
}

void print_pkt_info(const u_char *packet, struct pcap_pkthdr pkt_header)
{
	printf("Packet capture length: %d\n", pkt_header.caplen);
	printf("Packet total length %d\n", pkt_header.len);
}
