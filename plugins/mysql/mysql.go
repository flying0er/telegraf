package mysql

import (
	"database/sql"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/influxdb/tivan/plugins"
)

type Server struct {
	Address string
}

type Mysql struct {
	Disabled bool
	Servers  []*Server
}

var localhost = &Server{}

func (m *Mysql) Gather(acc plugins.Accumulator) error {
	if m.Disabled {
		return nil
	}

	if len(m.Servers) == 0 {
		// if we can't get stats in this case, thats fine, don't report
		// an error.
		m.gatherServer(localhost, acc)
		return nil
	}

	for _, serv := range m.Servers {
		err := m.gatherServer(serv, acc)
		if err != nil {
			return err
		}
	}

	return nil
}

type mapping struct {
	onServer string
	inExport string
}

var mappings = []*mapping{
	{
		onServer: "Bytes_",
		inExport: "mysql_bytes_",
	},
	{
		onServer: "Com_",
		inExport: "mysql_commands_",
	},
	{
		onServer: "Handler_",
		inExport: "mysql_handler_",
	},
	{
		onServer: "Innodb_",
		inExport: "mysql_innodb_",
	},
	{
		onServer: "Threads_",
		inExport: "mysql_threads_",
	},
}

func (m *Mysql) gatherServer(serv *Server, acc plugins.Accumulator) error {
	db, err := sql.Open("mysql", serv.Address)
	if err != nil {
		return err
	}

	defer db.Close()

	rows, err := db.Query(`SHOW /*!50002 GLOBAL */ STATUS`)
	if err != nil {
		return nil
	}

	for rows.Next() {
		var name string
		var val interface{}

		err = rows.Scan(&name, &val)
		if err != nil {
			return err
		}

		var found bool

		for _, mapped := range mappings {
			if strings.HasPrefix(name, mapped.onServer) {
				i, _ := strconv.Atoi(string(val.([]byte)))
				acc.Add(mapped.inExport+name[len(mapped.onServer):], i, nil)
				found = true
			}
		}

		if found {
			continue
		}

		switch name {
		case "Queries":
			i, err := strconv.ParseInt(string(val.([]byte)), 10, 64)
			if err != nil {
				return err
			}

			acc.Add("mysql_queries", i, nil)
		case "Slow_queries":
			i, err := strconv.ParseInt(string(val.([]byte)), 10, 64)
			if err != nil {
				return err
			}

			acc.Add("mysql_slow_queries", i, nil)
		}
	}

	return nil
}

func init() {
	plugins.Add("mysql", func() plugins.Plugin {
		return &Mysql{}
	})
}