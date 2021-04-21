package gorm

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/micro/micro/v3/service/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type Helper struct {
	sync.RWMutex
	gormConns  map[string]*gorm.DB
	dbConn     *sql.DB
	migrations []interface{}
}

func (h *Helper) Migrations(migrations ...interface{}) *Helper {
	h.migrations = migrations
	return h
}

func (h *Helper) DBConn(conn *sql.DB) *Helper {
	h.dbConn = conn
	h.gormConns = map[string]*gorm.DB{}
	return h
}

func getTenancyKey(acc *auth.Account) string {
	owner := acc.Metadata["apikey_owner"]
	if len(owner) == 0 {
		owner = acc.ID
	}
	return fmt.Sprintf("%s_%s", acc.Issuer, owner)
}

func (h *Helper) GetDBConn(ctx context.Context) (*gorm.DB, error) {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing account from context")
	}
	h.RLock()
	tenancyKey := getTenancyKey(acc)
	if conn, ok := h.gormConns[tenancyKey]; ok {
		h.RUnlock()
		return conn, nil
	}
	h.RUnlock()
	h.Lock()
	// double check
	if conn, ok := h.gormConns[tenancyKey]; ok {
		h.Unlock()
		return conn, nil
	}
	defer h.Unlock()
	ns := schema.NamingStrategy{
		TablePrefix: fmt.Sprintf("%s_", strings.ReplaceAll(tenancyKey, "-", "")),
	}
	db, err := gorm.Open(
		newGormDialector(postgres.Config{
			Conn: h.dbConn,
		}, ns),
		&gorm.Config{
			NamingStrategy: ns,
		})
	if err != nil {
		return nil, err
	}
	if len(h.migrations) == 0 {
		// record success
		h.gormConns[tenancyKey] = db
		return db, nil
	}

	if err := db.AutoMigrate(h.migrations...); err != nil {
		return nil, err
	}

	// record success
	h.gormConns[tenancyKey] = db
	return db, nil
}

func newGormDialector(config postgres.Config, ns schema.NamingStrategy) gorm.Dialector {
	return &postgresDial{
		Dialector: postgres.Dialector{Config: &config},
		namer:     &ns,
	}
}

// postgresDial is a postgres dialector that prefixes index names with the table prefix when doing migrations.
// NOTE, it does not support the gorm tag priority option
type postgresDial struct {
	postgres.Dialector
	namer schema.Namer
}

func (p postgresDial) Migrator(db *gorm.DB) gorm.Migrator {
	return gormMigrator{
		postgres.Migrator{
			migrator.Migrator{Config: migrator.Config{
				DB:                          db,
				Dialector:                   p,
				CreateIndexAfterCreateTable: true,
			}},
		},
		p.namer,
	}
}

type gormMigrator struct {
	postgres.Migrator
	namer schema.Namer
}

