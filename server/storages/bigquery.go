package storages

import (
	"errors"
	"fmt"
	"github.com/jitsucom/jitsu/server/identifiers"

	"github.com/hashicorp/go-multierror"
	"github.com/jitsucom/jitsu/server/adapters"
	"github.com/jitsucom/jitsu/server/caching"
	"github.com/jitsucom/jitsu/server/events"
	"github.com/jitsucom/jitsu/server/logging"
	"github.com/jitsucom/jitsu/server/schema"
)

var disabledRecognitionConfiguration = &UserRecognitionConfiguration{enabled: false}

//BigQuery stores files to google BigQuery in two modes:
//batch: via google cloud storage in batch mode (1 file = 1 operation)
//stream: via events queue in stream mode (1 object = 1 operation)
type BigQuery struct {
	destinationID   string
	gcsAdapter      *adapters.GoogleCloudStorage
	bqAdapter       *adapters.BigQuery
	tableHelper     *TableHelper
	processor       *schema.Processor
	streamingWorker *StreamingWorker
	fallbackLogger  *logging.AsyncLogger
	eventsCache     *caching.EventsCache
	uniqueIDField   *identifiers.UniqueID
	staged          bool
}

func init() {
	RegisterStorage(BigQueryType, NewBigQuery)
}

//NewBigQuery returns BigQuery configured instance
func NewBigQuery(config *Config) (Storage, error) {
	gConfig := config.destination.Google
	if err := gConfig.Validate(config.streamMode); err != nil {
		return nil, err
	}

	if gConfig.Project == "" {
		return nil, errors.New("BigQuery project(bq_project) is required parameter")
	}

	//enrich with default parameters
	if gConfig.Dataset == "" {
		gConfig.Dataset = "default"
		logging.Warnf("[%s] dataset wasn't provided. Will be used default one: %s", config.destinationID, gConfig.Dataset)
	}

	var gcsAdapter *adapters.GoogleCloudStorage
	if !config.streamMode {
		var err error
		gcsAdapter, err = adapters.NewGoogleCloudStorage(config.ctx, gConfig)
		if err != nil {
			return nil, err
		}
	}

	queryLogger := config.loggerFactory.CreateSQLQueryLogger(config.destinationID)
	bigQueryAdapter, err := adapters.NewBigQuery(config.ctx, gConfig, queryLogger, config.sqlTypes)
	if err != nil {
		return nil, err
	}

	//create dataset if doesn't exist
	err = bigQueryAdapter.CreateDataset(gConfig.Dataset)
	if err != nil {
		bigQueryAdapter.Close()
		if gcsAdapter != nil {
			gcsAdapter.Close()
		}
		return nil, err
	}

	tableHelper := NewTableHelper(bigQueryAdapter, config.monitorKeeper, config.pkFields, adapters.SchemaToBigQueryString, config.streamMode, config.maxColumns)

	bq := &BigQuery{
		destinationID:  config.destinationID,
		gcsAdapter:     gcsAdapter,
		bqAdapter:      bigQueryAdapter,
		tableHelper:    tableHelper,
		processor:      config.processor,
		fallbackLogger: config.loggerFactory.CreateFailedLogger(config.destinationID),
		eventsCache:    config.eventsCache,
		uniqueIDField:  config.uniqueIDField,
		staged:         config.destination.Staged,
	}

	if config.streamMode {
		bq.streamingWorker = newStreamingWorker(config.eventQueue, config.processor, bq, config.eventsCache, config.loggerFactory.CreateStreamingArchiveLogger(config.destinationID), tableHelper)
		bq.streamingWorker.start()
	}

	return bq, nil
}

func (bq *BigQuery) DryRun(payload events.Event) ([]adapters.TableField, error) {
	return dryRun(payload, bq.processor, bq.tableHelper)
}

//Insert event in BigQuery
func (bq *BigQuery) Insert(dataSchema *adapters.Table, event events.Event) (err error) {
	dbTable, err := bq.tableHelper.EnsureTable(bq.ID(), dataSchema)
	if err != nil {
		return err
	}

	err = bq.bqAdapter.Insert(dbTable, event)

	//renew current db schema and retry
	if err != nil {
		dbTable, err := bq.tableHelper.RefreshTableSchema(bq.ID(), dataSchema)
		if err != nil {
			return err
		}

		dbTable, err = bq.tableHelper.EnsureTable(bq.ID(), dataSchema)
		if err != nil {
			return err
		}

		return bq.bqAdapter.Insert(dbTable, event)
	}

	return nil
}

