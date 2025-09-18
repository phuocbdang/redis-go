package core

import (
	"syscall"
)

func ExecuteAndResponse(cmd *Command, connFd int) error {
	var res []byte
	switch cmd.Cmd {
	case "PING":
		res = cmdPING(cmd.Args)
	/* Dict */
	case "SET":
		res = cmdSET(cmd.Args)
	case "GET":
		res = cmdGET(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	case "INFO":
		res = cmdINFO(cmd.Args)
	/* Set */
	case "SADD":
		res = cmdSADD(cmd.Args)
	case "SREM":
		res = cmdSREM(cmd.Args)
	case "SMEMBERS":
		res = cmdSMEMBERS(cmd.Args)
	case "SISMEMBER":
		res = cmdSISMEMBER(cmd.Args)
	/* Sorted Set */
	case "ZADD":
		res = cmdZADD(cmd.Args)
	case "ZSCORE":
		res = cmdZSCORE(cmd.Args)
	case "ZRANK":
		res = cmdZRANK(cmd.Args)
	/* Count-min Sketch */
	case "CMS.INITBYDIM":
		res = cmdCMSINITBYDIM(cmd.Args)
	case "CMS.INITBYPROB":
		res = cmdCMSINITBYPROB(cmd.Args)
	case "CMS.INCRBY":
		res = cmdCMSINCRBY(cmd.Args)
	case "CMS.QUERY":
		res = cmdCMSQUERY(cmd.Args)
	/* Bloom filter */
	case "BF.RESERVE":
		res = cmdBFRESERVE(cmd.Args)
	case "BF.MADD":
		res = cmdBFMADD(cmd.Args)
	case "BF.EXISTS":
		res = cmdBFEXISTS(cmd.Args)
	/* Default case */
	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}
	_, err := syscall.Write(connFd, res)
	return err
}
