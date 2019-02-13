package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
	"strings"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron"
	"github.com/jamiealquiza/envy"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
var instance = flag.String("instance", "localhost:5060", "Instance")
var every = flag.String("every", "1m", "Update time")
var myClient = &http.Client{Timeout: 10 * time.Second}

var (
	CoreBadUrisRcvd = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_bad_uris_rcvd",
	    Help: "kam_core_bad_uris_rcvd",
	  },
	  []string{"instance"},
	)
	CoreBadMsgHdr = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_bad_msg_hdr",
	    Help: "kam_core_bad_msg_hdr",
	  },
	  []string{"instance"},
	)
	CoreDropReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_drop_reply",
	    Help: "kam_core_drop_reply",
	  },
	  []string{"instance"},
	)
	CoreDropRequest = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_drop_request",
	    Help: "kam_core_drop_request",
	  },
	  []string{"instance"},
	)
	CoreErrReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_err_reply",
	    Help: "kam_core_err_reply",
	  },
	  []string{"instance"},
	)
	CoreErrRequest = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_err_request",
	    Help: "kam_core_err_request",
	  },
	  []string{"instance"},
	)
	CoreFwdReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_fwd_reply",
	    Help: "kam_core_fwd_reply",
	  },
	  []string{"instance"},
	)
	CoreFwdRequest = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_fwd_request",
	    Help: "kam_core_fwd_request",
	  },
	  []string{"instance"},
	)
	CoreRcvReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_reply",
	    Help: "kam_core_rcv_reply",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies18x = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies18x",
	    Help: "kam_core_rcv_replies18x",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies1xx = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies1xx",
	    Help: "kam_core_rcv_replies1xx",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies2xx = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies2xx",
	    Help: "kam_core_rcv_replies2xx",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies3xx = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies3xx",
	    Help: "kam_core_rcv_replies3xx",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies401 = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies401",
	    Help: "kam_core_rcv_replies401",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies404 = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies404",
	    Help: "kam_core_rcv_replies404",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies407 = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies407",
	    Help: "kam_core_rcv_replies407",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies480 = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies480",
	    Help: "kam_core_rcv_replies480",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies486 = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies486",
	    Help: "kam_core_rcv_replies486",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies4xx = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies4xx",
	    Help: "kam_core_rcv_replies4xx",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies5xx = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies5xx",
	    Help: "kam_core_rcv_replies5xx",
	  },
	  []string{"instance"},
	)
	CoreRcvReplies6xx = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_replies6xx",
	    Help: "kam_core_rcv_replies6xx",
	  },
	  []string{"instance"},
	)
	CoreRcvRequest = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_request",
	    Help: "kam_core_rcv_request",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsAck = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_ack",
	    Help: "kam_core_rcv_requests_ack",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsBye = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_bye",
	    Help: "kam_core_rcv_requests_bye",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsCancel = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_cancel",
	    Help: "kam_core_rcv_requests_cancel",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsInfo = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_info",
	    Help: "kam_core_rcv_requests_info",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsInvite = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_invite",
	    Help: "kam_core_rcv_requests_invite",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsMessage = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_message",
	    Help: "kam_core_rcv_requests_message",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsNotify = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_notify",
	    Help: "kam_core_rcv_requests_notify",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsOption = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_option",
	    Help: "kam_core_rcv_requests_option",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsPrack = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_prack",
	    Help: "kam_core_rcv_requests_prack",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsPublish = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_publish",
	    Help: "kam_core_rcv_requests_publish",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsRefer = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_refer",
	    Help: "kam_core_rcv_requests_refer",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsRegister = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_register",
	    Help: "kam_core_rcv_requests_register",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsSubscribe = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_subscribe",
	    Help: "kam_core_rcv_requests_subscribe",
	  },
	  []string{"instance"},
	)
	CoreRcvRequestsUpdate = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_rcv_requests_update",
	    Help: "kam_core_rcv_requests_update",
	  },
	  []string{"instance"},
	)
	CoreUnsupportedMethod = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_core_unsupported_method",
	    Help: "kam_core_unsupported_method",
	  },
	  []string{"instance"},
	)
	DialogActiveDialog = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_dialog_active_dialog",
	    Help: "kam_dialog_active_dialog",
	  },
	  []string{"instance"},
	)
	DialogEarlyDialog = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_dialog_early_dialog",
	    Help: "kam_dialog_early_dialog",
	  },
	  []string{"instance"},
	)
	DialogExpiredDialog = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_dialog_expired_dialog",
	    Help: "kam_dialog_expired_dialog",
	  },
	  []string{"instance"},
	)
	DialogFailedDialog = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_dialog_failed_dialog",
	    Help: "kam_dialog_failed_dialog",
	  },
	  []string{"instance"},
	)
	DialogProcessedDialog = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_dialog_processed_dialog",
	    Help: "kam_dialog_processed_dialog",
	  },
	  []string{"instance"},
	)
	DnsFailedDnsRequest = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_dns_failed_dns_request",
	    Help: "kam_dns_failed_dns_request",
	  },
	  []string{"instance"},
	)
	HttpclientConnection = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_httpclient_connection",
	    Help: "kam_httpclient_connection",
	  },
	  []string{"instance"},
	)
	HttpclientConnfail = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_httpclient_connfail",
	    Help: "kam_httpclient_connfail",
	  },
	  []string{"instance"},
	)
	HttpclientConnok = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_httpclient_connok",
	    Help: "kam_httpclient_connok",
	  },
	  []string{"instance"},
	)
	MysqlDriverError = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_mysql_driver_error",
	    Help: "kam_mysql_driver_error",
	  },
	  []string{"instance"},
	)
	RegistrarAcceptedReg = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_registrar_accepted_reg",
	    Help: "kam_registrar_accepted_reg",
	  },
	  []string{"instance"},
	)
	RegistrarDefaultExpire = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_registrar_default_expire",
	    Help: "kam_registrar_default_expire",
	  },
	  []string{"instance"},
	)
	RegistrarDefaultExpiresRange = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_registrar_default_expires_range",
	    Help: "kam_registrar_default_expires_range",
	  },
	  []string{"instance"},
	)
	RegistrarExpiresRange = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_registrar_expires_range",
	    Help: "kam_registrar_expires_range",
	  },
	  []string{"instance"},
	)
	RegistrarMaxContact = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_registrar_max_contact",
	    Help: "kam_registrar_max_contact",
	  },
	  []string{"instance"},
	)
	RegistrarMaxExpire = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_registrar_max_expire",
	    Help: "kam_registrar_max_expire",
	  },
	  []string{"instance"},
	)
	RegistrarRejectedReg = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_registrar_rejected_reg",
	    Help: "kam_registrar_rejected_reg",
	  },
	  []string{"instance"},
	)
	ShmemFragment = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_shmem_fragment",
	    Help: "kam_shmem_fragment",
	  },
	  []string{"instance"},
	)
	ShmemFreeSize = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_shmem_free_size",
	    Help: "kam_shmem_free_size",
	  },
	  []string{"instance"},
	)
	ShmemMaxUsedSize = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_shmem_max_used_size",
	    Help: "kam_shmem_max_used_size",
	  },
	  []string{"instance"},
	)
	ShmemRealUsedSize = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_shmem_real_used_size",
	    Help: "kam_shmem_real_used_size",
	  },
	  []string{"instance"},
	)
	ShmemTotalSize = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_shmem_total_size",
	    Help: "kam_shmem_total_size",
	  },
	  []string{"instance"},
	)
	ShmemUsedSize = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_shmem_used_size",
	    Help: "kam_shmem_used_size",
	  },
	  []string{"instance"},
	)
	Sl1xxReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl1xx_reply",
	    Help: "kam_sl1xx_reply",
	  },
	  []string{"instance"},
	)
	Sl200Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl200_reply",
	    Help: "kam_sl200_reply",
	  },
	  []string{"instance"},
	)
	Sl202Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl202_reply",
	    Help: "kam_sl202_reply",
	  },
	  []string{"instance"},
	)
	Sl2xxReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl2xx_reply",
	    Help: "kam_sl2xx_reply",
	  },
	  []string{"instance"},
	)
	Sl300Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl300_reply",
	    Help: "kam_sl300_reply",
	  },
	  []string{"instance"},
	)
	Sl301Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl301_reply",
	    Help: "kam_sl301_reply",
	  },
	  []string{"instance"},
	)
	Sl302Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl302_reply",
	    Help: "kam_sl302_reply",
	  },
	  []string{"instance"},
	)
	Sl3xxReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl3xx_reply",
	    Help: "kam_sl3xx_reply",
	  },
	  []string{"instance"},
	)
	Sl400Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl400_reply",
	    Help: "kam_sl400_reply",
	  },
	  []string{"instance"},
	)
	Sl401Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl401_reply",
	    Help: "kam_sl401_reply",
	  },
	  []string{"instance"},
	)
	Sl403Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl403_reply",
	    Help: "kam_sl403_reply",
	  },
	  []string{"instance"},
	)
	Sl404Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl404_reply",
	    Help: "kam_sl404_reply",
	  },
	  []string{"instance"},
	)
	Sl407Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl407_reply",
	    Help: "kam_sl407_reply",
	  },
	  []string{"instance"},
	)
	Sl408Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl408_reply",
	    Help: "kam_sl408_reply",
	  },
	  []string{"instance"},
	)
	Sl483Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl483_reply",
	    Help: "kam_sl483_reply",
	  },
	  []string{"instance"},
	)
	Sl4xxReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl4xx_reply",
	    Help: "kam_sl4xx_reply",
	  },
	  []string{"instance"},
	)
	Sl500Reply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl500_reply",
	    Help: "kam_sl500_reply",
	  },
	  []string{"instance"},
	)
	Sl5xxReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl5xx_reply",
	    Help: "kam_sl5xx_reply",
	  },
	  []string{"instance"},
	)
	Sl6xxReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl6xx_reply",
	    Help: "kam_sl6xx_reply",
	  },
	  []string{"instance"},
	)
	SlFailure = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl_failure",
	    Help: "kam_sl_failure",
	  },
	  []string{"instance"},
	)
	SlReceivedAck = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl_received_ack",
	    Help: "kam_sl_received_ack",
	  },
	  []string{"instance"},
	)
	SlSentErrReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl_sent_err_reply",
	    Help: "kam_sl_sent_err_reply",
	  },
	  []string{"instance"},
	)
	SlSentReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl_sent_reply",
	    Help: "kam_sl_sent_reply",
	  },
	  []string{"instance"},
	)
	SlXxxReply = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_sl_xxx_reply",
	    Help: "kam_sl_xxx_reply",
	  },
	  []string{"instance"},
	)
	TcpConReset = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tcp_con_reset",
	    Help: "kam_tcp_con_reset",
	  },
	  []string{"instance"},
	)
	TcpConTimeout = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tcp_con_timeout",
	    Help: "kam_tcp_con_timeout",
	  },
	  []string{"instance"},
	)
	TcpConnectFailed = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tcp_connect_failed",
	    Help: "kam_tcp_connect_failed",
	  },
	  []string{"instance"},
	)
	TcpConnectSuccess = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tcp_connect_success",
	    Help: "kam_tcp_connect_success",
	  },
	  []string{"instance"},
	)
	TcpCurrentOpenedConnection = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tcp_current_opened_connection",
	    Help: "kam_tcp_current_opened_connection",
	  },
	  []string{"instance"},
	)
	TcpCurrentWriteQueueSize = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tcp_current_write_queue_size",
	    Help: "kam_tcp_current_write_queue_size",
	  },
	  []string{"instance"},
	)
	TcpEstablished = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tcp_established",
	    Help: "kam_tcp_established",
	  },
	  []string{"instance"},
	)
	TcpLocalReject = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tcp_local_reject",
	    Help: "kam_tcp_local_reject",
	  },
	  []string{"instance"},
	)
	TcpPassiveOpen = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tcp_passive_open",
	    Help: "kam_tcp_passive_open",
	  },
	  []string{"instance"},
	)
	TcpSendTimeout = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tcp_send_timeout",
	    Help: "kam_tcp_send_timeout",
	  },
	  []string{"instance"},
	)
	TcpSendqFull = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tcp_sendq_full",
	    Help: "kam_tcp_sendq_full",
	  },
	  []string{"instance"},
	)
	Tmx2xxTransaction = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx2xx_transaction",
	    Help: "kam_tmx2xx_transaction",
	  },
	  []string{"instance"},
	)
	Tmx3xxTransaction = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx3xx_transaction",
	    Help: "kam_tmx3xx_transaction",
	  },
	  []string{"instance"},
	)
	Tmx4xxTransaction = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx4xx_transaction",
	    Help: "kam_tmx4xx_transaction",
	  },
	  []string{"instance"},
	)
	Tmx5xxTransaction = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx5xx_transaction",
	    Help: "kam_tmx5xx_transaction",
	  },
	  []string{"instance"},
	)
	Tmx6xxTransaction = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx6xx_transaction",
	    Help: "kam_tmx6xx_transaction",
	  },
	  []string{"instance"},
	)
	TmxUacTransaction = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx_uac_transaction",
	    Help: "kam_tmx_uac_transaction",
	  },
	  []string{"instance"},
	)
	TmxUasTransaction = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx_uas_transaction",
	    Help: "kam_tmx_uas_transaction",
	  },
	  []string{"instance"},
	)
	TmxActiveTransaction = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx_active_transaction",
	    Help: "kam_tmx_active_transaction",
	  },
	  []string{"instance"},
	)
	TmxInuseTransaction = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx_inuse_transaction",
	    Help: "kam_tmx_inuse_transaction",
	  },
	  []string{"instance"},
	)
	TmxRplAbsorbed = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx_rpl_absorbed",
	    Help: "kam_tmx_rpl_absorbed",
	  },
	  []string{"instance"},
	)
	TmxRplGenerated = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx_rpl_generated",
	    Help: "kam_tmx_rpl_generated",
	  },
	  []string{"instance"},
	)
	TmxRplReceived = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx_rpl_received",
	    Help: "kam_tmx_rpl_received",
	  },
	  []string{"instance"},
	)
	TmxRplRelayed = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx_rpl_relayed",
	    Help: "kam_tmx_rpl_relayed",
	  },
	  []string{"instance"},
	)
	TmxRplSent = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_tmx_rpl_sent",
	    Help: "kam_tmx_rpl_sent",
	  },
	  []string{"instance"},
	)
	UsrlocLocationContact = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_usrloc_location_contact",
	    Help: "kam_usrloc_location_contact",
	  },
	  []string{"instance"},
	)
	UsrlocLocationExpire = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_usrloc_location_expire",
	    Help: "kam_usrloc_location_expire",
	  },
	  []string{"instance"},
	)
	UsrlocLocationUser = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_usrloc_location_user",
	    Help: "kam_usrloc_location_user",
	  },
	  []string{"instance"},
	)
	UsrlocRegisteredUser = prometheus.NewGaugeVec(
	  prometheus.GaugeOpts{
	    Name: "kam_usrloc_registered_user",
	    Help: "kam_usrloc_registered_user",
	  },
	  []string{"instance"},
	)
)

