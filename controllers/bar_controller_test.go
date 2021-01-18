package controllers_test

import (
	"testing"

	v1alpha1 "github.com/samze/reconciler-runtime-test/api/v1alpha1"
	"github.com/samze/reconciler-runtime-test/controllers"
	"github.com/vmware-labs/reconciler-runtime/reconcilers"
	rtesting "github.com/vmware-labs/reconciler-runtime/testing"
	"github.com/vmware-labs/reconciler-runtime/testing/factories"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestBarReconciler(t *testing.T) {

	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = v1alpha1.AddToScheme(scheme)

	testName := "bar-instance"
	testNamespace := "bar-test"
	testKey := types.NamespacedName{Namespace: testNamespace, Name: testName}

	bar := &v1alpha1.Bar{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName,
			Namespace: testNamespace,
			Labels:    map[string]string{"foo": "bar"},
		},
	}

	barFactory := rtesting.Wrapper(bar)

	expectedSecret := factories.Secret().
		ObjectMeta(func(om factories.ObjectMeta) {
			om.Namespace(testNamespace)
			om.Name("%s-secret", testName)
			om.AddLabel("foo", "bar")
			om.ControlledBy(barFactory, scheme)
		}).
		AddData("secret", "123")

	rts := rtesting.ReconcilerTestSuite{
		{
			Name: "bar creates secret",
			Key:  testKey,
			GivenObjects: []rtesting.Factory{
				//required objs
				barFactory,
			},
			ExpectTracks: []rtesting.TrackRequest{
				//rtesting.NewTrackRequest(inMemoryGatewayImagesConfigMap, inMemoryGateway, scheme),
			},
			ExpectEvents: []rtesting.Event{
				rtesting.NewEvent(barFactory, scheme, corev1.EventTypeNormal, "Created",
					`Created Secret "%s-secret"`, testName),
			},
			ExpectCreates: []rtesting.Factory{
				expectedSecret,
			},
			ExpectStatusUpdates: []rtesting.Factory{
				//barFactory.StatusObservedGeneration(1),
				// StatusConditions(
				// 	// the condition will be unknown since the child resource
				// 	// was just created and hasn't been reconciled by its
				// 	// controller yet
				// 	inMemoryGatewayConditionGatewayReady.Unknown(),
				// 	inMemoryGatewayConditionReady.Unknown(),
				// ),
			},
		},
		{
			Name: "error fetching bar",
			Key:  testKey,
			WithReactors: []rtesting.ReactionFunc{
				rtesting.InduceFailure("get", "Bar"),
			},
			GivenObjects: []rtesting.Factory{
				barFactory,
			},
			ShouldErr: true,
		},
	}

	rts.Test(t, scheme, func(t *testing.T, rtc *rtesting.ReconcilerTestCase, c reconcilers.Config) reconcile.Reconciler {
		return controllers.BarReconciler(c)
	})
}
