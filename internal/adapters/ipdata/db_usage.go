// Copyright (c) 2025 EterLine (Andrew)
// This file is part of micro-utils.
// Licensed under the MIT License. See the LICENSE file for details.

package ipdata

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/eterline/micro-utils/internal/models"
	"github.com/starskey-io/starskey"
)

const (
	maxCacheAge = 30 * time.Minute
)

type IpInfoSqlite struct {
	db       *sql.DB
	savePrep *sql.Stmt
	getPrep  *sql.Stmt
}

func NewIpInfoSqlite(ctx context.Context, db string) (*IpInfoSqlite, error) {
	self := &IpInfoSqlite{}

	sqlite, err := sql.Open("sqlite3", db)
	if err != nil {
		return nil, err
	}

	self.db = sqlite

	if err := self.checkMigrate(ctx); err != nil {
		return nil, err
	}

	if err := self.prepareExec(); err != nil {
		return nil, err
	}

	return self, nil
}

func (d *IpInfoSqlite) checkMigrate(ctx context.Context) error {
	_, err := d.db.ExecContext(ctx, `
    CREATE TABLE IF NOT EXISTS ip_data (      
		ip             TEXT PRIMARY KEY,
		status         TEXT,
		continent      TEXT,
		continent_code TEXT,
		country        TEXT,
		country_code   TEXT,
		region         TEXT,
		region_name    TEXT,
		city           TEXT,
		district       TEXT,
		zip            TEXT,
		lat            REAL,
		lon            REAL,
		timezone       TEXT,
		offset         INTEGER,
		currency       TEXT,
		isp            TEXT,
		org            TEXT,
		as_field       TEXT,
		asname         TEXT,
		reverse        TEXT,
		mobile         BOOLEAN,
		proxy          BOOLEAN,
		hosting        BOOLEAN,
		request_time   TIMESTAMP
    );`)

	if err != nil {
		return fmt.Errorf("failed migrate ip_data: %w", err)
	}

	return nil
}

func (d *IpInfoSqlite) prepareExec() error {
	savePrep, err := d.db.Prepare(`
	INSERT INTO ip_data (
		ip, status, continent, continent_code, country, country_code,
		region, region_name, city, district, zip,
		lat, lon, timezone, offset, currency, isp, org, as_field, asname,
		reverse, mobile, proxy, hosting, request_time
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(ip) DO UPDATE SET
		status         = excluded.status,
		continent      = excluded.continent,
		continent_code = excluded.continent_code,
		country        = excluded.country,
		country_code   = excluded.country_code,
		region         = excluded.region,
		region_name    = excluded.region_name,
		city           = excluded.city,
		district       = excluded.district,
		zip            = excluded.zip,
		lat            = excluded.lat,
		lon            = excluded.lon,
		timezone       = excluded.timezone,
		offset         = excluded.offset,
		currency       = excluded.currency,
		isp            = excluded.isp,
		org            = excluded.org,
		as_field       = excluded.as_field,
		asname         = excluded.asname,
		reverse        = excluded.reverse,
		mobile         = excluded.mobile,
		proxy          = excluded.proxy,
		hosting        = excluded.hosting,
		request_time   = excluded.request_time;
	`)

	if err != nil {
		return err
	}
	d.savePrep = savePrep

	getPrep, err := d.db.Prepare(`
        SELECT status, continent, continent_code, country, country_code,
               region, region_name, city, district, zip,
               lat, lon, timezone, offset, currency, isp, org, as_field, asname,
               reverse, mobile, proxy, hosting, request_time
        FROM ip_data
        WHERE ip = ?
        AND (? - request_time) <= ?;
    `)

	if err != nil {
		return err
	}
	d.getPrep = getPrep

	return err
}

func (d *IpInfoSqlite) Save(ctx context.Context, ip net.IP, obj models.AboutIPobject) error {
	_, err := d.savePrep.ExecContext(ctx,
		ip, &obj.Status, &obj.Continent, &obj.ContinentCode, &obj.Country, &obj.CountryCode,
		&obj.Region, &obj.RegionName, &obj.City, &obj.District, &obj.Zip, &obj.Lat, &obj.Lon,
		&obj.Timezone, &obj.Offset, &obj.Currency, &obj.Isp, &obj.Org, &obj.As, &obj.Asname,
		&obj.Reverse, &obj.Mobile, &obj.Proxy, &obj.Hosting, &obj.RequestTime,
	)
	return err
}

func (d *IpInfoSqlite) Get(ctx context.Context, ip net.IP) (*models.AboutIPobject, error) {
	now := time.Now().Unix()
	row := d.getPrep.QueryRowContext(ctx, ip, now, maxCacheAge.Seconds())

	var ipNew net.IP = ip

	obj := &models.AboutIPobject{}
	err := row.Scan(
		&obj.Status, &obj.Continent, &obj.ContinentCode, &obj.Country, &obj.CountryCode,
		&obj.Region, &obj.RegionName, &obj.City, &obj.District, &obj.Zip, &obj.Lat, &obj.Lon,
		&obj.Timezone, &obj.Offset, &obj.Currency, &obj.Isp, &obj.Org, &obj.As, &obj.Asname,
		&obj.Reverse, &obj.Mobile, &obj.Proxy, &obj.Hosting, &obj.RequestTime, &ipNew,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return obj, nil
}

type IpInfoStarskey struct {
	db *starskey.Starskey
}

func NewIpInfoStarskey(ctx context.Context, db string) (*IpInfoStarskey, error) {

	stubLogCh := make(chan string)

	go func() {
		for range stubLogCh {
		}
	}()

	skey, err := starskey.Open(&starskey.Config{
		Permission:        0755,
		Directory:         db,
		FlushThreshold:    (1024 * 1024) * 64,
		MaxLevel:          3,
		SizeFactor:        10,
		BloomFilter:       false,
		SuRF:              false,
		Logging:           false,
		CompressionOption: starskey.NoCompression,
		Optional: &starskey.OptionalConfig{
			LogChannel: stubLogCh,
		},
	})

	if ctx.Err() != nil {
		return nil, fmt.Errorf("failed to init cache: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to init cache: %w", err)
	}

	go func() {
		<-ctx.Done()
		skey.Close()
		close(stubLogCh)
	}()

	self := &IpInfoStarskey{db: skey}

	return self, nil
}

func (d *IpInfoStarskey) Get(ctx context.Context, ip net.IP) (*models.AboutIPobject, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	payload, err := d.db.Get(ip)
	if err != nil {
		return nil, err
	}

	obj := &models.AboutIPobject{}
	if err := json.Unmarshal(payload, obj); err != nil {
		return nil, err
	}

	expireTime := obj.RequestTime.Add(maxCacheAge)
	if time.Now().After(expireTime) {
		return nil, nil
	}

	return obj, err
}

func (d *IpInfoStarskey) Save(ctx context.Context, ip net.IP, obj models.AboutIPobject) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if ip == nil {
		return errors.New("ip is nil")
	}

	payload, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return d.db.Put(ip, payload)
}
