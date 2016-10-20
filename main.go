package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"

	log "github.com/cihub/seelog"
	"github.com/google/go-github/github"
	"github.com/jeremyletang/amish/conf"
	"github.com/jeremyletang/amish/domain"
	"github.com/jeremyletang/amish/gmail"
	"github.com/jeremyletang/amish/slack"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db                *gorm.DB
	wg                sync.WaitGroup
	quit              = make(chan struct{}, 2)
	updateUserBase    = make(chan struct{}, 2)
	updateUserBaseMtx = sync.Mutex{}
	githubClient      *github.Client
)

func init() {
	logger, err := log.LoggerFromConfigAsString(conf.Seelog)
	if err != nil {
		panic("unable to find logger configuration")
	}
	log.ReplaceLogger(logger)
}

func getConfig() conf.Conf {
	c, err := ioutil.ReadFile(".amish.conf.json")
	if err != nil {
		panic(fmt.Sprintf("[getConfig] unable to read config file: %s", err.Error()))
	}

	var config conf.Conf
	if err := json.Unmarshal(c, &config); err != nil {
		panic(fmt.Sprintf("[getConfig] invalid config format: %s", err.Error()))
	}

	return config
}

func initGithubClient(conf conf.Conf) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.Github.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	githubClient = github.NewClient(tc)
}

func main() {
	// read config
	config := getConfig()

	// init github client
	initGithubClient(config)

	// init db
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True",
		config.Mysql.User,
		config.Mysql.Password,
		config.Mysql.Ip,
		config.Mysql.Port,
		config.Mysql.Database)

	if db, err = gorm.Open("mysql", dsn); err != nil {
		panic(fmt.Sprintf("[main] unable to initialize gorm: %s", err.Error()))
	}
	defer db.Close()

	if config.UseGmail {
		// init gmail service
		gmail.InitGmailService(config)
		gmail.SendMail()
	}

	slack.InitSlack(config)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// create repositories if not exists
	repos := createRepositories(config)

	go startTasks(repos, config)

	_ = <-c
	fmt.Println("\nasked to exit, cleaning tasks...")
	quit <- struct{}{}
	wg.Wait()
	fmt.Println("see you later.")
}

func createRepositories(conf conf.Conf) []*domain.Repository {
	dao := domain.NewRepositoryDao(db)

	for _, r := range conf.Github.Repositories {
		arr := strings.Split(r, "/")
		if len(arr) != 2 || arr[0] == "" || arr[1] == "" {
			panic(fmt.Sprintf("invalid repository format name %s", r))
		}
		if _, err := dao.GetOrCreate(&domain.Repository{Owner: arr[0], Name: arr[1]}); err != nil {
			panic(fmt.Sprintf("database error: %s", err.Error))
		}
	}

	repos, err := dao.List()
	if err != nil {
		panic(fmt.Sprintf("unable to get all repositories list: %s", err.Error()))
	}

	return repos
}

func getFmtDate(in string) time.Duration {
	d, err := time.ParseDuration(in)
	if err != nil {
		panic("invalid date format")
	}
	if d > time.Hour*24 {
		panic("invalid date format (>24h)")
	}
	return d
}

func startTasks(repos []*domain.Repository, conf conf.Conf) {
	// first start repositories checkers
	for _, r := range repos {
		go repositoryChecker(r, conf)
	}

	// start notifiers
	for _, l := range conf.Listeners {
		go listenerNotifier(l, conf)
	}

	// start user info collector
	go collectUserInfos()
}

func collectUserInfos() {
	wg.Add(1)
	for {
		select {
		case <-updateUserBase:
			go getUsersInfos()
		case v := <-quit:
			quit <- v
			wg.Done()
			return
		}
	}
}

func getUsersInfos() {
	// add lock to not try to update user base in the multiple time at the same moment
	updateUserBaseMtx.Lock()
	defer updateUserBaseMtx.Unlock()

	log.Infof("[collectUserInfos] starting user infos updates")
	userDao := domain.NewUserDao(db)
	users, _ := userDao.GetContentNotUpdated()
	for _, u := range users {
		id, _ := strconv.Atoi(u.Id)
		githubUser, _, err := githubClient.Users.GetByID(id)
		if err != nil {
			log.Errorf("[getUsersInfos] unable to reach github: %s", err.Error())
		} else {
			log.Infof("[getUsersInfos] updating user %s", u.Login)
			domain.UpdateUserWithGithubUser(u, githubUser)
			u.ContentUpdated = 1
			userDao.Update(u)
		}
	}
	log.Infof("[collectUserInfos] collect user infos ended")
}

