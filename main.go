package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sandro-h/sibylgo/backup"
	"github.com/sandro-h/sibylgo/extsources"
	"github.com/sandro-h/sibylgo/outlook"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/popup"
	"github.com/sandro-h/sibylgo/reminder"
	"github.com/sandro-h/sibylgo/util"
	log "github.com/sirupsen/logrus"
)

var ascii = `.▄▄ · ▪  ▄▄▄▄·  ▄· ▄▌▄▄▌  
▐█ ▀. ██ ▐█ ▀█▪▐█▪██▌██•  
▄▀▀▀█▄▐█·▐█▀▀█▄▐█▌▐█▪██▪  
▐█▄▪▐█▐█▌██▄▪▐█ ▐█▀·.▐█▌▐▌
 ▀▀▀▀ ▀▀▀·▀▀▀▀   ▀ • .▀▀▀ `

var buildVersion = "0.0.0"
var buildNumber = "0"
var buildRevision = "-"

var configFile = flag.String("config", "", "Path to config yml file. By default uses sibylgo.yml in same directory as this executable, if it exists.")
var doEncrypt = flag.Bool("encrypt", false, "Encrypt stdin and write to stdout")
var doDecrypt = flag.Bool("decrypt", false, "Decrypt stdin and write to stdout")
var files *util.FileConfig
var extSourcesProcess *extsources.ExternalSourcesProcess

func main() {
	flag.Parse()

	if *doEncrypt || *doDecrypt {
		log.SetOutput(ioutil.Discard)
		cfg := loadConfig()
		err := cryptContent(cfg.GetSubConfig("backup"))
		if err != nil {
			panic(err)
		}
		return
	}

	log.SetFormatter(&SimpleFormatter{})

	fmt.Printf("%s\n", ascii)
	log.Infof("Version %s.%s (%s)\n", buildVersion, buildNumber, buildRevision)

	cfg := loadConfig()

	log.SetLevel(getConfigLogLevel(cfg))

	parse.ParseConfig.BackingCfg = cfg.GetSubConfig("parse")

	files = util.NewFileConfigFromConfig(cfg)
	if files.TodoFile != "" {
		log.Infof("Using todo file %s\n", files.TodoFile)
		startBackups(cfg.GetSubConfig("backup"))
	}

	if cfg.HasKey("mailTo") {
		startMailReminders(cfg)
	}

	if cfg.HasKey("external_sources") {
		extSrcConfig := cfg.GetSubConfig("external_sources")
		if files.TodoFile == "" {
			panic("Cannot run external sources without todoFile set")
		}
		startExternalSources(files, extSrcConfig)
	}

	if cfg.HasKey("outlook_events") {
		outlookConfig := cfg.GetSubConfig("outlook_events")
		if files.TodoFile == "" {
			panic("Cannot run outlook events without todoFile set")
		}
		startOutlookEvents(files.TodoFile, outlookConfig)
	}

	startRestServer(cfg)

	if files.TodoFile != "" && cfg.HasKey("popup") {
		popup.Start(files, cfg.GetSubConfig("popup"))
	} else {
		// Wait forever
		select {}
	}
}

func loadConfig() *util.Config {
	absoluteCfgFile := *configFile
	if absoluteCfgFile == "" {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		absoluteCfgFile = filepath.Join(dir, "sibylgo.yml")
	} else {
		// We don't care if the default config file doesn't exist,
		// but if a user set a config file explicitly, we should inform them
		// if the file doesn't actually exist.
		if !util.Exists(absoluteCfgFile) {
			panic(fmt.Sprintf("Config file %s set with -config does not exist.\n", absoluteCfgFile))
		}
	}

	cfg := &util.Config{}
	log.Infof("%s\n", absoluteCfgFile)
	if util.Exists(absoluteCfgFile) {
		var err error
		cfg, err = util.LoadConfig(absoluteCfgFile)
		if err != nil {
			panic(err)
		}
	}

	return cfg
}

func startBackups(backupCfg *util.Config) {
	if backupCfg.HasKey("encrypt_password") {
		exec, err := os.Executable()
		if err != nil {
			panic(err)
		}
		backup.EnableGitEncryption(files.TodoDir, exec)
	}

	startDailyBackupProcess(backupCfg, files)
}

func cryptContent(backupCfg *util.Config) error {
	var cryptor backup.Cryptor = &backup.AnsibleCryptor{
		Password: backupCfg.GetStringOrFail("encrypt_password"),
	}

	if *doEncrypt {
		err := cryptor.EncryptContent(os.Stdin, os.Stdout)
		if err != nil {
			panic(err)
		}
	} else if *doDecrypt {
		in := os.Stdin
		var err error
		if flag.NArg() > 0 {
			in, err = os.Open(flag.Arg(0))
			if err != nil {
				panic(err)
			}
			defer in.Close()
		}

		err = cryptor.DecryptContent(in, os.Stdout)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func startMailReminders(cfg *util.Config) {
	cfg.GetStringOrFail("todoFile")
	mailHost := cfg.GetStringOrFail("mailHost")
	mailPort := cfg.GetIntOrFail("mailPort")
	mailFrom := cfg.GetStringOrFail("mailFrom")
	mailTo := cfg.GetStringOrFail("mailTo")
	mailUser := cfg.GetString("mailUser", "")
	mailPassword := cfg.GetString("mailPassword", "")

	host := reminder.MailHostProperties{Host: mailHost, Port: mailPort, User: mailUser, Password: mailPassword}
	p := reminder.NewMailReminderProcessForSMTP(files.TodoFile, host, mailFrom, mailTo)
	go p.CheckInfinitely()
	log.Info("Started mail reminders\n")
}

func startExternalSources(files *util.FileConfig, extSrcConfig *util.Config) {
	extSourcesProcess = extsources.NewExternalSourcesProcess(files, extSrcConfig)
	go extSourcesProcess.CheckInfinitely()
	log.Info("Started external sources\n")
}

func startOutlookEvents(todoFile string, outlookConfig *util.Config) {
	if outlookConfig.GetBool("enabled", false) {
		go outlook.CheckInfinitely(todoFile, 5*time.Second)
		log.Info("Started outlook syncing\n")
	}
}

func startDailyBackupProcess(backupCfg *util.Config, files *util.FileConfig) {
	dailyBackupFunc := func() {
		for {
			doDailyBackup(backupCfg, files)
			time.Sleep(5 * time.Minute)
		}
	}
	go dailyBackupFunc()
	log.Info("Started daily backup\n")
}

func doDailyBackup(backupCfg *util.Config, files *util.FileConfig) {
	newBackup, err := backup.CheckAndMakeDailyBackup(files)
	if err != nil {
		log.Errorf("Error making the daily backup: %s\n", err)
		return
	}

	if backupCfg.HasKey("remote_url") && newBackup != nil {
		err = backup.SyncToRemote(backupCfg, files)
		if err != nil {
			log.Errorf("Error syncing backup to remote: %s\n", err)
		}
	}
}
