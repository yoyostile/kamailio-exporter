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
	"github.com/iancoleman/strcase"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
var instances = flag.String("instances", "localhost:5060,127.0.0.1:5060", "Instances")
var every = flag.String("every", "1m", "Update time")
var myClient = &http.Client{Timeout: 10 * time.Second}

func NewMetricGauge(name string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("kam_%s", strcase.ToSnake(name)),
		},
		[]string{"instance"},
	)
}

var measurementMapping = map[string] *prometheus.GaugeVec {
	"core:bad_URIs_rcvd": NewMetricGauge("CoreBadUrisRcvd"),
	"core:bad_msg_hdr": NewMetricGauge("CoreBadMsgHdr"),
	"core:drop_replies": NewMetricGauge("CoreDropReply"),
	"core:drop_requests": NewMetricGauge("CoreDropRequest"),
	"core:err_replies": NewMetricGauge("CoreErrReply"),
	"core:err_requests": NewMetricGauge("CoreErrRequest"),
	"core:fwd_replies": NewMetricGauge("CoreFwdReply"),
	"core:fwd_requests": NewMetricGauge("CoreFwdRequest"),
	"core:rcv_replies": NewMetricGauge("CoreRcvReply"),
	"core:rcv_replies_18x": NewMetricGauge("CoreRcvReplies18x"),
	"core:rcv_replies_1xx": NewMetricGauge("CoreRcvReplies1xx"),
	"core:rcv_replies_2xx": NewMetricGauge("CoreRcvReplies2xx"),
	"core:rcv_replies_3xx": NewMetricGauge("CoreRcvReplies3xx"),
	"core:rcv_replies_401": NewMetricGauge("CoreRcvReplies401"),
	"core:rcv_replies_404": NewMetricGauge("CoreRcvReplies404"),
	"core:rcv_replies_407": NewMetricGauge("CoreRcvReplies407"),
	"core:rcv_replies_480": NewMetricGauge("CoreRcvReplies480"),
	"core:rcv_replies_486": NewMetricGauge("CoreRcvReplies486"),
	"core:rcv_replies_4xx": NewMetricGauge("CoreRcvReplies4xx"),
	"core:rcv_replies_5xx": NewMetricGauge("CoreRcvReplies5xx"),
	"core:rcv_replies_6xx": NewMetricGauge("CoreRcvReplies6xx"),
	"core:rcv_requests": NewMetricGauge("CoreRcvRequest"),
	"core:rcv_requests_ack": NewMetricGauge("CoreRcvRequestsAck"),
	"core:rcv_requests_bye": NewMetricGauge("CoreRcvRequestsBye"),
	"core:rcv_requests_cancel": NewMetricGauge("CoreRcvRequestsCancel"),
	"core:rcv_requests_info": NewMetricGauge("CoreRcvRequestsInfo"),
	"core:rcv_requests_invite": NewMetricGauge("CoreRcvRequestsInvite"),
	"core:rcv_requests_message": NewMetricGauge("CoreRcvRequestsMessage"),
	"core:rcv_requests_notify": NewMetricGauge("CoreRcvRequestsNotify"),
	"core:rcv_requests_options": NewMetricGauge("CoreRcvRequestsOption"),
	"core:rcv_requests_prack": NewMetricGauge("CoreRcvRequestsPrack"),
	"core:rcv_requests_publish": NewMetricGauge("CoreRcvRequestsPublish"),
	"core:rcv_requests_refer": NewMetricGauge("CoreRcvRequestsRefer"),
	"core:rcv_requests_register": NewMetricGauge("CoreRcvRequestsRegister"),
	"core:rcv_requests_subscribe": NewMetricGauge("CoreRcvRequestsSubscribe"),
	"core:rcv_requests_update": NewMetricGauge("CoreRcvRequestsUpdate"),
	"core:unsupported_methods": NewMetricGauge("CoreUnsupportedMethod"),
	"dialog:active_dialogs": NewMetricGauge("DialogActiveDialog"),
	"dialog:early_dialogs": NewMetricGauge("DialogEarlyDialog"),
	"dialog:expired_dialogs": NewMetricGauge("DialogExpiredDialog"),
	"dialog:failed_dialogs": NewMetricGauge("DialogFailedDialog"),
	"dialog:processed_dialogs": NewMetricGauge("DialogProcessedDialog"),
	"dns:failed_dns_request": NewMetricGauge("DnsFailedDnsRequest"),
	"httpclient:connections": NewMetricGauge("HttpclientConnection"),
	"httpclient:connfail": NewMetricGauge("HttpclientConnfail"),
	"httpclient:connok": NewMetricGauge("HttpclientConnok"),
	"mysql:driver_errors": NewMetricGauge("MysqlDriverError"),
	"registrar:accepted_regs": NewMetricGauge("RegistrarAcceptedReg"),
	"registrar:default_expire": NewMetricGauge("RegistrarDefaultExpire"),
	"registrar:default_expires_range": NewMetricGauge("RegistrarDefaultExpiresRange"),
	"registrar:expires_range": NewMetricGauge("RegistrarExpiresRange"),
	"registrar:max_contacts": NewMetricGauge("RegistrarMaxContact"),
	"registrar:max_expires": NewMetricGauge("RegistrarMaxExpire"),
	"registrar:rejected_regs": NewMetricGauge("RegistrarRejectedReg"),
	"shmem:fragments": NewMetricGauge("ShmemFragment"),
	"shmem:free_size": NewMetricGauge("ShmemFreeSize"),
	"shmem:max_used_size": NewMetricGauge("ShmemMaxUsedSize"),
	"shmem:real_used_size": NewMetricGauge("ShmemRealUsedSize"),
	"shmem:total_size": NewMetricGauge("ShmemTotalSize"),
	"shmem:used_size": NewMetricGauge("ShmemUsedSize"),
	"sl:1xx_replies": NewMetricGauge("Sl1xxReply"),
	"sl:200_replies": NewMetricGauge("Sl200Reply"),
	"sl:202_replies": NewMetricGauge("Sl202Reply"),
	"sl:2xx_replies": NewMetricGauge("Sl2xxReply"),
	"sl:300_replies": NewMetricGauge("Sl300Reply"),
	"sl:301_replies": NewMetricGauge("Sl301Reply"),
	"sl:302_replies": NewMetricGauge("Sl302Reply"),
	"sl:3xx_replies": NewMetricGauge("Sl3xxReply"),
	"sl:400_replies": NewMetricGauge("Sl400Reply"),
	"sl:401_replies": NewMetricGauge("Sl401Reply"),
	"sl:403_replies": NewMetricGauge("Sl403Reply"),
	"sl:404_replies": NewMetricGauge("Sl404Reply"),
	"sl:407_replies": NewMetricGauge("Sl407Reply"),
	"sl:408_replies": NewMetricGauge("Sl408Reply"),
	"sl:483_replies": NewMetricGauge("Sl483Reply"),
	"sl:4xx_replies": NewMetricGauge("Sl4xxReply"),
	"sl:500_replies": NewMetricGauge("Sl500Reply"),
	"sl:5xx_replies": NewMetricGauge("Sl5xxReply"),
	"sl:6xx_replies": NewMetricGauge("Sl6xxReply"),
	"sl:failures": NewMetricGauge("SlFailure"),
	"sl:received_ACKs": NewMetricGauge("SlReceivedAck"),
	"sl:sent_err_replies": NewMetricGauge("SlSentErrReply"),
	"sl:sent_replies": NewMetricGauge("SlSentReply"),
	"sl:xxx_replies": NewMetricGauge("SlXxxReply"),
	"tcp:con_reset": NewMetricGauge("TcpConReset"),
	"tcp:con_timeout": NewMetricGauge("TcpConTimeout"),
	"tcp:connect_failed": NewMetricGauge("TcpConnectFailed"),
	"tcp:connect_success": NewMetricGauge("TcpConnectSuccess"),
	"tcp:current_opened_connections": NewMetricGauge("TcpCurrentOpenedConnection"),
	"tcp:current_write_queue_size": NewMetricGauge("TcpCurrentWriteQueueSize"),
	"tcp:established": NewMetricGauge("TcpEstablished"),
	"tcp:local_reject": NewMetricGauge("TcpLocalReject"),
	"tcp:passive_open": NewMetricGauge("TcpPassiveOpen"),
	"tcp:send_timeout": NewMetricGauge("TcpSendTimeout"),
	"tcp:sendq_full": NewMetricGauge("TcpSendqFull"),
	"tmx:2xx_transactions": NewMetricGauge("Tmx2xxTransaction"),
	"tmx:3xx_transactions": NewMetricGauge("Tmx3xxTransaction"),
	"tmx:4xx_transactions": NewMetricGauge("Tmx4xxTransaction"),
	"tmx:5xx_transactions": NewMetricGauge("Tmx5xxTransaction"),
	"tmx:6xx_transactions": NewMetricGauge("Tmx6xxTransaction"),
	"tmx:UAC_transactions": NewMetricGauge("TmxUacTransaction"),
	"tmx:UAS_transactions": NewMetricGauge("TmxUasTransaction"),
	"tmx:active_transactions": NewMetricGauge("TmxActiveTransaction"),
	"tmx:inuse_transactions": NewMetricGauge("TmxInuseTransaction"),
	"tmx:rpl_absorbed": NewMetricGauge("TmxRplAbsorbed"),
	"tmx:rpl_generated": NewMetricGauge("TmxRplGenerated"),
	"tmx:rpl_received": NewMetricGauge("TmxRplReceived"),
	"tmx:rpl_relayed": NewMetricGauge("TmxRplRelayed"),
	"tmx:rpl_sent": NewMetricGauge("TmxRplSent"),
	"usrloc:location-contacts": NewMetricGauge("UsrlocLocationContact"),
	"usrloc:location-expires": NewMetricGauge("UsrlocLocationExpire"),
	"usrloc:location-users": NewMetricGauge("UsrlocLocationUser"),
	"usrloc:registered_users": NewMetricGauge("UsrlocRegisteredUser"),
}

type Measurement struct {
	Jsonrpc  string
	Result 	 []string
}

func init() {
  for _, metric := range measurementMapping {
  	prometheus.MustRegister(metric)
  }
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
	arr := strings.Split(*instances, ",")

	for _, instance := range arr {
		kamMeasurement := new(Measurement)
		getMeasurement(instance, kamMeasurement)

	  for _, el := range kamMeasurement.Result {
	  	arr := strings.Split(el, " = ")
	  	metric := arr[0]
	  	value, err := strconv.ParseFloat(arr[1], 64)
	  	if err != nil {
	  		log.Println(err)
	  		continue
	  	}
	  	if measurementMapping[metric] != nil {
		  	measurementMapping[metric].With(prometheus.Labels{"instance": instance}).Set(value)
	  	}
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
