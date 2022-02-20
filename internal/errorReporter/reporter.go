package errorReporter

import (
	"fmt"
	ErrorReporter "github.com/zytell3301/tg-error-reporter"
)

type reportError func(message string, parameters ...string)

type reporter struct {
	instanceId    string
	serviceId     string
	errorReporter ErrorReporter.Reporter
	reportError
}

var Reporter reporter

func NewReporter(instanceId string, serviceId string, errorReporter ErrorReporter.Reporter) reporter {
	r := reporter{
		instanceId:    instanceId,
		serviceId:     serviceId,
		errorReporter: errorReporter,
	}
	r.reportError = func(message string, parameters ...string) {
		r.errorReporter.Report(ErrorReporter.Error{
			ServiceId:  r.serviceId,
			InstanceId: r.instanceId,
			Message:    fmt.Sprintf(message, parameters),
		})
	}
	return r
}

func ReportError(message string, parameters ...string) {
	go Reporter.reportError(message, parameters...)
}
