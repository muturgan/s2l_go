package dal

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/muturgan/s2l_go/src/config"
	"github.com/muturgan/s2l_go/src/models"
)

const _DB_TYPE = "mysql"

type IDal interface {
	GetLinkByHash(hash string) (*models.LinkEntry, error)
	CreateNewLink(link string) (*models.Result, error)
}

type Dal struct {
	db     *sql.DB
	config *config.Config
}

func New(conf *config.Config) (*Dal, error) {
	dbUrl := conf.GetDbUrl()
	db, err := sql.Open(_DB_TYPE, dbUrl)
	if err != nil {
		return nil, err
	}

	err = testConnection(db)
	if err != nil {
		return nil, err
	}
	fmt.Println("DB successfully connected...")

	dal := &Dal{
		db:     db,
		config: conf,
	}
	return dal, nil
}

func (dal *Dal) GetLinkByHash(hash string) (*models.LinkEntry, error) {
	l := &models.LinkEntry{}

	err := dal.db.QueryRow(
		"SELECT link FROM links WHERE hash = ?",
		hash,
	).Scan(
		&l.Link,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return l, nil
}

func (dal *Dal) CreateNewLink(link string) (*models.Result, error) {
	urlStr := dal.config.APP_DOMAIN()

	existingHash, err := dal.getHashByLink(link)
	if err != nil {
		return nil, err
	}

	if existingHash != nil {
		urlStr = urlStr + "/" + existingHash.Hash
	} else {
		newHash, err := dal.generateUniqueHash()
		if err != nil {
			return nil, err
		}

		err = dal.insertNewLinkEntry(link, newHash)
		if err != nil {
			return nil, err
		}

		urlStr = urlStr + "/" + newHash
	}

	r := &models.Result{
		Link: urlStr,
	}
	return r, nil
}

func (dal *Dal) getHashByLink(link string) (*models.LinkEntry, error) {
	l := &models.LinkEntry{}

	rows := dal.db.QueryRow(
		"SELECT hash FROM links WHERE link = ?",
		link,
	)
	err := rows.Scan(
		&l.Hash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return l, nil
}

func (dal *Dal) getEntryByHash(hash string) (*models.LinkEntry, error) {
	l := &models.LinkEntry{}

	err := dal.db.QueryRow(
		"SELECT ID, link, hash FROM links WHERE hash = ?",
		hash,
	).Scan(
		&l.ID,
		&l.Link,
		&l.Hash,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return l, nil
}

func (dal *Dal) Stop() {
	dal.db.Close()
}

func (dal *Dal) generateUniqueHash() (string, error) {
	newHash, err := generateHash()
	if err != nil {
		return "", err
	}

	entry, err := dal.getEntryByHash(newHash)
	if err != nil {
		return "", err
	}
	if entry == nil {
		return newHash, nil
	} else {
		return dal.generateUniqueHash()
	}
}

func (dal *Dal) insertNewLinkEntry(link string, hash string) error {
	_, err := dal.db.Exec(
		"INSERT INTO links (hash, link) VALUES (?, ?)",
		hash,
		link,
	)
	return err
}

func testConnection(db *sql.DB) error {
	_, err := db.Exec(
		"SELECT 1 + 1",
	)
	return err
}
