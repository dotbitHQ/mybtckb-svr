package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/scorpiotzh/toolib"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v74"
	"mybtckb-svr/lib/common"
)

var (
	Cfg CfgServer
)

func InitCfg(configFilePath string) error {
	if configFilePath == "" {
		configFilePath = "../config/config.yaml"
	}
	log.Debug("config file path：", configFilePath)
	if err := toolib.UnmarshalYamlFile(configFilePath, &Cfg); err != nil {
		return fmt.Errorf("UnmarshalYamlFile err:%s", err.Error())
	}
	initStripe()
	log.Debug("config file：", toolib.JsonString(Cfg))
	return nil
}

func AddCfgFileWatcher(configFilePath string) (*fsnotify.Watcher, error) {
	if configFilePath == "" {
		configFilePath = "../config/config.yaml"
	}
	return toolib.AddFileWatcher(configFilePath, func() {
		log.Debug("config file path：", configFilePath)
		if err := toolib.UnmarshalYamlFile(configFilePath, &Cfg); err != nil {
			log.Error("UnmarshalYamlFile err:", err.Error())
		}
		log.Debug("config file：", toolib.JsonString(Cfg))
	})
}

type CfgServer struct {
	Server struct {
		IsUpdate         bool           `json:"is_update" yaml:"is_update"`
		HttpServerAddr   string         `json:"http_server_addr" yaml:"http_server_addr"`
		Name             string         `json:"name" yaml:"name"`
		Net              common.NetType `json:"net" yaml:"net"`
		HttpPort         string         `json:"http_port" yaml:"http_port"`
		HttpPortInternal string         `json:"http_port_internal" yaml:"http_port_internal"`
		PayServerAddress string         `json:"pay_server_address" yaml:"pay_server_address"`
		PayPrivate       string         `json:"pay_private" yaml:"pay_private"`
		RemoteSignApiUrl string         `json:"remote_sign_api_url" yaml:"remote_sign_api_url"`
		TxTeeRate        uint64         `json:"tx_fee_rate" yaml:"tx_fee_rate"`
	} `json:"server" yaml:"server"`
	Origins []string `json:"origins" yaml:"origins"`
	Notify  struct {
		MinBalance            decimal.Decimal `json:"min_balance" yaml:"min_balance"`
		SentryDsn             string          `json:"sentry_dsn" yaml:"sentry_dsn"`
		LarkKey               string          `json:"lark_key" yaml:"lark_key"`
		LarkErrKey            string          `json:"lark_err_key" yaml:"lark_err_key"`
		PrometheusPushGateway string          `json:"prometheus_push_gateway" yaml:"prometheus_push_gateway"`
		LarkStripeErrKey      string          `json:"lark_stripe_err_key" yaml:"lark_stripe_err_key"`
	} `json:"notify" yaml:"notify"`
	DB struct {
		Mysql DbMysql `json:"mysql" yaml:"mysql"`
		Redis struct {
			Addr     string `json:"addr" yaml:"addr"`
			Password string `json:"password" yaml:"password"`
			DbNum    int    `json:"db_num" yaml:"db_num"`
		} `json:"redis" yaml:"redis"`
	} `json:"db" yaml:"db"`
	Chain struct {
		Ckb struct {
			Url                string `json:"url" yaml:"url"`
			IndexUrl           string `json:"index_url" yaml:"index_url"`
			CurrentBlockNumber uint64 `json:"current_block_number" yaml:"current_block_number"`
			ConfirmNum         uint64 `json:"confirm_num" yaml:"confirm_num"`
			ConcurrencyNum     uint64 `json:"concurrency_num" yaml:"concurrency_num"`
			PrivateKey         string `json:"privateKey" yaml:"privateKey"`
		} `json:"ckb" yaml:"ckb"`
		Btc struct {
			Url    string `json:"url" yaml:"url"`
			ApiKey string `json:"api_key" yaml:"api_key"`
		}
	} `json:"chain" yaml:"chain"`
	Stripe struct {
		PremiumPercentage decimal.Decimal `json:"premium_percentage" yaml:"premium_percentage"`
		PremiumBase       decimal.Decimal `json:"premium_base" yaml:"premium_base"`
		Key               string          `json:"key" yaml:"key"`
	} `json:"stripe" yaml:"stripe"`
}

type DbMysql struct {
	Addr        string `json:"addr" yaml:"addr"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"password"`
	DbName      string `json:"db_name" yaml:"db_name"`
	MaxOpenConn int    `json:"max_open_conn" yaml:"max_open_conn"`
	MaxIdleConn int    `json:"max_idle_conn" yaml:"max_idle_conn"`
}

func initStripe() {
	stripe.Key = Cfg.Stripe.Key
}
