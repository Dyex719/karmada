package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/util/gclient"
)

func TestUpdateFailoverStatus(t *testing.T) {
	tests := []struct {
		name         string
		binding      *workv1alpha2.ResourceBinding
		cluster      string
		failoverType string
		wantErr      bool
	}{
		{
			name: "application failover",
			binding: &workv1alpha2.ResourceBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "binding",
					Namespace: "default",
				},
			},
			cluster:      "cluster1",
			failoverType: workv1alpha2.EvictionReasonApplicationFailure,
			wantErr:      false,
		},
		{
			name: "cluster failover",
			binding: &workv1alpha2.ResourceBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "binding",
					Namespace: "default",
				},
			},
			cluster:      "cluster2",
			failoverType: workv1alpha2.EvictionReasonTaintUntolerated,
			wantErr:      false,
		},
		{
			name: "invalid failover type",
			binding: &workv1alpha2.ResourceBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "binding",
					Namespace: "default",
				},
			},
			cluster:      "cluster3",
			failoverType: "InvalidType",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := fake.NewClientBuilder().WithScheme(gclient.NewSchema()).Build()

			err := UpdateFailoverStatus(fakeClient, tt.binding, tt.cluster, tt.failoverType)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
