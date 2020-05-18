package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/amoghe/distillog"
	"github.com/mjavier2k/solidfire-exporter/pkg/prom"
	"github.com/mjavier2k/solidfire-exporter/pkg/solidfire"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

func init() {
	flag.CommandLine.SortFlags = false
	flag.IntP(solidfire.ListenPortFlag, "l", 9987, fmt.Sprintf("Port for the exporter to listen on. May also be set by environment variable %v.", solidfire.ListenPortFlagEnv))
	flag.StringP(solidfire.UsernameFlag, "u", "my_solidfire_user", fmt.Sprintf("User with which to authenticate to the Solidfire API. May also be set by environment variable %v.", solidfire.UsernameFlagEnv))
	flag.StringP(solidfire.PasswordFlag, "p", "my_solidfire_password", fmt.Sprintf("Password with which to authenticate to the Solidfire API. May also be set by environment variable %v.", solidfire.PasswordFlagEnv))
	flag.StringP(solidfire.EndpointFlag, "e", "https://192.168.1.2/json-rpc/11.3", fmt.Sprintf("Endpoint for the Solidfire API. May also be set by environment variable %v.", solidfire.EndpointFlagEnv))
	flag.BoolP(solidfire.InsecureSSLFlag, "i", false, fmt.Sprintf("Whether to disable TLS validation when calling the Solidfire API. May also be set by environment variable %v.", solidfire.InsecureSSLFlagEnv))
	flag.Parse()

	// PORT environment variable takes precedence in order to be backwards-compatible
	if legacyPort, legacyPortFlagExists := os.LookupEnv("PORT"); legacyPortFlagExists {
		viper.BindEnv(solidfire.ListenPortFlag, "PORT")
		log.Warningf("Found environment variable PORT=%v, skipping lookup of %v", legacyPort, solidfire.ListenPortFlagEnv)
	} else {
		viper.BindEnv(solidfire.ListenPortFlag, solidfire.ListenPortFlagEnv)
	}

	viper.BindEnv(solidfire.UsernameFlag, solidfire.UsernameFlagEnv)
	viper.BindEnv(solidfire.PasswordFlag, solidfire.PasswordFlagEnv)
	viper.BindEnv(solidfire.EndpointFlag, solidfire.EndpointFlagEnv)
	viper.BindEnv(solidfire.InsecureSSLFlag, solidfire.InsecureSSLFlagEnv)
	viper.BindPFlags(flag.CommandLine)
}
func main() {
	log.Infof("Version: %v", sha1ver)
	log.Infof("Built: %v", buildTime)
	listenAddr := fmt.Sprintf("0.0.0.0:%v", viper.GetInt(solidfire.ListenPortFlag))
	solidfireExporter, _ := prom.NewCollector()
	prometheus.MustRegister(solidfireExporter)
	http.Handle("/metrics", promhttp.Handler())
	log.Infof("Booted and listening on %v/metrics\n", listenAddr)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "UP")
	})

	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Errorln(err)
	}
}
