package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/medatechnology/suresql"

	"github.com/medatechnology/goutil/metrics"
	"github.com/medatechnology/goutil/simplelog"
	"github.com/medatechnology/simplehttp"
)

const (
	SUCCESS_EVENT = "success"
	ERROR_EVENT   = "error"
)

// Standardized handler state. NOTE: maybe next time move this to simplehttp. This is easier way to make end-points
// Which will always have:
// - binding the request parameter/body/file attachment?
// - logic
// - logging (optional) : can be console/log-files OR insert into DB
// - response back to caller
type HandlerState struct {
	Context             simplehttp.Context
	Header              *simplehttp.RequestHeader // Please make sure that the HeaderParser middleware is used! If not it's nil.
	Label               string                    // This is used for logging, it's like the title/function name
	User                string                    // Mostly user that call make the request, could be also user that connect to DB?
	LogMessage          string                    // for success event logs
	ErrorMessage        string                    // for error event logs + Err?
	ResponseMessage     string                    // for handler response (API/Endpoint response)
	Status              int                       // http status usually
	Err                 error                     // original error
	Data                interface{}
	DBLogging           bool
	ConsoleLogging      bool
	DBLoggingEvent      string              // success,error - can be multiple. Which event is logged
	ConsoleLoggingEvent string              // success,error - can be multiple. Which event is logged
	TableNames          string              // table in DB that is related, if applicable
	TimerID             int64               // if using timer, ie from Meda metrics
	Duration            float64             // if using timer, ie from Meda metrics
	Token               *suresql.TokenTable // for specific handlers that requires token
	LogTable            AccessLogTable      // TODO: put them here but somewhat abstract?
}

// This is the configuration for logging for the project
// This is state status for basic handlers that do not require status
func NewHandlerState(ctx simplehttp.Context, user, label, table string) HandlerState {
	if user == "" {
		user = suresql.CurrentNode.InternalConfig.Username
	}
	// This is the default setting, make sure the HeaderParser middleware is in used!
	return HandlerState{
		Context:             ctx,
		Label:               label,
		User:                user,
		DBLogging:           true,
		ConsoleLogging:      true,
		TableNames:          table,
		DBLoggingEvent:      SUCCESS_EVENT,
		ConsoleLoggingEvent: ERROR_EVENT + ", " + SUCCESS_EVENT,
		Header:              ctx.Get(simplehttp.REQUEST_HEADER_PARSED_STRING).(*simplehttp.RequestHeader),
		TimerID:             metrics.StartTimeIt("", 0),
	}
}

// This is the state status for handlers that requires token.
func NewHandlerTokenState(ctx simplehttp.Context, label, table string) HandlerState {
	// This is the default setting, make sure the HeaderParser middleware is in used!
	state := HandlerState{
		Context:             ctx,
		Label:               label,
		User:                suresql.CurrentNode.InternalConfig.Username,
		DBLogging:           true,
		ConsoleLogging:      true,
		TableNames:          table,
		DBLoggingEvent:      SUCCESS_EVENT,
		ConsoleLoggingEvent: ERROR_EVENT + ", " + SUCCESS_EVENT,
		Header:              ctx.Get(simplehttp.REQUEST_HEADER_PARSED_STRING).(*simplehttp.RequestHeader),
		Token:               ctx.Get(TOKEN_TABLE_STRING).(*suresql.TokenTable),
		TimerID:             metrics.StartTimeIt("", 0),
	}
	// This is important, if not it will get the real username used to connect to DBMS
	if state.Token != nil {
		state.User = state.Token.UserName
	}
	return state
}

// This is the configuration for logging for the middleware
func NewMiddlewareState(ctx simplehttp.Context, name string) HandlerState {
	// This is the default setting, make sure the HeaderParser middleware is in used!
	return HandlerState{
		Context:             ctx,
		User:                suresql.CurrentNode.InternalConfig.Username,
		Label:               "middleware",
		TableNames:          name,
		DBLogging:           false,                              // no DB logging
		ConsoleLogging:      true,                               // has console logging
		ConsoleLoggingEvent: ERROR_EVENT + ", " + SUCCESS_EVENT, // Production: only ERROR_EVENTS for hacking checks
		Header:              ctx.Get(simplehttp.REQUEST_HEADER_PARSED_STRING).(*simplehttp.RequestHeader),
		// TimerID:             metrics.StartTimeIt("", 0),
	}
}

// Readibility for the state logging configuration
func (h *HandlerState) IsErrorLoggedInConsole() bool {
	return strings.Contains(h.ConsoleLoggingEvent, ERROR_EVENT) && h.ConsoleLogging
}
func (h *HandlerState) IsSuccessLoggedInConsole() bool {
	return strings.Contains(h.ConsoleLoggingEvent, SUCCESS_EVENT) && h.ConsoleLogging
}
func (h *HandlerState) IsSuccessLoggedInDB() bool {
	return strings.Contains(h.DBLoggingEvent, SUCCESS_EVENT) && h.DBLogging
}
func (h *HandlerState) IsErrorLoggedInDB() bool {
	return strings.Contains(h.DBLoggingEvent, ERROR_EVENT) && h.DBLogging
}

