// Copyright The Shipwright Contributors
//
// SPDX-License-Identifier: Apache-2.0
package resources

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	buildv1alpha1 "github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	"github.com/shipwright-io/build/pkg/config"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

// TestGenerateTaskRun_Params tests param handling scenarios and how they
// generate TaskRunSpec.TaskSpec.Params specs, and TaskRunSpec.Params values.
func TestGenerateTaskRun_Params(t *testing.T) {
	cfg := &config.Config{}
	const serviceAccountName = "build-bot"

	// These params are always defined, whether or not they're used.
	standardParams := []v1beta1.ParamSpec{{
		Name:        "DOCKERFILE",
		Description: "Path to the Dockerfile",
		Default:     v1beta1.NewArrayOrString("Dockerfile"),
	}, {
		Name:        "CONTEXT_DIR",
		Description: "The root of the code",
		Default:     v1beta1.NewArrayOrString("."),
	}}

	// strPtr returns a pointer to the string, to satisfy param.default.
	strPtr := func(s string) *string { return &s }

	for _, c := range []struct {
		name            string
		strategy        buildv1alpha1.BuildStrategyInterface
		buildSpec       buildv1alpha1.BuildSpec
		buildRunSpec    buildv1alpha1.BuildRunSpec
		wantParams      []v1beta1.ParamSpec
		wantParamValues []v1beta1.Param
	}{{
		name: "build strategy specifies no params",
		strategy: &buildv1alpha1.BuildStrategy{
			Spec: buildv1alpha1.BuildStrategySpec{},
		},
		wantParams:      standardParams,
		wantParamValues: nil,
	}, {
		name: "build strategy with param with default",
		strategy: &buildv1alpha1.BuildStrategy{
			Spec: buildv1alpha1.BuildStrategySpec{
				Params: []buildv1alpha1.Param{{
					Name:    "param",
					Default: strPtr("default value"),
				}},
			},
		},
		wantParams: append(standardParams, v1alpha1.ParamSpec{
			Name:    "param",
			Default: v1beta1.NewArrayOrString("default value"),
		}),
		wantParamValues: nil,
	}, {
		name: "build strategy with param, build sets",
		strategy: &buildv1alpha1.BuildStrategy{
			Spec: buildv1alpha1.BuildStrategySpec{
				Params: []buildv1alpha1.Param{{
					Name: "param",
				}},
			},
		},
		buildSpec: buildv1alpha1.BuildSpec{
			Params: []buildv1alpha1.ParamValue{{Name: "param", Value: "build value"}},
		},
		wantParams: append(standardParams, v1alpha1.ParamSpec{
			Name: "param",
		}),
		wantParamValues: []v1alpha1.Param{{
			Name:  "param",
			Value: *v1beta1.NewArrayOrString("build value"),
		}},
	}, {
		name: "build strategy with param, buildrun sets",
		strategy: &buildv1alpha1.BuildStrategy{
			Spec: buildv1alpha1.BuildStrategySpec{
				Params: []buildv1alpha1.Param{{
					Name: "param",
				}},
			},
		},
		buildSpec: buildv1alpha1.BuildSpec{},
		buildRunSpec: buildv1alpha1.BuildRunSpec{
			Params: []buildv1alpha1.ParamValue{{Name: "param", Value: "buildrun value"}},
		},
		wantParams: append(standardParams, v1alpha1.ParamSpec{
			Name: "param",
		}),
		wantParamValues: []v1alpha1.Param{{
			Name:  "param",
			Value: *v1beta1.NewArrayOrString("buildrun value"),
		}},
	}, {
		name: "build strategy with param, build sets, buildrun overrides",
		strategy: &buildv1alpha1.BuildStrategy{
			Spec: buildv1alpha1.BuildStrategySpec{
				Params: []buildv1alpha1.Param{{
					Name: "param",
				}},
			},
		},
		buildSpec: buildv1alpha1.BuildSpec{
			Params: []buildv1alpha1.ParamValue{{Name: "param", Value: "build value"}},
		},
		buildRunSpec: buildv1alpha1.BuildRunSpec{
			Params: []buildv1alpha1.ParamValue{{Name: "param", Value: "buildrun value"}},
		},
		wantParams: append(standardParams, v1alpha1.ParamSpec{
			Name: "param",
		}),
		wantParamValues: []v1alpha1.Param{{
			Name:  "param",
			Value: *v1beta1.NewArrayOrString("buildrun value"),
		}},
	}, {
		// For completeness, below are examples of how other Build
		// features turn into TaskRun params:

		// Specifying build.Dockerfile sets a non-default param value
		// for the DOCKERFILE param.
		name: "build specifies dockerfile -> taskrun param",
		strategy: &buildv1alpha1.BuildStrategy{
			Spec: buildv1alpha1.BuildStrategySpec{},
		},
		buildSpec: buildv1alpha1.BuildSpec{
			Dockerfile: strPtr("path/to/Dockerfile"),
		},
		wantParams: standardParams,
		wantParamValues: []v1alpha1.Param{{
			Name:  "DOCKERFILE",
			Value: *v1beta1.NewArrayOrString("path/to/Dockerfile"),
		}},
	}, {
		// Specifying build.builder adds a param for BUILDER_IMAGE with
		// the default set to the value, _and_ adds a param value to
		// specify the value (it doesn't need to).
		name: "build specifies builder -> taskrun param",
		strategy: &buildv1alpha1.BuildStrategy{
			Spec: buildv1alpha1.BuildStrategySpec{},
		},
		buildSpec: buildv1alpha1.BuildSpec{
			BuilderImage: &buildv1alpha1.Image{
				ImageURL: "builder/image",
			},
		},
		wantParams: append(standardParams, v1alpha1.ParamSpec{
			Name:        "BUILDER_IMAGE",
			Description: "Image containing the build tools/logic",
			Default:     v1beta1.NewArrayOrString("builder/image"),
		}),
		wantParamValues: []v1beta1.Param{{
			Name:  "BUILDER_IMAGE",
			Value: *v1beta1.NewArrayOrString("builder/image"),
		}},
	}, {
		// Specifying build.source.contextDir sets a non-default param
		// value for the CONTEXT_DIR param.
		name: "build specifies contextDir -> taskrun param",
		strategy: &buildv1alpha1.BuildStrategy{
			Spec: buildv1alpha1.BuildStrategySpec{},
		},
		buildSpec: buildv1alpha1.BuildSpec{
			Source: buildv1alpha1.GitSource{
				ContextDir: strPtr("path/to/context"),
			},
		},
		wantParams: standardParams,
		wantParamValues: []v1beta1.Param{{
			Name:  "CONTEXT_DIR",
			Value: *v1beta1.NewArrayOrString("path/to/context"),
		}},
	}} {
		t.Run(c.name, func(t *testing.T) {
			got, err := GenerateTaskRun(cfg, &buildv1alpha1.Build{Spec: c.buildSpec}, &buildv1alpha1.BuildRun{Spec: c.buildRunSpec}, serviceAccountName, c.strategy)
			if err != nil {
				t.Fatalf("GenerateTaskRun: %v", err)
			}

			if d := cmp.Diff(c.wantParams, got.Spec.TaskSpec.Params); d != "" {
				t.Errorf("Params Diff(-want,+got): %s", d)
			}
			if d := cmp.Diff(c.wantParamValues, got.Spec.Params); d != "" {
				t.Errorf("Param values Diff(-want,+got): %s", d)
			}
		})
	}
}
