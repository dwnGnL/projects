package db

import (
	"fmt"
	"log"

	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"projects/internal/database/actionPlan"
	"projects/internal/database/epics"
	"projects/internal/database/milestone"
	"projects/internal/database/processes"
	"projects/internal/database/projects"
	"projects/internal/database/stage"
	"projects/internal/database/tasks"
	"projects/internal/database/workspace"
	"projects/pkg/config"
)

var Module = fx.Provide(Setup)

type Params struct {
	fx.In
	*logrus.Logger
	*config.Tuner
}

type DbInter interface {
	Close() error
	GetDB() *gorm.DB
}

type gormDB struct {
	db  *gorm.DB
	log *logrus.Logger
}

func Setup(param Params) DbInter {
	var err error
	dbrUri := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		param.Tuner.DB.Host,
		param.Tuner.DB.Port,
		param.Tuner.DB.User,
		param.Tuner.DB.Password,
		param.Tuner.DB.Database, param.Tuner.DB.SSlMode,
	)
	log.Println(dbrUri)
	db, err := gorm.Open(postgres.Open(dbrUri), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		param.Logger.Fatal("db.Setup err:", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	// AutoMigrate
	if err := autoMigrate(db); err != nil {

		param.Logger.Fatal("create model migrate err: ", err)

	}
	param.Logger.Println("DB successfully connected! ")

	return &gormDB{
		db:  db,
		log: param.Logger,
	}

}

func (d gormDB) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (d gormDB) GetDB() *gorm.DB {
	return d.db
}

func autoMigrate(db *gorm.DB) error {
	if err := db.Exec(fmt.Sprintf(`
		DO
		$$
		BEGIN
			IF NOT EXISTS (SELECT * FROM pg_type typ
				INNER JOIN pg_namespace nsp ON nsp.oid = typ.typnamespace
				WHERE nsp.nspname = current_schema() AND typ.typname = 'category') THEN
				CREATE TYPE category AS ENUM('%s', '%s', '%s', '%s');
			END IF;
			IF NOT EXISTS (SELECT * FROM pg_type typ
				INNER JOIN pg_namespace nsp ON nsp.oid = typ.typnamespace
				WHERE nsp.nspname = current_schema() AND typ.typname = 'ws_type') THEN
				CREATE TYPE ws_type AS ENUM('development', 'marketing', 'legal', 'support','back office');
			END IF;
		END;
		$$
		LANGUAGE plpgsql;
	`, projects.Project, projects.Product, projects.Platform, projects.BackOffice)).Error; err != nil {
		return err
	}
	if err := db.Exec(fmt.Sprintf(`
		DO
		$$
		BEGIN
			IF NOT EXISTS (SELECT * FROM pg_type typ
				INNER JOIN pg_namespace nsp ON nsp.oid = typ.typnamespace
				WHERE nsp.nspname = current_schema() AND typ.typname = 'enum_type') THEN
				CREATE TYPE enum_type AS ENUM('%s', '%s', '%s', '%s', '%s' );
			END IF;
		END;
		$$
		LANGUAGE plpgsql;
	`, projects.VentureBuilding, projects.ServiceDevt, projects.ServiceVB, projects.Social, projects.Internal)).Error; err != nil {
		return err
	}
	if err := addStageEnum(db); err != nil {
		return err
	}

	if err := addStatusEnum(db); err != nil {
		return err
	}
	if err := addAcStatusEnum(db); err != nil {
		return err
	}

	projectTypes, err := db.Migrator().ColumnTypes(projects.ProjectEntity{})
	if err != nil {
		return err
	}
	for _, pType := range projectTypes {
		if pType.Name() == "type" && pType.DatabaseTypeName() != "enum_type" {
			if err := db.Exec("ALTER TABLE projects DROP COLUMN type").Error; err != nil {
				return err
			}
		}
		if pType.Name() == "phase" && pType.DatabaseTypeName() != "enum_phases" {
			if err := db.Exec("ALTER TABLE projects DROP COLUMN phase").Error; err != nil {
				return err
			}
		}
		if pType.Name() == "stage" && pType.DatabaseTypeName() != "enum_stages" {
			if err := db.Exec("ALTER TABLE projects DROP COLUMN stage").Error; err != nil {
				return err
			}
		}
	}
	milestoneTypes, err := db.Migrator().ColumnTypes(milestone.MilestoneEntity{})
	if err != nil {
		return err
	}
	for _, mType := range milestoneTypes {
		if mType.Name() == "status" && mType.DatabaseTypeName() != "enum_status" {
			if err := db.Exec("ALTER TABLE milestone DROP COLUMN status").Error; err != nil {
				return err
			}
		}
	}
	actionPlanTypes, err := db.Migrator().ColumnTypes(actionPlan.ActionPlanEntity{})
	if err != nil {
		return err
	}
	for _, acType := range actionPlanTypes {
		if acType.Name() == "status" && acType.DatabaseTypeName() != "enum_ac_status" {
			if err := db.Exec("ALTER TABLE action_plan DROP COLUMN status").Error; err != nil {
				return err
			}
		}
	}

	for _, model := range []interface{}{
		(*workspace.WorkspaceEntity)(nil),
		(*actionPlan.ActionPlanEntity)(nil),
		(*projects.PhaseEntity)(nil),
		(*projects.ProjectEntity)(nil),
		(*stage.StageEntity)(nil),
		(*milestone.MilestoneEntity)(nil),
		(*epics.EpicEntity)(nil),
		(*tasks.TaskEntity)(nil),
		(*processes.ProcessEntity)(nil),
	} {
		dbSilent := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})
		if err := dbSilent.AutoMigrate(model); err != nil {
			return err
		}
	}
	if err := db.Where("name = 'building&launch'").FirstOrCreate(&projects.PhaseEntity{Name: "building&launch"}).Error; err != nil {
		return err
	}
	if err := db.Where("name = 'scale&growth'").FirstOrCreate(&projects.PhaseEntity{Name: "scale&growth"}).Error; err != nil {
		return err
	}

	return nil
}

