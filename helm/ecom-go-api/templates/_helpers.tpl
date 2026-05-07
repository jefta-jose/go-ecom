{{/*
    =================================================================================
        1. define
            define is like creating a function or template snippet in Helm.
            You give it a name (e.g., "go-api-chart.name") and a block of template code.
            You can then reuse this block anywhere in your Helm chart.
                Example:
                {{- define "go-api-chart.name" -}}
                {{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
                {{- end }}
            This defines a template called "go-api-chart.name".
            It will return:
            .Values.nameOverride if the user set it, else
            .Chart.Name from Chart.yaml
            trunc 63 | trimSuffix "-" ensures it’s a valid Kubernetes name (≤63 chars, no trailing dash).
            Think of define as “here’s a mini-template I can call later”.

        2. include
            include is how you call a template defined with define.
            You give it the name of the template and the current context (.).
                Example:
                helm.sh/chart: {{ include "go-api-chart.name" . }}-{{ .Chart.Version }}
            Here, it calls "go-api-chart.name" and gets the chart name.
            . passes all the current values, like .Chart and .Values, into the template.
            Think of include as “run that mini-template and give me the result”.
    =================================================================================
*/}}



{{/*
    =================================================================================
        Returns chart name from Chart.yaml or .Values.nameOverride if set
    =================================================================================
*/}}

{{- define "go-api-chart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}



{{/*
    =========================================================================================
        Returns a unique fully qualified app name by combining release name with chart name.
    =========================================================================================
*/}}

{{- define "go-api-chart" -}}
{{- if .ValuesOverride }}
{{- .ValuesOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}



{{/*
    =================================================================================
    COMMON LABELS HELPER
    =================================================================================
        Generates standard Kubernetes labels that should be applied to ALL resources
        (Deployments, Services, ConfigMaps, etc.) for consistency and tracking.

        Labels generated:
        - helm.sh/chart: Chart name and version (e.g., "go-api-chart-0.1.0")
        - app.kubernetes.io/name: Application name
        - app.kubernetes.io/instance: Release name (identifies this specific installation)
        - app.kubernetes.io/managed-by: Always "Helm" (indicates this is Helm-managed)

        These follow Kubernetes recommended labels:
        https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
    =================================================================================
*/}}

{{- define "go-api-chart.labels" -}}
helm.sh/chart: {{ include "go-api-chart.name" . }}-{{ .Chart.Version | replace "+" "_" }}
{{ include "go-api-chart.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}



{{/*
    =================================================================================
    SELECTOR LABELS HELPER
    =================================================================================
        Generates IMMUTABLE labels used for matching Pods to Services.

        WARNING: These labels should NEVER change after initial deployment, otherwise
        Services won't be able to match Pods. Only use these in:
        - Deployment spec.selector.matchLabels
        - Service spec.selector
        - DO NOT include version or other changing values here

        Labels generated:
        - app.kubernetes.io/name: Application name (stable)
        - app.kubernetes.io/instance: Release name (stable)

        Usage in templates:
        Deployment spec.selector:
            selector:
            matchLabels:
                {{- include "go-api-chart.selectorLabels" . | nindent 8 }}
        
        Service spec.selector:
            selector:
            {{- include "go-api-chart.selectorLabels" . | nindent 4 }}
    =================================================================================
*/}}

{{- define "go-api-chart.selectorLabels" -}}
app.kubernetes.io/name: {{ include "go-api-chart.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}