// Stopping the timer if not already stopped. This function is saved to be
// called multiple times!
func (h *HandlerState) SaveStopTimer() float64 {
	if h.TimerID != 0 {
		h.Duration = float64(metrics.StopTimeIt(h.TimerID))
		h.TimerID = 0
	}
	return h.Duration
}

// Implement this method according to the project logTable AND console logging
// Separate the Logging part and Response part (though it can be combined to make it faster)
// because sometimes the HandlerState is only for logging before returning anything.
// Like in the handler it has multiple log entries then at the end send back response
func (h *HandlerState) OnlyLog(message string, data interface{}, restartTimer bool) error {
	// if restartTimer=true then we will stop the timer, then later start it again. The handler will have multiple
	if restartTimer {
		h.SaveStopTimer()
	}
	// fmt.Println("called onlylog message:", message, ", data:", data, ", timer:", restartTimer)
	// Implement the logging for project. Most projects have different logging
	logEntry := AccessLogTable{
		Username:   h.User,
		Duration:   h.Duration,
		ActionType: h.Label,
		Occurred:   time.Now(),
		Table:      h.TableNames,
		// Note:        note,
		// Description: description,
		// Result:      result,
		// ResultStatus:  ERROR_EVENT,
		Method:        h.Context.Request().Method,
		ClientIP:      h.Header.RemoteIP, // NOTE: is this accurate??
		ClientBrowser: h.Header.UserAgent,
		ClientDevice:  h.Header.Device,
		NodeNumber:    suresql.CurrentNode.Config.NodeNumber,
		// Error:         h.ErrorMessage,
		// RawQuery: ,
	}
	// if data is passed, use this is for the RAW_QUERY_LOG. NOTE: this is a bit ambiguous
	if data != nil && LOG_RAW_QUERY {
		logEntry.RawQuery = fmt.Sprintf("%v", data)
		// if logEntry.Description == "" {
		// 	logEntry.Description = fmt.Sprintf("%v", data)
		// } else {
		// 	logEntry.Note = fmt.Sprintf("%v", data)
		// }
	}
	var err error
	// Error event
	if h.Err != nil {
		// if ErrorMessage is not set (maybe didn't call SetError) then errorMessage = message (parameter)
		if message != "" && h.ErrorMessage == "" {
			h.ErrorMessage = message
		}
		logEntry.ResultStatus = ERROR_EVENT
		logEntry.Result = h.ErrorMessage
		logEntry.Error = h.Err.Error()
		if h.IsErrorLoggedInDB() {
			// fmt.Println("error logged in DB")
			err = logEntry.DBLogging(&suresql.CurrentNode.InternalConnection)
		}
		if h.IsErrorLoggedInConsole() {
			// fmt.Println("error logged in Console")
			// NOTE: data is not logged in console
			simplelog.LogErrorAny(h.Label, h.Err, h.ErrorMessage)
		}
	} else {
		// Success event
		if message != "" {
			h.LogMessage = message
		}
		logEntry.ResultStatus = SUCCESS_EVENT
		logEntry.Result = h.LogMessage
		if h.IsSuccessLoggedInDB() {
			// fmt.Println("success logged in DB")
			err = logEntry.DBLogging(&suresql.CurrentNode.InternalConnection)
		}
		if h.IsSuccessLoggedInConsole() {
			// fmt.Println("success logged in Console")
			// NOTE: data is not logged in console
			simplelog.LogThis(h.Label, h.LogMessage)
		}
	}
	if restartTimer {
		h.TimerID = metrics.StartTimeIt("", 0)
	}
	return err
}

// For chaining calls, this message in parameter is used for response
func (h *HandlerState) SetError(msg string, err error, status int) *HandlerState {
	h.Err = err
	if status == 0 {
		status = http.StatusBadRequest
	}
	h.ErrorMessage = msg
	h.Status = status
	h.Data = err
	return h
}

// For chaining calls, this message parameter is for response
func (h *HandlerState) SetSuccess(msg string, data interface{}) *HandlerState {
	h.ResponseMessage = msg
	h.Data = data
	h.Status = http.StatusOK
	return h
}

// Return response based on context, message and data is used for logging (not response)
// logAgain if want to run the logging once more and return (if not yet called)
func (h *HandlerState) LogAndResponse(message string, data interface{}, logAgain bool) error {
	h.SaveStopTimer()
	if logAgain {
		h.OnlyLog(message, data, false) // always ignore the error for logging, DO NOT restart timer, we are giving response and exit
	}
	// If state.responsemessage is not set then set it the same as this.
	if h.ResponseMessage == "" {
		h.ResponseMessage = message
	}
	if h.Data == nil {
		h.Data = data
	}
	resp := suresql.StandardResponse{
		Status:  h.Status,
		Message: h.ResponseMessage,
		Data:    h.Data,
	}
	if h.Err != nil {
		// Error Event
		if resp.Status == 0 {
			resp.Status = http.StatusBadRequest
		}
		resp.Data = h.Err
	}
	return h.Context.JSON(resp.Status, resp)
}
