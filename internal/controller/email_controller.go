package controllers

import (
	"context"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	emailsv1alpha1 "github.com/yourusername/email-operator/api/v1alpha1"
)

// EmailReconciler reconciles an Email object
type EmailReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=emails.example.com,resources=emails,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=emails.example.com,resources=emails/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=emails.example.com,resources=emails/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=secrets,verbs=get;list;watch

func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Fetch the Email instance
	email := &emailsv1alpha1.Email{}
	err := r.Get(ctx, req.NamespacedName, email)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return. Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Fetch the related EmailSenderConfig
	emailSenderConfig := &emailsv1alpha1.EmailSenderConfig{}
	err = r.Get(ctx, client.ObjectKey{Name: email.Spec.SenderConfigRef, Namespace: req.Namespace}, emailSenderConfig)
	if err != nil {
		r.Log.Error(err, "Failed to get EmailSenderConfig")
		return ctrl.Result{}, err
	}

	// Get the API token from the secret
	secret := &corev1.Secret{}
	err = r.Get(ctx, client.ObjectKey{Name: emailSenderConfig.Spec.ApiTokenSecretRef, Namespace: req.Namespace}, secret)
	if err != nil {
		r.Log.Error(err, "Failed to get Secret")
		return ctrl.Result{}, err
	}
	apiToken := string(secret.Data["apiToken"])

	// Send email using MailerSend or Mailgun
	if err := sendEmail(apiToken, emailSenderConfig.Spec.SenderEmail, email.Spec.RecipientEmail, email.Spec.Subject, email.Spec.Body); err != nil {
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = err.Error()
	} else {
		email.Status.DeliveryStatus = "Sent"
		email.Status.MessageId = "some-message-id" // Capture real message ID from provider
	}

	if err := r.Status().Update(ctx, email); err != nil {
		r.Log.Error(err, "Failed to update Email status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func sendEmail(apiToken, senderEmail, recipientEmail, subject, body string) error {
	// Implement the logic to send email using MailerSend or Mailgun APIs
	// Use apiToken for authentication
	// Return error if the email sending fails
	return nil
}

func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&emailsv1alpha1.Email{}).
		Complete(r)
}
