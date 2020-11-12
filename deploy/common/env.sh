#!/usr/bin/env bash
set -e

export CLUSTER_API_ENDPOINT="${CLUSTER_API_ENDPOINT:-central.stackrox:443}"
echo "In-cluster Central endpoint set to $CLUSTER_API_ENDPOINT"

export COLLECTION_METHOD="${COLLECTION_METHOD:-${RUNTIME_SUPPORT:-ebpf}}"
echo "COLLECTION_METHOD set to $COLLECTION_METHOD"

export SCANNER_SUPPORT=${SCANNER_SUPPORT:-true}
echo "SCANNER_SUPPORT set to $SCANNER_SUPPORT"

export OFFLINE_MODE=${OFFLINE_MODE:-false}
echo "OFFLINE_MODE set to $OFFLINE_MODE"

export ROX_HTPASSWD_AUTH=${ROX_HTPASSWD_AUTH:-true}
echo "ROX_HTPASSWD_AUTH set to $ROX_HTPASSWD_AUTH"

echo "MONITORING_SUPPORT set to ${MONITORING_SUPPORT}"

export CLUSTER=${CLUSTER:-remote}
echo "CLUSTER set to $CLUSTER"

export STORAGE="${STORAGE:-none}"
echo "STORAGE set to ${STORAGE}"

export STORAGE_SIZE="${STORAGE_SIZE:-10}"
echo "STORAGE_SIZE set to ${STORAGE_SIZE}"

export OUTPUT_FORMAT="${OUTPUT_FORMAT:-kubectl}"
echo "OUTPUT_FORMAT set to ${OUTPUT_FORMAT}"

export LOAD_BALANCER="${LOAD_BALANCER:-none}"
echo "LOAD_BALANCER set to ${LOAD_BALANCER}"

export MONITORING_LOAD_BALANCER="${MONITORING_LOAD_BALANCER:-none}"
echo "MONITORING_LOAD_BALANCER set to ${MONITORING_LOAD_BALANCER}"

export ADMISSION_CONTROLLER="${ADMISSION_CONTROLLER:-false}"
echo "ADMISSION_CONTROLLER set to ${ADMISSION_CONTROLLER}"

export ADMISSION_CONTROLLER_UPDATES="${ADMISSION_CONTROLLER_UPDATES:-false}"
echo "ADMISSION_CONTROLLER_UPDATES set to ${ADMISSION_CONTROLLER_UPDATES}"

export ROX_NETWORK_ACCESS_LOG=${ROX_NETWORK_ACCESS_LOG:-false}
echo "ROX_NETWORK_ACCESS_LOG set to $ROX_NETWORK_ACCESS_LOG"

export ROX_DEVELOPMENT_BUILD=true
echo "ROX_DEVELOPMENT_BUILD is set to ${ROX_DEVELOPMENT_BUILD}"

export API_ENDPOINT="${API_ENDPOINT:-localhost:8000}"
echo "API_ENDPOINT is set to ${API_ENDPOINT}"

export AUTH0_SUPPORT="${AUTH0_SUPPORT:-true}"
echo "AUTH0_SUPPORT is set to ${AUTH0_SUPPORT}"

export TRUSTED_CA_FILE="${TRUSTED_CA_FILE:-}"
if [[ -n "${TRUSTED_CA_FILE}" ]]; then
  [[ -f "${TRUSTED_CA_FILE}" ]] || { echo "Trusted CA file ${TRUSTED_CA_FILE} not found"; return 1; }
  echo "TRUSTED_CA_FILE is set to ${TRUSTED_CA_FILE}"
else
  echo "No TRUSTED_CA_FILE provided"
fi

export ROX_DEFAULT_TLS_CERT_FILE="${ROX_DEFAULT_TLS_CERT_FILE:-}"
export ROX_DEFAULT_TLS_KEY_FILE="${ROX_DEFAULT_TLS_KEY_FILE:-}"

if [[ -n "$ROX_DEFAULT_TLS_CERT_FILE" ]]; then
	[[ -f "$ROX_DEFAULT_TLS_CERT_FILE" ]] || { echo "Default TLS certificate ${ROX_DEFAULT_TLS_CERT_FILE} not found"; return 1; }
	[[ -f "$ROX_DEFAULT_TLS_KEY_FILE" ]] || { echo "Default TLS key ${ROX_DEFAULT_TLS_KEY_FILE} not found"; return 1; }
	echo "Using default TLS certificate/key material from $ROX_DEFAULT_TLS_CERT_FILE, $ROX_DEFAULT_TLS_KEY_FILE"
elif [[ -n "$ROX_DEFAULT_TLS_KEY_FILE" ]]; then
	echo "ROX_DEFAULT_TLS_KEY_FILE is nonempty, but ROX_DEFAULT_TLS_CERT_FILE is"
	return 1
else
	echo "No default TLS certificates provided"
fi
