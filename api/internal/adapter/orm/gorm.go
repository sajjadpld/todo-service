package orm

import (
	"fmt"
	"microservice/config"
	"microservice/internal/adapter/locale"
	"microservice/internal/adapter/registry"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type sql struct {
	service *config.Service
	config  config.Database
	l       locale.ILocale
	db      *gorm.DB
	tx      *gorm.DB
}

func New(service *config.Service, registry registry.IRegistry, locale locale.ILocale) ISql {
	db := new(sql)

	if service.Debug == false {
		debug, _ := strconv.ParseBool(os.Getenv("DB_DEBUG"))
		maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
		maxOpenConn, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTIONS"))
		maxLifetimeSec, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_SECONDS"))
		slowSqlThreshold, _ := strconv.Atoi(os.Getenv("DB_SLOW_SQL_THRESHOLD"))

		db.config.Debug = debug
		db.config.Host = os.Getenv("DB_HOST")
		db.config.Port = os.Getenv("DB_PORT")
		db.config.Username = os.Getenv("DB_USERNAME")
		db.config.Password = os.Getenv("DB_PASSWORD")
		db.config.Database = os.Getenv("DB_DATABASE")
		db.config.Ssl = os.Getenv("DB_SSL")
		db.config.MaxIdleConnections = maxIdleConn
		db.config.MaxOpenConnections = maxOpenConn
		db.config.MaxLifetimeSeconds = maxLifetimeSec
		db.config.SlowSqlThreshold = slowSqlThreshold
	} else {
		registry.Parse(&db.config)
	}

	db.l = locale
	return db
}

func (s *sql) Init() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		s.config.Host,
		s.config.Username,
		s.config.Password,
		s.config.Database,
		s.config.Port,
		s.config.Ssl,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 s.newGormLog(s.config.SlowSqlThreshold),
		NowFunc:                func() time.Time { return time.Now().UTC() },
	})
	if err != nil {
		log.Fatalf("[sql] connect err: %s", err)
	}

	sqlDatabase, err := database.DB()
	if err != nil {
		log.Fatalf("[sql] init err: %s", err)
	}

	if s.config.MaxIdleConnections != 0 {
		sqlDatabase.SetMaxIdleConns(s.config.MaxIdleConnections)
	}

	if s.config.MaxOpenConnections != 0 {
		sqlDatabase.SetMaxOpenConns(s.config.MaxOpenConnections)
	}

	if s.config.MaxLifetimeSeconds != 0 {
		sqlDatabase.SetConnMaxLifetime(time.Second * time.Duration(s.config.MaxLifetimeSeconds))
	}

	if s.config.Debug {
		database = database.Debug()
		log.Print("[sql] debug is enabled\n\n")
	}

	s.db = database
}

func (s *sql) C() *gorm.DB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *sql) Stop() {
	closeOnce := sync.Once{}

	sqlDatabase, err := s.db.DB()
	if err != nil {
		log.Printf("[sql] connection close retrieve err: %s", err)
	}

	closeOnce.Do(func() {
		err = sqlDatabase.Close()
		if err != nil {
			log.Printf("[sql] connection close err: %s", err)
		}

		log.Printf("[sql] daatabase stopped successfully")
	})
}

func (s *sql) Migrate(path string) {
	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		log.Fatalf("[sql] migrations dir scan err: %s", err)
	}

	defer func() {
		_ = dir.Close()
	}()

	// Read the directory contents
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Fatalf("[sql] reading directory contents err: %s", err)
	}

	if len(fileInfos) == 0 {
		log.Fatalf("[sql] no migration file found in the current path:\n%s", path)
	}

	// sort the entries alphabetically by name - sql file order by numeric(01, 02, etc)
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Name() < fileInfos[j].Name()
	})

	// Iterate over the file info slice and print the file names
	for _, fileInfo := range fileInfos {
		if fileInfo.Mode().IsRegular() {
			if err = s.db.Exec(s.parseSqlFile(path, fileInfo)).Error; err != nil {
				log.Fatalf("[sql] migrate err: %s", err)
			}
		}
	}
}

func (s *sql) Seed() {
	// Consider desired seeder data here
}

func (s *sql) Begin() {
	s.tx = s.db.Begin()
}

func (s *sql) Commit() error {
	if s.tx != nil {
		return s.tx.Commit().Error
	}

	return nil
}

func (s *sql) Rollback() (err error) {
	if s.tx != nil {
		return s.tx.Rollback().Error
	}

	return nil
}

func (s *sql) Resolve(dbErr error) (err error) {
	if dbErr != nil {
		err = s.Rollback()
		s.tx = nil //resets the DB transactional instance to avoid being reused by the C() method.
		return
	}

	err = s.Commit()
	s.tx = nil
	return
}

// HELPER METHODS

func (s *sql) newGormLog(SlowSqlThreshold int) logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Duration(SlowSqlThreshold) * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn,                                   // Log level
			IgnoreRecordNotFoundError: false,                                         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                                          // Disable color
		})
}

func (s *sql) parseSqlFile(path string, fileInfo os.FileInfo) string {
	sqlFile := fmt.Sprintf("%s/%s", path, fileInfo.Name())
	sqlBytes, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("[sql] SQL file parse err: %s", err)
	}
	// Convert SQL file contents to string
	return string(sqlBytes)
}
