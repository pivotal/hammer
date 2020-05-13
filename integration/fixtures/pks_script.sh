prods="$(OM_CLIENT_ID='fakeClientID' OM_CLIENT_SECRET='fakeClientSecret' OM_USERNAME='pivotalcf' OM_PASSWORD='fakePassword' om -t https://pcf.manatee.cf-app.com -k curl -s -p /api/v0/staged/products)"
guid="$(echo "$prods" | jq -r '.[] | select(.type == "pivotal-container-service") | .guid')"
creds="$(OM_CLIENT_ID='fakeClientID' OM_CLIENT_SECRET='fakeClientSecret' OM_USERNAME='pivotalcf' OM_PASSWORD='fakePassword' om -t https://pcf.manatee.cf-app.com -k curl -s -p /api/v0/deployed/products/"$guid"/credentials/.properties.uaa_admin_password)"
pass="$(echo "$creds" | jq -r .credential.value.secret)"
pks login -a https://api.pks.manatee.cf-app.com -u admin -p "$pass" --skip-ssl-validation
