package unstructured

import (
	"encoding/json"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/google/go-cmp/cmp"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"
)

var _ resource.Composite = &Composite{}

func emptyXR() *Composite {
	return &Composite{Unstructured: unstructured.Unstructured{Object: map[string]interface{}{}}}
}

func TestCompositeCondition(t *testing.T) {
	cases := map[string]struct {
		reason string
		u      *Composite
		set    []xpv1.Condition
		get    xpv1.ConditionType
		want   xpv1.Condition
	}{
		"NewCondition": {
			reason: "It should be possible to set a condition of an emptyXR Unstructured.",
			u:      emptyXR(),
			set:    []xpv1.Condition{xpv1.Available(), xpv1.ReconcileSuccess()},
			get:    xpv1.TypeReady,
			want:   xpv1.Available(),
		},
		"ExistingCondition": {
			reason: "It should be possible to overwrite a condition that is already set.",
			u: func() *Composite {
				c := emptyXR()
				c.SetConditions(xpv1.Creating())
				return c
			}(),
			set:  []xpv1.Condition{xpv1.Available()},
			get:  xpv1.TypeReady,
			want: xpv1.Available(),
		},
		"WeirdStatus": {
			reason: "It should not be possible to set a condition when status is not an object.",
			u: &Composite{unstructured.Unstructured{Object: map[string]interface{}{
				"status": "wat",
			}}},
			set:  []xpv1.Condition{xpv1.Available()},
			get:  xpv1.TypeReady,
			want: xpv1.Condition{},
		},
		"WeirdStatusConditions": {
			reason: "Conditions should be overwritten if they are not an object.",
			u: &Composite{unstructured.Unstructured{Object: map[string]interface{}{
				"status": map[string]interface{}{
					"conditions": "wat",
				},
			}}},
			set:  []xpv1.Condition{xpv1.Available()},
			get:  xpv1.TypeReady,
			want: xpv1.Available(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.u.SetConditions(tc.set...)
			got := tc.u.GetCondition(tc.get)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("\n%s\nu.GetCondition(%s): -want, +got:\n%s", tc.reason, tc.get, diff)
			}
		})
	}
}

func TestCompositeConditions(t *testing.T) {
	cases := map[string]struct {
		reason string
		u      *Composite
		want   []xpv1.Condition
	}{
		"NewCondition": {
			reason: "It should be possible to get conditions.",
			u: func() *Composite {
				c := emptyXR()
				c.SetConditions(xpv1.Available(), xpv1.ReconcileSuccess())
				return c
			}(),
			want: []xpv1.Condition{xpv1.Available(), xpv1.ReconcileSuccess()},
		},
		"WeirdStatus": {
			reason: "It should not be possible to get conditions when status is not an object.",
			u: &Composite{unstructured.Unstructured{Object: map[string]interface{}{
				"status": "wat",
			}}},
			want: nil,
		},
		"WeirdStatusConditions": {
			reason: "It should notbe possible to get conditions when they are not an array.",
			u: &Composite{unstructured.Unstructured{Object: map[string]interface{}{
				"status": map[string]interface{}{
					"conditions": "wat",
				},
			}}},
			want: nil,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := tc.u.GetConditions()
			if diff := cmp.Diff(tc.want, got, test.EquateConditions()); diff != "" {
				t.Errorf("\n%s\nu.GetConditions(): -want, +got:\n%s", tc.reason, diff)
			}
		})
	}
}

func TestCompositeCompositionSelector(t *testing.T) {
	sel := &metav1.LabelSelector{MatchLabels: map[string]string{"cool": "very"}}
	cases := map[string]struct {
		u    *Composite
		set  *metav1.LabelSelector
		want *metav1.LabelSelector
	}{
		"NewSel": {
			u:    emptyXR(),
			set:  sel,
			want: sel,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.u.SetCompositionSelector(tc.set)
			got := tc.u.GetCompositionSelector()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("\nu.GetCompositionSelector(): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestCompositeCompositionReference(t *testing.T) {
	ref := &corev1.ObjectReference{Namespace: "ns", Name: "cool"}
	cases := map[string]struct {
		u    *Composite
		set  *corev1.ObjectReference
		want *corev1.ObjectReference
	}{
		"NewRef": {
			u:    emptyXR(),
			set:  ref,
			want: ref,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.u.SetCompositionReference(tc.set)
			got := tc.u.GetCompositionReference()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("\nu.GetCompositionReference(): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestCompositeClaimReference(t *testing.T) {
	ref := &corev1.ObjectReference{Namespace: "ns", Name: "cool"}
	cases := map[string]struct {
		u    *Composite
		set  *corev1.ObjectReference
		want *corev1.ObjectReference
	}{
		"NewRef": {
			u:    emptyXR(),
			set:  ref,
			want: ref,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.u.SetClaimReference(tc.set)
			got := tc.u.GetClaimReference()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("\nu.GetClaimReference(): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestCompositeResourceReferences(t *testing.T) {
	ref := corev1.ObjectReference{Namespace: "ns", Name: "cool"}
	cases := map[string]struct {
		u    *Composite
		set  []corev1.ObjectReference
		want []corev1.ObjectReference
	}{
		"NewRef": {
			u:    emptyXR(),
			set:  []corev1.ObjectReference{ref},
			want: []corev1.ObjectReference{ref},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.u.SetResourceReferences(tc.set)
			got := tc.u.GetResourceReferences()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("\nu.GetResourceReferences(): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestCompositeWriteConnectionSecretToReference(t *testing.T) {
	ref := &xpv1.SecretReference{Namespace: "ns", Name: "cool"}
	cases := map[string]struct {
		u    *Composite
		set  *xpv1.SecretReference
		want *xpv1.SecretReference
	}{
		"NewRef": {
			u:    emptyXR(),
			set:  ref,
			want: ref,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.u.SetWriteConnectionSecretToReference(tc.set)
			got := tc.u.GetWriteConnectionSecretToReference()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("\nu.GetWriteConnectionSecretToReference(): -want, +got:\n%s", diff)
			}
		})
	}
}

func TestCompositeConnectionDetailsLastPublishedTime(t *testing.T) {
	now := &metav1.Time{Time: time.Now()}

	// The timestamp loses a little resolution when round-tripped through JSON
	// encoding.
	lores := func(t *metav1.Time) *metav1.Time {
		out := &metav1.Time{}
		j, _ := json.Marshal(t)
		_ = json.Unmarshal(j, out)
		return out
	}

	cases := map[string]struct {
		u    *Composite
		set  *metav1.Time
		want *metav1.Time
	}{
		"NewTime": {
			u:    emptyXR(),
			set:  now,
			want: lores(now),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.u.SetConnectionDetailsLastPublishedTime(tc.set)
			got := tc.u.GetConnectionDetailsLastPublishedTime()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("\nu.GetConnectionDetailsLastPublishedTime(): -want, +got:\n%s", diff)
			}
		})
	}
}
