package secret

import (
	"context"
	"encoding/json"

	admissionv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/openyurtio/pool-coordinator/utils/certs"
)

func CreateSecret(certset *certs.Certs, namespace string, secretname string) {
	config := ctrl.GetConfigOrDie()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Info(err.Error())
	}
	secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretname, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			secret = &corev1.Secret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Secret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      secretname,
				},
			}
		} else {
			klog.Error(err)
			return
		}
	}
	secret.Data = map[string][]byte{
		certs.CAKeyName:       certset.CAKey,
		certs.CACertName:      certset.CACert,
		certs.ServerKeyName:   certset.Key,
		certs.ServerKeyName2:  certset.Key,
		certs.ServerCertName:  certset.Cert,
		certs.ServerCertName2: certset.Cert,
	}
	_, err = clientset.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		klog.Error(err)
	}
}

func UpdateCABundle(caBundle []byte, webhookconfig string) error {
	config := ctrl.GetConfigOrDie()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Info(err.Error())
		return err
	}
	validatingConfig, err := clientset.AdmissionregistrationV1().ValidatingWebhookConfigurations().Get(context.TODO(), webhookconfig, metav1.GetOptions{})
	if err != nil {
		klog.Info(err.Error())
		return err
	}

	validatingTemplate, err := parseValidatingTemplate(validatingConfig)
	if err != nil {
		return err
	}

	for i := range validatingTemplate {
		wh := &validatingTemplate[i]
		wh.ClientConfig.CABundle = caBundle
	}

	_, err = clientset.AdmissionregistrationV1().ValidatingWebhookConfigurations().Update(context.TODO(), validatingConfig, metav1.UpdateOptions{})
	if err != nil {
		klog.Error(err)
	}

	return nil
}

func parseValidatingTemplate(validatingConfig *admissionv1.ValidatingWebhookConfiguration) ([]admissionv1.ValidatingWebhook, error) {
	if templateStr := validatingConfig.Annotations["template"]; len(templateStr) > 0 {
		var validatingWHs []admissionv1.ValidatingWebhook
		if err := json.Unmarshal([]byte(templateStr), &validatingWHs); err != nil {
			return nil, err
		}
		return validatingWHs, nil
	}

	templateBytes, err := json.Marshal(validatingConfig.Webhooks)
	if err != nil {
		return nil, err
	}
	if validatingConfig.Annotations == nil {
		validatingConfig.Annotations = make(map[string]string, 1)
	}
	validatingConfig.Annotations["template"] = string(templateBytes)
	return validatingConfig.Webhooks, nil
}
