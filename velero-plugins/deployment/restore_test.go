// unit tests for restore.go in deployment

package deployment

import (
	"testing"
	"github.com/konveyor/openshift-velero-plugin/velero-plugins/util/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
        appsv1API "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"encoding/json"
	"reflect"
        velerov1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	//"fmt"
	//"k8s.io/apimachinery/pkg/runtime"
)

func TestRestorePluginAppliesTo(t *testing.T) {
	restorePlugin := &RestorePlugin{Log: test.NewLogger()}
	actual, err := restorePlugin.AppliesTo()
	require.NoError(t, err)
	assert.Equal(t, velero.ResourceSelector{IncludedResources: []string{"deployments.apps"}}, actual)
}

func TestRestorePluginExecute(t *testing.T) {
	restorePlugin := &RestorePlugin{Log: test.NewLogger()}

	tests := map[string]struct{
		deployment appsv1API.Deployment
		exp	   appsv1API.Deployment
	}{
		"1": {
			deployment: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "foo",
						"openshift.io/restore-registry-hostname": "bar",
					},
				},
				Spec: appsv1API.DeploymentSpec {
					Template: apiv1.PodTemplateSpec {
						Spec: apiv1.PodSpec {
							Containers: []apiv1.Container {
								apiv1.Container{Image: "foo/cat"},
							},
						},
					},
				},
			},
			exp: appsv1API.Deployment {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "foo",
						"openshift.io/restore-registry-hostname": "bar",
					},
				},
                                Spec: appsv1API.DeploymentSpec {
                                        Template: apiv1.PodTemplateSpec {
                                                Spec: apiv1.PodSpec {
                                                        Containers: []apiv1.Container {
                                                                apiv1.Container{Image: "bar/cat"},
                                                        },
                                                },
                                        },
                                },
                        },
		},
	}


	for name, tc := range tests {
                t.Run(name, func(t *testing.T) {
			var out map[string]interface{}
			item := unstructured.Unstructured{}
			deploymentRec, _ := json.Marshal(tc.deployment) // Marshal it to JSON
			json.Unmarshal(deploymentRec, &out) // Unmarshal into the proper format
			item.SetUnstructuredContent(out) // Set unstructured object
			restore := velerov1.Restore{}
			input := &velero.RestoreItemActionExecuteInput{Item: &item, Restore: &restore}

			output, _ := restorePlugin.Execute(input)

			deployment := appsv1API.Deployment{}
			itemMarshal, _ := json.Marshal(output.UpdatedItem)
			json.Unmarshal(itemMarshal, &deployment)

			if !reflect.DeepEqual(deployment, tc.exp) {
                                t.Fatalf("expected: %v, got: %v", tc.exp, deployment)
                        }
		})
        }
}

func int32Ptr(i int32) *int32 { return &i }
