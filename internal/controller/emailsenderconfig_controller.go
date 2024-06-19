package controllers

import (
	"context"
	"github.com/go-logr/logr"
	emailsv1alpha1 "github.com/yourusername/email-operator/api/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EmailSenderConfigReconciler reconciles an EmailSenderConfig object
type EmailSenderConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=emails.example.com,resources=emailsenderconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=emails.example.com,resources=emailsenderconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=emails.example.com,resources=emailsenderconfigs/finalizers,verbs=update

func (r *EmailSenderConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("emailsenderconfig", req.NamespacedName)

	// Fetch the EmailSenderConfig instance
	var emailSenderConfig emailsv1alpha1.EmailSenderConfig
	if err := r.Get(ctx, req.NamespacedName, &emailSenderConfig); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("EmailSenderConfig resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get EmailSenderConfig")
		return ctrl.Result{}, err
	}

	// Log creation or update
	log.Info("Reconciling EmailSenderConfig", "EmailSenderConfig", emailSenderConfig)

	// Add logic to confirm the email sending settings (optional)

	return ctrl.Result{}, nil
}

func (r *EmailSenderConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&emailsv1alpha1.EmailSenderConfig{}).
		Complete(r)
}
