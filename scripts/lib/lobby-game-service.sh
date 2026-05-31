# Lobby LOBBY_GAME_SERVICE_TOKEN (sent to games as lobby.serviceToken on provision).
# Source from deploy scripts after NAMESPACE is set.

apply_lobby_game_service_secret() {
  local secrets_file="${1:-k8s/secrets/lobby-game-service.yaml}"
  if [ ! -f "$secrets_file" ]; then
    echo "Game service token: no $secrets_file — provision/player lookup auth disabled"
    return 0
  fi

  echo "Applying $secrets_file (namespace ${NAMESPACE})..."
  kubectl apply -f "$secrets_file"

  if kubectl get secret lobby-game-service -n "$NAMESPACE" >/dev/null 2>&1; then
    kubectl set env deployment/lobby-backend -n "$NAMESPACE" \
      --from=secret/lobby-game-service \
      --containers=backend
  fi
}

run_patch_rps_handoff_urls_job() {
  local ns="${1:-joinquest}"
  kubectl delete job lobby-patch-rps-handoff-urls -n "$ns" --ignore-not-found
  sed 's/namespace: playhub/namespace: '"$ns"'/g' k8s/jobs/patch-rps-handoff-urls.yaml | kubectl apply -f -
  kubectl wait --for=condition=complete --timeout=120s "job/lobby-patch-rps-handoff-urls" -n "$ns"
  kubectl logs "job/lobby-patch-rps-handoff-urls" -n "$ns"
}