func makeNotifyFirstBuffer(notifyTime time.Duration) time.Duration {
	var buffer time.Duration
	daySeconds := time.Hour * 24
	// first get Now seconds of the day
	h, m, s := time.Now().Clock()
	now := time.Hour*time.Duration(h) + time.Duration(m)*time.Minute + time.Duration(s)*time.Second
	// if we are before the notifyTime
	// just make a buffer of the difference
	if now < notifyTime {
		buffer = notifyTime - now
	} else {
		buffer = daySeconds - now
		buffer = buffer + notifyTime
	}

	return buffer
}

func repositoryChecker(repository *domain.Repository, conf conf.Conf) {
	// init refresh shit
	refreshRate := getFmtDate(conf.Refresh)
	log.Infof("[repositoryChecker] (%s/%s) updates will be done every %s",
		repository.Owner, repository.Name, conf.Refresh)
	refreshTicker := time.NewTicker(refreshRate).C
	// start it at least one time
	go updateRepositoryStargazers(repository, conf.Slack.Channels)

	wg.Add(1)
	for {
		select {
		case <-refreshTicker:
			log.Infof("[repositoryChecker] (%s/%s) starting update",
				repository.Owner, repository.Name)
			go updateRepositoryStargazers(repository, conf.Slack.Channels)
		case v := <-quit:
			quit <- v
			wg.Done()
			return
		}
	}
}

func listenerNotifier(listener string, conf conf.Conf) {
	// init notification shit
	notifyTime := getFmtDate(conf.Notify)
	notifyBuffer := makeNotifyFirstBuffer(notifyTime)
	log.Infof("[listenerNotifier] (%s) notification will be send every day at %s",
		listener, conf.Notify)
	log.Infof("[listenerNotifier] (%s) next notification in %s", listener, notifyBuffer.String())
	notifyTimer := time.NewTimer(notifyBuffer)

	wg.Add(1)
	for {
		select {
		case <-notifyTimer.C:
			notifyTimer.Reset(time.Hour * 24)
			log.Info("[listernerNotifier] (%s) send notification", listener)
		case v := <-quit:
			quit <- v
			wg.Done()
			return
		}
	}
}

func addUsersAndStarsToRepository(repository *domain.Repository, stargazers []*github.Stargazer, channels []string) {
	userDao := domain.NewUserDao(db)
	starDao := domain.NewStarDao(db)
	userIds := []string{}
	newStars := []*domain.User{}

	for _, gazer := range stargazers {
		// first create user if not exists
		u, err := userDao.GetById(strconv.Itoa(*(gazer.User.ID)))
		if err != nil {
			log.Infof("[addUserAndStarsToRepository] (%s/%s) create user %s",
				repository.Owner, repository.Name, *(gazer.User.Login))
			// user do not exists
			u = domain.NewUserFromGithubUser(gazer.User)
			_ = userDao.Create(u)
		}
		userIds = append(userIds, u.Id)
		// create the new star
		star := domain.Star{
			RepositoryId: repository.Id,
			UserId:       u.Id,
			StarredAt:    gazer.StarredAt.Time,
			Valid:        1,
		}

		if b, _ := starDao.CreateIfNotExists(&star); b {
			newStars = append(newStars, u)
		}
	}

	if len(newStars) != 0 {
		slack.Notify(slack.Star, newStars, repository, channels)
	}

	// set all other as invalid
	toInvalidateList, _ := starDao.NotOneOfUsers(repository.Id, userIds)
	starDao.SetInvalids(toInvalidateList)

	removedStars := GetUserFromStars(toInvalidateList)
	// slack notify unstar
	slack.Notify(slack.UnStar, removedStars, repository, channels)
}

func GetUserFromStars(stars []domain.Star) []*domain.User {
	ids := []string{}
	for _, s := range stars {
		ids = append(ids, s.UserId)
	}
	userDao := domain.NewUserDao(db)
	users, _ := userDao.GetByIds(ids)
	return users
}

func updateRepositoryStargazers(repository *domain.Repository, channels []string) {
	perPage := 100
	results := []*github.Stargazer{}
	lo := github.ListOptions{Page: 1, PerPage: perPage}

	// get one time at least
	for len(results)%perPage == 0 {
		stargazers, _, err := githubClient.Activity.ListStargazers(
			repository.Owner, repository.Name, &lo)
		if err != nil {
			log.Errorf("[updateRepositoryStargazers] %s", err.Error())
		}
		results = append(results, stargazers...)
		lo.Page += 1
	}

	addUsersAndStarsToRepository(repository, results, channels)
	updateUserBase <- struct{}{}

	log.Infof("[updateRepositoryStargazers] (%s/%s) update finished",
		repository.Owner, repository.Name)
}
