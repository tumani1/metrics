package metrics

import (
	log "github.com/sirupsen/logrus"

	"github.com/trafficstars/fastmetrics/worker"
)

func runAndRegisterWorkerWrapper(key string, worker Worker, tags AnyTags) error {
	err := runAndRegister(key, worker, tags)
	if err != nil {
		worker.Stop()
		log.WithFields(log.Fields{
			"metric_key": key,
			"tags":       tags,
		}).Errorf(`Cannot register a metric "%s": %v`, key, err)
	}
	return err
}

func createWorkerCount(key string, tags AnyTags) (WorkerCount, error) {
	keyBuf := generateStorageKey("", key, tags)
	statsdKey := keyBuf.result.String()
	keyBuf.Unlock()
	worker := metricworker.NewWorkerCount(metrics.GetSender(), statsdKey)
	worker.SetGCEnabled(true)
	return worker, runAndRegisterWorkerWrapper(key, worker, tags)
}

func createWorkerGauge(key string, tags AnyTags) (WorkerGauge, error) {
	keyBuf := generateStorageKey("", key, tags)
	statsdKey := keyBuf.result.String()
	keyBuf.Unlock()
	worker := metricworker.NewWorkerGauge(metrics.GetSender(), statsdKey)
	worker.SetGCEnabled(true)
	return worker, runAndRegisterWorkerWrapper(key, worker, tags)
}

func createWorkerGaugeFunc(key string, tags AnyTags, fn func() int64) (WorkerGaugeFunc, error) {
	keyBuf := generateStorageKey("", key, tags)
	statsdKey := keyBuf.result.String()
	keyBuf.Unlock()
	worker := metricworker.NewWorkerGaugeFunc(metrics.GetSender(), statsdKey, fn)
	worker.SetGCEnabled(true)
	return worker, runAndRegisterWorkerWrapper(key, worker, tags)
}

func createWorkerTiming(key string, tags AnyTags) (WorkerTiming, error) {
	keyBuf := generateStorageKey("", key, tags)
	statsdKey := keyBuf.result.String()
	keyBuf.Unlock()
	worker := metricworker.NewWorkerTiming(metrics.GetSender(), statsdKey)
	worker.SetGCEnabled(true)
	return worker, runAndRegisterWorkerWrapper(key, worker, tags)
}

func CreateOrGetWorkerCountWithError(key string, tags AnyTags) (WorkerCount, error) {
	m := Get(MetricTypeCount, key, tags)
	if m != nil {
		return m.worker.(WorkerCount), nil
	}
	return createWorkerCount(key, tags)
}

func CreateOrGetWorkerCount(key string, tags AnyTags) WorkerCount {
	worker, _ := CreateOrGetWorkerCountWithError(key, tags)
	return worker
}

func CreateOrGetWorkerGaugeWithError(key string, tags AnyTags) (WorkerGauge, error) {
	m := Get(MetricTypeGauge, key, tags)
	if m != nil {
		return m.worker.(WorkerGauge), nil
	}
	return createWorkerGauge(key, tags)
}

func CreateOrGetWorkerGauge(key string, tags AnyTags) WorkerGauge {
	worker, _ := CreateOrGetWorkerGaugeWithError(key, tags)
	return worker
}

func CreateOrGetWorkerGaugeFuncWithError(key string, tags AnyTags, fn func() int64) (WorkerGaugeFunc, error) {
	m := Get(MetricTypeGauge, key, tags)
	if m != nil {
		return m.worker.(WorkerGaugeFunc), nil
	}
	return createWorkerGaugeFunc(key, tags, fn)
}

func CreateOrGetWorkerGaugeFunc(key string, tags AnyTags, fn func() int64) WorkerGaugeFunc {
	worker, _ := CreateOrGetWorkerGaugeFuncWithError(key, tags, fn)
	return worker
}

func CreateOrGetWorkerTimingWithError(key string, tags AnyTags) (WorkerTiming, error) {
	m := Get(MetricTypeTiming, key, tags)
	if m != nil {
		return m.worker.(WorkerTiming), nil
	}
	return createWorkerTiming(key, tags)
}

func CreateOrGetWorkerTiming(key string, tags AnyTags) WorkerTiming {
	worker, _ := CreateOrGetWorkerTimingWithError(key, tags)
	return worker
}
