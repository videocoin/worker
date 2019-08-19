#!/bin/bash

readonly CHART_NAME=transcoder
readonly CHART_DIR=./helm/transcoder

CONSUL_ADDR=${CONSUL_ADDR:=127.0.0.1:8500}
ENV=${ENV:=dev}
DOCKER_REGISTRY=us.gcr.io
VERSION=$(git describe --abbrev=0)-$(git rev-parse --short HEAD)
PROJECT=${PROJECT:=`gcloud config list --format 'value(core.project)' 2>/dev/null`}

function log {
  local readonly level="$1"
  local readonly message="$2"
  local readonly timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  >&2 echo -e "${timestamp} [${level}] [$SCRIPT_NAME] ${message}"
}

function log_info {
  local readonly message="$1"
  log "INFO" "$message"
}

function log_warn {
  local readonly message="$1"
  log "WARN" "$message"
}

function log_error {
  local readonly message="$1"
  log "ERROR" "$message"
}

function update_deps() {
    log_info "Syncing dependencies..."
    helm dependencies update --kube-context ${KUBE_CONTEXT} ${CHART_DIR}
}

function has_jq {
  [ -n "$(command -v jq)" ]
}

function has_consul {
  [ -n "$(command -v consul)" ]
}

function has_helm {
  [ -n "$(command -v helm)" ]
}

function get_vars() {
    log_info "Getting variables..."
    readonly KUBE_CONTEXT=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/common/kube_context`
    readonly BASE_STREAM_URL=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/vars/baseStreamUrl`
    readonly BASE_STORAGE_URL=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/vars/baseStorageUrl`
    readonly SMCA=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/vars/smca`
    readonly BUCKET=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/vars/bucket`
    readonly VERIFIER_RPC_ADDR=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/vars/verifierRpcAddr`
    readonly LOG_LEVEL=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/vars/logLevel` 

    readonly SECRET=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/secrets/secret`
    readonly CONSUL_URL=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/secrets/consulUrl`
    readonly CONSUL_USERNAME=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/secrets/consulUsername`
    readonly CONSUL_PASSWORD=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/secrets/consulPassword`
    readonly BLOCKCHAIN_URL=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/secrets/blockchainUrl`
    readonly KEY=`consul kv get -http-addr=${CONSUL_ADDR} config/${ENV}/services/${CHART_NAME}/secrets/key`

}

function deploy() {
    log_info "Deploying ${CHART_NAME} version ${VERSION}"
    helm upgrade \
        --kube-context "${KUBE_CONTEXT}" \
        --install \
        --set image.tag="${VERSION}" \
        --set image.repo="${DOCKER_REGISTRY}/${PROJECT}/${CHART_NAME}" \
        --set config.baseStreamUrl="${BASE_STREAM_URL}" \
        --set config.baseStorageUrl="${BASE_STORAGE_URL}" \
        --set config.smca="${SMCA}" \
        --set config.logLevel="${LOG_LEVEL}" \
        --set config.verifierRpcAddr="${VERIFIER_RPC_ADDR}" \
        --set config.bucket="${BUCKET}" \
        --set secrets.blockchainUrl="${BLOCKCHAIN_URL}" \
        --set secrets.secret="${SECRET}" \
        --set secrets.key="${KEY}" \
        --set secrets.consulPassword="${CONSUL_PASSWORD}" \
        --set secrets.consulUsername="${CONSUL_USERNAME}" \
        --set secrets.consulUrl="${CONSUL_URL}" \
        --wait ${CHART_NAME} ${CHART_DIR}
}

if ! $(has_jq); then
    log_error "Could not find jq"
    exit 1
fi

if ! $(has_consul); then
    log_error "Could not find consul"
    exit 1
fi

if ! $(has_helm); then
    log_error "Could not find helm"
    exit 1
fi

get_vars
update_deps
deploy

exit $?