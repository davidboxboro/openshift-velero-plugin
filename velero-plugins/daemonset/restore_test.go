// unit tests for restore.go in deployment

package daemonset

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
	assert.Equal(t, velero.ResourceSelector{IncludedResources: []string{"daemonsets.apps"}}, actual)
}

func TestRestorePluginExecute(t *testing.T) {
	restorePlugin := &RestorePlugin{Log: test.NewLogger()}

	tests := map[string]struct {
		daemonSet appsv1API.DaemonSet
		exp	  appsv1API.DaemonSet
	}{
		"1": {
			daemonSet: appsv1API.DaemonSet {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "foo",
						"openshift.io/restore-registry-hostname": "bar",
					},
				},
				Spec: appsv1API.DaemonSetSpec {
					Template: apiv1.PodTemplateSpec {
						Spec: apiv1.PodSpec {
							Containers: []apiv1.Container {
								apiv1.Container{Image: "foo/cat"},
							},
						},
					},
				},
			},
			exp: appsv1API.DaemonSet {
				ObjectMeta: metav1.ObjectMeta {
					Annotations: map[string]string{
						"openshift.io/backup-registry-hostname": "foo",
						"openshift.io/restore-registry-hostname": "bar",
					},
				},
                                Spec: appsv1API.DaemonSetSpec {
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
			daemonSetRec, _ := json.Marshal(tc.daemonSet) // Marshal it to JSON
			json.Unmarshal(daemonSetRec, &out) // Unmarshal into the proper format
			item.SetUnstructuredContent(out) // Set unstructured object
			restore := velerov1.Restore{}
			input := &velero.RestoreItemActionExecuteInput{Item: &item, Restore: &restore}

			output, _ := restorePlugin.Execute(input)

			outDaemonSet := appsv1API.DaemonSet{}
			itemMarshal, _ := json.Marshal(output.UpdatedItem)
			json.Unmarshal(itemMarshal, &outDaemonSet)

			if !reflect.DeepEqual(outDaemonSet, tc.exp) {
                                t.Fatalf("expected: %v, got: %v", tc.exp, outDaemonSet)
                        }
		})
        }
}

