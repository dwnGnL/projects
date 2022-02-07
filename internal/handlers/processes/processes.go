package processes

import (
	"net/http"
	"projects/internal/database/processes"
	"projects/internal/models"
	"projects/pkg/config"
	"projects/pkg/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

var Module = fx.Provide(NewProcessesHandler)

type ProcessesHandler interface {
	CreateProcess(c *gin.Context)
	ReadProcesses(c *gin.Context)
	UpdateProcess(c *gin.Context)
	DeleteProcess(c *gin.Context)
}

type Params struct {
	fx.In
	db.DbInter
	*config.Tuner
	*logrus.Logger
}

type processesHandler struct {
	db   db.DbInter
	log  *logrus.Logger
	conf *config.Tuner
}

func NewProcessesHandler(params Params) ProcessesHandler {
	return &processesHandler{db: params.DbInter, log: params.Logger, conf: params.Tuner}
}

func (p processesHandler) CreateProcess(c *gin.Context) {
	var proc models.ProcessReq
	if err := c.ShouldBindJSON(&proc); err != nil {
		p.log.Warnln(err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": "bind error"})
		return
	}

	ps := processes.New(p.db.GetDB())
	procEntity := processes.ProcessEntity{
		Name: proc.Name,
	}
	if err := ps.Create(&procEntity); err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, &models.ProcessResp{
		ProcessID: procEntity.ProcessID,
		Name:      procEntity.Name,
	})
}

func (p processesHandler) ReadProcesses(c *gin.Context) {
	ps := processes.New(p.db.GetDB())
	procs, err := ps.GetAll()
	if err != nil {
		p.log.Warnln("Get processes err", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	var procResp []models.ProcessResp
	for _, proc := range procs {
		procResp = append(procResp, models.ProcessResp{
			ProcessID: proc.ProcessID,
			Name:      proc.Name,
		})
	}

	c.JSON(http.StatusOK, procResp)
}

func (p processesHandler) UpdateProcess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	var procReq models.ProcessReq
	if err := c.ShouldBindJSON(&procReq); err != nil {
		p.log.Warnln("bind error")
		c.JSON(http.StatusBadGateway, gin.H{"error": "bind error"})
		return
	}
	ps := processes.New(p.db.GetDB())
	proc := processes.ProcessEntity{
		ProcessID: id,
		Name:      procReq.Name,
	}
	if err := ps.Update(&proc); err != nil {
		p.log.Warnln("process update err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusOK, models.ProcessResp{
		ProcessID: proc.ProcessID,
		Name:      proc.Name,
	})
}

func (p processesHandler) DeleteProcess(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		p.log.Warnln("Param err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ps := processes.New(p.db.GetDB())
	if err := ps.Delete(id); err != nil {
		p.log.Warnln("Can't delete process entity with err: ", err.Error())
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
