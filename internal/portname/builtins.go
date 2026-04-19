package portname

type entry struct {
	port  int
	proto string
	name  string
}

var builtins = []entry{
	{20, "tcp", "ftp-data"},
	{21, "tcp", "ftp"},
	{22, "tcp", "ssh"},
	{23, "tcp", "telnet"},
	{25, "tcp", "smtp"},
	{53, "tcp", "dns"},
	{53, "udp", "dns"},
	{80, "tcp", "http"},
	{110, "tcp", "pop3"},
	{143, "tcp", "imap"},
	{443, "tcp", "https"},
	{465, "tcp", "smtps"},
	{587, "tcp", "submission"},
	{993, "tcp", "imaps"},
	{995, "tcp", "pop3s"},
	{3306, "tcp", "mysql"},
	{5432, "tcp", "postgresql"},
	{6379, "tcp", "redis"},
	{8080, "tcp", "http-alt"},
	{8443, "tcp", "https-alt"},
	{27017, "tcp", "mongodb"},
}