var measurementMapping = map[string] *prometheus.GaugeVec {
	"core:bad_URIs_rcvd": CoreBadUrisRcvd,
	"core:bad_msg_hdr": CoreBadMsgHdr,
	"core:drop_replies": CoreDropReply,
	"core:drop_requests": CoreDropRequest,
	"core:err_replies": CoreErrReply,
	"core:err_requests": CoreErrRequest,
	"core:fwd_replies": CoreFwdReply,
	"core:fwd_requests": CoreFwdRequest,
	"core:rcv_replies": CoreRcvReply,
	"core:rcv_replies_18x": CoreRcvReplies18x,
	"core:rcv_replies_1xx": CoreRcvReplies1xx,
	"core:rcv_replies_2xx": CoreRcvReplies2xx,
	"core:rcv_replies_3xx": CoreRcvReplies3xx,
	"core:rcv_replies_401": CoreRcvReplies401,
	"core:rcv_replies_404": CoreRcvReplies404,
	"core:rcv_replies_407": CoreRcvReplies407,
	"core:rcv_replies_480": CoreRcvReplies480,
	"core:rcv_replies_486": CoreRcvReplies486,
	"core:rcv_replies_4xx": CoreRcvReplies4xx,
	"core:rcv_replies_5xx": CoreRcvReplies5xx,
	"core:rcv_replies_6xx": CoreRcvReplies6xx,
	"core:rcv_requests": CoreRcvRequest,
	"core:rcv_requests_ack": CoreRcvRequestsAck,
	"core:rcv_requests_bye": CoreRcvRequestsBye,
	"core:rcv_requests_cancel": CoreRcvRequestsCancel,
	"core:rcv_requests_info": CoreRcvRequestsInfo,
	"core:rcv_requests_invite": CoreRcvRequestsInvite,
	"core:rcv_requests_message": CoreRcvRequestsMessage,
	"core:rcv_requests_notify": CoreRcvRequestsNotify,
	"core:rcv_requests_options": CoreRcvRequestsOption,
	"core:rcv_requests_prack": CoreRcvRequestsPrack,
	"core:rcv_requests_publish": CoreRcvRequestsPublish,
	"core:rcv_requests_refer": CoreRcvRequestsRefer,
	"core:rcv_requests_register": CoreRcvRequestsRegister,
	"core:rcv_requests_subscribe": CoreRcvRequestsSubscribe,
	"core:rcv_requests_update": CoreRcvRequestsUpdate,
	"core:unsupported_methods": CoreUnsupportedMethod,
	"dialog:active_dialogs": DialogActiveDialog,
	"dialog:early_dialogs": DialogEarlyDialog,
	"dialog:expired_dialogs": DialogExpiredDialog,
	"dialog:failed_dialogs": DialogFailedDialog,
	"dialog:processed_dialogs": DialogProcessedDialog,
	"dns:failed_dns_request": DnsFailedDnsRequest,
	"httpclient:connections": HttpclientConnection,
	"httpclient:connfail": HttpclientConnfail,
	"httpclient:connok": HttpclientConnok,
	"mysql:driver_errors": MysqlDriverError,
	"registrar:accepted_regs": RegistrarAcceptedReg,
	"registrar:default_expire": RegistrarDefaultExpire,
	"registrar:default_expires_range": RegistrarDefaultExpiresRange,
	"registrar:expires_range": RegistrarExpiresRange,
	"registrar:max_contacts": RegistrarMaxContact,
	"registrar:max_expires": RegistrarMaxExpire,
	"registrar:rejected_regs": RegistrarRejectedReg,
	"shmem:fragments": ShmemFragment,
	"shmem:free_size": ShmemFreeSize,
	"shmem:max_used_size": ShmemMaxUsedSize,
	"shmem:real_used_size": ShmemRealUsedSize,
	"shmem:total_size": ShmemTotalSize,
	"shmem:used_size": ShmemUsedSize,
	"sl:1xx_replies": Sl1xxReply,
	"sl:200_replies": Sl200Reply,
	"sl:202_replies": Sl202Reply,
	"sl:2xx_replies": Sl2xxReply,
	"sl:300_replies": Sl300Reply,
	"sl:301_replies": Sl301Reply,
	"sl:302_replies": Sl302Reply,
	"sl:3xx_replies": Sl3xxReply,
	"sl:400_replies": Sl400Reply,
	"sl:401_replies": Sl401Reply,
	"sl:403_replies": Sl403Reply,
	"sl:404_replies": Sl404Reply,
	"sl:407_replies": Sl407Reply,
	"sl:408_replies": Sl408Reply,
	"sl:483_replies": Sl483Reply,
	"sl:4xx_replies": Sl4xxReply,
	"sl:500_replies": Sl500Reply,
	"sl:5xx_replies": Sl5xxReply,
	"sl:6xx_replies": Sl6xxReply,
	"sl:failures": SlFailure,
	"sl:received_ACKs": SlReceivedAck,
	"sl:sent_err_replies": SlSentErrReply,
	"sl:sent_replies": SlSentReply,
	"sl:xxx_replies": SlXxxReply,
	"tcp:con_reset": TcpConReset,
	"tcp:con_timeout": TcpConTimeout,
	"tcp:connect_failed": TcpConnectFailed,
	"tcp:connect_success": TcpConnectSuccess,
	"tcp:current_opened_connections": TcpCurrentOpenedConnection,
	"tcp:current_write_queue_size": TcpCurrentWriteQueueSize,
	"tcp:established": TcpEstablished,
	"tcp:local_reject": TcpLocalReject,
	"tcp:passive_open": TcpPassiveOpen,
	"tcp:send_timeout": TcpSendTimeout,
	"tcp:sendq_full": TcpSendqFull,
	"tmx:2xx_transactions": Tmx2xxTransaction,
	"tmx:3xx_transactions": Tmx3xxTransaction,
	"tmx:4xx_transactions": Tmx4xxTransaction,
	"tmx:5xx_transactions": Tmx5xxTransaction,
	"tmx:6xx_transactions": Tmx6xxTransaction,
	"tmx:UAC_transactions": TmxUacTransaction,
	"tmx:UAS_transactions": TmxUasTransaction,
	"tmx:active_transactions": TmxActiveTransaction,
	"tmx:inuse_transactions": TmxInuseTransaction,
	"tmx:rpl_absorbed": TmxRplAbsorbed,
	"tmx:rpl_generated": TmxRplGenerated,
	"tmx:rpl_received": TmxRplReceived,
	"tmx:rpl_relayed": TmxRplRelayed,
	"tmx:rpl_sent": TmxRplSent,
	"usrloc:location-contacts": UsrlocLocationContact,
	"usrloc:location-expires": UsrlocLocationExpire,
	"usrloc:location-users": UsrlocLocationUser,
	"usrloc:registered_users": UsrlocRegisteredUser,
}

