package util

import (
	"context"

	"github.com/onsi/ginkgo"
	ginko "github.com/onsi/ginkgo"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	jiraservicedeskv1alpha1 "github.com/stakater/jira-service-desk-operator/api/v1alpha1"
)

// TestUtil contains necessary objects required to perform operations during tests
type TestUtil struct {
	ctx       context.Context
	k8sClient client.Client
	r         reconcile.Reconciler
}

// New creates new TestUtil
func New(ctx context.Context, k8sClient client.Client, r reconcile.Reconciler) *TestUtil {
	return &TestUtil{
		ctx:       ctx,
		k8sClient: k8sClient,
		r:         r,
	}
}

// CreateNamespace creates a namespace in the kubernetes server
func (t *TestUtil) CreateNamespace(name string) {
	namespaceObject := t.CreateNamespaceObject(name)
	err := t.k8sClient.Create(t.ctx, namespaceObject)

	if err != nil {
		ginkgo.Fail(err.Error())
	}
}

// CreateNamespaceObject creates a namespace object
func (t *TestUtil) CreateNamespaceObject(name string) *v1.Namespace {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

// CreateProjectObject creates a jira project custom resource object
func (t *TestUtil) CreateProjectObject(project jiraservicedeskv1alpha1.Project, namespace string) *jiraservicedeskv1alpha1.Project {
	return &jiraservicedeskv1alpha1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name:      project.Spec.Name,
			Namespace: namespace,
		},
		Spec: jiraservicedeskv1alpha1.ProjectSpec{
			Name:                project.Spec.Name,
			Key:                 project.Spec.Key,
			ProjectTypeKey:      project.Spec.ProjectTypeKey,
			ProjectTemplateKey:  project.Spec.ProjectTemplateKey,
			Description:         project.Spec.Description,
			AssigneeType:        project.Spec.AssigneeType,
			LeadAccountId:       project.Spec.LeadAccountId,
			URL:                 project.Spec.URL,
			AvatarId:            project.Spec.AvatarId,
			IssueSecurityScheme: project.Spec.IssueSecurityScheme,
			PermissionScheme:    project.Spec.PermissionScheme,
			NotificationScheme:  project.Spec.NotificationScheme,
			CategoryId:          project.Spec.CategoryId,
		},
	}
}

// CreateProject creates and submits a Project object to the kubernetes server
func (t *TestUtil) CreateProject(project jiraservicedeskv1alpha1.Project, namespace string) *jiraservicedeskv1alpha1.Project {
	projectObject := t.CreateProjectObject(project, namespace)
	err := t.k8sClient.Create(t.ctx, projectObject)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: project.Spec.Name, Namespace: namespace}}

	_, err = t.r.Reconcile(req)
	if err != nil {
		ginkgo.Fail(err.Error())
	}

	return projectObject
}

// GetProject fetches a project object from kubernetes
func (t *TestUtil) GetProject(name string, namespace string) *jiraservicedeskv1alpha1.Project {
	projectObject := &jiraservicedeskv1alpha1.Project{}
	err := t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, projectObject)

	if err != nil {
		ginko.Fail(err.Error())
	}

	return projectObject
}

// DeleteProject deletes the project resource
func (t *TestUtil) DeleteProject(name string, namespace string) {
	projectObject := &jiraservicedeskv1alpha1.Project{}
	err := t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, projectObject)
	if err != nil {
		ginko.Fail(err.Error())
	}
	err = t.k8sClient.Delete(t.ctx, projectObject)
	if err != nil {
		ginko.Fail(err.Error())
	}
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}
	_, err = t.r.Reconcile(req)
	if err != nil {
		ginko.Fail(err.Error())
	}
}

// TryDeleteProject - Tries to delete Project if it exists, does not fail on any error
func (t *TestUtil) TryDeleteProject(name string, namespace string) {
	projectObject := &jiraservicedeskv1alpha1.Project{}
	_ = t.k8sClient.Get(t.ctx, types.NamespacedName{Name: name, Namespace: namespace}, projectObject)
	_ = t.k8sClient.Delete(t.ctx, projectObject)
	req := reconcile.Request{NamespacedName: types.NamespacedName{Name: name, Namespace: namespace}}
	_, _ = t.r.Reconcile(req)
}
