{{- define "neo4j.backup.checkIfSecretExistsOrNot" -}}
    {{- if (.Values.backup.secretName | trim) -}}
        {{- if (not .Values.disableLookups) -}}

            {{- include "neo4j.backup.checkIfSecretKeyNameExistsOrNot" . -}}
            {{- $secret := (lookup "v1" "Secret" .Release.Namespace .Values.backup.secretName) }}
            {{- $secretExists := $secret | all }}

            {{- if (not $secretExists) -}}
                {{ fail (printf "Secret %s configured in 'backup.secretName' not found" .Values.backup.secretName) }}
             {{- else if not (hasKey $secret.data .Values.backup.secretKeyName) -}}
                {{ fail (printf "Secret %s must contain key %s" .Values.backup.secretName .Values.backup.secretKeyName) }}
            {{- end -}}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{- define "neo4j.backup.checkAzureStorageAccountName" -}}
    {{- if eq .Values.backup.cloudProvider "azure" }}
        {{- if and (or (empty .Values.backup.secretName) (empty .Values.backup.secretKeyName)) (empty .Values.backup.azureStorageAccountName) -}}
            {{ fail (printf "Both secretName|secretKeyName and azureStorageAccountName key cannot be empty. Please set one of them via --set backup.secretName or --set backup.azureStorageAccountName") }}
        {{- end -}}

        {{- if and (or (.Values.backup.secretName) (.Values.backup.secretKeyName)) (.Values.backup.azureStorageAccountName) -}}
            {{ fail (printf "Both secretName|secretKeyName and azureStorageAccountName key cannot be present. Please set only one of them via --set backup.secretName or --set backup.azureStorageAccountName") }}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{/*check for secretKeyName existence only when secretName is provided*/}}
{{- define "neo4j.backup.checkIfSecretKeyNameExistsOrNot" -}}
   {{- if .Values.backup.secretName -}}
    {{- if kindIs "invalid" .Values.backup.secretKeyName -}}
        {{- fail (printf "Missing secretKeyName !!") -}}
    {{- else if (not (.Values.backup.secretKeyName | trim)) -}}
        {{- fail (printf "Empty secretKeyName") -}}
    {{- end -}}
   {{- end -}}

{{- end -}}

{{/* checks if serviceAccountName is provided or not  when secretName is missing */}}
{{- define "neo4j.backup.checkServiceAccountName" -}}
    {{- if and (empty .Values.serviceAccountName) (empty .Values.backup.secretName) (not (empty .Values.backup.cloudProvider)) -}}
        {{ fail (printf "Please provide either secretName or serviceAccountName. Both cannot be empty. Please set only one of them via --set backup.secretName or --set serviceAccountName") }}
    {{- end -}}
{{- end -}}

{{/* checks if serviceAccountName is provided or not  when secretName is missing */}}
{{- define "neo4j.backup.checkBucketName" -}}
    {{- if or (kindIs "invalid" .Values.backup.aggregate) (not .Values.backup.aggregate.enabled) -}}
        {{- if .Values.backup.cloudProvider -}}
            {{- if empty .Values.backup.bucketName -}}
                {{ fail (printf "Empty bucketName. Please set bucketName via --set backup.bucketName") }}
            {{- end -}}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{- define "neo4j.backup.checkDatabaseIPAndServiceName" -}}
    {{/* Skip validation for aggregate backup */}}
    {{- if and .Values.backup.aggregateBackup .Values.backup.aggregateBackup.enabled -}}
        {{/* Skip all validation for aggregate backup */}}
    {{- else if .Values.disableLookups -}}
        {{- if and (empty (.Values.backup.databaseAdminServiceName | trim)) (empty (.Values.backup.databaseAdminServiceIP | trim)) -}}
            {{- fail (printf "Empty fields. Please set databaseAdminServiceName") -}}
        {{- end -}}
    {{- else -}}
        {{- if and (empty (.Values.backup.databaseBackupEndpoints | trim)) (empty (.Values.backup.databaseAdminServiceName | trim)) (empty (.Values.backup.databaseAdminServiceIP | trim)) -}}
            {{- fail (printf "Empty fields. Please set either databaseBackupEndpoints or databaseAdminServiceName/databaseAdminServiceIP") -}}
        {{- end -}}
        {{- if and (.Values.backup.databaseBackupEndpoints) (or .Values.backup.databaseAdminServiceName .Values.backup.databaseAdminServiceIP) -}}
            {{- fail (printf "Cannot specify both databaseBackupEndpoints and databaseAdminServiceName/databaseAdminServiceIP") -}}
        {{- end -}}
    {{- end -}}
{{- end -}}
