prods="$(om -t https://pcf.manatee.cf-app.com -k -u pivotalcf -p fakePassword curl -s -p /api/v0/staged/products)"
guid="$(echo "$prods" | jq -r '.[] | select(.type == "cf") | .guid')"
creds="$(om -t https://pcf.manatee.cf-app.com -k -u pivotalcf -p fakePassword curl -s -p /api/v0/deployed/products/"$guid"/credentials/.uaa.admin_credentials)"
user="$(echo "$creds" | jq -r .credential.value.identity)"
pass="$(echo "$creds" | jq -r .credential.value.password)"
cf login -a "api.sys.manatee.cf-app.com" -u "$user" -p "$pass" --skip-ssl-validation