//Store process events and stores with storeTable() func
//returns store result per table, failed events (group of events which are failed to process) and err
func (bq *BigQuery) Store(fileName string, objects []map[string]interface{}, alreadyUploadedTables map[string]bool) (map[string]*StoreResult, *events.FailedEvents, error) {
	flatData, failedEvents, err := bq.processor.ProcessEvents(fileName, objects, alreadyUploadedTables)
	if err != nil {
		return nil, nil, err
	}

	//update cache with failed events
	for _, failedEvent := range failedEvents.Events {
		bq.eventsCache.Error(bq.ID(), failedEvent.EventID, failedEvent.Error)
	}

	storeFailedEvents := true
	tableResults := map[string]*StoreResult{}
	for _, fdata := range flatData {
		table := bq.tableHelper.MapTableSchema(fdata.BatchHeader)
		err := bq.storeTable(fdata, table)
		tableResults[table.Name] = &StoreResult{Err: err, RowsCount: fdata.GetPayloadLen(), EventsSrc: fdata.GetEventsPerSrc()}
		if err != nil {
			storeFailedEvents = false
		}

		//events cache
		for _, object := range fdata.GetPayload() {
			if err != nil {
				bq.eventsCache.Error(bq.ID(), bq.uniqueIDField.Extract(object), err.Error())
			} else {
				bq.eventsCache.Succeed(bq.ID(), bq.uniqueIDField.Extract(object), object, table)
			}
		}
	}

	//store failed events to fallback only if other events have been inserted ok
	if storeFailedEvents {
		return tableResults, failedEvents, nil
	}

	return tableResults, nil, nil
}

//check table schema
//and store data into one table via google cloud storage
func (bq *BigQuery) storeTable(fdata *schema.ProcessedFile, table *adapters.Table) error {
	dbTable, err := bq.tableHelper.EnsureTable(bq.ID(), table)
	if err != nil {
		return err
	}

	b := fdata.GetPayloadBytes(schema.JSONMarshallerInstance)
	if err := bq.gcsAdapter.UploadBytes(fdata.FileName, b); err != nil {
		return err
	}

	if err := bq.bqAdapter.Copy(fdata.FileName, dbTable.Name); err != nil {
		return fmt.Errorf("Error copying file [%s] from gcp to bigquery: %v", fdata.FileName, err)
	}

	if err := bq.gcsAdapter.DeleteObject(fdata.FileName); err != nil {
		logging.SystemErrorf("[%s] file %s wasn't deleted from gcs: %v", bq.ID(), fdata.FileName, err)
	}

	return nil
}

//Update isn't supported
func (bq *BigQuery) Update(object map[string]interface{}) error {
	return errors.New("BigQuery doesn't support updates")
}

//SyncStore isn't supported
func (bq *BigQuery) SyncStore(overriddenDataSchema *schema.BatchHeader, objects []map[string]interface{}, timeIntervalValue string) error {
	return errors.New("BigQuery doesn't support sync store")
}

//GetUsersRecognition returns disabled users recognition configuration
func (bq *BigQuery) GetUsersRecognition() *UserRecognitionConfiguration {
	return disabledRecognitionConfiguration
}

//Fallback logs event with error to fallback logger
func (bq *BigQuery) Fallback(failedEvents ...*events.FailedEvent) {
	for _, failedEvent := range failedEvents {
		bq.fallbackLogger.ConsumeAny(failedEvent)
	}
}

//ID returns destination ID
func (bq *BigQuery) ID() string {
	return bq.destinationID
}

//Type returns BigQuery type
func (bq *BigQuery) Type() string {
	return BigQueryType
}

func (bq *BigQuery) IsStaging() bool {
	return bq.staged
}

//GetUniqueIDField returns unique ID field configuration
func (bq *BigQuery) GetUniqueIDField() *identifiers.UniqueID {
	return bq.uniqueIDField
}

//Close closes BigQuery adapter, fallback logger and streaming worker
func (bq *BigQuery) Close() (multiErr error) {
	if bq.gcsAdapter != nil {
		if err := bq.gcsAdapter.Close(); err != nil {
			multiErr = multierror.Append(multiErr, fmt.Errorf("[%s] Error closing google cloud storage client: %v", bq.ID(), err))
		}
	}

	if err := bq.bqAdapter.Close(); err != nil {
		multiErr = multierror.Append(multiErr, fmt.Errorf("[%s] Error closing BigQuery client: %v", bq.ID(), err))
	}

	if bq.streamingWorker != nil {
		if err := bq.streamingWorker.Close(); err != nil {
			multiErr = multierror.Append(multiErr, fmt.Errorf("[%s] Error closing streaming worker: %v", bq.ID(), err))
		}
	}

	if err := bq.fallbackLogger.Close(); err != nil {
		multiErr = multierror.Append(multiErr, fmt.Errorf("[%s] Error closing fallback logger: %v", bq.ID(), err))
	}

	return
}
