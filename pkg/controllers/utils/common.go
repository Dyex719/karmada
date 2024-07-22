package utils

import (
	"context"
	"fmt"

	"github.com/karmada-io/karmada/pkg/util/helper"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func UpdateFailoverStatus(client client.Client, binding *workv1alpha2.ResourceBinding, cluster string, failoverType string) (err error) {
	message := fmt.Sprintf("Failover triggered for replica on cluster %s", cluster)

	var reason string
	if failoverType == workv1alpha2.EvictionReasonApplicationFailure {
		reason = "ApplicationFailover"
	} else if failoverType == workv1alpha2.EvictionReasonTaintUntolerated {
		reason = "ClusterFailover"
	} else {
		errMsg := "Invalid failover type passed into updateFailoverStatus"
		klog.Errorf(errMsg)
		return fmt.Errorf(errMsg)
	}

	newFailoverAppliedCondition := metav1.Condition{
		Type:               failoverType,
		Status:             metav1.ConditionTrue,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: metav1.Now(),
	}

	err = retry.RetryOnConflict(retry.DefaultRetry, func() (err error) {
		_, err = helper.UpdateStatus(context.Background(), client, binding, func() error {
			// set binding status with the newest condition
			currentTime := metav1.Now()
			failoverHistoryItem := workv1alpha2.FailoverHistoryItem{
				FailoverTime:  &currentTime,
				OriginCluster: cluster,
				Reason:        reason,
			}
			binding.Status.FailoverHistory = append(binding.Status.FailoverHistory, failoverHistoryItem)
			klog.V(4).Infof("Failover history is %+v", binding.Status.FailoverHistory)
			meta.SetStatusCondition(&binding.Status.Conditions, newFailoverAppliedCondition)
			return nil
		})
		return err
	})

	if err != nil {
		klog.Errorf("Failed to update condition of binding %s/%s: %s", binding.Namespace, binding.Name, err.Error())
		return err
	}
	return nil
}
