# Config

The app service currently defaults to Google Cloud Run as a provider. In future it will support more. 
This requires setting up a config key `app` with service accounts and related info. See handler/google.go 
to understand what is required. Additionally we require a dummy service account with no permissions 
for the deployed apps.
