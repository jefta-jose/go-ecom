{{/*
    Returns chart name from Chart.yaml or .Values.nameOverride if set
*/}}
{{- define "postgres-db.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}
