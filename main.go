package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/oxakromax/Backend_UipathMonitor/ORM"
	"github.com/oxakromax/Backend_UipathMonitor/UipathAPI"
)

var (
	BearerToken string
	APIUrl      = ""
	JobLastTime = make(map[uint]time.Time)
)

func pingAuth() bool {
	url := APIUrl + "/pingAuth"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return false
	}
	req.Header.Add("Authorization", "Bearer "+BearerToken)

	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return res.StatusCode == http.StatusOK
}

func Auth() string {
	APIUrl = os.Getenv("API_URL")
	url := APIUrl + "/auth"
	method := "POST"
	payload := strings.NewReader("email=" + os.Getenv("MONITOR_USER") + "&password=" + os.Getenv("MONITOR_PASS"))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	if res.StatusCode != 200 {
		fmt.Println("Error: ", res.StatusCode)
		fmt.Println(string(body))
		return ""
	}

	jsonStruct := make(map[string]interface{})
	json.Unmarshal(body, &jsonStruct)

	if token, ok := jsonStruct["token"].(string); ok {
		return token
	}

	fmt.Println("Token not found in response")
	return ""
}

func newIncident(incidente *ORM.TicketsProceso) {
	// POST /monitor/:id/newIncident (id = id del proceso)
	url := APIUrl + "/monitor/" + strconv.FormatUint(uint64(incidente.ProcesoID), 10) + "/newTicket"
	method := "POST"
	body, err := json.Marshal(incidente)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(string(body)))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+BearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.StatusCode != 200 {
		fmt.Println("Error: ", res.StatusCode)
		fmt.Println(string(body))
		return
	}
}

func JobKeyException(JobKey string) {
	// PUT /monitor/UpdateExceptionJob
	// Query: JobKey
	url := APIUrl + "/monitor/UpdateExceptionJob?JobKey=" + JobKey
	method := "PUT"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+BearerToken)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
}

func refreshTokenRoutine() {
	RefreshTime := time.Now().Add(24 * time.Hour)
	for {
		if !pingAuth() {
			BearerToken = Auth()
		}
		if RefreshTime.Before(time.Now()) {
			RefreshTime = time.Now().Add(24 * time.Hour)
			BearerToken = Auth()
		}
		time.Sleep(15 * time.Second)
	}
}

