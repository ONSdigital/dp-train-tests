package train

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ONSdigital/log.go/log"
)

var (
	healthCheckAttempts  = 3
	backOffIncrementSecs = 2
)

type Instance struct {
	cmd               *exec.Cmd
	ErrC              chan error
	CompletedC        chan bool
	ThreadPoolSize    int
	Port              int
	MaxFileUploadSize int
	MaxRequestSize    int
	FileThresholdSize int
	TransactionsDir   string
	WebsiteDir        string
}

var (
	instance *Instance
)

func NewRunner() (*Instance, error) {
	transactionsDir, err := ioutil.TempDir(".", "*-transactions")
	if err != nil {
		return nil, err
	}

	websiteDir, err := ioutil.TempDir(".", "*-website")
	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	t := &Instance{
		ErrC:              make(chan error, 1),
		CompletedC:        make(chan bool, 1),
		ThreadPoolSize:    100,
		Port:              8084,
		MaxFileUploadSize: -1,
		MaxRequestSize:    -1,
		FileThresholdSize: 10,
		TransactionsDir:   filepath.Join(wd, transactionsDir),
		WebsiteDir:        filepath.Join(wd, websiteDir),
	}

	instance = t
	return t, nil
}

func GetInstance() *Instance {
	return instance
}

func (r *Instance) Start() error {
	r.cmd = exec.Command("java", r.getArgs()...)
	r.cmd.Stderr = os.Stderr
	//r.cmd.Stdout = os.Stdout
	r.cmd.Dir = "/Users/Dave/Development/java/The-Train/target"

	go func() {
		if err := r.cmd.Run(); err != nil {
			r.ErrC <- err
		}
	}()

	if err := healthCheck(); err != nil {
		r.ErrC <- err
	}

	return nil
}

func (r *Instance) Stop() {
	log.Event(nil, "stopping train instance", log.INFO)
	if err := r.cmd.Process.Kill(); err != nil {
		log.Event(nil, "error terminating train process", log.Error(err), log.FATAL)
		os.Exit(1)
	}

	if r.TransactionsDir != "" {
/*		if err := os.RemoveAll(r.TransactionsDir); err != nil {
			log.Event(nil, "error removing transactions dir", log.Error(err), log.ERROR)
		}*/
	}

	if r.WebsiteDir != "" {
		if err := os.RemoveAll(r.WebsiteDir); err != nil {
			log.Event(nil, "error terminating train process", log.Error(err), log.ERROR)
		}
	}
}

func (r *Instance) getArgs() []string {
	return []string{
		fmt.Sprintf("-DPUBLISHING_THREAD_POOL_SIZE=%d", r.ThreadPoolSize),
		fmt.Sprintf("-DPORT=%d", r.Port),
		fmt.Sprintf("-DMAX_FILE_UPLOAD_SIZE_MB=%d", r.MaxFileUploadSize),
		fmt.Sprintf("-DMAX_REQUEST_SIZE_MB=%d", r.MaxRequestSize),
		fmt.Sprintf("-DFILE_THRESHOLD_SIZE_MB=%d", r.FileThresholdSize),
		fmt.Sprintf("-DTRANSACTION_STORE=%s", r.TransactionsDir),
		fmt.Sprintf("-DWEBSITE=%s", r.WebsiteDir),
		"-Xmx512m",
		"-Xms512m",
		"-jar",
		"the-train-0.0.1-SNAPSHOT-jar-with-dependencies.jar",
	}
}

func healthCheck() error {
	httpCLi := &http.Client{Timeout: time.Second * 3}
	attempts := 0
	healthy := false
	backoff := 0

	for attempts <= healthCheckAttempts {
		if err := doHealthCheck(httpCLi); err == nil {
			healthy = true
			break
		}
		attempts += 1
		backoff += backOffIncrementSecs
		<-time.After(time.Second * time.Duration(backoff))
	}

	if !healthy {
		return fmt.Errorf("failed to start up instance")
	}

	return nil
}

func doHealthCheck(httpCli *http.Client) error {
	req, err := http.NewRequest("GET", "http://localhost:8084/health", nil)
	if err != nil {
		return err
	}

	resp, err := httpCli.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("incorrect choo choo health check status")
	}

	log.Event(nil, "train is a gogo", log.INFO)
	return nil
}