type Measurement struct {
	Jsonrpc  string
	Result 	 []string
}

func init() {
	prometheus.MustRegister(CoreBadUrisRcvd)
	prometheus.MustRegister(CoreBadMsgHdr)
	prometheus.MustRegister(CoreDropReply)
	prometheus.MustRegister(CoreDropRequest)
	prometheus.MustRegister(CoreErrReply)
	prometheus.MustRegister(CoreErrRequest)
	prometheus.MustRegister(CoreFwdReply)
	prometheus.MustRegister(CoreFwdRequest)
	prometheus.MustRegister(CoreRcvReply)
	prometheus.MustRegister(CoreRcvReplies18x)
	prometheus.MustRegister(CoreRcvReplies1xx)
	prometheus.MustRegister(CoreRcvReplies2xx)
	prometheus.MustRegister(CoreRcvReplies3xx)
	prometheus.MustRegister(CoreRcvReplies401)
	prometheus.MustRegister(CoreRcvReplies404)
	prometheus.MustRegister(CoreRcvReplies407)
	prometheus.MustRegister(CoreRcvReplies480)
	prometheus.MustRegister(CoreRcvReplies486)
	prometheus.MustRegister(CoreRcvReplies4xx)
	prometheus.MustRegister(CoreRcvReplies5xx)
	prometheus.MustRegister(CoreRcvReplies6xx)
	prometheus.MustRegister(CoreRcvRequest)
	prometheus.MustRegister(CoreRcvRequestsAck)
	prometheus.MustRegister(CoreRcvRequestsBye)
	prometheus.MustRegister(CoreRcvRequestsCancel)
	prometheus.MustRegister(CoreRcvRequestsInfo)
	prometheus.MustRegister(CoreRcvRequestsInvite)
	prometheus.MustRegister(CoreRcvRequestsMessage)
	prometheus.MustRegister(CoreRcvRequestsNotify)
	prometheus.MustRegister(CoreRcvRequestsOption)
	prometheus.MustRegister(CoreRcvRequestsPrack)
	prometheus.MustRegister(CoreRcvRequestsPublish)
	prometheus.MustRegister(CoreRcvRequestsRefer)
	prometheus.MustRegister(CoreRcvRequestsRegister)
	prometheus.MustRegister(CoreRcvRequestsSubscribe)
	prometheus.MustRegister(CoreRcvRequestsUpdate)
	prometheus.MustRegister(CoreUnsupportedMethod)
	prometheus.MustRegister(DialogActiveDialog)
	prometheus.MustRegister(DialogEarlyDialog)
	prometheus.MustRegister(DialogExpiredDialog)
	prometheus.MustRegister(DialogFailedDialog)
	prometheus.MustRegister(DialogProcessedDialog)
	prometheus.MustRegister(DnsFailedDnsRequest)
	prometheus.MustRegister(HttpclientConnection)
	prometheus.MustRegister(HttpclientConnfail)
	prometheus.MustRegister(HttpclientConnok)
	prometheus.MustRegister(MysqlDriverError)
	prometheus.MustRegister(RegistrarAcceptedReg)
	prometheus.MustRegister(RegistrarDefaultExpire)
	prometheus.MustRegister(RegistrarDefaultExpiresRange)
	prometheus.MustRegister(RegistrarExpiresRange)
	prometheus.MustRegister(RegistrarMaxContact)
	prometheus.MustRegister(RegistrarMaxExpire)
	prometheus.MustRegister(RegistrarRejectedReg)
	prometheus.MustRegister(ShmemFragment)
	prometheus.MustRegister(ShmemFreeSize)
	prometheus.MustRegister(ShmemMaxUsedSize)
	prometheus.MustRegister(ShmemRealUsedSize)
	prometheus.MustRegister(ShmemTotalSize)
	prometheus.MustRegister(ShmemUsedSize)
	prometheus.MustRegister(Sl1xxReply)
	prometheus.MustRegister(Sl200Reply)
	prometheus.MustRegister(Sl202Reply)
	prometheus.MustRegister(Sl2xxReply)
	prometheus.MustRegister(Sl300Reply)
	prometheus.MustRegister(Sl301Reply)
	prometheus.MustRegister(Sl302Reply)
	prometheus.MustRegister(Sl3xxReply)
	prometheus.MustRegister(Sl400Reply)
	prometheus.MustRegister(Sl401Reply)
	prometheus.MustRegister(Sl403Reply)
	prometheus.MustRegister(Sl404Reply)
	prometheus.MustRegister(Sl407Reply)
	prometheus.MustRegister(Sl408Reply)
	prometheus.MustRegister(Sl483Reply)
	prometheus.MustRegister(Sl4xxReply)
	prometheus.MustRegister(Sl500Reply)
	prometheus.MustRegister(Sl5xxReply)
	prometheus.MustRegister(Sl6xxReply)
	prometheus.MustRegister(SlFailure)
	prometheus.MustRegister(SlReceivedAck)
	prometheus.MustRegister(SlSentErrReply)
	prometheus.MustRegister(SlSentReply)
	prometheus.MustRegister(SlXxxReply)
	prometheus.MustRegister(TcpConReset)
	prometheus.MustRegister(TcpConTimeout)
	prometheus.MustRegister(TcpConnectFailed)
	prometheus.MustRegister(TcpConnectSuccess)
	prometheus.MustRegister(TcpCurrentOpenedConnection)
	prometheus.MustRegister(TcpCurrentWriteQueueSize)
	prometheus.MustRegister(TcpEstablished)
	prometheus.MustRegister(TcpLocalReject)
	prometheus.MustRegister(TcpPassiveOpen)
	prometheus.MustRegister(TcpSendTimeout)
	prometheus.MustRegister(TcpSendqFull)
	prometheus.MustRegister(Tmx2xxTransaction)
	prometheus.MustRegister(Tmx3xxTransaction)
	prometheus.MustRegister(Tmx4xxTransaction)
	prometheus.MustRegister(Tmx5xxTransaction)
	prometheus.MustRegister(Tmx6xxTransaction)
	prometheus.MustRegister(TmxUacTransaction)
	prometheus.MustRegister(TmxUasTransaction)
	prometheus.MustRegister(TmxActiveTransaction)
	prometheus.MustRegister(TmxInuseTransaction)
	prometheus.MustRegister(TmxRplAbsorbed)
	prometheus.MustRegister(TmxRplGenerated)
	prometheus.MustRegister(TmxRplReceived)
	prometheus.MustRegister(TmxRplRelayed)
	prometheus.MustRegister(TmxRplSent)
	prometheus.MustRegister(UsrlocLocationContact)
	prometheus.MustRegister(UsrlocLocationExpire)
	prometheus.MustRegister(UsrlocLocationUser)
	prometheus.MustRegister(UsrlocRegisteredUser)
}

func getMeasurement(instance string, target interface{}) error {
	var url = fmt.Sprintf("%s/JSONRPC", instance)

	var body = "{\"jsonrpc\": \"2.0\", \"method\":\"stats.get_statistics\",\"params\":[\"all\"]}"
	r, err := myClient.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		log.Println(err)
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func collectSample() {
	log.Println("Collecting sample...")
	kamMeasurement := new(Measurement)
	getMeasurement(*instance, kamMeasurement)

  for _, el := range kamMeasurement.Result {
  	arr := strings.Split(el, " = ")
  	metric := arr[0]
  	value, err := strconv.ParseFloat(arr[1], 64)
  	if err != nil {
  		log.Println(err)
  		continue
  	}
  	if measurementMapping[metric] != nil {
	  	measurementMapping[metric].With(prometheus.Labels{"instance": *instance}).Set(value)
  	}
	}
}

func main() {
	envy.Parse("KAM")
	flag.Parse()
	http.Handle("/metrics", prometheus.Handler())

	collectSample()
	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %s", *every), collectSample)
	c.Start()

	log.Printf("Listening on %s!", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
