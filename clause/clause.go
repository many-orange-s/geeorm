package clause

import "strings"

type Type int

type Clause struct {
	sql     map[Type]string
	sqlvars map[Type][]interface{}
}

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlvars = make(map[Type][]interface{})
	}

	sql, values := generators[name](vars...)
	c.sql[name] = sql
	c.sqlvars[name] = values
}

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sql []string
	var vars []interface{}
	for _, order := range orders {
		if _, ok := c.sql[order]; ok {
			sql = append(sql, c.sql[order])
			vars = append(vars, c.sqlvars[order]...)
		}
	}
	return strings.Join(sql, " "), vars
}
