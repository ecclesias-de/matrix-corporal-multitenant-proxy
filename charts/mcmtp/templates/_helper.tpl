{{/*
Basename for chart resources
{{ include "basename" . }}
*/}}
{{- define "basename" -}}
{{ if hasSuffix .Chart.Name .Release.Name  }}
{{- .Release.Name -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
Kubernetes standard labels
labels: {{ include "labels" . | nindent 4 }}
*/}}
{{- define "labels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
app.metaways.net/team: {{ .Values.metadata.team }}
app.metaways.net/name: {{ .Values.metadata.name }}
{{- end -}}

{{/*
Labels used to match standard labels
selector: {{- include "matchLabels" . | nindent 4 }}
*/}}
{{- define "matchLabels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Secret name
secretName: {{ template "secret" . }}
*/}}
{{- define "secret" -}}
{{- if .Values.secrets.existingSecret -}}
{{ .Values.secrets.existingSecret }}
{{- else -}}
{{ include "basename" . }}
{{- end -}}
{{- end -}}

{{/*
Pvc name
claimName: {{ template "pvc" . }}
*/}}
{{- define "pvc" -}}
{{- if .Values.persistence.existingPvc -}}
{{ .Values.persistence.existingPvc }}
{{- else -}}
{{ include "basename" . }}
{{- end -}}
{{- end -}}