#!/bin/bash
# Deployment helper script for payment-service
# This script helps deploy the payment service to remote k3s cluster

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="payment-service"
K8S_HOST="root@37.27.40.86"

echo -e "${GREEN}Payment Service Deployment Helper${NC}"
echo "===================================="
echo "Remote Host: $K8S_HOST"
echo "Namespace: $NAMESPACE"
echo ""

# Function to print colored messages
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to run kubectl commands on remote host
remote_kubectl() {
    ssh $K8S_HOST "kubectl $@"
}

# Check if SSH is configured
if ! ssh -q $K8S_HOST exit; then
    print_error "Cannot connect to $K8S_HOST"
    print_info "Please ensure SSH access is configured"
    exit 1
fi

# Menu
echo "Select an option:"
echo "1) Create namespace"
echo "2) Create image pull secret"
echo "3) Create application secrets (interactive)"
echo "4) Check deployment status"
echo "5) View logs"
echo "6) Test health endpoint"
echo "7) View ArgoCD application"
echo "8) Manual sync in ArgoCD"
echo "9) SSH into cluster"
echo ""
read -p "Enter option (1-9): " option

case $option in
    1)
        print_info "Creating namespace on remote cluster..."
        remote_kubectl create namespace $NAMESPACE --dry-run=client -o yaml | remote_kubectl apply -f -
        print_info "Namespace created successfully"
        ;;

    2)
        print_info "Creating image pull secret..."
        read -p "Enter GitHub username (frallan97): " gh_user
        gh_user=${gh_user:-frallan97}
        read -sp "Enter GitHub Personal Access Token: " gh_token
        echo ""

        remote_kubectl create secret docker-registry ghcr-pull-secret \
            --docker-server=ghcr.io \
            --docker-username=$gh_user \
            --docker-password=$gh_token \
            --namespace=$NAMESPACE \
            --dry-run=client -o yaml | remote_kubectl apply -f -

        print_info "Image pull secret created"
        ;;

    3)
        print_info "Creating application secrets..."
        print_warning "You will be prompted for sensitive values"
        echo ""

        read -p "Stripe API Key (sk_test_... or sk_live_...): " stripe_key
        read -p "Stripe Webhook Secret (whsec_...): " stripe_webhook
        read -p "Database Password: " db_password

        # Create secret on remote
        ssh $K8S_HOST "kubectl create secret generic payment-service-secrets \
            --namespace=$NAMESPACE \
            --from-literal=stripe-api-key='$stripe_key' \
            --from-literal=stripe-webhook-secret='$stripe_webhook' \
            --from-literal=database-url='postgresql://paymentuser:$db_password@payment-service-postgresql:5432/paymentdb?sslmode=disable' \
            --from-literal=redis-url='redis://payment-service-redis-master:6379' \
            --dry-run=client -o yaml | kubectl apply -f -"

        print_info "Application secrets created"
        ;;

    4)
        print_info "Checking deployment status on remote cluster..."
        echo ""

        echo "=== Pods ==="
        remote_kubectl get pods -n $NAMESPACE
        echo ""

        echo "=== Services ==="
        remote_kubectl get svc -n $NAMESPACE
        echo ""

        echo "=== Ingress ==="
        remote_kubectl get ingress -n $NAMESPACE
        echo ""

        echo "=== Recent Events ==="
        remote_kubectl get events -n $NAMESPACE --sort-by='.lastTimestamp' | tail -10
        ;;

    5)
        print_info "Fetching logs from remote cluster..."
        POD=$(remote_kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=payment-service -o jsonpath='{.items[0].metadata.name}')

        if [ -z "$POD" ]; then
            print_error "No payment-service pods found"
            exit 1
        fi

        print_info "Showing logs for pod: $POD"
        print_info "Press Ctrl+C to stop following logs"
        remote_kubectl logs -n $NAMESPACE $POD --tail=100 -f
        ;;

    6)
        print_info "Testing health endpoint..."
        HOST=$(remote_kubectl get ingress -n $NAMESPACE payment-service -o jsonpath='{.spec.rules[0].host}' 2>/dev/null || echo "")

        if [ -z "$HOST" ]; then
            print_warning "Ingress not found, checking if deployment exists..."
            remote_kubectl get deployment -n $NAMESPACE payment-service
        else
            print_info "Testing https://$HOST/health"
            curl -s https://$HOST/health
            echo ""

            print_info "Testing https://$HOST/metrics"
            curl -s https://$HOST/metrics | head -20
            echo "..."
        fi
        ;;

    7)
        print_info "Checking ArgoCD for payment-service application..."
        echo ""
        print_info "ArgoCD URL: https://argocd.vibeoholic.com"
        print_info "Looking for application: payment-service"
        echo ""

        # Try to get ArgoCD app status
        remote_kubectl get application payment-service -n argocd -o wide 2>/dev/null || \
            print_warning "Application not found in ArgoCD. Have you added it to app-of-apps.yaml?"
        ;;

    8)
        print_info "Manually syncing payment-service in ArgoCD..."
        remote_kubectl patch application payment-service -n argocd \
            --type merge \
            -p '{"operation":{"initiatedBy":{"username":"admin"},"sync":{"revision":"HEAD"}}}' 2>/dev/null || \
            print_error "Failed to sync. Use ArgoCD UI: https://argocd.vibeoholic.com"
        ;;

    9)
        print_info "Opening SSH session to $K8S_HOST"
        print_info "Useful commands:"
        echo "  kubectl get pods -n $NAMESPACE"
        echo "  kubectl logs -n $NAMESPACE -l app.kubernetes.io/name=payment-service"
        echo "  kubectl describe pod -n $NAMESPACE <pod-name>"
        echo ""
        ssh $K8S_HOST
        ;;

    *)
        print_error "Invalid option"
        exit 1
        ;;
esac
