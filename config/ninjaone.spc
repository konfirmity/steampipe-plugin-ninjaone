connection "ninjaone" {
  plugin = "ninjaone"

  # `make install` puts the binary under plugins/local/ninjaone. For a build
  # installed that way, change the line above to: plugin = "local/ninjaone"

  # Client ID of an API Services (machine-to-machine) client app.
  # Can also be set with the NINJAONE_CLIENT_ID environment variable.
  # A value here wins over the environment variable, so delete this line to use the variable.
  client_id = "example-client-id-not-real"

  # Client secret of the same client app. NinjaOne shows it once, when the app is created.
  # Can also be set with the NINJAONE_CLIENT_SECRET environment variable.
  # A value here wins over the environment variable, so delete this line to use the variable.
  client_secret = "example-client-secret-not-real"

  # Regional instance, including https://. Defaults to https://app.ninjarmm.com.
  # https://oc.ninjarmm.com is verified. NinjaOne also has eu, ca and us2 regions:
  # for those, copy the host from your browser's address bar while signed in to NinjaOne.
  # Can also be set with the NINJAONE_BASE_URL environment variable.
  # base_url = "https://oc.ninjarmm.com"

  # Limits every table to one NinjaOne organization. Write it as a quoted whole number.
  # A client app cannot be limited to one organization in NinjaOne, so a managed service
  # provider's credential sees every tenant. This setting is what limits the plugin.
  # Leave it unset and every table returns every organization the credential can see.
  # Can also be set with the NINJAONE_ORGANIZATION_ID environment variable.
  # organization_id = "12345"
}
