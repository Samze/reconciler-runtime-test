/*
Copyright 2019 VMware, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package factories

import (
	"fmt"

	v1alpha1 "github.com/samze/reconciler-runtime-test/api/v1alpha1"
	"github.com/vmware-labs/reconciler-runtime/apis"
	rtesting "github.com/vmware-labs/reconciler-runtime/testing"
	testingfactories "github.com/vmware-labs/reconciler-runtime/testing/factories"
)

type bar struct {
	target *v1alpha1.Bar
}

var (
	_ rtesting.Factory = (*bar)(nil)
)

func Bar(seed ...*v1alpha1.Bar) *bar {
	var target *v1alpha1.Bar
	switch len(seed) {
	case 0:
		target = &v1alpha1.Bar{}
	case 1:
		target = seed[0]
	default:
		panic(fmt.Errorf("expected exactly zero or one seed, got %v", seed))
	}
	return &bar{
		target: target,
	}
}

func (f *bar) deepCopy() *bar {
	return Bar(f.target.DeepCopy())
}

func (f *bar) Create() *v1alpha1.Bar {
	return f.deepCopy().target
}

func (f *bar) CreateObject() apis.Object {
	return f.Create()
}

func (f *bar) mutation(m func(*v1alpha1.Bar)) *bar {
	f = f.deepCopy()
	m(f.target)
	return f
}

func (f *bar) NamespaceName(namespace, name string) *bar {
	return f.mutation(func(s *v1alpha1.Bar) {
		s.ObjectMeta.Namespace = namespace
		s.ObjectMeta.Name = name
	})
}

func (f *bar) ObjectMeta(nf func(testingfactories.ObjectMeta)) *bar {
	return f.mutation(func(s *v1alpha1.Bar) {
		omf := testingfactories.ObjectMetaFactory(s.ObjectMeta)
		nf(omf)
		s.ObjectMeta = omf.Create()
	})
}
func (f *bar) StatusFoo(foo string) *bar {
	return f.mutation(func(s *v1alpha1.Bar) {
		s.Status.FooStatus = foo
	})
}

func (f *bar) SpecFoo(foo string) *bar {
	return f.mutation(func(s *v1alpha1.Bar) {
		s.Spec.Foo = foo
	})
}
