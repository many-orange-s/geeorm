package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}

func _insert(value ...interface{}) (string, []interface{}) {
	// INSERT INTO $tableName ($fields)
	tablename := value[0]
	fields := strings.Join(value[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tablename, fields), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	// VALUES ($v1), ($v2), ...
	var sql strings.Builder
	var bindstr string
	var vars []interface{}
	sql.WriteString("VALUES")
	for i, value := range values {
		//这是确定传进来的是一个数组
		v := value.([]interface{})
		if bindstr == "" {
			bindstr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindstr))

		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _select(value ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tablename := value[0]
	fields := strings.Join(value[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tablename), []interface{}{}
}

func _limit(value ...interface{}) (string, []interface{}) {
	return "LIMIT ?", value
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderBy(value ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", value), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{}) {
	tablename := values[0]
	m := values[1].(map[string]interface{})
	var key []string
	var vars []interface{}
	for name, value := range m {
		key = append(key, fmt.Sprintf(name+" =?"))
		vars = append(vars, value)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tablename, strings.Join(key, ",")), vars
}

func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}
