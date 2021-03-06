/*
Copyright 2019 The Tekton Authors

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

package reconciler

import (
	"context"
	"reflect"
	"strings"
	"testing"

	tb "github.com/tektoncd/pipeline/internal/builder/v1beta1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/pipeline/pkg/reconciler/pipelinerun/resources"
	ttesting "github.com/tektoncd/pipeline/pkg/reconciler/testing"
	test "github.com/tektoncd/pipeline/test"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"knative.dev/pkg/apis"
)

// Test case for providing recorder in the option
func TestRecorderOptions(t *testing.T) {

	prs := []*v1beta1.PipelineRun{tb.PipelineRun("test-pipeline-run-completed",
		tb.PipelineRunSpec("test-pipeline", tb.PipelineRunServiceAccountName("test-sa")),
		tb.PipelineRunStatus(tb.PipelineRunStatusCondition(apis.Condition{
			Type:    apis.ConditionSucceeded,
			Status:  corev1.ConditionTrue,
			Reason:  resources.ReasonSucceeded,
			Message: "All Tasks have completed executing",
		})),
	)}
	ps := []*v1beta1.Pipeline{tb.Pipeline("test-pipeline", tb.PipelineSpec(
		tb.PipelineTask("hello-world-1", "hellow-world"),
	))}
	ts := []*v1beta1.Task{tb.Task("hello-world")}
	d := test.Data{
		PipelineRuns: prs,
		Pipelines:    ps,
		Tasks:        ts,
	}
	ctx, _ := ttesting.SetupFakeContext(t)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	c, _ := test.SeedTestData(t, ctx, d)

	observer, _ := observer.New(zap.InfoLevel)

	// recorder is ont provided in the option
	b := NewBase(Options{
		Logger:            zap.New(observer).Sugar(),
		KubeClientSet:     c.Kube,
		PipelineClientSet: c.Pipeline,
	}, "test", pipeline.Images{})

	if strings.Compare(reflect.TypeOf(b.Recorder).String(), "*record.recorderImpl") != 0 {
		t.Errorf("Expected recorder type '*record.recorderImpl' but actual type is: %s", reflect.TypeOf(b.Recorder).String())
	}

	fr := record.NewFakeRecorder(1)

	// recorder is provided in the option
	b = NewBase(Options{
		Logger:            zap.New(observer).Sugar(),
		KubeClientSet:     c.Kube,
		PipelineClientSet: c.Pipeline,
		Recorder:          fr,
	}, "test", pipeline.Images{})

	if strings.Compare(reflect.TypeOf(b.Recorder).String(), "*record.FakeRecorder") != 0 {
		t.Errorf("Expected recorder type '*record.FakeRecorder' but actual type is: %s", reflect.TypeOf(b.Recorder).String())
	}
}
