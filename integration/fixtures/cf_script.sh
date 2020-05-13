prods="$(OM_CLIENT_ID='fakeClientID' OM_CLIENT_SECRET='fakeClientSecret' OM_USERNAME='pivotalcf' OM_PASSWORD='fakePassword' om -t https://pcf.manatee.cf-app.com -k curl -s -p /api/v0/staged/products)"
guid="$(echo "$prods" | jq -r '.[] | select(.type == "cf") | .guid')"
creds="$(OM_CLIENT_ID='fakeClientID' OM_CLIENT_SECRET='fakeClientSecret' OM_USERNAME='pivotalcf' OM_PASSWORD='fakePassword' om -t https://pcf.manatee.cf-app.com -k curl -s -p /api/v0/deployed/products/"$guid"/credentials/.uaa.admin_credentials)"
user="$(echo "$creds" | jq -r .credential.value.identity)"
pass="$(echo "$creds" | jq -r .credential.value.password)"
cf login -a "api.sys.manatee.cf-app.com" -u "$user" -p "$pass" --skip-ssl-validation