// AutoMigrate
func (m gormMigrator) AutoMigrate(values ...interface{}) error {
	for _, value := range m.ReorderModels(values, true) {
		tx := m.DB.Session(&gorm.Session{NewDB: true})
		if !tx.Migrator().HasTable(value) {
			if err := tx.Migrator().CreateTable(value); err != nil {
				return err
			}
		} else {
			if err := m.RunWithValue(value, func(stmt *gorm.Statement) (errr error) {
				columnTypes, _ := m.DB.Migrator().ColumnTypes(value)

				for _, field := range stmt.Schema.FieldsByDBName {
					var foundColumn gorm.ColumnType

					for _, columnType := range columnTypes {
						if columnType.Name() == field.DBName {
							foundColumn = columnType
							break
						}
					}

					if foundColumn == nil {
						// not found, add column
						if err := tx.Migrator().AddColumn(value, field.DBName); err != nil {
							return err
						}
					} else if err := m.DB.Migrator().MigrateColumn(value, field, foundColumn); err != nil {
						// found, smart migrate
						return err
					}
				}

				for _, rel := range stmt.Schema.Relationships.Relations {
					if !m.DB.Config.DisableForeignKeyConstraintWhenMigrating {
						if constraint := rel.ParseConstraint(); constraint != nil {
							if constraint.Schema == stmt.Schema {
								if !tx.Migrator().HasConstraint(value, constraint.Name) {
									if err := tx.Migrator().CreateConstraint(value, constraint.Name); err != nil {
										return err
									}
								}
							}
						}
					}

					for _, chk := range stmt.Schema.ParseCheckConstraints() {
						if !tx.Migrator().HasConstraint(value, chk.Name) {
							if err := tx.Migrator().CreateConstraint(value, chk.Name); err != nil {
								return err
							}
						}
					}
				}

				for _, idx := range m.ParseIndexes(stmt.Schema) {
					if !tx.Migrator().HasIndex(value, idx.Name) {
						if err := tx.Migrator().CreateIndex(value, idx.Name); err != nil {
							return err
						}
					}
				}

				return nil
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m gormMigrator) CreateIndex(value interface{}, name string) error {
	return m.RunWithValue(value, func(stmt *gorm.Statement) error {
		if idx := m.LookIndex(stmt.Schema, name); idx != nil {
			opts := m.BuildIndexOptions(idx.Fields, stmt)
			values := []interface{}{clause.Column{Name: idx.Name}, m.CurrentTable(stmt), opts}

			createIndexSQL := "CREATE "
			if idx.Class != "" {
				createIndexSQL += idx.Class + " "
			}
			createIndexSQL += "INDEX ?"

			createIndexSQL += " ON ?"

			if idx.Type != "" {
				createIndexSQL += " USING " + idx.Type + "(?)"
			} else {
				createIndexSQL += " ?"
			}

			if idx.Where != "" {
				createIndexSQL += " WHERE " + idx.Where
			}

			return m.DB.Exec(createIndexSQL, values...).Error
		}

		return fmt.Errorf("failed to create index with name %v", name)
	})
}

func (m gormMigrator) LookIndex(sch *schema.Schema, name string) *schema.Index {
	if sch != nil {
		indexes := m.ParseIndexes(sch)
		for _, index := range indexes {
			if index.Name == name {
				return &index
			}

			for _, field := range index.Fields {
				if field.Name == name {
					return &index
				}
			}
		}
	}

	return nil
}

func (g gormMigrator) parseFieldIndexes(field *schema.Field) (indexes []schema.Index) {
	for _, value := range strings.Split(field.Tag.Get("gorm"), ";") {
		if value != "" {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if k == "INDEX" || k == "UNIQUEINDEX" {
				var (
					name      string
					tag       = strings.Join(v[1:], ":")
					idx       = strings.Index(tag, ",")
					settings  = schema.ParseTagSetting(tag, ",")
					length, _ = strconv.Atoi(settings["LENGTH"])
				)

				if idx == -1 {
					idx = len(tag)
				}

				if idx != -1 {
					name = tag[0:idx]
				}

				if name == "" {
					name = g.namer.IndexName(field.Schema.Table, field.Name)
				} else {
					ns := g.namer.(*schema.NamingStrategy)
					name = fmt.Sprintf("%s%s", ns.TablePrefix, name)
				}

				if (k == "UNIQUEINDEX") || settings["UNIQUE"] != "" {
					settings["CLASS"] = "UNIQUE"
				}

				//priority, err := strconv.Atoi(settings["PRIORITY"])
				//if err != nil {
				//	priority = 10
				//}

				indexes = append(indexes, schema.Index{
					Name:    name,
					Class:   settings["CLASS"],
					Type:    settings["TYPE"],
					Where:   settings["WHERE"],
					Comment: settings["COMMENT"],
					Option:  settings["OPTION"],
					Fields: []schema.IndexOption{{
						Field:      field,
						Expression: settings["EXPRESSION"],
						Sort:       settings["SORT"],
						Collate:    settings["COLLATE"],
						Length:     length,
						//priority:   priority, // TODO does not support priority
					}},
				})
			}
		}
	}

	return
}

func (m gormMigrator) CreateTable(values ...interface{}) error {
	for _, value := range m.ReorderModels(values, false) {
		tx := m.DB.Session(&gorm.Session{NewDB: true})
		if err := m.RunWithValue(value, func(stmt *gorm.Statement) (errr error) {
			var (
				createTableSQL          = "CREATE TABLE ? ("
				values                  = []interface{}{m.CurrentTable(stmt)}
				hasPrimaryKeyInDataType bool
			)

			for _, dbName := range stmt.Schema.DBNames {
				field := stmt.Schema.FieldsByDBName[dbName]
				createTableSQL += "? ?"
				hasPrimaryKeyInDataType = hasPrimaryKeyInDataType || strings.Contains(strings.ToUpper(string(field.DataType)), "PRIMARY KEY")
				values = append(values, clause.Column{Name: dbName}, m.DB.Migrator().FullDataTypeOf(field))
				createTableSQL += ","
			}

			if !hasPrimaryKeyInDataType && len(stmt.Schema.PrimaryFields) > 0 {
				createTableSQL += "PRIMARY KEY ?,"
				primaryKeys := []interface{}{}
				for _, field := range stmt.Schema.PrimaryFields {
					primaryKeys = append(primaryKeys, clause.Column{Name: field.DBName})
				}

				values = append(values, primaryKeys)
			}

			for _, idx := range m.ParseIndexes(stmt.Schema) {
				if m.CreateIndexAfterCreateTable {
					defer func(value interface{}, name string) {
						errr = tx.Migrator().CreateIndex(value, name)
					}(value, idx.Name)
				} else {
					if idx.Class != "" {
						createTableSQL += idx.Class + " "
					}
					createTableSQL += "INDEX ? ?"

					if idx.Option != "" {
						createTableSQL += " " + idx.Option
					}

					createTableSQL += ","
					values = append(values, clause.Expr{SQL: idx.Name}, tx.Migrator().(migrator.BuildIndexOptionsInterface).BuildIndexOptions(idx.Fields, stmt))
				}
			}

			for _, rel := range stmt.Schema.Relationships.Relations {
				if !m.DB.DisableForeignKeyConstraintWhenMigrating {
					if constraint := rel.ParseConstraint(); constraint != nil {
						if constraint.Schema == stmt.Schema {
							sql, vars := buildConstraint(constraint)
							createTableSQL += sql + ","
							values = append(values, vars...)
						}
					}
				}
			}

			for _, chk := range stmt.Schema.ParseCheckConstraints() {
				createTableSQL += "CONSTRAINT ? CHECK (?),"
				values = append(values, clause.Column{Name: chk.Name}, clause.Expr{SQL: chk.Constraint})
			}

			createTableSQL = strings.TrimSuffix(createTableSQL, ",")

			createTableSQL += ")"

			if tableOption, ok := m.DB.Get("gorm:table_options"); ok {
				createTableSQL += fmt.Sprint(tableOption)
			}

			errr = tx.Exec(createTableSQL, values...).Error
			return errr
		}); err != nil {
			return err
		}
	}
	return nil
}

type gormIndexOption struct {
	schema.IndexOption
	priority int
}

func (g gormMigrator) ParseIndexes(sch *schema.Schema) map[string]schema.Index {
	var indexes = map[string]schema.Index{}

	for _, field := range sch.Fields {
		if field.TagSettings["INDEX"] != "" || field.TagSettings["UNIQUEINDEX"] != "" {
			for _, index := range g.parseFieldIndexes(field) {
				idx := indexes[index.Name]
				idx.Name = index.Name
				if idx.Class == "" {
					idx.Class = index.Class
				}
				if idx.Type == "" {
					idx.Type = index.Type
				}
				if idx.Where == "" {
					idx.Where = index.Where
				}
				if idx.Comment == "" {
					idx.Comment = index.Comment
				}
				if idx.Option == "" {
					idx.Option = index.Option
				}

				idx.Fields = append(idx.Fields, index.Fields...)
				// TODO priority not supported
				//sort.Slice(idx.Fields, func(i, j int) bool {
				//	return idx.Fields[i].priority < idx.Fields[j].priority
				//})

				indexes[index.Name] = idx
			}
		}
	}

	return indexes
}

func buildConstraint(constraint *schema.Constraint) (sql string, results []interface{}) {
	sql = "CONSTRAINT ? FOREIGN KEY ? REFERENCES ??"
	if constraint.OnDelete != "" {
		sql += " ON DELETE " + constraint.OnDelete
	}

	if constraint.OnUpdate != "" {
		sql += " ON UPDATE " + constraint.OnUpdate
	}

	var foreignKeys, references []interface{}
	for _, field := range constraint.ForeignKeys {
		foreignKeys = append(foreignKeys, clause.Column{Name: field.DBName})
	}

	for _, field := range constraint.References {
		references = append(references, clause.Column{Name: field.DBName})
	}
	results = append(results, clause.Table{Name: constraint.Name}, foreignKeys, clause.Table{Name: constraint.ReferenceSchema.Table}, references)
	return
}