func addStageEnum(db *gorm.DB) error {
	stages := []interface{}{
		projects.Ideation, projects.Concept, projects.Business, projects.Development, projects.Pilot,
		projects.CL, projects.EG, projects.Expansion, projects.ShakeOut, projects.Mature, projects.Ecosystem,
	}
	return db.Exec(fmt.Sprintf(`
		DO
		$$
		BEGIN
			IF NOT EXISTS (SELECT * FROM pg_type typ
				INNER JOIN pg_namespace nsp ON nsp.oid = typ.typnamespace
				WHERE nsp.nspname = current_schema() AND typ.typname = 'enum_stages') THEN
				CREATE TYPE enum_stages AS ENUM('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');
			END IF;
		END;
		$$
		LANGUAGE plpgsql;
	`, stages...)).Error
}

func addStatusEnum(db *gorm.DB) error {
	return db.Exec(fmt.Sprintf(`
		DO
		$$
		BEGIN
			IF NOT EXISTS (SELECT * FROM pg_type typ
				INNER JOIN pg_namespace nsp ON nsp.oid = typ.typnamespace
				WHERE nsp.nspname = current_schema() AND typ.typname = 'enum_status') THEN
				CREATE TYPE enum_status AS ENUM('%s', '%s', '%s', '%s', '%s');
			END IF;
		END;
		$$
		LANGUAGE plpgsql;
	`, milestone.NewStatus, milestone.Hold, milestone.Cancelled, milestone.InProgress, milestone.Completed)).Error
}

func addAcStatusEnum(db *gorm.DB) error {
	return db.Exec(fmt.Sprintf(`
		DO
		$$
		BEGIN
			IF NOT EXISTS (SELECT * FROM pg_type typ
				INNER JOIN pg_namespace nsp ON nsp.oid = typ.typnamespace
				WHERE nsp.nspname = current_schema() AND typ.typname = 'enum_ac_status') THEN
				CREATE TYPE enum_ac_status AS ENUM('%s', '%s', '%s');
			END IF;
		END;
		$$
		LANGUAGE plpgsql;
	`, actionPlan.Active, actionPlan.Archived, actionPlan.Draft)).Error
}
