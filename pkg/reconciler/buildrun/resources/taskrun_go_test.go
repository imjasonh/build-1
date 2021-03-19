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

	standardParams := []v1beta1.ParamSpec{{
		Name:        "DOCKERFILE",
		Description: "Path to the Dockerfile",
		Default:     v1beta1.NewArrayOrString("Dockerfile"),
	}, {
		Name:        "CONTEXT_DIR",
		Description: "The root of the code",
		Default:     v1beta1.NewArrayOrString("."),
	}}

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
				t.Errorf("Params Diff(-want,+got): %s", d)
			}
		})
	}
}
