package ipdata

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"github.com/eterline/micro-utils/internal/models"
)

type IpInfoSqlite struct {
	db       *sql.DB
	savePrep *sql.Stmt
	getPrep  *sql.Stmt
}

func NewIpInfoSqlite(ctx context.Context, db *sql.DB) (*IpInfoSqlite, error) {
	self := &IpInfoSqlite{
		db: db,
	}

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

	return fmt.Errorf("failed migrate ip_data: %w", err)
}

func (d *IpInfoSqlite) prepareExec() error {
	savePrep, err := d.db.Prepare(`
	INSERT INTO ip_info (
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
        FROM ip_info
        WHERE ip = ?
        AND (? - request_time) <= ?;
    `)

	if err != nil {
		return err
	}
	d.getPrep = getPrep

	return err
}

func (d *IpInfoSqlite) Save(ctx context.Context, obj models.AboutIPobject) error {
	_, err := d.savePrep.ExecContext(ctx,
		&obj.Status, &obj.Continent, &obj.ContinentCode, &obj.Country, &obj.CountryCode,
		&obj.Region, &obj.RegionName, &obj.City, &obj.District, &obj.Zip, &obj.Lat, &obj.Lon,
		&obj.Timezone, &obj.Offset, &obj.Currency, &obj.Isp, &obj.Org, &obj.As, &obj.Asname,
		&obj.Reverse, &obj.Mobile, &obj.Proxy, &obj.Hosting, &obj.RequestTime,
	)
	return err
}

func (d *IpInfoSqlite) Get(ctx context.Context, ip net.IP) (*models.AboutIPobject, error) {
	now := time.Now().Unix()
	maxAge := 30 * time.Minute
	row := d.getPrep.QueryRowContext(ctx, ip, now, maxAge.Seconds())

	obj := &models.AboutIPobject{}
	err := row.Scan(
		&obj.Status, &obj.Continent, &obj.ContinentCode, &obj.Country, &obj.CountryCode,
		&obj.Region, &obj.RegionName, &obj.City, &obj.District, &obj.Zip, &obj.Lat, &obj.Lon,
		&obj.Timezone, &obj.Offset, &obj.Currency, &obj.Isp, &obj.Org, &obj.As, &obj.Asname,
		&obj.Reverse, &obj.Mobile, &obj.Proxy, &obj.Hosting, &obj.RequestTime,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return obj, nil
}
