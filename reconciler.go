/*
Copyright 2021 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/logging"
	kreconciler "knative.dev/pkg/reconciler"
)

// TODO: move this to its own package, with tests.
type Reconciler struct{}

// ReconcileKind implements Interface.ReconcileKind.
func (c *Reconciler) ReconcileKind(ctx context.Context, r *v1beta1.PipelineRun) kreconciler.Event {
	logger := logging.FromContext(ctx)
	logger.Infof("Reconciling %s/%s", r.Namespace, r.Name)

	cfg := FromContextOrDefaults(ctx)
	logger.Infof("key=%q, limit=%d", cfg.Concurrency.Key, cfg.Concurrency.Limit)

	if r.IsPending() {
		logger.Info("PipelineRun is Pending, setting it to start")
		r.Spec.Status = ""
	}

	return kreconciler.NewEvent(corev1.EventTypeNormal, "PipelineRunReconciled", "PipelineRun reconciled: \"%s/%s\"", r.Namespace, r.Name)
}