func getOrgs() []*ORM.Organizacion {
	url := APIUrl + "/monitor/Orgs"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	req.Header.Add("Authorization", "Bearer "+BearerToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if res.StatusCode != 200 {
		fmt.Println("Error: ", res.StatusCode)
		fmt.Println(string(body))
		return nil
	}

	var orgs []*ORM.Organizacion
	json.Unmarshal(body, &orgs)
	return orgs
}

func RefreshOrgs() {
	// PATCH /monitor/RefreshOrgs
	url := APIUrl + "/monitor/RefreshOrgs"
	method := "PATCH"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+BearerToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.StatusCode != 200 {
		fmt.Println("Error: ", res.StatusCode)
		fmt.Println(string(body))
		return
	}
}

func RefreshJobHistory() {
	// PATCH /monitor/PatchJobHistory
	url := APIUrl + "/monitor/PatchJobHistory"
	method := "PATCH"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+BearerToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.StatusCode != 200 {
		fmt.Println("Error: ", res.StatusCode)
		fmt.Println(string(body))
		return
	}
}

func refreshOrgsRoutine() {
	for {
		if pingAuth() {
			RefreshOrgs()
		}
		time.Sleep(1 * time.Minute)
	}
}

func Monitor() {
	if pingAuth() {
		RefreshJobHistory()
		wg := new(sync.WaitGroup)
		for _, org := range getOrgs() {
			wg.Add(1)
			go orgMonitor(org, wg)
		}
		wg.Wait()
	}
	time.Sleep(30 * time.Second)
}

func orgMonitor(org *ORM.Organizacion, wg *sync.WaitGroup) {
	defer wg.Done()
	FoldersAndProcesses := make(map[uint][]*ORM.Proceso)
	for _, process := range org.Procesos {
		if !process.ActiveMonitoring {
			continue
		}
		FoldersAndProcesses[process.Folderid] = append(FoldersAndProcesses[process.Folderid], process)
		incidentsProcess := process.TicketsProcesos
		sort.Slice(incidentsProcess, func(i, j int) bool {
			// Ordenamos los incidentes por fecha de creación, de más reciente a más antiguo
			return incidentsProcess[i].CreatedAt.After(incidentsProcess[j].CreatedAt)
		})
		if len(incidentsProcess) == 0 {
			continue
		}
		// buscar el incidente más reciente, de tipo 1 para establecerlo como fecha base
		for _, incident := range incidentsProcess {
			if incident.Tipo.ID == 1 { // incidente
				if JobLastTime[process.ID].Before(incident.UpdatedAt) {
					JobLastTime[process.ID] = incident.UpdatedAt // Se establece la fecha base, desde la ultima vez que se atendió la alerta
				}
				break
			}
		}

	}
	subwg := new(sync.WaitGroup)
	for folderid, processes := range FoldersAndProcesses {
		subwg.Add(1)
		go orgFolderRoutine(subwg, org, folderid, processes)
	}
	subwg.Wait()
}

func orgFolderRoutine(subwg *sync.WaitGroup, org *ORM.Organizacion, folderid uint, processes []*ORM.Proceso) {
	defer subwg.Done()
	JobsResponse := new(UipathAPI.JobsResponse)
	org.GetFromApi(JobsResponse, int(folderid))
	LogResponse := new(UipathAPI.LogResponse)
	org.GetFromApi(LogResponse, int(folderid))
	sort.Slice(JobsResponse.Value, func(i, j int) bool {
		// Ordenamos los jobs por fecha de creación, de más reciente a más antiguo
		return JobsResponse.Value[i].CreationTime.After(JobsResponse.Value[j].CreationTime)
	})

	ProcessJobKeyMap := make(map[uint]UipathAPI.JobsValue)
	for _, job := range JobsResponse.Value {
		var jobProcess ORM.Proceso
		for _, process := range processes {
			if process.Nombre == job.ReleaseName {
				jobProcess = *process
				break
			}
		}
		if !jobProcess.ActiveMonitoring {
			continue
		}

		if job.State == "Pending" {
			JobLastTime[jobProcess.ID] = job.CreationTime
			Now := time.Now()
			// if job is pending for more than 30 minutes, we report it
			if Now.Sub(job.CreationTime).Minutes() > float64(jobProcess.MaxQueueTime) {
				go ReportIncident(&jobProcess, "Ejecución pendiente no realizada", "El proceso está pendiente de ejecución desde hace más de "+strconv.Itoa(jobProcess.MaxQueueTime)+" minutos")
			}
			continue
		}

		// Check if start or end time of job is newer than last time
		if job.StartTime.Before(JobLastTime[jobProcess.ID]) && job.EndTime.Before(JobLastTime[jobProcess.ID]) {
			continue
		}

		switch job.State {
		case "Running", "Successful":
			JobLastTime[jobProcess.ID] = job.StartTime // we want to check if everything is ok until the job is finished
		case "Stopped":
			JobLastTime[jobProcess.ID] = job.EndTime.Add(1 * time.Second) // manually stopped, so we don't want to report it
		case "Faulted", "Error":
			go ReportIncident(&jobProcess, "Error de ejecución", *job.Info)
			JobLastTime[jobProcess.ID] = job.EndTime.Add(1 * time.Second) // it's already reported, so we don't want to report it again
		default:
			JobLastTime[jobProcess.ID] = time.Now()
		}
		ProcessJobKeyMap[jobProcess.ID] = job
	}
	for _, process := range processes {
		warnCounter := 0
		errorCounter := 0
		fatalCounter := 0
		var Reason string
		var Message string
		needStop := false
		if !process.ActiveMonitoring {
			continue
		}
		LastTimeLog := time.Time{}
		for _, log := range LogResponse.Value {
			if needStop {
				break
			}
			if log.ProcessName == process.Nombre {
				if log.TimeStamp.After(JobLastTime[process.ID]) {
					switch log.Level {
					case "Warn":
						warnCounter++
					case "Error":
						errorCounter++
					case "Fatal":
						fatalCounter++
					}
					switch {
					case fatalCounter >= process.FatalTolerance:
						Reason = "[" + log.TimeStamp.Format("2006-01-02 15:04:05") + "] \n" + log.Message
						Message = "Se ha superado el umbral de errores fatales en el proceso"
						go ReportIncident(process, Message, Reason)
						needStop = true
					case errorCounter >= process.ErrorTolerance:
						Reason = "[" + log.TimeStamp.Format("2006-01-02 15:04:05") + "] \n" + log.Message
						Message = "Se ha superado el umbral de errores en el proceso"
						go ReportIncident(process, Message, Reason)
						needStop = true
					case warnCounter >= process.WarningTolerance:
						Reason = "[" + log.TimeStamp.Format("2006-01-02 15:04:05") + "] \n" + log.Message
						Message = "Se ha superado el umbral de advertencias en el proceso"
						go ReportIncident(process, Message, Reason)
						needStop = true
					}
					if needStop {
						go JobKeyException(log.JobKey)
					}
				}
				if LastTimeLog.Before(log.TimeStamp) {
					LastTimeLog = log.TimeStamp
				}
			}
		}
		if needStop {
			continue // Go to next process
		}
		MaxTimeForLog := process.MaxQueueTime
		if MaxTimeForLog == 0 {
			continue
		}
		Now := time.Now()
		if Now.Sub(LastTimeLog).Minutes() > float64(MaxTimeForLog) {
			if val, ok := ProcessJobKeyMap[process.ID]; ok {
				if val.State != "Running" {
					continue
				}
				go JobKeyException(val.Key)
			} else {
				continue // Doesn't have a job running
			}
			TimeDifference := int(Now.Sub(LastTimeLog).Minutes())
			Reason = "El proceso no ha registrado logs en los últimos " + strconv.Itoa(TimeDifference) + " minutos"
			Message = "El proceso no ha registrado logs en los últimos " + strconv.Itoa(TimeDifference) + " minutos"
			go ReportIncident(process, Message, Reason)
		}

	}
}

func ReportIncident(process *ORM.Proceso, Message string, Reason string) {
	Incident := &ORM.TicketsProceso{
		Proceso:     process,
		ProcesoID:   process.ID,
		TipoID:      1,
		Descripcion: Message,
	}
	Detail := &ORM.TicketsDetalle{
		Detalle: Reason,
	}
	Incident.Detalles = append(Incident.Detalles, Detail)
	newIncident(Incident)
}

func main() {
	println("Starting...")
	previousToClose()
	envVars()
	if !initializeConnection() {
		return
	}
	go refreshTokenRoutine()
	go refreshOrgsRoutine()
	fmt.Println("Running...")
	for {
		Monitor()
	}
}

func initializeConnection() bool {
	BearerToken = Auth()
	if BearerToken == "" {
		fmt.Println("Error getting token")
		return false
	}
	return true
}

func envVars() {
	var err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file, please be sure ENV variables are set")
	}
}

func previousToClose() {
	// Captura señales de terminación.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Closing...")
		os.Exit(0)
	}()
}
